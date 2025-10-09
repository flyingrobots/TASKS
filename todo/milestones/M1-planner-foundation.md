---
id: M1
title: Planner Foundation
status: finished
features:
  - F001
  - F002
  - F003
reference: docs/formal-spec.md
---

# Milestone M1 – Planner Foundation

## User Story
As a platform owner I need deterministic serialization, schemas, and a codebase census so later planning stages inherit trusted building blocks.

## Objectives
- deliver canonical JSON serialization and SHA-256 hashing helpers
- define canonical data models plus embedded JSON Schemas
- inventory the repository and shared resources to seed planning metadata

## Deliverables
- `canonjson` and `hash` packages with CLI entry points
- embedded schemas and validation wiring in `tasksd validate`
- census tooling that populates `tasks.json.meta.codebase_analysis`

## Exit Criteria
- canonicalizer + hash helpers pass regression suites for numbers and key ordering
- validator rejects malformed artifacts using embedded schemas
- census handles permission-denied directories without aborting

## Evidence & References
- docs/formal-spec.md §§3–4,12 (deterministic artifacts, census prerequisites)

Features: F001, F002, F003

<!-- PROGRESS:START M1 -->
4/4 done (100%)
<!-- PROGRESS:END M1 -->

## Links
- Features: [F001](../features/F001-canonical-json.md), [F002](../features/F002-model-and-schemas.md), [F003](../features/F003-codebase-census.md)
- Tasks: see `../tasks/m1/*`
