---
id: T030
milestone: m2
features: [F005]
title: Generate bounded tasks with acceptance + autonormalization
status: finished
deps: [T020]
---

## User Story
As a planner engineer I need automatically generated, bounded tasks with machine-verifiable acceptance so execution teams inherit actionable work units.

## Summary
Transform parsed spec tasks into normalized `m.Task` entries sized within the 2–8 hour target, attach execution logging defaults, autonormalization notes, and acceptance/evidence stubs.

## Scope
### In Scope
- apply split/merge heuristics and record actions in `meta.autonormalization`
- enforce verb-first titles and resource annotations for each task
- seed acceptance checks with structured command specs validated by tooling
### Out of Scope
- dependency inference (handled by T040)
- wave simulation (T060)

## Execution Guidance
- maintain deterministic ordering when generating IDs (`T001` pattern)
- provide helper functions for constructing default logging + compensation structs
- fail fast when acceptance checks or evidence are missing

## Acceptance (machine-verifiable)
```acceptance
[
  {"type":"command","cmd":"cd planner && go test ./internal/model","timeoutSeconds":180},
  {"type":"command","cmd":"cd planner && go test ./cmd/tasksd","timeoutSeconds":300}
]
```

## Evidence & References
- docs/formal-spec.md §§5–6,11 (task sizing, acceptance, logging)
- Feature F005 Task Generation
