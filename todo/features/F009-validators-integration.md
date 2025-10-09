---
id: F009
title: Validator Integration
depends: [F005, F002]
priority: P0
status: backlog
---

## User Story
As a compliance reviewer I need the planner to invoke acceptance, evidence, and interface validators so every task and dependency carries machine-verified backing before artifacts are finalized.

## Outcome
Add validator clients that run during `plan`, collect structured diagnostics, halt the pipeline on missing acceptance checks or invalid evidence, and persist validator fingerprints into plan metadata.

## Scope & Boundaries
- integrate acceptance, evidence, and interface validators described in the spec
- surface validator output in Plan.md and `tasks.json.analysis`
- cache validator results for deterministic reruns when artifacts are unchanged
- exclude execution-time health checks (handled by F013/F018)

## Acceptance Criteria
- planner run stops when validator gating fails and explains the remediation path
- acceptance tool generates runnable command sequences per task
- evidence tool validates ≥95% of tasks/edges with machine-confirmed sources

## Evidence & References
- docs/formal-spec.md §§6,10–11 (validator orchestration requirements)

## Linked Tasks
- T070-wire-validators
