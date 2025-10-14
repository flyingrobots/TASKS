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
	df.Nodes = []m.DagNode{{ID: "T001", Depth: 0}}
	df.Edges = []m.DagEdge{}

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
