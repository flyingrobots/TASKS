---
id: M3
title: Validators Integration
status: finished
features:
  - F009
reference: docs/formal-spec.md
---

# Milestone M3 – Validators Integration

## User Story
As a compliance lead I need validator gating wired into the planner so every task, edge, and interface carries machine-verified acceptance and evidence before artifacts hash out.

## Objectives
- integrate acceptance, evidence, and interface validators into the `plan` workflow
- halt planning when acceptance, evidence, or interface contracts fail
- surface validator diagnostics with actionable remediation steps

## Deliverables
- validator client package with caching + deterministic logs
- updates to `tasksd plan` and Plan.md summarizing validator run results
- regression fixtures demonstrating success and failure cases

## Exit Criteria
- planning aborts when validator coverage drops below the ≥95% evidence threshold
- acceptance checks render executable command sequences stored alongside tasks
- Plan.md documents validator fingerprints for audit

## Evidence & References
- docs/formal-spec.md §§6,10–11 (validator orchestration and evidence requirements)

Features: F009

<!-- PROGRESS:START M3 -->
1/1 done (100%)
<!-- PROGRESS:END M3 -->

## Links
- Features: [F009](../features/F009-validators-integration.md)
- Tasks: see `../tasks/m3/*`
