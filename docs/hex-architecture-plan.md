# Hexagonal Architecture Adoption Plan

## Why
- Preserve pure domain logic (planner + executor) while isolating environment-specific adapters (CLI, filesystem, validators).
- Improve unit-test coverage by swapping external services (validators, Graphviz export, todo scripts) with in-memory fakes.
- Enable future multi-platform runtimes (daemon, API) without rewriting business logic.

## Current State Snapshot
- `cmd/tasksd` mixes CLI flag parsing, planning orchestration, artifact writing, and validator invocation.
- Planner internals already have separable responsibilities (`analysis`, `docparse`, `dag`, `validators`, `emitter`), but dependencies are direct package imports.
- External concerns (I/O, subprocess exec, cache) are invoked inline; swapping them requires integration tests.

## Target Architecture Overview
```
+---------------------------+
| Application Services     |
|  (Orchestrators)         |
|   - PlanService          |
|   - ExportService        |
+-------------+-------------+
              |
      +-------+-------+
      |   Domain Core |
      |  (pure Go)    |
      +---+-------+---+
          |       |
      Ports    Ports
    (interfaces)
```
Ports / adapters:
- **DocPort**: parse requests (file, string, remote spec) -> `[]Feature`, `[]TaskSpec`.
- **ValidatorPort**: acceptance/evidence/interface validators; default adapter uses subprocess, future adapters may call HTTP services.
- **ArtifactPort**: persistence for artifacts (filesystem today, future memory/remote).
- **AnalysisPort**: wraps repo census + interface scanners.
- **Logging/Telemetry Port**: attaches structured output for plan/executor.

## Incremental Steps
1. **Introduce Service Layer (PlanningService)**
   - Move orchestration from `cmd/tasksd` into `internal/app/plan/service.go`.
   - Inject current helpers (analysis, docparse, dag builder, validators) via interfaces.
2. **Define Port Interfaces**
   - Create small interfaces in `internal/ports` for validators, doc sources, artifact sink.
   - Provide default adapters wrapping existing packages.
3. **Adjust CLI to use Services**
   - `tasksd plan` resolves adapters (filesystem, subprocess validators) and calls service.
4. **Executor Alignment**
   - Mirror approach for `slapsd` once planner is hexagonal.
5. **Testing Story**
   - Add table-driven tests for PlanningService using in-memory adapters.
6. **Migration Cleanup**
   - Remove directly-coupled helper functions from CLI layer once services stable.

## Open Questions
- Do we version ports (e.g., validator API) to track upstream external tool changes?
- Should artifact writer remain a shared adapter or become part of the domain core?
- How to stage executor hexagonalization without blocking validator-focused milestone?

## References
- docs/formal-spec.md §§3-15 — deterministic artifacts, validators, execution contract.
- docs/go-architecture.md §"Validators" / §"Resource modeling".
- README.md (stub CLI usage / validator flags).
