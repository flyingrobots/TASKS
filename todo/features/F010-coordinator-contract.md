---
id: F010
title: Coordinator Contract
depends: [F007, F008, F009]
priority: P0
status: backlog
---

## User Story
As an executor operator I need a coordinator contract that enumerates structural nodes, edges, and resource policies so S.L.A.P.S. can schedule work safely and reproducibly.

## Outcome
Generate `coordinator.json` from planner artifacts, encoding resource catalogs, lock ordering, concurrency ceilings, and telemetry requirements aligned with the DAG and waves preview.

## Scope & Boundaries
- translate DAG nodes/edges and resource metadata into coordinator payloads
- define default resource profiles and lock sequences per spec
- include runtime policy defaults (concurrency max, breaker thresholds) while allowing overrides
- exclude the executor runtime implementation (covered by F011)

## Acceptance Criteria
- `coordinator.json` passes schema validation and includes node/edge parity with `dag.json`
- resource catalog entries specify capacity, mode, and lock ordering tokens
- Plan.md cross-references the coordinator hash for auditor traceability

## Evidence & References
- docs/formal-spec.md §§3,10,12 (coordinator contract and resource policies)

## Linked Tasks
- T080-generate-coordinator-contract
