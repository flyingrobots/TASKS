package dag

import (
	"strings"
	"testing"

	m "github.com/james/tasks-planner/internal/model"
)

func TestBuildRejectsDuplicateTaskIDs(t *testing.T) {
	tasks := []m.Task{{ID: "T001"}, {ID: "T001"}}
	df, err := Build(tasks, nil, 0.7)
	if err == nil {
		t.Fatalf("expected error for duplicate task ids")
	}
	if df.Analysis.OK {
		t.Fatalf("analysis should be marked not OK on duplicate ids")
	}
	const want = "duplicate task id T001"
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("expected error to contain %q, got %v", want, err)
	}
	found := false
	for _, msg := range df.Analysis.Errors {
		if strings.Contains(msg, want) {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("analysis errors missing duplicate id message: %#v", df.Analysis.Errors)
	}
}

func TestBuildRejectsEmptyTasks(t *testing.T) {
	df, err := Build(nil, nil, 0.7)
	if err == nil {
		t.Fatalf("expected error for empty task list")
	}
	if df.Analysis.OK {
		t.Fatalf("analysis should be marked not OK for empty task list")
	}
	const want = "no tasks to build DAG"
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("expected error to contain %q, got %v", want, err)
	}
	found := false
	for _, msg := range df.Analysis.Errors {
		if strings.Contains(msg, want) {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("analysis errors missing empty task message: %#v", df.Analysis.Errors)
	}
}
