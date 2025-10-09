---
id: F007
title: DAG Build + Reduction
depends: [F006]
priority: P0
status: backlog
---

## User Story
As a release manager I need a deterministic precedence DAG with cycle diagnostics, layering, and transitive reduction so downstream scheduling and visualization stay trustworthy.

## Outcome
Compile hard, high-confidence edges into a DAG, detect cycles with DFS, perform Kahn layering, compute width and critical path, and apply transitive reduction before emitting `dag.json` metrics.

## Scope & Boundaries
- include only `isHard` edges meeting the minimum confidence threshold
- emit analysis diagnostics when duplicate tasks or cycles are detected
- compute metrics (nodes, edges, width, critical path) and populate `dag.json.metrics`
- exclude wave simulation or executor projections (covered by F008/F010)

## Acceptance Criteria
- DAG build fails loudly on duplicate task IDs or remaining cycles
- transitive reduction eliminates redundant edges on fixture graphs
- metrics match expectations for sample plans and are referenced in Plan.md hashes

## Evidence & References
- docs/formal-spec.md §§7–8,16 (DAG purity, layering, transitive reduction)

## Linked Tasks
- T050-build-dag-and-transitive-reduction
