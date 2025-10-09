---
id: F013
title: Resilience – Breakers and Hot Patch
depends: [F011]
priority: P1
status: backlog
---

## User Story
As an SRE I need automated circuit breakers and hot-patch tooling so the executor can pause failing fronts, apply remediation tasks, and resume safely.

## Outcome
Detect repeated failure fingerprints, trigger configurable breakers, and apply scripted mitigation tasks (e.g., dependency installs) while journaling the intervention in provenance logs.

## Scope & Boundaries
- track failure fingerprints and thresholds per task/resource profile
- inject remediation tasks or lock adjustments when breakers fire
- persist breaker state and interventions in coordinator/provenance outputs
- exclude long-term observability/storage export (handled by F014/F018)

## Acceptance Criteria
- synthetic failure suites trip breakers and append remediation tasks before resuming
- hot patches update plan metadata without invalidating existing hashes incorrectly
- provenance ledger records breaker activations with operator-friendly context

## Evidence & References
- docs/formal-spec.md §§2,15 (rolling frontier resilience and human-in-loop checkpoints)

## Linked Tasks
- T130-circuit-breakers
- T140-hot-patch-engine
