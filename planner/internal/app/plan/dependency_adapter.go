package plan

import (
	"sort"

	m "github.com/james/tasks-planner/internal/model"
)

// DependencyResolver computes precedence edges and resource conflicts from tasks + doc edges.
type DependencyResolver interface {
	Resolve(tasks []m.Task, docEdges []m.Edge) (edges []m.Edge, resourceConflicts map[string]any)
}

// DefaultDependencyResolver implements the existing inference logic.
type DefaultDependencyResolver struct{}

func (DefaultDependencyResolver) Resolve(tasks []m.Task, baseEdges []m.Edge) ([]m.Edge, map[string]any) {
	deps := append([]m.Edge{}, baseEdges...)
	if len(deps) == 0 && len(tasks) >= 2 {
		for i := 0; i < len(tasks)-1; i++ {
			typ := "sequential"
			if i > 0 {
				typ = "technical"
			}
			deps = append(deps, m.Edge{From: tasks[i].ID, To: tasks[i+1].ID, Type: typ, IsHard: true, Confidence: 1.0})
		}
	}

	resToTasks := map[string][]string{}
	for _, task := range tasks {
		for _, r := range task.Resources.Exclusive {
			resToTasks[r] = append(resToTasks[r], task.ID)
		}
	}

	resourceConflicts := map[string]any{}
	keys := make([]string, 0, len(resToTasks))
	for k := range resToTasks {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, r := range keys {
		ids := resToTasks[r]
		sort.Strings(ids)
		if len(ids) < 2 {
			continue
		}
		resourceConflicts[r] = map[string]any{"type": "exclusive", "tasks": ids}
		for i := 0; i < len(ids); i++ {
			for j := i + 1; j < len(ids); j++ {
				deps = append(deps, m.Edge{From: ids[i], To: ids[j], Type: "resource", Subtype: "mutual_exclusion", IsHard: true, Confidence: 1.0})
			}
		}
	}

	return deps, resourceConflicts
}
