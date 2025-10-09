---
id: T001
milestone: m1
features: [F001]
title: Normalize canonical JSON numbers to minimal decimals
status: finished
deps: []
---

## User Story
As a planner developer I need deterministic minimal-decimal rendering so canonical JSON outputs stay stable across runs and platforms.

## Summary
Implement numeric normalization inside `canonjson` that removes redundant signs, trims trailing zeros, and scrubs exponent noise while preserving semantic equality.

## Scope
### In Scope
- update `planner/internal/canonjson` encoder helpers
- add regression fixtures for zero, scientific notation, and negative handling
- expose helper utilities to other packages via internal API
### Out of Scope
- hashing logic (covered by T002)
- higher-level CLI wiring (handled by T002/T180)

## Execution Guidance
- follow Go 1.25 formatting rules; avoid locale-specific formatting
- expand table-driven tests under `planner/internal/canonjson`
- update docs to clarify minimal decimal expectations

## Acceptance (machine-verifiable)
```acceptance
[
  {"type":"command","cmd":"cd planner && go test ./internal/canonjson","timeoutSeconds":120},
  {"type":"command","cmd":"cd planner && go test ./...","timeoutSeconds":300}
]
```

## Evidence & References
- docs/formal-spec.md §§3,12 (deterministic serialization requirements)
- Feature F001 Canonical JSON + Hashing
