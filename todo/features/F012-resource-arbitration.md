---
id: F012
title: Resource Arbitration
depends: [F010]
priority: P1
status: backlog
---

## User Story
As an executor engineer I need deterministic lock ordering and quota enforcement so concurrent tasks respect exclusive resources without polluting the DAG structure.

## Outcome
Implement resource catalogs, lock ordering, quota enforcement, and profile-based concurrency controls within the executor, informed by `tasks.json.meta.codebase_analysis` and coordinator metadata.

## Scope & Boundaries
- honor exclusive resources with all-or-nothing acquisition using a global lock order
- enforce per-profile quotas and concurrency ceilings defined in coordinator files
- expose metrics and logs when arbitration delays tasks or detects saturation
- exclude resilience playbooks and circuit breakers (handled by F013)

## Acceptance Criteria
- integration tests show deadlock-free execution under competing resource demands
- lock order is deterministic and recorded in the coordinator artifact
- telemetry captures arbitration wait times for Plan.md summaries

## Evidence & References
- docs/formal-spec.md §§2,7,10 (resource separation and runtime arbitration)

## Linked Tasks
- T110-locks-global-order
- T120-quotas-and-profiles
