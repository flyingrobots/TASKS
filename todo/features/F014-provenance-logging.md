---
id: F014
title: Provenance and Logging
depends: [F011]
priority: P1
status: backlog
---

## User Story
As an auditor I need append-only provenance logs with structured task telemetry so I can trace every execution decision and artifact mutation.

## Outcome
Emit JSONL provenance with immutable identifiers, task status transitions, resource arbitration notes, and hash references, and surface summaries in Plan.md and the admin interfaces.

## Scope & Boundaries
- define provenance schema and rotation strategy
- ensure executor logs include required fields and correlation IDs
- provide tooling to query provenance by task, resource, or wave
- exclude admin surfaces (handled by F015)

## Acceptance Criteria
- every executed task produces provenance entries with timestamps, node IDs, and hashes
- provenance writer prevents truncation or rewrite of existing entries
- Plan.md references provenance location and retention policy

## Evidence & References
- docs/formal-spec.md §§2,11,14 (telemetry and provenance governance)

## Linked Tasks
- T150-provenance-ledger
