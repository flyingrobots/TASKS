---
id: F002
title: Canonical Data Model and Schemas
depends: []
priority: P0
status: shaping
---

## User Story
As a planner reviewer I need strongly typed artifact models and embedded schemas so that validation catches drift automatically and every artifact documents its contract.

## Outcome
Provide Go structs for every artifact, embed the authoritative JSON Schemas, and expose validate routines so the CLI and CI can ensure files match the spec before hashing or shipping.

## Scope & Boundaries
- define canonical Go structs aligned with `features.json`, `tasks.json`, `dag.json`, `waves.json`, and `coordinator.json`
- embed and version JSON Schemas alongside the planner for offline validation
- wire schema + struct validation into `tasksd validate` and unit tests
- exclude executor runtime schemas (covered under feature F010)

## Acceptance Criteria
- `tasksd validate` rejects malformed fixture artifacts and passes the happy path
- schemas embed build-time hash/version annotations for audit trails
- Go structs expose helper constructors that enforce required defaults (e.g., meta version)
- documentation explains how downstream automation consumes the definitions

## Evidence & References
- docs/formal-spec.md §§3,12 (artifact contracts and schema determinism)

## Linked Tasks
- T003-embed-schemas-and-validate
