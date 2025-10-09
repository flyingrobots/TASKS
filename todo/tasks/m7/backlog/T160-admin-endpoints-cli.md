---
id: T160
milestone: m7
features: [F015]
title: Ship admin HTTP + CLI for graph, provenance, and patch workflows
status: backlog
deps: [T150]
---

## User Story
As an operator I need authenticated HTTP endpoints and CLI shims to inspect the graph, tail provenance, and submit safe patches without manual DB edits.

## Summary
Expose admin REST endpoints and CLI wrappers for graph export, breaker state, provenance queries, and patch submission with audit logging.

## Scope
### In Scope
- implement authenticated REST handlers for read/write admin verbs
- mirror handlers via CLI commands for local workflows
- log admin actions (user, verb, payload summary) to provenance
### Out of Scope
- visualization rendering (handled by F017)
- evidence redaction controls (T170)

## Execution Guidance
- enforce RBAC/role enforcement for admin verbs
- return deterministic JSON responses for DOT export requests
- integrate rate limiting to protect runtime control plane

## Acceptance (machine-verifiable)
```acceptance
[
  {"type":"command","cmd":"cd executor && go test ./cmd/admin","timeoutSeconds":360},
  {"type":"command","cmd":"cd executor && go test ./...","timeoutSeconds":480}
]
```

## Evidence & References
- docs/formal-spec.md §§11,15 (admin checkpoints and human oversight)
- Feature F015 Admin CLI/HTTP
