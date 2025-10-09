---
id: T100
milestone: m4
features: [F011]
title: Deliver rolling-frontier executor stub with telemetry
status: backlog
deps: [T080]
---

## User Story
As an operator I need a frontier executor stub that consumes coordinator contracts so I can simulate execution, retries, and telemetry before full resource arbitration is ready.

## Summary
Implement a Go service that ingests `coordinator.json`, maintains a ready frontier respecting precedence, dispatches tasks with retries, and emits structured JSONL telemetry to the provenance pipeline.

## Scope
### In Scope
- ready-queue management and DAG-based unlock logic
- pluggable worker shims that simulate steps via shell or Go callbacks
- structured telemetry with the required fields (`timestamp`, `task_id`, `step`, `status`, `message`)
### Out of Scope
- advanced resource locking (T110+)
- circuit breakers or hot patches (T130/T140)

## Execution Guidance
- implement bounded retry/backoff with idempotency awareness
- ensure telemetry writes are append-only and survive restarts
- supply sample coordinator + tasks fixtures for smoke tests

## Acceptance (machine-verifiable)
```acceptance
[
  {"type":"command","cmd":"cd executor && go test ./...","timeoutSeconds":360},
  {"type":"command","cmd":"cd executor && go run ./cmd/slaps --plan ../examples/contracts/basic","timeoutSeconds":360}
]
```

## Evidence & References
- docs/formal-spec.md §§2,11 (rolling frontier execution and telemetry)
- Feature F011 Executor Frontier Stub
