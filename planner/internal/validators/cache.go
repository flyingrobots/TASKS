package validators

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type cacheStore struct {
	dir string
	mu  sync.RWMutex
}

func newCacheStore(dir string) (*cacheStore, error) {
	if dir == "" {
		return nil, fmt.Errorf("validators: empty cache directory")
	}
	info, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("validators: cache dir: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("validators: cache path is not a directory: %s", dir)
	}
	if err := ensureWritable(dir); err != nil {
		return nil, err
	}
	return &cacheStore{dir: dir}, nil
}

func ensureWritable(dir string) error {
	probe, err := os.CreateTemp(dir, ".validator-cache-*")
	if err != nil {
		return fmt.Errorf("validators: cache dir not writable: %w", err)
	}
	name := probe.Name()
	probe.Close()
	_ = os.Remove(name)
	return nil
}

func (c *cacheStore) cachePath(name string) string {
	return filepath.Join(c.dir, fmt.Sprintf("%s.json", name))
}

func (c *cacheStore) Load(name, key string) (Report, bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	path := c.cachePath(name)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Report{}, false, nil
		}
		return Report{}, false, fmt.Errorf("validators: read cache: %w", err)
	}
	entries := map[string]Report{}
	if err := json.Unmarshal(data, &entries); err != nil {
		return Report{}, false, fmt.Errorf("validators: corrupt cache file %s: %w", path, err)
	}
	rep, ok := entries[key]
	return rep, ok, nil
}

func (c *cacheStore) Store(name, key string, rep Report) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	path := c.cachePath(name)
	entries := map[string]Report{}
	if data, err := os.ReadFile(path); err == nil {
		if err := json.Unmarshal(data, &entries); err != nil {
			return fmt.Errorf("validators: corrupt cache file %s: %w", path, err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("validators: read cache %s: %w", path, err)
	}
	entries[key] = rep
	encoded, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("validators: encode cache: %w", err)
	}
	tmp, err := os.CreateTemp(c.dir, ".validator-cache-write-*")
	if err != nil {
		return fmt.Errorf("validators: temp cache: %w", err)
	}
	tmpPath := tmp.Name()
	if _, err := tmp.Write(encoded); err != nil {
		tmp.Close()
		_ = os.Remove(tmpPath)
		return fmt.Errorf("validators: write temp cache: %w", err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("validators: close temp cache: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("validators: rename cache: %w", err)
	}
	return nil
}
