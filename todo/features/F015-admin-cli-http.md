---
id: F015
title: Admin CLI/HTTP
depends: [F014]
priority: P2
status: backlog
---

## User Story
As an operator I need introspection endpoints and CLI verbs to inspect graphs, apply safe patches, and review provenance without SSHing into workers.

## Outcome
Expose admin HTTP routes and CLI commands for graph export, breaker status, provenance queries, and patch submission, all authenticated and auditable.

## Scope & Boundaries
- implement read and write endpoints for DAG, coordinator, and provenance inspection
- mirror endpoints through `npm run`/CLI wrappers for local workflow
- enforce authN/Z and log administrative actions into provenance
- exclude visualization rendering (handled by F017)

## Acceptance Criteria
- `GET /admin/graph`, `GET /admin/provenance`, and `POST /admin/patch` endpoints documented and tested
- CLI wrappers provide ergonomic shorthands for the same verbs
- admin operations append provenance entries capturing operator ID and payload summary

## Evidence & References
- docs/formal-spec.md §§2,11,15 (human-in-loop checkpoints and telemetry)

## Linked Tasks
- T160-admin-endpoints-cli
