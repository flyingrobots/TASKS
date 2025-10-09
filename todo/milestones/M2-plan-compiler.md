---
id: M2
title: Plan Compiler
status: finished
features:
  - F004
  - F005
  - F006
  - F007
  - F008
reference: docs/formal-spec.md
---

# Milestone M2 – Plan Compiler

## User Story
As the planning lead I need an automated pipeline that turns a repo+spec into canonical features, tasks, DAG, and waves so downstream executors inherit a deterministic contract.

## Objectives
- parse the charter into prioritized features and sized tasks
- infer structural dependencies and compile the precedence DAG
- generate waves previews without polluting DAG purity
- document autonormalization and validator gating policies

## Deliverables
- updated `tasksd plan` wiring for doc parsing → tasks → dependencies → DAG → waves
- validated artifacts (`features.json`, `tasks.json`, `dag.json`, `waves.json`, `Plan.md`)
- CI fixtures that exercise plan compilation end to end

## Exit Criteria
- dag.json reports `analysis.ok == true` for baseline plans
- waves preview references the `tasks.json` hash and passes schema validation
- validator integration blocks plans missing acceptance, evidence, or confidence thresholds

## Evidence & References
- docs/formal-spec.md §§3–8 (planning pipeline, DAG purity, wave preview)

Features: F004, F005, F006, F007, F008

<!-- PROGRESS:START M2 -->
5/5 done (100%)
<!-- PROGRESS:END M2 -->

## Links
- Features: [F004](../features/F004-feature-extraction.md), [F005](../features/F005-task-generation.md), [F006](../features/F006-dependency-inference.md), [F007](../features/F007-dag-build-and-reduction.md), [F008](../features/F008-wave-simulation.md)
- Tasks: see `../tasks/m2/*`
