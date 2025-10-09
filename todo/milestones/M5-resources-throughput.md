---
id: M5
title: Resources & Throughput
status: backlog
features:
  - F012
reference: docs/formal-spec.md
---

# Milestone M5 – Resources & Throughput

## User Story
As a scheduler engineer I need deterministic lock ordering and quota enforcement so the executor maximizes throughput without violating DAG purity.

## Objectives
- implement global lock ordering and all-or-nothing acquisition for exclusive resources
- enforce quota and profile-based concurrency ceilings
- surface arbitration wait metrics and saturation signals in telemetry

## Deliverables
- resource arbitration module wired into the executor frontier
- coordinator profile examples illustrating concurrency tuning
- regression tests exercising contention scenarios without deadlock

## Exit Criteria
- executor handles competing exclusive tasks without starvation or deadlock
- arbitration metrics recorded in provenance and surfaced in Plan.md
- documentation guides operators on defining resource profiles

## Evidence & References
- docs/formal-spec.md §§2,7,10 (resource arbitration and executor responsibilities)

Features: F012

<!-- PROGRESS:START M5 -->
0/2 done (0%)
<!-- PROGRESS:END M5 -->

## Links
- Features: [F012](../features/F012-resource-arbitration.md)
- Tasks: see `../tasks/m5/*`
