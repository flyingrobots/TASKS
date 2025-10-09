---
id: F016
title: Security & Redaction
depends: [F009]
priority: P0
status: backlog
---

## User Story
As a security reviewer I need deterministic redaction of secrets in evidence so plans can be shared without leaking credentials while still proving assertions.

## Outcome
Implement redaction rules for API keys, tokens, JWTs, and PEM blocks in evidence excerpts, integrate checks into validation, and record redaction markers in artifacts.

## Scope & Boundaries
- define regex and heuristic matchers for sensitive substrings
- apply redaction before hashing artifacts and before storing provenance
- ensure validators accept redacted spans and link to original evidence references
- exclude runtime secret rotation (out of scope)

## Acceptance Criteria
- evidence validator confirms redacted excerpts still match sources with placeholders
- acceptance pipeline fails if secrets remain after redaction checks
- documentation outlines operator workflow for supplying sanitized evidence

## Evidence & References
- docs/formal-spec.md §14 (security and redaction rules)

## Linked Tasks
- T170-evidence-redaction
