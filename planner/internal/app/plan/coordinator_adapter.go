package plan

import m "github.com/james/tasks-planner/internal/model"

// CoordinatorBuilder abstracts coordinator artifact construction.
type CoordinatorBuilder interface {
	Build(tasks []m.Task, deps []m.Edge) m.Coordinator
}

// DefaultCoordinatorBuilder constructs the stub coordinator artifact.
type DefaultCoordinatorBuilder struct{}

func (DefaultCoordinatorBuilder) Build(tasks []m.Task, deps []m.Edge) m.Coordinator {
	coord := m.Coordinator{}
	coord.Version = "v8"
	coord.Graph.Nodes = tasks
	coord.Graph.Edges = deps
	coord.Config.Resources.Catalog = map[string]struct {
		Capacity  int    `json:"capacity"`
		Mode      string `json:"mode"`
		LockOrder int    `json:"lock_order"`
	}{}
	coord.Config.Resources.Profiles = map[string]map[string]int{"default": {}}
	coord.Config.Policies.LockOrdering = []string{}
	return coord
}
