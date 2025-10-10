package exec

import (
	"context"
	"errors"

	m "github.com/james/tasks-planner/internal/model"
)

// Service orchestrates executor initialization and runtime.
type Service struct {
	LoadCoordinator func(path string) (m.Coordinator, error)
	InitRuntime     func(ctx context.Context, coord m.Coordinator) error
	RunLoop         func(ctx context.Context) error
}

// Run loads the coordinator contract, initializes runtime components, then enters the execution loop.
func (s Service) Run(ctx context.Context, coordPath string) error {
	if s.LoadCoordinator == nil || s.InitRuntime == nil || s.RunLoop == nil {
		return errors.New("exec service: missing adapters")
	}
	coord, err := s.LoadCoordinator(coordPath)
	if err != nil {
		return err
	}
	if err := s.InitRuntime(ctx, coord); err != nil {
		return err
	}
	return s.RunLoop(ctx)
}
