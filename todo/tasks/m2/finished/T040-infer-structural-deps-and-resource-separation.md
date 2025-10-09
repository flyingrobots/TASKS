---
id: T040
milestone: m2
features: [F006]
title: Infer structural dependencies and record resource conflicts
status: finished
deps: [T030]
---

## User Story
As a planner architect I need structural dependencies separated from resource conflicts so the DAG remains pure and resource metadata stays auditable.

## Summary
Analyze tasks, interfaces, and acceptance hints to infer technical, sequential, infrastructure, and knowledge edges, while cataloging resource exclusivity relationships outside the DAG.

## Scope
### In Scope
- compute edge confidence and `isHard` fields per inference
- maintain `tasks.json.resource_conflicts` for exclusivity groups
- expose diagnostics for unresolved references or confidence below thresholds
### Out of Scope
- DAG compilation (T050)
- wave simulation (T060)

## Execution Guidance
- ensure resource edges never enter the precedence build path
- support configurable minimum confidence in `tasks.json.meta`
- create regression tests mixing structural and resource edges

## Acceptance (machine-verifiable)
```acceptance
[
  {"type":"command","cmd":"cd planner && go test ./internal/planner","timeoutSeconds":240},
  {"type":"command","cmd":"cd planner && go test ./cmd/tasksd","timeoutSeconds":300}
]
```

## Evidence & References
- docs/formal-spec.md §§5,7,73 (dependency inference and resource separation)
- Feature F006 Dependency Inference
