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

func TestServiceRunPropagatesInitError(t *testing.T) {
    initErr := errors.New("init failed")
    svc := execapp.Service{
        LoadCoordinator: func(path string) (m.Coordinator, error) {
            return m.Coordinator{Version: "v8"}, nil
        },
        InitRuntime: func(ctx context.Context, coord m.Coordinator) error {
            return initErr
        },
        RunLoop: func(ctx context.Context) error {
            return nil
        },
    }
    if err := svc.Run(context.Background(), "coord.json"); err == nil {
        t.Fatalf("expected init error")
    } else if !errors.Is(err, initErr) {
        t.Fatalf("expected 'init failed', got: %v", err)
    }
}

func TestServiceRunPropagatesLoopError(t *testing.T) {
    errLoop := errors.New("loop failed")
    svc := execapp.Service{
        LoadCoordinator: func(path string) (m.Coordinator, error) {
            return m.Coordinator{Version: "v8"}, nil
        },
        InitRuntime: func(ctx context.Context, coord m.Coordinator) error {
            return nil
        },
        RunLoop: func(ctx context.Context) error {
            return errLoop
        },
    }
    if err := svc.Run(context.Background(), "coord.json"); err == nil {
        t.Fatalf("expected loop error")
    } else if !errors.Is(err, errLoop) {
        t.Fatalf("expected loop error sentinel, got: %v", err)
    }
}

func TestServiceRunFailsOnZeroValueCoordinatorAndSkipsInit(t *testing.T) {
    called := struct{ init, loop bool }{}
    svc := execapp.Service{
        LoadCoordinator: func(string) (m.Coordinator, error) { return m.Coordinator{}, nil },
        InitRuntime: func(ctx context.Context, coord m.Coordinator) error { called.init = true; return nil },
        RunLoop: func(ctx context.Context) error { called.loop = true; return nil },
    }
    err := svc.Run(context.Background(), "coord.json")
    if err == nil || err.Error() != "invalid coordinator: missing version" {
        t.Fatalf("expected invalid coordinator error, got %v", err)
    }
    if called.init || called.loop {
        t.Fatalf("expected InitRuntime/RunLoop not called on validation failure: %+v", called)
    }
}

func TestServiceRunInitErrorPrecedenceOverLoop(t *testing.T) {
    called := struct{ loop bool }{}
    svc := execapp.Service{
        LoadCoordinator: func(string) (m.Coordinator, error) { return m.Coordinator{Version: "v8"}, nil },
        InitRuntime: func(context.Context, m.Coordinator) error { return errors.New("init failed") },
        RunLoop: func(context.Context) error { called.loop = true; return errors.New("loop failed") },
    }
    err := svc.Run(context.Background(), "coord.json")
    if err == nil || err.Error() != "init failed" {
        t.Fatalf("expected init error precedence, got %v", err)
    }
    if called.loop {
        t.Fatalf("expected RunLoop not to run when init fails")
    }
}

func TestServiceRun_ContextCanceledBefore(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    cancel()
    svc := execapp.Service{
        LoadCoordinator: func(path string) (m.Coordinator, error) { return m.Coordinator{Version: "v8"}, nil },
        InitRuntime: func(ctx context.Context, coord m.Coordinator) error {
            // honor context
            return ctx.Err()
        },
        RunLoop: func(ctx context.Context) error { return nil },
    }
    if err := svc.Run(ctx, "coord.json"); !errors.Is(err, context.Canceled) {
        t.Fatalf("expected context.Canceled, got %v", err)
    }
}

func TestServiceRun_ContextCanceledDuringInitAndLoop(t *testing.T) {
    // During InitRuntime
    ctx1, cancel1 := context.WithCancel(context.Background())
    svc1 := execapp.Service{
        LoadCoordinator: func(string) (m.Coordinator, error) { return m.Coordinator{Version: "v8"}, nil },
        InitRuntime: func(ctx context.Context, coord m.Coordinator) error {
            cancel1() // cancel while initializing
            return ctx.Err()
        },
        RunLoop: func(context.Context) error { return nil },
    }
    if err := svc1.Run(ctx1, "coord.json"); !errors.Is(err, context.Canceled) {
        t.Fatalf("expected canceled from init, got %v", err)
    }

    // During RunLoop
    ctx2, cancel2 := context.WithCancel(context.Background())
    svc2 := execapp.Service{
        LoadCoordinator: func(string) (m.Coordinator, error) { return m.Coordinator{Version: "v8"}, nil },
        InitRuntime: func(context.Context, m.Coordinator) error { return nil },
        RunLoop: func(ctx context.Context) error {
            cancel2()
            <-ctx.Done()
            return ctx.Err()
        },
    }
    if err := svc2.Run(ctx2, "coord.json"); !errors.Is(err, context.Canceled) {
        t.Fatalf("expected canceled from loop, got %v", err)
    }
}
