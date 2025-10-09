---
id: T003
milestone: m1
features: [F002]
title: Embed schemas and wire planner validation
status: finished
deps: [T001, T002]
---

## User Story
As a planner maintainer I need embedded JSON Schemas and a validate command so artifacts fail fast when they drift from the normative contract.

## Summary
Define Go structs for every artifact, embed matching schemas, and update `tasksd validate` to run schema + struct validation with hash verification.

## Scope
### In Scope
- generate schemas for features, tasks, dag, waves, and coordinator artifacts
- embed schemas via `embed` directives and expose validator helpers
- extend `tasksd validate` to check hashes, parse, and schema-validate every artifact
### Out of Scope
- generating plan artifacts (handled elsewhere)
- executor runtime schemas (future milestones)

## Execution Guidance
- keep schema references stable and include version metadata
- ensure validation distinguishes between hash mismatch and schema violations
- cover schema changes with regression tests using fixture artifacts

## Acceptance (machine-verifiable)
```acceptance
[
  {"type":"command","cmd":"cd planner && go test ./internal/validate","timeoutSeconds":180},
  {"type":"command","cmd":"cd planner && go run ./cmd/tasksd validate --dir ./plans","timeoutSeconds":180}
]
```

## Evidence & References
- docs/formal-spec.md §§3,12 (artifact contract and schema validation)
- Feature F002 Canonical Data Model and Schemas
