---
id: M9
title: CI & Documentation
status: backlog
features:
  - F018
  - F019
reference: docs/formal-spec.md
---

# Milestone M9 – CI & Documentation

## User Story
As a release captain I need CI enforcement and refreshed documentation so contributors trust the pipeline and can reproduce plans from scratch.

## Objectives
- wire deterministic planner + validator runs into CI with hash drift detection
- package binaries/containers and publish status badges
- overhaul documentation and samples aligned with the normative spec

## Deliverables
- CI workflows covering canonicalization, validation, unit/integration tests
- release automation that publishes `tasksd` binaries or container images
- documentation suite with runnable samples and troubleshooting guides

## Exit Criteria
- CI fails on schema/hash drift, validator regressions, or nondeterministic outputs
- reproducibility job replays plan generation without hash changes
- sample doc walkthrough regenerates artifacts exactly as documented

## Evidence & References
- docs/formal-spec.md §§3,12,17–19 (determinism, evaluation, scaffolding)

Features: F018, F019

<!-- PROGRESS:START M9 -->
0/2 done (0%)
<!-- PROGRESS:END M9 -->

## Links
- Features: [F018](../features/F018-ci-quality-gates.md), [F019](../features/F019-docs-and-samples.md)
- Tasks: see `../tasks/m9/*`
