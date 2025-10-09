---
id: F008
title: Wave Simulation (Preview Only)
depends: [F007]
priority: P1
status: backlog
---

## User Story
As a delivery lead I need a preview of execution waves derived from the DAG so I can reason about concurrency opportunities while keeping resource constraints outside the graph.

## Outcome
Simulate Kahn-based waves, annotate each wave with depth, resource exclusivity hints, and readiness checks, and emit `waves.json` keyed off the `tasks.json` hash without mutating the DAG.

## Scope & Boundaries
- derive waves strictly from DAG topology with optional resource annotations
- flag missing tasks/resources as validation errors
- keep simulation outputs informative only—no feedback edges into DAG build
- exclude runtime scheduling (handled by executor features F010–F013)

## Acceptance Criteria
- `waves.json` references the `tasks.json` hash and passes schema validation
- simulation fails on missing task IDs, ensuring traceability
- wave preview reports estimated width per depth for the plan summary

## Evidence & References
- docs/formal-spec.md §§7–8 (Kahn layering, wave preview separation)

## Linked Tasks
- T060-generate-waves-preview
