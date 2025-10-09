---
id: T080
milestone: m4
features: [F010]
title: Generate coordinator.json from planner artifacts
status: backlog
deps: [T050, T060]
---

## User Story
As an execution lead I need the planner to emit a coordinator contract describing nodes, edges, and resource policies so the runtime can schedule work deterministically.

## Summary
Compose `coordinator.json` using tasks, DAG, and waves data, embed resource catalogs, profiles, and policies, and validate the output against the coordinator schema.

## Scope
### In Scope
- map tasks and DAG edges into coordinator nodes/edges arrays
- publish resource catalogs with lock ordering and capacity settings
- include runtime policy defaults (concurrency, breaker hints, telemetry)
### Out of Scope
- executor implementation (T100+)
- admin/provenance tooling (future milestones)

## Execution Guidance
- keep coordinator generation deterministic and sorted for hashing
- fail fast when tasks or edges referenced in coordinator are missing
- update Plan.md to summarize coordinator metrics

## Acceptance (machine-verifiable)
```acceptance
[
  {"type":"command","cmd":"cd planner && go test ./cmd/tasksd","timeoutSeconds":300},
  {"type":"command","cmd":"cd planner && go run ./cmd/tasksd plan --out /tmp/tasks-plan && cd planner && go run ./cmd/tasksd validate --dir /tmp/tasks-plan","timeoutSeconds":420}
]
```

## Evidence & References
- docs/formal-spec.md §§3,10,12 (coordinator contract requirements)
- Feature F010 Coordinator Contract
