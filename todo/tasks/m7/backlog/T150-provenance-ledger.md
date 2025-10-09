---
id: T150
milestone: m7
features: [F014]
title: Build append-only provenance ledger with structured search
status: backlog
deps: [T100, T130]
---

## User Story
As an auditor I need an append-only provenance ledger so I can trace every task execution, resource arbitration decision, and patch event.

## Summary
Implement a JSONL provenance writer that records execution steps, resource waits, breaker activations, and hash references, paired with tooling to query and tail the ledger.

## Scope
### In Scope
- define provenance schema and rotation policy
- emit structured entries from executor tasks and arbitration modules
- provide CLI helpers for filtering by task, resource, or time window
### Out of Scope
- admin HTTP/CLI surfaces (T160)
- external log shipping (future work)

## Execution Guidance
- ensure writes are atomic and append-only even under concurrent workers
- include hash references to artifacts and remediation actions
- integrate provenance path/location into Plan.md summary

## Acceptance (machine-verifiable)
```acceptance
[
  {"type":"command","cmd":"cd executor && go test ./internal/provenance","timeoutSeconds":360},
  {"type":"command","cmd":"cd executor && go test ./...","timeoutSeconds":480}
]
```

## Evidence & References
- docs/formal-spec.md §§11,14 (telemetry and provenance governance)
- Feature F014 Provenance and Logging
