package plan

import (
	"context"

	m "github.com/james/tasks-planner/internal/model"
	wavesim "github.com/james/tasks-planner/internal/planner/wavesim"
)

// WaveBuilder abstracts wave simulation.
type WaveBuilder interface {
	Build(ctx context.Context, df m.DagFile, tasks []m.Task) (map[string]any, error)
}

// DefaultWaveBuilder uses the wavesim package to layer tasks.
type DefaultWaveBuilder struct{}

func (DefaultWaveBuilder) Build(ctx context.Context, df m.DagFile, tasks []m.Task) (map[string]any, error) {
	waves := map[string]any{"meta": map[string]any{"version": "v8", "planId": "", "artifact_hash": ""}}
	ids, err := wavesim.Generate(df, tasks)
	if err != nil {
		return nil, err
	}
	waves["waves"] = ids
	return waves, nil
}
