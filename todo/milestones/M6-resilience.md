---
id: M6
title: Resilience
status: backlog
features:
  - F013
reference: docs/formal-spec.md
---

# Milestone M6 – Resilience

## User Story
As an SRE I need circuit breakers and hot-patch tooling so the executor pauses failing fronts, remediates safely, and resumes with provenance accountability.

## Objectives
- detect repeated failure fingerprints and trigger configurable breakers
- inject remediation tasks or lock adjustments under breaker control
- journal interventions into provenance for audit

## Deliverables
- breaker evaluation and state persistence modules
- hot-patch executor helpers that append remediation tasks and restart affected waves
- docs for operators explaining breaker tuning and rollback

## Exit Criteria
- synthetic failure suites trigger breakers and apply remediation before resuming execution
- provenance ledger records breaker activations and operator confirmations
- tests cover rollback of misapplied patches without corrupting plan state

## Evidence & References
- docs/formal-spec.md §§2,15 (rolling frontier resilience and human-in-loop checkpoints)

Features: F013

<!-- PROGRESS:START M6 -->
0/2 done (0%)
<!-- PROGRESS:END M6 -->

## Links
- Features: [F013](../features/F013-resilience-breakers-hotpatch.md)
- Tasks: see `../tasks/m6/*`
