# Tasks-of-TASKS (Planner + Executor)

This table maps the TASKS feature plan to existing GitHub issues (where they exist today), and records planner/executor precedence plus runtime resource notes. “Blocking” lists downstream tasks this task unlocks; “Blocked By” lists prerequisites. Runtime resource columns describe executor-time constraints (planner-time tasks are N/A).

| Task Name | Feature Group | GH Issue # | Blocking | Blocked By | Runtime Resource Dependencies | Runtime Resource Exclusivity Policy |
| --- | --- | --- | --- | --- | --- | --- |
| T1 — Canonical JSON minimal-number normalization | F-PLN (Planner Core) | — | T2, T8, T10 | — | — | — |
| T2 — Preimage hashing helper + validator | F-PLN | — | T8, T12, T20 | T1 | — | — |
| T3 — DAG builder (filter, cycle detect, transitive reduction) | F-PLN | — | T4, T5, T9, T10 | T1 | — | — |
| T4 — Metrics (width, density, critical path, verb-first, evidence) | F-PLN | — | T9, T12, T20 | T3 | — | — |
| T5 — Waves preview (Kahn layers; preview only) | F-PLN | — | T10 | T3 | — | — |
| T6 — Validator integration (acceptance/evidence/interface) | F-PLN | #21, #27 | T19, T20, T10 | T1, T2 | — | — |
| T7 — Evidence redaction rules + tests | F-PLN | #29 | T19 | — | — | — |
| T8 — JSON Schemas (features/tasks/dag/waves/coordinator) | F-CRT (Contracts) | #28 | T10, T12 | T1, T2 | — | — |
| T9 — Coordinator estimator consistency (p50, width, longest path) | F-CRT | — | — | T3, T4 | — | — |
| T10 — tasksd CLI (plan, validate, export-dot, canonical) | F-CLI | #31, #30 | T12 | T1–T6, T8 | — | — |
| T11 — DOT export styling options | F-CLI | #20 | T12 | T10 | — | — |
| T12 — CI jobs: plan/validate/comment (inline SVG) | F-CLI | #31, #39 | — | T8, T10, T11 | — | — |
| T13 — Ports package (planner/executor) | F-HEX (Ports) | #15, #18, #19 | T14, T15–T18 | — | — | — |
| T14 — Adapters (doc parse, census, validators, artifacts) | F-HEX | #20 | T10, T15–T18 | T13 | — | — |
| T15 — Coordinator loader + topo gating | F-EXE (Executor) | #45 | T17 | T13, T14 | — | — |
| T16 — Dispatcher & Worker components (skeletons) | F-EXE | #46 (interface) | T17 | T13 | workers | up to concurrencyMax (e.g., 5) |
| T17 — Rolling frontier loop + arbitration hooks | F-EXE | #40, #22–#25, #48, #47 | T18 | T15, T16 | locks, quotas, workers | locks: exclusive (global order); quotas: capacity; workers: up to concurrencyMax |
| T18 — Circuit breakers + retry/backoff | F-EXE | #52–#55 | — | T17 | — | — |
| T19 — Quality gate policy (≥95% evidence; ≥80% verb-first) | F-GOV (Gates) | #27, #17 | T12 (gating) | T6, T7 | — | — |
| T20 — Plan.md synthesis (Hashes + Validators + Gaps) | F-GOV | #31 | — | T4, T6 | — | — |
| T21 — Planner→GitHub sync (mirror plan to issues/projects) | F-GH (Sync) | #32–#39 | T22 | T3, T4, T8, T10 | — | — |
| T22 — Overview DAG improvements (Issues) | F-GH | #35, #39 | — | T21 | — | — |

Notes
- GH Issue # column maps to currently visible issues in this repo’s project boards (e.g., #28 “Tighten schemas consistently”, #46 “X3 — Dispatcher + Worker interface”, etc.). Where no specific issue exists yet, the cell is left blank (—).
- Blocking/Blocked By are task-level (Txx) precedence edges per T.A.S.K.S. DAG purity; resource constraints appear only in the runtime columns and do not affect precedence.
- Runtime resource entries describe executor-time behavior. Planner/contract/schema/CLI tasks show “—”.
- Example policies: “exclusive” (one-at-a-time), “capacity” (limited parallelism), “up to concurrencyMax” (global runtime setting).

