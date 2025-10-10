package exec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	m "github.com/james/tasks-planner/internal/model"
)

// ErrLoopNotImplemented is returned by the default loop placeholder.
var ErrLoopNotImplemented = errors.New("executor loop not implemented")

// FilesystemCoordinatorLoader reads coordinator contracts from disk.
type FilesystemCoordinatorLoader struct {
	ReadFile func(string) ([]byte, error)
}

// Load loads and decodes a coordinator artifact.
func (l FilesystemCoordinatorLoader) Load(path string) (m.Coordinator, error) {
	bs, err := l.read(path)
	if err != nil {
		return m.Coordinator{}, fmt.Errorf("read coordinator: %w", err)
	}
	var coord m.Coordinator
	if err := json.Unmarshal(bs, &coord); err != nil {
		return m.Coordinator{}, fmt.Errorf("decode coordinator: %w", err)
	}
	return coord, nil
}

func (l FilesystemCoordinatorLoader) read(path string) ([]byte, error) {
	if l.ReadFile != nil {
		return l.ReadFile(path)
	}
	return os.ReadFile(path)
}

// NewDefaultService assembles the executor service with default adapters.
func NewDefaultService() Service {
	loader := FilesystemCoordinatorLoader{ReadFile: os.ReadFile}
	return Service{
		LoadCoordinator: loader.Load,
		InitRuntime:     func(ctx context.Context, coord m.Coordinator) error { return nil },
		RunLoop:         func(ctx context.Context) error { return ErrLoopNotImplemented },
	}
}
