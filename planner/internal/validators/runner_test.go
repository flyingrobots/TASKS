package validators

import (
	"context"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	m "github.com/james/tasks-planner/internal/model"
)

func TestRunnerExecutesAndCaches(t *testing.T) {
	tmp := t.TempDir()
	cmdPath := validatorCmd(t)
	cfg := Config{AcceptanceCmd: cmdPath, CacheDir: tmp, Timeout: 5 * time.Second}
	runner, err := NewRunner(cfg)
	if err != nil {
		t.Fatalf("runner: %v", err)
	}
	payload := Payload{Tasks: stubTasks()}
	ctx := context.Background()
	reports, err := runner.Run(ctx, payload)
	if err != nil {
		t.Fatalf("first run: %v", err)
	}
	if len(reports) != 1 {
		t.Fatalf("expected 1 report, got %d", len(reports))
	}
	if reports[0].Cached {
		t.Fatalf("expected fresh run, got cached")
	}
	reports2, err := runner.Run(ctx, payload)
	if err != nil {
		t.Fatalf("second run: %v", err)
	}
	if !reports2[0].Cached {
		t.Fatalf("expected cached run")
	}
}

func stubTasks() *m.TasksFile {
	tf := &m.TasksFile{}
	tf.Meta.Version = "v8"
	tf.Meta.MinConfidence = 0.7
	tf.Tasks = []m.Task{{ID: "T001", FeatureID: "F1", Title: "Do thing", Duration: m.DurationPERT{Optimistic: 1, MostLikely: 2, Pessimistic: 3}, DurationUnit: "hours", AcceptanceChecks: []m.AcceptanceCheck{{Type: "command", Cmd: "echo ok"}}}}
	return tf
}

func validatorCmd(t *testing.T) string {
	t.Helper()
	script := filepath.Join("testdata", "mockvalidator", "main.go")
	if _, err := exec.LookPath("go"); err != nil {
		t.Fatalf("go command required: %v", err)
	}
	return "go run " + script
}
