---
id: T110
milestone: m5
features: [F012]
title: Implement deterministic global lock ordering in executor
status: backlog
deps: [T100]
---

## User Story
As a scheduler engineer I need deterministic lock ordering and acquisition logic so parallel tasks acquire exclusive resources without deadlocks.

## Summary
Add a resource arbitration module that orders lock acquisition globally, performs all-or-nothing grabs, and records wait telemetry for provenance and Plan.md reporting.

## Scope
### In Scope
- define resource descriptors with priority and lock-order rank
- ensure executor acquires locks using sorted token lists
- emit metrics for wait times and failed acquisitions
### Out of Scope
- quota enforcement (T120)
- breaker interventions (T130)

## Execution Guidance
- rely on coordinator resource catalog for lock ordering definitions
- maintain context cancellation pathways to release locks cleanly on failure
- cover cycle detection with fuzz/integration tests under high contention

## Acceptance (machine-verifiable)
```acceptance
[
  {"type":"command","cmd":"cd executor && go test ./internal/arbitration","timeoutSeconds":300},
  {"type":"command","cmd":"cd executor && go test ./...","timeoutSeconds":420}
]
```

## Evidence & References
- docs/formal-spec.md §§2,7,10 (resource arbitration and executor safety)
- Feature F012 Resource Arbitration
