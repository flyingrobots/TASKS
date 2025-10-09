---
id: T060
milestone: m2
feature: F008
title: Generate waves preview from DAG layering
status: finished
deps: [T050]
---

## User Story
As a delivery lead I need a preview of execution waves derived from the DAG so I can judge concurrency and staffing before execution.

## Summary
Use Kahn layering to group tasks into waves, annotate each wave with resource exclusivity collisions, and emit `waves.json` keyed to the `tasks.json` hash without mutating the DAG.

## Scope
### In Scope
- derive waves strictly from `dag.json` ordering
- enrich wave metadata with resource conflict summaries and estimates
- validate that every DAG node maps to a wave entry
### Out of Scope
- runtime scheduling (handled by executor milestones)
- feedback edges into DAG build

## Execution Guidance
- ensure `waves.json.meta.planId` equals the `tasks.json` hash
- fail fast when task IDs referenced in DAG are missing from tasks list
- provide unit tests for empty graph, diamond, and resource-heavy shapes

## Acceptance (machine-verifiable)
```acceptance
[
  {"type":"command","cmd":"cd planner && go test ./internal/planner/wavesim","timeoutSeconds":180},
  {"type":"command","cmd":"rm -rf /tmp/tasks-plan && cd planner && go run ./cmd/tasksd plan --out /tmp/tasks-plan","timeoutSeconds":300}
]
```

## Evidence & References
- docs/formal-spec.md §§7–8 (wave simulation preview contract)
- Feature F008 Wave Simulation
