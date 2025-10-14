package exec_test

import (
	"context"
	"errors"
	"testing"

	execapp "github.com/james/tasks-planner/internal/app/exec"
	m "github.com/james/tasks-planner/internal/model"
)

func TestServiceRunHappyPath(t *testing.T) {
	called := struct {
		load, init, loop bool
	}{}

	svc := execapp.Service{
		LoadCoordinator: func(path string) (m.Coordinator, error) {
			called.load = true
			if path != "coord.json" {
				t.Fatalf("unexpected coord path: %s", path)
			}
			coord := m.Coordinator{}
			coord.Version = "v8"
			return coord, nil
		},
		InitRuntime: func(ctx context.Context, coord m.Coordinator) error {
			called.init = true
			if coord.Version != "v8" {
				t.Fatalf("unexpected coordinator version: %s", coord.Version)
			}
			return nil
		},
		RunLoop: func(ctx context.Context) error {
			called.loop = true
			return nil
		},
	}

	if err := svc.Run(context.Background(), "coord.json"); err != nil {
		t.Fatalf("run: %v", err)
	}
	if !called.load || !called.init || !called.loop {
		t.Fatalf("expected all adapters invoked, got %+v", called)
	}
}

func TestServiceRunPropagatesErrors(t *testing.T) {
	svc := execapp.Service{
		LoadCoordinator: func(path string) (m.Coordinator, error) {
			return m.Coordinator{}, errors.New("boom")
		},
		InitRuntime: func(ctx context.Context, coord m.Coordinator) error {
			return nil
		},
		RunLoop: func(ctx context.Context) error {
			return nil
		},
	}
	if err := svc.Run(context.Background(), "coord.json"); err == nil {
		t.Fatalf("expected error")
	}
}
