---
id: F018
title: CI + Quality Gates
depends: [F001, F002, F007, F009]
priority: P0
status: backlog
---

## User Story
As a release captain I need CI pipelines that gate on schema, hash, validator, and DAG checks so regressions are caught before merge.

## Outcome
Wire planner commands and validator suites into CI, enforce hash determinism, publish build artifacts, and expose status badges documenting gate coverage.

## Scope & Boundaries
- run `go test`, schema validation, and canonicalization checks in CI
- confirm artifact hashes match committed Plan.md entries
- publish release binaries or container images that bundle `tasksd`
- exclude operational deployment automation (future work)

## Acceptance Criteria
- CI fails on schema/hash drift or validator failures and links to remediation docs
- nightly reproducibility job replays plan generation without hash deltas
- README badges show CI and lint status for default branch

## Evidence & References
- docs/formal-spec.md §§3,12,17 (determinism and evaluation gates)

## Linked Tasks
- T190-ci-quality-gates
