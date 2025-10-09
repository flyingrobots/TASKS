---
id: M8
title: Security & Visualization
status: backlog
features:
  - F016
  - F017
reference: docs/formal-spec.md
---

# Milestone M8 – Security & Visualization

## User Story
As a program steward I need enforced evidence redaction and polished visualization tooling so I can safely share plans and explain scope to stakeholders.

## Objectives
- implement deterministic evidence redaction with validator hooks
- extend DOT export helpers and docs for sharing DAG/runtime views
- ensure visualization tooling respects canonical ordering and critical-path cues

## Deliverables
- redaction utility + validation wiring with regression fixtures
- updated `tasksd export-dot` with label/format toggles and examples
- documentation covering secure evidence handling and visualization workflows

## Exit Criteria
- evidence validator passes redacted excerpts and fails unredacted secrets
- DOT exports render deterministic outputs in Graphviz + sample SVGs
- README/tutorial demonstrates secure sharing workflow referencing redaction policy

## Evidence & References
- docs/formal-spec.md §§7–8,14 (visualization guidance and security redaction)

Features: F016, F017

<!-- PROGRESS:START M8 -->
1/2 done (50%)
<!-- PROGRESS:END M8 -->

## Links
- Features: [F016](../features/F016-security-redaction.md), [F017](../features/F017-visualization-export.md)
- Tasks: see `../tasks/m8/*`
