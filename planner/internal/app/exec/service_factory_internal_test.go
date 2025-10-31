package exec

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestFilesystemCoordinatorLoaderReadOverride(t *testing.T) {
    // Provide a realistic coordinator.json payload to exercise decoding
    payload := []byte(`{
      "version": "v8",
      "graph": {
        "nodes": [
          {
            "id": "T001",
            "feature_id": "F001",
            "title": "Setup DB",
            "duration": {"optimistic": 1, "mostLikely": 2, "pessimistic": 3},
            "durationUnits": "hours",
            "acceptance_checks": [{"type": "command", "cmd": "echo ok", "timeout": 5}],
            "execution_logging": {"format": "JSONL", "required_fields": ["timestamp","task_id","step","status","message"]},
            "compensation": {"idempotent": true}
          }
        ],
        "edges": [
          {"from": "T001", "to": "T001", "type": "technical", "transitive": false}
        ]
      },
      "config": {
        "resources": {
          "catalog": {"db": {"capacity": 1, "mode": "exclusive", "lock_order": 10}},
          "profiles": {"default": {"db": 1}}
        },
        "policies": {
          "concurrency_max": 4,
          "lock_ordering": ["db"],
          "circuit_breaker_thresholds": {"task_fail_rate": 0.5}
        }
      },
      "metrics": {"estimates": {"p50_total_hours": 2.0, "longest_path_length": 1, "width_approx": 1}}
    }`)

    loader := FilesystemCoordinatorLoader{
        ReadFile: func(path string) ([]byte, error) { return payload, nil },
    }
    coord, err := loader.Load("custom.json")
    if err != nil {
        t.Fatalf("load: %v", err)
    }
    if coord.Version != "v8" {
        t.Fatalf("unexpected version: %s", coord.Version)
    }
    if len(coord.Graph.Nodes) != 1 || coord.Graph.Nodes[0].ID != "T001" {
        t.Fatalf("unexpected nodes: %+v", coord.Graph.Nodes)
    }
    if coord.Config.Resources.Catalog["db"].Capacity != 1 {
        t.Fatalf("unexpected resource catalog: %+v", coord.Config.Resources.Catalog)
    }
    if coord.Metrics.Estimates.P50TotalHours <= 0 {
        t.Fatalf("estimates not decoded: %+v", coord.Metrics.Estimates)
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
    // Write a complete-ish coordinator payload
    if err := os.WriteFile(path, []byte(`{"version":"v8","graph":{"nodes":[],"edges":[]}}`), 0o644); err != nil {
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

func TestFilesystemCoordinatorLoaderNegativeCases(t *testing.T) {
    // (1) syntactically invalid JSON
    loader := FilesystemCoordinatorLoader{ReadFile: func(string) ([]byte, error) { return []byte("{"), nil }}
    if _, err := loader.Load("x.json"); err == nil {
        t.Fatalf("expected decode error for invalid JSON")
    }

    // (2) structurally valid but missing fields (no version) — current behavior: no error, zero values
    loader = FilesystemCoordinatorLoader{ReadFile: func(string) ([]byte, error) { return []byte(`{"graph":{"nodes":[],"edges":[]}}`), nil }}
    c, err := loader.Load("y.json")
    if err != nil {
        t.Fatalf("unexpected error for missing fields: %v", err)
    }
    if c.Version != "" {
        t.Fatalf("expected empty version for missing field, got %q", c.Version)
    }

    // (3) wrong types — number for version -> decode error
    loader = FilesystemCoordinatorLoader{ReadFile: func(string) ([]byte, error) { return []byte(`{"version": 123}`), nil }}
    if _, err := loader.Load("z.json"); err == nil {
        t.Fatalf("expected type error for version number")
    }

    // (4) empty file — decode error
    loader = FilesystemCoordinatorLoader{ReadFile: func(string) ([]byte, error) { return []byte(""), nil }}
    if _, err := loader.Load("empty.json"); err == nil {
        t.Fatalf("expected error for empty payload")
    }
}

func TestFilesystemCoordinatorLoaderReadFallbackError(t *testing.T) {
	loader := FilesystemCoordinatorLoader{}
	if _, err := loader.Load("missing.json"); err == nil {
		t.Fatalf("expected error for missing coordinator")
	}
}
