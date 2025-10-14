package plan

import (
	"context"

	m "github.com/james/tasks-planner/internal/model"
	wavesim "github.com/james/tasks-planner/internal/planner/wavesim"
)

// WaveBuilder abstracts wave simulation.
type WaveBuilder interface {
	Build(ctx context.Context, df *m.DagFile, tasks []m.Task) (*m.WavesArtifact, error)
}

// DefaultWaveBuilder uses the wavesim package to layer tasks.
type DefaultWaveBuilder struct{}

func (DefaultWaveBuilder) Build(ctx context.Context, df *m.DagFile, tasks []m.Task) (*m.WavesArtifact, error) {
	ids, err := wavesim.Generate(*df, tasks)
	if err != nil {
		return nil, err
	}
	return &m.WavesArtifact{
		Meta:  m.WavesMeta{Version: schemaVersion, PlanID: "", ArtifactHash: ""},
		Waves: ids,
	}, nil
}
