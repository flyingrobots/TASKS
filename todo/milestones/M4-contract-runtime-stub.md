---
id: M4
title: Contract + Runtime Stub
status: backlog
features:
  - F010
  - F011
reference: docs/formal-spec.md
---

# Milestone M4 – Contract + Runtime Stub

## User Story
As an execution lead I need a complete coordinator contract and a rolling-frontier stub so the runtime can consume deterministic planner outputs and exercise scheduling flows.

## Objectives
- emit `coordinator.json` aligned with planner DAG and resource metadata
- bootstrap an executor stub that respects precedence and logs telemetry
- validate interoperability by replaying sample plans end to end

## Deliverables
- coordinator generation module with schema validation
- executor stub with ready-queue management, retries, and JSONL telemetry
- documentation describing contract fields and runtime expectations

## Exit Criteria
- end-to-end sample runs produce matching hashes and JSONL logs
- coordinator contract enumerates lock ordering and runtime policies per spec
- executor stub surfaces critical path and wave insights in logs

## Evidence & References
- docs/formal-spec.md §§2,3,10–11 (planner/executor contract and telemetry)

Features: F010, F011

<!-- PROGRESS:START M4 -->
0/2 done (0%)
<!-- PROGRESS:END M4 -->

## Links
- Features: [F010](../features/F010-coordinator-contract.md), [F011](../features/F011-executor-frontier.md)
- Tasks: see `../tasks/m4/*`
