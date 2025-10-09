---
id: T170
milestone: m8
features: [F016]
title: Implement deterministic evidence redaction pipeline
status: backlog
deps: [T070]
---

## User Story
As a security reviewer I need evidence excerpts redacted deterministically so artifacts can be shared without leaking secrets while validators still confirm provenance.

## Summary
Add redaction utilities that detect sensitive substrings, replace them with `[REDACTED_*]` markers, integrate with evidence validation, and document the policy.

## Scope
### In Scope
- regex/heuristic matchers for keys, tokens, JWTs, PEM blocks
- integration hooks in validator and artifact writer workflows
- regression fixtures covering before/after redaction
### Out of Scope
- admin UI for secrets rotation (future milestone)
- runtime secret management (executor scope)

## Execution Guidance
- perform redaction before canonicalization and hashing to keep determinism
- store original evidence references securely for validator replay if authorized
- fail validation when redaction cannot sanitize detected secrets

## Acceptance (machine-verifiable)
```acceptance
[
  {"type":"command","cmd":"cd planner && go test ./internal/validate","timeoutSeconds":240},
  {"type":"command","cmd":"cd planner && go test ./cmd/tasksd","timeoutSeconds":360}
]
```

## Evidence & References
- docs/formal-spec.md §14 (security and redaction guidance)
- Feature F016 Security & Redaction
