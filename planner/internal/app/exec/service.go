package exec

import (
    "context"
    "errors"

    m "github.com/james/tasks-planner/internal/model"
)

// ErrInvalidCoordinator indicates the loaded coordinator contract is invalid
// (e.g., missing required fields like version). Use errors.Is with this
// sentinel in tests and callers.
var ErrInvalidCoordinator = errors.New("invalid coordinator: missing version")

// Service orchestrates executor initialization and runtime.
type Service struct {
	LoadCoordinator func(path string) (m.Coordinator, error)
	InitRuntime     func(ctx context.Context, coord m.Coordinator) error
	RunLoop         func(ctx context.Context) error
}

// ErrMissingAdapters indicates the executor Service is misconfigured (one or
// more adapter functions are nil). Use errors.Is with this sentinel in tests.
var ErrMissingAdapters = errors.New("exec service: missing adapters")

// Run loads the coordinator contract, initializes runtime components, then enters the execution loop.
func (s Service) Run(ctx context.Context, coordPath string) error {
    if s.LoadCoordinator == nil || s.InitRuntime == nil || s.RunLoop == nil {
        return ErrMissingAdapters
    }
    coord, err := s.LoadCoordinator(coordPath)
    if err != nil {
        return err
    }
    // Minimal validation: require a non-empty version to guard zero-value coordinators.
    if coord.Version == "" {
        return ErrInvalidCoordinator
    }
    if err := s.InitRuntime(ctx, coord); err != nil {
        return err
    }
    return s.RunLoop(ctx)
}
