---
id: F011
title: Executor Frontier Stub
depends: [F010]
priority: P1
status: backlog
---

## User Story
As an operator I need a rolling-frontier executor stub that consumes `coordinator.json` so we can validate scheduling, telemetry, and retries before full resource arbitration lands.

## Outcome
Stand up a Go (or compatible) stub that loads coordinator contracts, maintains a ready frontier, dispatches tasks respecting precedence, and emits JSONL telemetry compatible with downstream provenance logging.

## Scope & Boundaries
- implement ready-queue management honoring DAG precedence
- emit structured logs per task (`timestamp`, `task_id`, `step`, `status`, `message`)
- provide simple retry and backoff behavior for transient failures
- exclude advanced resource locks or circuit breakers (handled by F012–F013)

## Acceptance Criteria
- sample coordinator contract executes end-to-end with JSONL telemetry stored under `logs/`
- retries happen with capped exponential backoff for recoverable failures
- operator README explains how to run the stub with sample artifacts

## Evidence & References
- docs/formal-spec.md §§2,11 (rolling frontier execution and telemetry requirements)

## Linked Tasks
- T100-executor-frontier-stub
