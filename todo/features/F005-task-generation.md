---
id: F005
title: Task Generation
depends: [F004]
priority: P0
status: backlog
---

## User Story
As an implementing engineer I need bounded, evidence-backed tasks with machine-verifiable acceptance so work can flow in 2–8 hour slices without reinterpreting scope.

## Outcome
Transform parsed feature specs into verb-first, bounded tasks with acceptance checks, execution logging guidance, compensation stance, resource annotations, and autonormalization notes stored in `tasks.json`.

## Scope & Boundaries
- enforce task sizing caps and merge/split heuristics with audit notes
- attach execution logging and telemetry requirements per spec §11
- create structured acceptance check payloads validated by the acceptance tool
- exclude dependency inference (handled by F006)

## Acceptance Criteria
- generator rejects tasks lacking acceptance checks or compensation metadata
- autonormalization log records merges/splits triggered during shaping
- generated tasks achieve ≥80% verb-first title compliance and include execution logging fields

## Evidence & References
- docs/formal-spec.md §§5–7,11 (task sizing, acceptance, logging)

## Linked Tasks
- T030-generate-bounded-tasks-with-acceptance
