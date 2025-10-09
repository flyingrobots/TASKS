package main

import (
	"os"
	"path/filepath"
	"testing"

	m "github.com/james/tasks-planner/internal/model"
)

func TestInferDependenciesSequentialFallback(t *testing.T) {
	tasks := []m.Task{
		{ID: "T001", Title: "One"},
		{ID: "T002", Title: "Two"},
		{ID: "T003", Title: "Three"},
	}
	deps, conflicts := inferDependencies(tasks, nil)
	if len(deps) != 2 {
		t.Fatalf("expected 2 fallback edges, got %d", len(deps))
	}
	if len(conflicts) != 0 {
		t.Fatalf("expected no resource conflicts, got %v", conflicts)
	}
	if deps[0].From != "T001" || deps[0].To != "T002" || deps[0].Type != "sequential" {
		t.Fatalf("unexpected first edge: %+v", deps[0])
	}
	if deps[1].From != "T002" || deps[1].To != "T003" || deps[1].Type != "technical" {
		t.Fatalf("unexpected second edge: %+v", deps[1])
	}
}

func TestInferDependenciesResourceEdges(t *testing.T) {
	tasks := []m.Task{
		{ID: "T001"},
		{ID: "T002"},
		{ID: "T003"},
	}
	tasks[0].Resources.Exclusive = []string{"db"}
	tasks[1].Resources.Exclusive = []string{"db"}
	tasks[2].Resources.Exclusive = []string{"db"}

	deps, conflicts := inferDependencies(tasks, nil)
	if len(deps) < 3 {
		t.Fatalf("expected resource edges to be added, got %d total edges", len(deps))
	}
	entry, ok := conflicts["db"].(map[string]any)
	if !ok {
		t.Fatalf("expected conflict entry for db, got %T", conflicts["db"])
	}
	tasksList, ok := entry["tasks"].([]string)
	if !ok || len(tasksList) != 3 {
		t.Fatalf("unexpected tasks list: %#v", entry)
	}
	resourceCount := 0
	for _, edge := range deps {
		if edge.Type == "resource" {
			resourceCount++
		}
	}
	if resourceCount != 3 {
		t.Fatalf("expected 3 resource edges, got %d", resourceCount)
	}
}

func TestBuildTasksFromDocFallback(t *testing.T) {
	tasks, features, edges, docProvided, err := buildTasksFromDoc("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if docProvided {
		t.Fatalf("expected fallback stub when no doc present")
	}
	if len(tasks) == 0 || len(features) == 0 {
		t.Fatalf("expected stub tasks and features, got %d tasks %d features", len(tasks), len(features))
	}
	if len(edges) != 0 {
		t.Fatalf("expected no edges from stub doc, got %d", len(edges))
	}
}

func TestBuildTasksFromDocReadsFile(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "plan.md")
	content := "## Feature Alpha\n- First task"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write doc: %v", err)
	}
	tasks, features, edges, docProvided, err := buildTasksFromDoc(path)
	if err != nil {
		t.Fatalf("unexpected error parsing doc: %v", err)
	}
	if !docProvided {
		t.Fatalf("expected docProvided to be true")
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	if len(features) == 0 {
		t.Fatalf("expected at least one feature summary")
	}
	if len(edges) != 0 {
		t.Fatalf("expected no edges for single task, got %d", len(edges))
	}
}
