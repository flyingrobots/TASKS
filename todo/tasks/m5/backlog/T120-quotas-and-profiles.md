---
id: T120
milestone: m5
features: [F012]
title: Enforce resource quotas and execution profiles
status: backlog
deps: [T110]
---

## User Story
As an operator I need concurrency profiles and quotas so I can balance throughput and contention when executing planner tasks.

## Summary
Extend the executor arbitration module with profile-based quotas, concurrency ceilings, and fairness policies derived from coordinator resource profiles.

## Scope
### In Scope
- define execution profiles that map to coordinator catalog entries
- enforce per-profile concurrency ceilings with backpressure metrics
- export telemetry for wait queues and quota utilization
### Out of Scope
- lock acquisition ordering (covered by T110)
- breaker-based adjustments (T130/T140)

## Execution Guidance
- support configuration reload or command-line overrides without restart
- integrate quota stats into provenance and admin endpoints
- ensure fairness between profiles via round-robin or weighted scheduling

## Acceptance (machine-verifiable)
```acceptance
[
  {"type":"command","cmd":"cd executor && go test ./internal/arbitration","timeoutSeconds":300},
  {"type":"command","cmd":"cd executor && go test ./...","timeoutSeconds":420}
]
```

## Evidence & References
- docs/formal-spec.md §§2,7,10 (profile-driven arbitration requirements)
- Feature F012 Resource Arbitration
