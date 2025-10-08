# Project: T.A.S.K.S. + S.L.A.P.S. v8.0

## Project Overview

This project, "T.A.S.K.S. + S.L.A.P.S.", is a sophisticated system for project planning and execution. It consists of two main components:

*   **T.A.S.K.S.**: A planning engine that analyzes a project specification and codebase to generate an optimized and audit-ready plan. The plan is represented as a directed acyclic graph (DAG) of tasks, along with a set of deterministic JSON artifacts.
*   **S.L.A.P.S.**: An execution runtime that consumes the plan from T.A.S.K.S. and executes it using a rolling frontier scheduler, resource arbitration, and self-healing mechanisms.

The project's core principle is **DAG purity**, which means that the dependency graph only represents precedence, while resource constraints are managed separately by the runtime. This approach maximizes parallelism and ensures the plan's mathematical validity.

The system is designed to be used as a command-line tool, through an API, or as a Claude command.

## Key Artifacts

The T.A.S.K.S. engine produces the following canonical artifacts:

*   `features.json`: A breakdown of features with priorities.
*   `tasks.json`: A list of tasks with detailed information, including acceptance criteria, resource requirements, and dependencies.
*   `dag.json`: The dependency graph of the tasks.
*   `waves.json`: A simulation of parallel execution waves.
*   `Plan.md`: A human-readable roadmap.
*   `coordinator.json`: The execution contract for S.L.A.P.S., which includes the dependency graph, resource catalog, and runtime policies.

## Building and Running

This repository is Go‑first for the planner and shared model.

### Planner (Go)

Basic demo (canonicalize + hash JSON):

```bash
cd planner
go run ./cmd/tasksd ./test.json
```

Run tests (where present):

```bash
cd planner
go test ./...
```

As the planner grows, the `tasksd` CLI will add subcommands per docs/go-architecture.md (e.g., `plan`, `validate`).

### DAG Viewer

The DAG viewer has moved to its own repository and is not included here. Use the external viewer to visualize `dag.json` and `waves.json` outputs produced by the planner.

## Development Conventions

- Go version: see `planner/go.mod`.
- Packages live under `planner/internal/*` to avoid public API drift.
- Canonical JSON serialization must sort object keys, preserve array order, and emit LF‑terminated output.
- Normative validators (acceptance, evidence, interface) are treated as external tools/services with strict JSON I/O; the planner/runtime integrate with them but do not re‑implement their logic.
