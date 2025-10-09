---
id: F003
title: Codebase Census (Reuse-first)
depends: []
priority: P0
status: in-flight
---

## User Story
As a planner I need an automated census of the repository so decomposition reuses existing interfaces, resources, and tooling per the reuse-first directive.

## Outcome
Implement `analysis.RunCensus` to inventory languages, directories, shared resources, and permissions, persisting the summary into `tasks.json.meta.codebase_analysis` for downstream planning and executor resource profiles.

## Scope & Boundaries
- traverse the repository respecting ignore files and permissions
- capture file counts, language buckets, shared resource hints, and chokepoints
- surface permission or traversal errors without aborting the census
- exclude semantic parsing of features/tasks (handled in docparse)

## Acceptance Criteria
- census CLI succeeds on macOS/Linux sample runs and logs permission denials without crashing
- `tasks.json.meta.codebase_analysis` includes shared resource suggestions consumed by resource arbitration
- unit tests cover empty, nested, and permission-restricted directories

## Evidence & References
- docs/formal-spec.md §4 (codebase-first planning requirements)

## Linked Tasks
- T010-repo-census
