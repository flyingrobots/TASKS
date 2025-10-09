package validators

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type cacheStore struct {
	dir string
}

func newCacheStore(dir string) (*cacheStore, error) {
	return &cacheStore{dir: dir}, nil
}

func (c *cacheStore) cachePath(name string) string {
	return filepath.Join(c.dir, fmt.Sprintf("%s.json", name))
}

func (c *cacheStore) Load(name, key string) (Report, bool) {
	path := c.cachePath(name)
	data, err := os.ReadFile(path)
	if err != nil {
		return Report{}, false
	}
	var entries map[string]Report
	if err := json.Unmarshal(data, &entries); err != nil {
		return Report{}, false
	}
	rep, ok := entries[key]
	return rep, ok
}

func (c *cacheStore) Store(name, key string, rep Report) {
	path := c.cachePath(name)
	entries := map[string]Report{}
	if data, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(data, &entries)
	}
	entries[key] = rep
	enc, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(path, enc, 0o644)
}
