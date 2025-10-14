package plan

import (
	"context"
	"testing"

	m "github.com/james/tasks-planner/internal/model"
)

func TestDefaultWaveBuilderBuildsWaves(t *testing.T) {
	builder := DefaultWaveBuilder{}
	df := &m.DagFile{}
	df.Meta.Version = "v8"
	df.Nodes = []struct {
		ID                  string `json:"id"`
		Depth               int    `json:"depth"`
		CriticalPath        bool   `json:"critical_path"`
		ParallelOpportunity int    `json:"parallel_opportunity"`
	}{{ID: "T001", Depth: 0}}
	df.Edges = []struct {
		From       string `json:"from"`
		To         string `json:"to"`
		Type       string `json:"type"`
		Transitive bool   `json:"transitive"`
	}{}

	tasks := []m.Task{{ID: "T001"}}
	waves, err := builder.Build(context.Background(), df, tasks)
	if err != nil {
		t.Fatalf("build waves: %v", err)
	}
	if waves.Meta.Version != schemaVersion {
		t.Fatalf("expected meta version %s, got %+v", schemaVersion, waves.Meta)
	}
	if len(waves.Waves) != 1 {
		t.Fatalf("expected waves array, got %+v", waves.Waves)
	}
}
