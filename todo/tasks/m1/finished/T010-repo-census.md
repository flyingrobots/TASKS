---
id: T010
milestone: m1
features: [F003]
title: Implement repository census with resource hints
status: finished
deps: []
---

## User Story
As a planner architect I need an automated repository census so the planner understands languages, directories, and shared resources before decomposition.

## Summary
Implement `analysis.RunCensus` to traverse the repo, collect file statistics, capture shared resource hints, and surface permission issues without aborting.

## Scope
### In Scope
- walk directories with configurable ignore rules
- record totals for files, Go files, and other language families
- persist census output into `tasks.json.meta.codebase_analysis`
### Out of Scope
- semantic interface extraction (future milestones)
- resource arbitration runtime wiring (handled later)

## Execution Guidance
- handle permission-denied directories gracefully by logging and continuing
- provide reusable helper for CLI commands to print census summaries
- update documentation describing reuse-first expectations

## Acceptance (machine-verifiable)
```acceptance
[
  {"type":"command","cmd":"cd planner && go test ./internal/analysis","timeoutSeconds":180},
  {"type":"command","cmd":"cd planner && go test ./...","timeoutSeconds":300}
]
```

## Evidence & References
- docs/formal-spec.md §4 (codebase-first planning and census output)
- Feature F003 Codebase Census
