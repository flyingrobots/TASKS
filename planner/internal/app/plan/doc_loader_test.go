package plan

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	m "github.com/james/tasks-planner/internal/model"
)

func TestMarkdownDocLoaderFallbackWhenMissing(t *testing.T) {
	loader := NewMarkdownDocLoader()
	res, err := loader.Load(context.Background(), "")
	if err != nil {
		t.Fatalf("expected stub plan, got error: %v", err)
	}
	if res.DocProvided {
		t.Fatalf("expected DocProvided=false")
	}
	if len(res.Tasks) == 0 {
		t.Fatalf("expected stub tasks, got none")
	}
	if len(res.Features) == 0 {
		t.Fatalf("expected stub feature entries, got none")
	}
}

func TestMarkdownDocLoaderParsesFile(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "plan.md")
	content := strings.Join([]string{
		"## Login",
		"- Implement login (3h) after: Setup DB",
		"## Data",
		"- Setup DB",
	}, "\n")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	loader := NewMarkdownDocLoader()
	res, err := loader.Load(context.Background(), path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if !res.DocProvided {
		t.Fatalf("expected DocProvided=true")
	}
	if len(res.Tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(res.Tasks))
	}
	if res.Tasks[0].ID != "T001" || res.Tasks[1].ID != "T002" {
		t.Fatalf("unexpected task IDs: %v", res.Tasks)
	}
	if len(res.Dependencies) == 0 {
		t.Fatalf("expected dependencies from doc")
	}
}

func TestMarkdownDocLoaderPropagatesParseErrors(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "plan.md")
	content := "## Feature\n- Task One\n\n```acceptance\nnot-json\n```\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	loader := NewMarkdownDocLoader()
	_, err := loader.Load(context.Background(), path)
	if err == nil {
		t.Fatalf("expected parse error")
	}
}

func TestApplyTaskDefaults(t *testing.T) {
	task := m.Task{}
	applyTaskDefaults(&task)
	if len(task.AcceptanceChecks) == 0 {
		t.Fatalf("expected default acceptance check")
	}
	if task.DurationUnit != "hours" {
		t.Fatalf("expected duration unit hours")
	}
	if task.ExecutionLogging.Format != "JSONL" {
		t.Fatalf("expected JSONL logging format")
	}
}
