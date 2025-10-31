package exec

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestFilesystemCoordinatorLoaderReadOverride(t *testing.T) {
	loader := FilesystemCoordinatorLoader{
		ReadFile: func(path string) ([]byte, error) {
			if path != "custom.json" {
				t.Fatalf("unexpected path: %s", path)
			}
			return []byte(`{"version":"v8"}`), nil
		},
	}
	coord, err := loader.Load("custom.json")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if coord.Version != "v8" {
		t.Fatalf("unexpected coordinator: %+v", coord)
	}
}

func TestFilesystemCoordinatorLoaderReadOverrideError(t *testing.T) {
	expected := errors.New("boom")
	loader := FilesystemCoordinatorLoader{
		ReadFile: func(string) ([]byte, error) {
			return nil, expected
		},
	}
	if _, err := loader.Load("anything.json"); !errors.Is(err, expected) {
		t.Fatalf("expected override error, got %v", err)
	}
}

func TestFilesystemCoordinatorLoaderReadFallback(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "coord.json")
	if err := os.WriteFile(path, []byte(`{"version":"v8"}`), 0o644); err != nil {
		t.Fatalf("write coord: %v", err)
	}

	loader := FilesystemCoordinatorLoader{}
	coord, err := loader.Load(path)
	if err != nil {
		t.Fatalf("load fallback: %v", err)
	}
	if coord.Version != "v8" {
		t.Fatalf("unexpected coordinator: %+v", coord)
	}
}

func TestFilesystemCoordinatorLoaderReadFallbackError(t *testing.T) {
	loader := FilesystemCoordinatorLoader{}
	if _, err := loader.Load("missing.json"); err == nil {
		t.Fatalf("expected error for missing coordinator")
	}
}
