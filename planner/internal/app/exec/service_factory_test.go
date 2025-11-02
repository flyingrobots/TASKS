package exec_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	execapp "github.com/james/tasks-planner/internal/app/exec"
)

func TestNewDefaultServiceLoadsCoordinator(t *testing.T) {
	dir := t.TempDir()
    coordPath := filepath.Join(dir, "coord.json")
    mini := `{"version":"v9","meta":{"version":"v9","artifactHash":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"graph":{"nodes":[],"edges":[]},"config":{"resources":{"catalog":{},"profiles":{"default":{}}},"policies":{"circuitBreakerThresholds":null,"concurrencyMax":0,"lockOrdering":[]}}}`
    if err := os.WriteFile(coordPath, []byte(mini), 0o644); err != nil {
        t.Fatalf("write coord: %v", err)
    }

	svc := execapp.NewDefaultService()
    err := svc.Run(context.Background(), coordPath)
    if !errors.Is(err, execapp.ErrLoopNotImplemented) {
        t.Fatalf("expected ErrLoopNotImplemented, got %v", err)
    }
}

func TestNewDefaultServiceMissingFile(t *testing.T) {
	svc := execapp.NewDefaultService()
	if err := svc.Run(context.Background(), "nope.json"); err == nil {
		t.Fatalf("expected error for missing coordinator")
	}
}
