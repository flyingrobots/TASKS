package plan

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
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
	byTitle := map[string]m.Task{}
	for _, task := range res.Tasks {
		byTitle[task.Title] = task
	}
	login, ok := byTitle["Implement login"]
	if !ok {
		t.Fatalf("missing Implement login task: %+v", res.Tasks)
	}
	setup, ok := byTitle["Setup DB"]
	if !ok {
		t.Fatalf("missing Setup DB task: %+v", res.Tasks)
	}
	if len(res.Dependencies) == 0 {
		t.Fatalf("expected dependencies from doc")
	}
	found := false
	for _, dep := range res.Dependencies {
		if dep.From == setup.ID && dep.To == login.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected dependency from %s to %s, got %+v", setup.ID, login.ID, res.Dependencies)
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
	task.Duration = m.DurationPERT{Optimistic: 1, MostLikely: 2, Pessimistic: 3}
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
	expectedFields := []string{"timestamp", "task_id", "step", "status", "message"}
	if !reflect.DeepEqual(task.ExecutionLogging.RequiredFields, expectedFields) {
		t.Fatalf("unexpected required fields: %+v", task.ExecutionLogging.RequiredFields)
	}
	if !task.Compensation.Idempotent {
		t.Fatalf("expected compensation to be idempotent")
	}
	if task.Duration.Optimistic > task.Duration.MostLikely || task.Duration.MostLikely > task.Duration.Pessimistic {
		t.Fatalf("duration PERT out of order: %+v", task.Duration)
	}
}
