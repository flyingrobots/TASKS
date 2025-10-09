---
id: T050
milestone: m2
features: [F007]
title: Build DAG, layering, and transitive reduction
status: finished
deps: [T040]
---

## User Story
As a release manager I need the planner to emit a minimal precedence DAG with cycle diagnostics so execution planning is dependable.

## Summary
Implement DAG compilation that filters eligible edges, detects cycles, computes Kahn layering and depth metrics, derives the critical path, and applies transitive reduction before emitting `dag.json`.

## Scope
### In Scope
- enforce duplicate task detection and diagnostic reporting
- compute metrics (nodes, edges, longest path length, width)
- produce critical path annotations for DAG nodes
### Out of Scope
- wave simulation (T060)
- coordinator emission (T080)

## Execution Guidance
- reuse adjacency data structures for DFS and reduction passes
- guard against empty task lists with explicit analysis errors
- maintain deterministic ordering of nodes and edges in output

## Acceptance (machine-verifiable)
```acceptance
[
  {"type":"command","cmd":"cd planner && go test ./internal/planner/dag","timeoutSeconds":240},
  {"type":"command","cmd":"cd planner && go test ./cmd/tasksd","timeoutSeconds":300}
]
```

## Evidence & References
- docs/formal-spec.md §§7–8,16 (DAG purity, layering, transitive reduction)
- Feature F007 DAG Build + Reduction
