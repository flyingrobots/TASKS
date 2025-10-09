---
id: F017
title: Visualization & Export
depends: [F007, F010]
priority: P2
status: finished
---

## User Story
As a program lead I need up-to-date DAG and runtime visualizations so I can reason about dependencies and resource flows during reviews.

## Outcome
Generate deterministic DOT exports for DAG and runtime graphs with stable ordering, critical-path highlighting, and optional titles, and provide helper scripts for SVG rendering.

## Scope & Boundaries
- render DAG and coordinator views via `tasksd export-dot`
- support configurable node/edge labels while preserving deterministic ordering
- keep viewer tooling separate (viewer repo lives elsewhere)
- exclude runtime UI bundling (future work)

## Acceptance Criteria
- `tasksd export-dot --dir plans` emits `dag.dot` and `runtime.dot` matching artifacts
- Graphviz renders highlight the critical path when requested
- README documents usage with sample outputs checked into `examples/`

## Evidence & References
- docs/formal-spec.md §§7–8,17 (visualization and evaluation guidance)

## Linked Tasks
- T180-dot-export
