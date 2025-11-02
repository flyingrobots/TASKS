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

func TestServiceRunPropagatesLoadCoordinatorError(t *testing.T) {
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

// The failure condition is a missing/empty Version field on the coordinator.
func TestServiceRunFailsOnMissingVersionAndSkipsInit(t *testing.T) {
    called := struct{ init, loop bool }{}
    svc := execapp.Service{
        LoadCoordinator: func(string) (m.Coordinator, error) {
            return m.Coordinator{}, nil
        },
        InitRuntime: func(ctx context.Context, coord m.Coordinator) error {
            called.init = true
            return nil
        },
        RunLoop: func(ctx context.Context) error {
            called.loop = true
            return nil
        },
    }
    err := svc.Run(context.Background(), "coord.json")
    if err == nil {
        t.Fatalf("expected invalid coordinator error")
    } else if !errors.Is(err, execapp.ErrInvalidCoordinator) {
        t.Fatalf("expected invalid coordinator error, got %v", err)
    }
    if called.init || called.loop {
        t.Fatalf("expected InitRuntime/RunLoop not called on validation failure: %+v", called)
    }
}

func TestServiceRunInitErrorPrecedenceOverLoop(t *testing.T) {
    called := struct{ loop bool }{}
    initErr := errors.New("init failed")
    svc := execapp.Service{
        LoadCoordinator: func(string) (m.Coordinator, error) {
            return m.Coordinator{Version: "v8"}, nil
        },
        InitRuntime: func(context.Context, m.Coordinator) error {
            return initErr
        },
        RunLoop: func(context.Context) error {
            called.loop = true
            return errors.New("loop failed")
        },
    }
    err := svc.Run(context.Background(), "coord.json")
    if err == nil || !errors.Is(err, initErr) {
        t.Fatalf("expected init error precedence, got %v", err)
    }
    if called.loop {
        t.Fatalf("expected RunLoop not to run when init fails")
    }
}

func TestServiceRunContextCanceledBefore(t *testing.T) {
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

func TestServiceRunContextCanceledDuringInitAndLoop(t *testing.T) {
    cases := []struct{
        name string
        cancelOn string // "init" or "loop"
    }{
        {name:"cancel during init", cancelOn:"init"},
        {name:"cancel during loop", cancelOn:"loop"},
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            ctx, cancel := context.WithCancel(context.Background())
            svc := execapp.Service{
                LoadCoordinator: func(string) (m.Coordinator, error) { return m.Coordinator{Version: "v8"}, nil },
                InitRuntime: func(c context.Context, _ m.Coordinator) error {
                    if tc.cancelOn == "init" { cancel(); return c.Err() }
                    return nil
                },
                RunLoop: func(c context.Context) error {
                    if tc.cancelOn == "loop" { cancel() }
                    return c.Err()
                },
            }
            if err := svc.Run(ctx, "coord.json"); !errors.Is(err, context.Canceled) {
                t.Fatalf("expected context.Canceled for %s, got %v", tc.name, err)
            }
        })
    }
}

func TestServiceRunFailsWithMissingAdapters(t *testing.T) {
    tests := []struct{
        name string
        svc  execapp.Service
    }{
        {
            name: "no LoadCoordinator",
            svc: execapp.Service{
                InitRuntime: func(context.Context, m.Coordinator) error { return nil },
                RunLoop:     func(context.Context) error { return nil },
            },
        },
        {
            name: "no InitRuntime",
            svc: execapp.Service{
                LoadCoordinator: func(string) (m.Coordinator, error) { return m.Coordinator{}, nil },
                RunLoop:         func(context.Context) error { return nil },
            },
        },
        {
            name: "no RunLoop",
            svc: execapp.Service{
                LoadCoordinator: func(string) (m.Coordinator, error) { return m.Coordinator{}, nil },
                InitRuntime:     func(context.Context, m.Coordinator) error { return nil },
            },
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if err := tt.svc.Run(context.Background(), "coord.json"); !errors.Is(err, execapp.ErrMissingAdapters) {
                t.Fatalf("expected missing adapters error, got %v", err)
            }
        })
    }
}
