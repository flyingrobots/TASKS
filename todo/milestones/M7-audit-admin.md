---
id: M7
title: Audit & Admin
status: backlog
features:
  - F014
  - F015
reference: docs/formal-spec.md
---

# Milestone M7 – Audit & Admin

## User Story
As an auditor and operator I need provenance logging plus admin tooling so I can trace execution, inspect graphs, and apply controlled patches.

## Objectives
- finalize append-only provenance logging with retention guidance
- expose admin HTTP/CLI endpoints for graph, provenance, and patch workflows
- log administrative actions into provenance for accountability

## Deliverables
- provenance writer module with rotation and retention policies
- admin HTTP service plus CLI wrappers covering read/write verbs
- documentation for admins on authentication and safe patching

## Exit Criteria
- provenance stores satisfy append-only guarantees and survive restarts
- admin operations produce provenance entries with operator identity and payload summary
- CLI + HTTP integration tests cover graph export and patch submission

## Evidence & References
- docs/formal-spec.md §§11,14–15 (telemetry, provenance, and human checkpoints)

Features: F014, F015

<!-- PROGRESS:START M7 -->
0/2 done (0%)
<!-- PROGRESS:END M7 -->

## Links
- Features: [F014](../features/F014-provenance-logging.md), [F015](../features/F015-admin-cli-http.md)
- Tasks: see `../tasks/m7/*`
