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
      "version": "v9",
      "meta": {"version":"v9", "artifactHash":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
      "graph": {
        "nodes": [
          {
            "id": "T001",
            "featureID": "F001",
            "title": "Setup DB",
            "duration": {"optimistic": 1, "mostLikely": 2, "pessimistic": 3},
            "durationUnits": "hours",
            "acceptanceChecks": [{"type": "command", "cmd": "echo ok", "timeout": 5}],
            "executionLogging": {"format": "JSONL", "requiredFields": ["timestamp","task_id","step","status","message"]},
            "compensation": {"idempotent": true}
          }
        ],
        "edges": [
          {"from": "T001", "to": "T001", "type": "technical", "transitive": false}
        ]
      },
      "config": {
        "resources": {
          "catalog": {"db": {"capacity": 1, "mode": "exclusive", "lockOrder": 10}},
          "profiles": {"default": {"db": 1}}
        },
        "policies": {
          "concurrencyMax": 4,
          "lockOrdering": ["db"],
          "circuitBreakerThresholds": {"task_fail_rate": 0.5}
        }
      },
      "metrics": {"estimates": {"p50TotalHours": 2.0, "longestPathLength": 1, "widthApprox": 1}}
    }`)

    loader := FilesystemCoordinatorLoader{
        ReadFile: func(path string) ([]byte, error) { return payload, nil },
    }
    coord, err := loader.Load("custom.json")
    if err != nil {
        t.Fatalf("load: %v", err)
    }
    if coord.Version == "" {
        t.Fatalf("expected non-empty version, got empty")
    }
    if len(coord.Graph.Nodes) != 1 || coord.Graph.Nodes[0].ID != "T001" {
        t.Fatalf("unexpected nodes: %+v", coord.Graph.Nodes)
    }
    if coord.Config.Resources.Catalog["db"].Capacity != 1 || coord.Config.Resources.Catalog["db"].LockOrder != 10 {
        t.Fatalf("unexpected resource catalog: %+v", coord.Config.Resources.Catalog)
    }
    if len(coord.Graph.Nodes[0].AcceptanceChecks) != 1 {
        t.Fatalf("expected acceptance checks, got %+v", coord.Graph.Nodes[0].AcceptanceChecks)
    }
    if len(coord.Graph.Nodes[0].ExecutionLogging.RequiredFields) == 0 || coord.Graph.Nodes[0].FeatureID != "F001" {
        t.Fatalf("unexpected logging fields or featureID: %+v / %s", coord.Graph.Nodes[0].ExecutionLogging, coord.Graph.Nodes[0].FeatureID)
    }
    if coord.Metrics.Estimates.P50TotalHours != 2.0 || coord.Metrics.Estimates.LongestPathLength != 1 || coord.Metrics.Estimates.WidthApprox != 1 {
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
    // Write a minimal valid coordinator payload
    mini := `{"version":"v9","meta":{"version":"v9","artifactHash":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"graph":{"nodes":[],"edges":[]},"config":{"resources":{"catalog":{},"profiles":{"default":{}}},"policies":{"circuitBreakerThresholds":null,"concurrencyMax":0,"lockOrdering":[]}}}`
    if err := os.WriteFile(path, []byte(mini), 0o644); err != nil {
        t.Fatalf("write coord: %v", err)
    }
    loader := FilesystemCoordinatorLoader{}
    coord, err := loader.Load(path)
    if err != nil {
        t.Fatalf("load fallback: %v", err)
    }
    if coord.Version == "" {
        t.Fatalf("unexpected empty version in coordinator: %+v", coord)
    }
}

func TestFilesystemCoordinatorLoaderNegativeCases(t *testing.T) {
    // (1) syntactically invalid JSON
    loader := FilesystemCoordinatorLoader{ReadFile: func(string) ([]byte, error) { return []byte("{"), nil }}
    if _, err := loader.Load("x.json"); err == nil {
        t.Fatalf("expected decode error for invalid JSON")
    }

    // (2) structurally valid but missing required fields -> schema error now
    loader = FilesystemCoordinatorLoader{ReadFile: func(string) ([]byte, error) { return []byte(`{"graph":{"nodes":[],"edges":[]}}`), nil }}
    if _, err := loader.Load("y.json"); err == nil {
        t.Fatalf("expected schema error for missing version")
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
