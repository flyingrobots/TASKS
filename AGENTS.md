<<<<<<< HEAD
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

## Task Execution Workflow (Git + todo CLI + PR)
Follow this sequence for every task. Make small, frequent commits.

1. Verify a clean working tree: `git status` (commit/stash before starting).
2. Pick a task ID from `todo/tasks/**` (e.g., `T070`). Read its frontmatter and linked feature/milestone.
3. Mark it active via the CLI: `npm run todo:task:set-active -- T070`.
4. Create a task branch off `origin/main`:
   - `git fetch origin`
   - `git checkout -B feat/{feature}-task-{taskid}` (e.g., `feat/validators-task-T070`).
5. Commit the move/state change: `git add -A && git commit -m "todo: set T070 active"`.
6. Re-read the task + feature/user story; update the test plan if needed (edit the task/feature doc).
7. Write tests first. Commit: `git add -A && git commit -m "T070: add tests"`.
8. Run tests. If they pass, proceed to step 10. If they fail, go to step 9.
9. Implement the task in small steps. Make micro-commits as you go. Re-run tests until green.
10. Update docs: task notes, feature rationale, how it was implemented, anything learned.
11. Commit docs: `git add -A && git commit -m "T070: docs + notes"`.
12. Mark the task finished: `npm run todo:task:set-finished -- T070`.
13. Commit + push + open PR:
    - `git add -A && git commit -m "todo: set T070 finished"`
    - `git push -u origin HEAD`
    - `gh pr create --fill --base main --head $(git branch --show-current)` (or `npm run todo:pr:create`)
14. Await review; address feedback with incremental commits. When merged, optionally `npm run todo:task:set-merged -- T070` on `main` and commit the move.

Notes
- Always keep branches focused on a single task.
- If the task reveals necessary subtasks, add new task files under `todo/tasks/...` and link them from the feature.
- Use `npm run todo:update` to refresh progress bars after state changes.
  
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

## PR #7 Debrief (feat/planner-dag-validate)

What we changed (high level)
- Planner core:
  - Canonical JSON minimal-number normalization (e.g., 1.0→1, 1.20e+03→1.2e3; 0e10→0).
  - DAG builder: cycle detect, Kahn layering, transitive reduction; guard duplicate task IDs.
  - Wave simulation (preview): layer by Kahn depth; split by exclusive resources; returns error on missing task IDs; integrated into tasksd plan.
  - Coordinator/waves/features/dag JSON Schemas tightened (required fields, minLength, additionalProperties policies; explicit metrics schema for dag).
  - Artifact writer: documented preimage hash policy; panic guard for setHash; writer errors are aggregated and reported after all writes.
  - Repo census: removed debug/test-only comments and unnecessary per-directory reads; tests relaxed for OS portability; added typed FileCensusCounts and logging on failure.
- CLI + DX:
  - export-dot refactored into helpers; directory/single-file modes reuse common loaders.
  - todo/ workflow scaffold with Node helpers: branch + PR commands; task indexing for O(1) lookup; root + ID validation; robust mkdir/IO handling; switched to gray-matter frontmatter.
  - Roadmap: milestones/features/tasks seeded; examples added (success/error/idempotent); markdown hygiene across docs/milestones.

Why these changes
- Made the planner outputs deterministic and auditable (preimage hashing, schema validation, canonical numbers).
- Enforced DAG purity (resource edges traced in tasks.json, not fed into dag.json) while still previewing waves.
- Improved safety and ergonomics for everyday workflow (todo/ + helpers) without requiring a full Hubless install.

Notable decisions
- Preimage hash policy is intentional: meta.artifact_hash reflects canonical bytes with the hash field left empty (documented here, enforced by writer helper and validated).
- Wave simulation is strictly a preview; the executor remains the authority on resources.

Pitfalls we found and fixed
- Silent duplicate task IDs in DAG builder (now fails fast with context).
- brittle frontmatter parsing in todo script (replaced with gray-matter; added validation & error handling).
- Over‑permissive schemas (tightened; added explicit metrics typing for dag).
- Tests assumed uniform permission semantics across OSes (relaxed to remain meaningful and portable).

Process notes
- Followed the task workflow: mark task active → branch → tests/impl → docs → mark finished → push → PR.
- Batched PR feedback in small commits, updated the feedback checklist as items were completed.

Next steps (suggested)
- Finish remaining PR #7 checklist items (continue small batches):
  - Further schema tightening where applicable and unify additionalProperties policies.
  - Minor code health passes (smaller helpers, comments, error messages).
- Start M3 (validators): adapters for acceptance/evidence/interface and gate the plan accordingly; add flags + validation cache.
- Wire a CI job to run tasksd validate + schema checks and block on drift; add a smoke test that plans a tiny spec and renders DOT.
- Optional: add a tiny JSON export from todo/ into a Hubless-compatible structure; later, replace with Hubless CLI once ready.

## Automation Session Log
{"date":"2025-10-07","time":"13:00","summary":"Created a pull request to add Gemini GitHub workflows to the repository.","topics":[{"topic":"Git Repository Setup for Gemini Integration","what":"Configured the git repository to integrate Gemini by updating .gitignore, adding GitHub workflows, and creating a new feature branch and pull request.","why":"The user wanted to set up the repository for Gemini integration.","context":"The user initiated a session to configure the git repository. This involved adding files to .gitignore, adding new GitHub workflows, and creating a new branch and pull request.","issue":"The repository was not configured for Gemini integration.","resolution":"I updated the .gitignore file, added the Gemini GitHub workflows, created a new branch 'feat/gemini-github', and opened a pull request with these changes.","future_work":"The pull request needs to be reviewed and merged.","time_percent":100}],"key_decisions":["Added .gemini/ and gha-creds-*.json to .gitignore.","Created a new branch 'feat/gemini-github'.","Added new GitHub workflows for Gemini.","Created a pull request for the changes."],"action_items":[]}
{"date":"2025-10-09","time":"12:44","summary":"Hardened validator reporting/caching and refreshed supporting tests while keeping PR #10 moving.","topics":[{"topic":"Validator pipeline hardening","what":"Aligned validator runner/cache with shared model types and error handling","why":"Address PR #10 feedback about inaccurate statuses and silent cache failures","context":"Validator integration for tasksd plan outputs","issue":"Reports fabricated pass statuses and cache silently dropped IO errors","resolution":"Aliased report type, normalized statuses, propagated cache read/write errors","future_work":"Monitor for additional validator schema changes or integration feedback","time_percent":60},{"topic":"Test + CLI reliability","what":"Stabilized validator integration tests and CLI validator harness","why":"Original tests were brittle regarding go toolchain paths and timeout detection","context":"planner/internal/validators and tasksd validator tests","issue":"Timeout checks misread signals and helper builds polluted logs","resolution":"Added deadline-aware assertions, repo-root discovery, and quieter build helpers","future_work":"Extend coverage for evidence/interface validator variants once fixtures exist","time_percent":40}],"key_decisions":["Alias validator reports to shared model","Persist cache errors instead of silent overwrites"],"action_items":[{"task":"Review remaining PR #10 feedback beyond validator scope","owner":"james"}]}
{"date":"2025-10-10","time":"05:48","summary":"Completed hexagonal migration for planner and executor fronts, exposing default services and adapter ports.","topics":[{"topic":"Planner Hex Ports","what":"Added default service factory with doc, census, deps, waves, coordinator, artifacts, validators","why":"Finalize planner-side hex architecture","context":"feat/hex-planning-service","issue":"CLI still owned orchestration and manual adapter wiring","resolution":"Introduced plan.NewDefaultService and rewired tasksd to use it","future_work":"Iterate on validator ergonomics and port adapters as they grow","time_percent":45},{"topic":"Executor Hex Skeleton","what":"Created exec service, default factory, and slapsd stub","why":"Extend hex architecture to executor","context":"feat/hex-planning-service","issue":"Executor lacked service layer and ports","resolution":"Added coordinator loader, runtime adapters, and tests; slapsd now delegates to service","future_work":"Implement real rolling-frontier loop and adapters","time_percent":35},{"topic":"Documentation & Tests","what":"Updated architecture docs and added factory coverage","why":"Keep guidance current with new wiring","context":"docs/go-architecture.md & tests","issue":"Docs/tests lagged behind new design","resolution":"Documented default service wiring and added factory tests","future_work":"Document executor runtime once loop lands","time_percent":20}],"key_decisions":["tasksd/slapsd rely on default service factories","Executor hex scaffold returns ErrLoopNotImplemented until runtime lands"],"action_items":[{"task":"Implement real slapsd execution loop","owner":"james"}]}
