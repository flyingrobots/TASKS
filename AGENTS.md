# AGENTS.md

This file guides agents working in this repository (T.A.S.K.S. + S.L.A.P.S.). Its scope is the entire repository. Follow these instructions when reading, editing, or running code here.

## Purpose
- Build a pure precedence DAG planner (T.A.S.K.S.) and a rolling‑frontier executor (S.L.A.P.S.).
- Keep the DAG pure: precedence only. Model and enforce resources in the runtime, not in the graph.
- Emit deterministic, hash‑verified artifacts as the plan contract between planner and executor.

## Read First
- Core spec: `docs/v8/v8.md` (normative v8 design and contracts).
- Go architecture: `docs/go-architecture.md` (repository layout, types, and binaries).
- Overview: `README.md` (high‑level concepts and CLI examples).

## Repository Layout (active parts)
- Planner (Go): `planner/`
  - Canonical JSON: `planner/internal/canonjson/`
  - Hashing: `planner/internal/hash/`
  - Codebase census: `planner/internal/analysis/`
  - Canonical model types: `planner/internal/model/`
  - CLI (starter): `planner/cmd/tasksd/`
- Legacy and prior specs: `commands/`, `docs/` (reference only)
  
Note: The DAG Viewer has been split into its own repository and is not included here.

## Deterministic Artifacts (contracts)
Planner emits canonical, newline‑terminated JSON with sorted object keys:
- `features.json`, `tasks.json`, `dag.json`, `waves.json`, and `Plan.md` (+ `Decisions.md` where applicable)
- Executor contract: `coordinator.json`
- Hashing: SHA‑256 of canonical bytes (newline‑terminated). Record the hash in `meta.artifact_hash`. List all hashes in Plan.md "Hashes" section.
  - Hash preimage: compute over the artifact with `meta.artifact_hash` present but set to an empty string; then embed the resulting hex into `meta.artifact_hash` and write the final canonical bytes.

Rules (from v8):
- Sort map keys lexicographically at all depths.
- Preserve array order as computed; do not sort arrays.
- Minimal decimal rendering for numbers; UTF‑8; LF line endings; no trailing whitespace.
- Redact secrets in evidence excerpts before hashing using the redaction rules.

## DAG Purity (non‑negotiable)
- Structural edges only (technical, sequential, infrastructure, knowledge) are included in the DAG and topo sort.
- Resource edges (exclusivity, quotas, time windows) are recorded in `tasks.json` for traceability but excluded from `dag.json`.
- Waves are simulated from the DAG (Kahn layering) for preview only; do not feed resource constraints back into the graph.

## Running Things
- Canonicalize + hash any JSON (starter CLI):
  - `cd planner && go run ./cmd/tasksd canonical ./test.json`
- Run Go tests (where present):
  - `cd planner && go test ./...`
- Export DOT from planner artifacts:
  - `cd planner && go run ./cmd/tasksd export-dot --dag ./plans/dag.json --tasks ./plans/tasks.json --out ./plans/dag.dot`
  - Render with Graphviz: `dot -Tsvg ./plans/dag.dot -o ./plans/dag.svg`
 - Default directory mode (emits both if present):
  - `cd planner && go run ./cmd/tasksd export-dot --dir ./plans`  (writes `dag.dot` and `runtime.dot` next to JSON)
  - Label options: `--node-label id|title|id-title` (default id-title), `--edge-label none|type` (default type)
  - Example: `cd planner && go run ./cmd/tasksd export-dot --dir ./plans --node-label title --edge-label none`

- Stub planner (emits artifacts + DOT automatically):
  - `cd planner && go run ./cmd/tasksd plan --out ./plans`
  - Then render with Graphviz as needed.

## Todo Workflow (Hubless-style, markdown-first)
- Roadmap lives in `todo/` with milestones, features, and tasks as markdown files.
- Scripts (Node) manage task state and update progress in place:
  - Set active: `npm run todo:task:set-active -- T001`
  - Set finished: `npm run todo:task:set-finished -- T001`
  - Set merged: `npm run todo:task:set-merged -- T001`
  - Recompute progress: `npm run todo:update`

Conventions
- Tasks live at `todo/tasks/{milestone}/{backlog|active|finished|merged}/T###-*.md` with frontmatter keys: `id`, `milestone`, `features`, `title`, `status`, `deps`.
- Milestones have PROGRESS markers; the script updates those and the roadmap block in `todo/README.md`.
  
Visualization: Use the external DAG Viewer repository (not part of this repo).

## Coding Conventions
- Go version: `go 1.25.1` (see `planner/go.mod`).
- Package boundaries: code under `planner/internal/*`; avoid export leakage and import cycles.
- Keep changes minimal and focused. Update documentation alongside code when you change contracts or behavior.
- Do not add resource constraints into the DAG logic. Keep them adjacent for the executor.
- Evidence and acceptance checks are machine‑verifiable only; include type, source anchor, confidence, rationale.
- Logging: execution logs are JSONL with `{timestamp, task_id, step, status, message}` and error fields on failure.

## Implementation Status (today)
- Implemented:
  - Canonical JSON serializer: `planner/internal/canonjson` (sorts map keys; preserves arrays; newline‑terminated via `encoding/json` encoder).
  - SHA‑256 hashing: `planner/internal/hash`.
  - Codebase census: `planner/internal/analysis` (+ tests).
  - Starter CLI: `planner/cmd/tasksd` (reads a JSON file, prints canonical JSON and its hash).
- Not yet implemented (follow `docs/go-architecture.md`):
  - Feature extraction, task generation, dependency inference, DAG build + transitive reduction, wave simulation, validators integration, and executor.

## Quality Gates (v8)
- Acyclic and transitive‑reduced DAG.
- ≥95% of tasks/edges carry validated evidence.
- All tasks have at least one machine‑verifiable acceptance check.
- Verb‑first task titles (≥80%).
- Isolated nodes eliminated or justified.

## Agent Workflow Expectations
- Before edits: skim `docs/v8/v8.md` and the relevant package(s). Prefer ripgrep for search (`rg`). Read files in ≤250‑line chunks.
- When changing Go code: run `go test ./...` within `planner/` where tests exist. Keep diffs surgical; preserve existing style.
- When changing JSON contracts: re‑serialize via `canonjson` and update `meta.artifact_hash` accordingly.
- When adding planner logic: keep resource edges out of DAG construction; include them only in `tasks.json` metadata and the executor contract.
- When unsure about evidence or acceptance shape: consult `docs/v8/v8.md` and mirror the types in `planner/internal/model`.

## Do / Don’t
- Do: keep the contract boundary firm (`coordinator.json` is the planner→executor interface).
- Do: write small, auditable functions with clear input/output and unit tests.
- Don’t: push resource constraints into precedence. Don’t reorder arrays when canonicalizing.
- Don’t: emit artifacts without hashes or with un‑redacted secrets.

## Handy Paths
- Spec: `docs/v8/v8.md:1`
- Go architecture: `docs/go-architecture.md:1`
- Canon JSON: `planner/internal/canonjson/canonjson.go:1`
- Hashing: `planner/internal/hash/hash.go:1`
- Census: `planner/internal/analysis/census.go:1`
- Model: `planner/internal/model/*.go`
- CLI: `planner/cmd/tasksd/main.go:1`

## Suggested Next Steps (if implementing)
- Add minimal‑decimal number normalization to `canonjson` to fully meet v8.
- Implement DAG package (cycle detection, Kahn topo, transitive reduction) under `planner/internal/planner/dag`.
- Add `wavesim` preview using antichains; keep output in `waves.json` only when DAG ok.
- Introduce `schemas/` with JSON Schema for artifacts and validate before hashing.
- Wrap acceptance/evidence/interface validators as subprocess clients.
 - Wire `tasksd plan` to call `export-dot --dir <out>` after artifacts are written so DOT files always accompany a plan.

---
This AGENTS.md is normative for agents modifying this repository. When in doubt, favor DAG purity, determinism, and auditability.
