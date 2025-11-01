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
    // We intentionally set both top-level version and meta.version for
    // backward-compatibility with existing consumers and schema validation.
    // Downstream executor validation checks the top-level field; meta.version
    // is kept as part of the artifact meta lineage.
    coord.Version = schemaVersion
    coord.Meta.Version = schemaVersion
    coord.Graph.Nodes = tasks
    coord.Graph.Edges = deps
    coord.Config.Resources.Catalog = map[string]m.ResourceSpec{}
    coord.Config.Resources.Profiles = map[string]map[string]int{"default": {}}
    coord.Config.Policies.LockOrdering = []string{}
    return coord
}
