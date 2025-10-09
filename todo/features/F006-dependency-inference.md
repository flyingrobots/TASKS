---
id: F006
title: Dependency Inference
depends: [F005]
priority: P0
status: backlog
---

## User Story
As a planner architect I need automatically inferred structural dependencies so the DAG reflects only precedence edges while resource relationships remain adjacent metadata.

## Outcome
Infer sequential, technical, infrastructure, and knowledge edges with confidence scores, record resource conflicts separately, and persist the results into `tasks.json` for DAG compilation.

## Scope & Boundaries
- parse acceptance and interface metadata to infer structural edges
- capture resource exclusivity edges in `tasks.json.resource_conflicts` while excluding them from DAG build
- annotate edges with confidence + evidence handles for validator tooling
- exclude DAG compilation (handled by F007)

## Acceptance Criteria
- structural edges meet the minimum confidence threshold configurable in meta
- resource edges appear in `tasks.json` while `dag.json` contains only structural edges
- unit tests cover mixed confidence edges and resource separation cases

## Evidence & References
- docs/formal-spec.md §§5,7,73 (dependency inference and resource separation)

## Linked Tasks
- T040-infer-structural-deps-and-resource-separation
