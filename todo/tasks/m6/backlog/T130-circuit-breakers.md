---
id: T130
milestone: m6
features: [F013]
title: Implement circuit breakers with failure fingerprinting
status: backlog
deps: [T100]
---

## User Story
As an SRE I need circuit breakers that pause failing fronts when error fingerprints exceed thresholds so we can protect throughput and trigger remediation playbooks.

## Summary
Add failure fingerprint tracking to the executor, configure breaker thresholds, and surface breaker status plus recommended remediation hooks to operators.

## Scope
### In Scope
- fingerprint failures based on task ID, error class, and resource profile
- evaluate rolling windows to trigger breaker open/half-open states
- expose breaker state via telemetry and admin APIs
### Out of Scope
- remediation task injection (handled by T140)
- provenance wiring beyond breaker state logging (covered by T150)

## Execution Guidance
- ensure breaker state survives restarts using on-disk or coordinator-backed storage
- allow manual override/close operations through admin CLI/HTTP
- simulate repeated transient errors vs. permanent faults in tests

## Acceptance (machine-verifiable)
```acceptance
[
  {"type":"command","cmd":"cd executor && go test ./internal/resilience","timeoutSeconds":360},
  {"type":"command","cmd":"cd executor && go test ./...","timeoutSeconds":480}
]
```

## Evidence & References
- docs/formal-spec.md §§2,15 (resilience patterns and human checkpoints)
- Feature F013 Resilience – Breakers and Hot Patch
