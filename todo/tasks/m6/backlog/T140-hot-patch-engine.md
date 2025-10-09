---
id: T140
milestone: m6
features: [F013]
title: Deliver hot-patch engine for remedial tasks
status: backlog
deps: [T130]
---

## User Story
As an operator I need to inject remediation tasks (e.g., dependency installs) when breakers fire so the executor can self-heal without editing source artifacts manually.

## Summary
Implement a hot-patch pipeline that, upon breaker activation, composes remediation tasks, updates coordinator/provenance, and resumes execution once the patch completes.

## Scope
### In Scope
- define remediation task templates and insertion points
- update coordinator/provenance metadata to reflect injected tasks
- provide operator controls for approving, rolling back, or replaying patches
### Out of Scope
- breaker detection (T130)
- admin surface UI (T160)

## Execution Guidance
- ensure injected tasks respect DAG purity by adding sequential edges only where necessary
- recompute hashes for modified artifacts and update Plan.md hash section deterministically
- record operator approvals with timestamps and rationale in provenance

## Acceptance (machine-verifiable)
```acceptance
[
  {"type":"command","cmd":"cd executor && go test ./internal/resilience","timeoutSeconds":420},
  {"type":"command","cmd":"cd executor && go test ./...","timeoutSeconds":480}
]
```

## Evidence & References
- docs/formal-spec.md §§2,15 (hot patching and human sign-off)
- Feature F013 Resilience – Breakers and Hot Patch
