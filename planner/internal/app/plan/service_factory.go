package plan

import (
	"context"

	m "github.com/james/tasks-planner/internal/model"
	dagbuild "github.com/james/tasks-planner/internal/planner/dag"
	"github.com/james/tasks-planner/internal/validate"
	"github.com/james/tasks-planner/internal/validators"
)

// NewDefaultService returns a Service wired with default adapters.
func NewDefaultService() Service {
	docLoader := NewMarkdownDocLoader()
	analyzer := CensusAnalyzer{}
	deps := DefaultDependencyResolver{}
	waves := DefaultWaveBuilder{}
	coord := DefaultCoordinatorBuilder{}
	artifacts := FileArtifactWriter{}

	return Service{
		BuildTasks:  docLoader.Load,
		AnalyzeRepo: analyzer.Analyze,
		ResolveDeps: deps.Resolve,
		BuildDAG: func(_ context.Context, tasks []m.Task, edges []m.Edge, minConfidence float64) (*m.DagFile, error) {
			return dagbuild.Build(tasks, edges, minConfidence)
		},
		BuildCoordinator: coord.Build,
		ValidateTasks:    validate.TasksFile,
		ValidateDAG:      validate.DagFile,
		BuildWaves:       waves.Build,
		WriteArtifacts:   artifacts.Write,
		NewValidatorRunner: func(cfg validators.Config) (ValidatorRunner, error) {
			return validators.NewRunner(cfg)
		},
	}
}
