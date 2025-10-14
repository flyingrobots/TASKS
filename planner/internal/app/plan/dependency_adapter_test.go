package plan

import (
	"testing"

	m "github.com/james/tasks-planner/internal/model"
)

func TestDefaultDependencyResolverSequentialFallback(t *testing.T) {
	tasks := []m.Task{
		{ID: "T001"},
		{ID: "T002"},
		{ID: "T003"},
	}
	deps, conflicts := DefaultDependencyResolver{}.Resolve(tasks, nil)
	if len(deps) != 2 {
		t.Fatalf("expected 2 fallback edges, got %d", len(deps))
	}
	if _, ok := conflicts["db"]; ok {
		t.Fatalf("unexpected resource conflict: %+v", conflicts)
	}
}

func TestDefaultDependencyResolverResources(t *testing.T) {
	tasks := []m.Task{{ID: "T001"}, {ID: "T002"}}
	tasks[0].Resources.Exclusive = []string{"db"}
	tasks[1].Resources.Exclusive = []string{"db"}
	deps, conflicts := DefaultDependencyResolver{}.Resolve(tasks, nil)
	if got := len(conflicts); got != 1 {
		t.Fatalf("expected 1 resource conflict, got %d (%+v)", got, conflicts)
	}
	if got := len(deps); got != 1 {
		t.Fatalf("expected only fallback precedence edge, got %d", got)
	}
}
