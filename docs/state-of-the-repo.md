# State of the Repo™ Report

Date: 2025-10-31

Scope: Comparison of the implemented codebase against these documents:
- docs/formal-spec.md
- docs/go-architecture.md
- docs/hex-architecture-plan.md
- docs/cli-demo.md

Repository root: `.` (planner + executor under `planner/`)

---

## Executive Summary

- Overall, the planner (T.A.S.K.S.) aligns strongly with the v8 specification on determinism, DAG purity, and artifact contracts. Canonical JSON with minimal-number normalization, preimage hashing, schema validation, DAG build + transitive reduction, and a preview waves generator are all in place and tested.
- The executor (S.L.A.P.S.) is scaffolded with a service and CLI stub and returns a not‑yet‑implemented error for the rolling‑frontier loop. Resource constraints are present only as task metadata and used for preview waves, consistent with DAG purity.
- The codebase has begun a hexagonal migration: a `plan.Service` composes adapters (doc loader, census, DAG builder, waves, validators, artifact writer). Ports are expressed as function fields rather than concrete interfaces, which achieves decoupling but diverges slightly from the “explicit ports module” described in docs.
- Documentation gaps remain (notably `docs/cli-demo.md`), and some schema constraints are intentionally looser than the “tighten all additionalProperties” goal.

Status (high level):
- Green: deterministic artifacts, DAG builder, wave preview, schemas + validation, validators integration, CLI (`tasksd`) surface, DOT exports, repo census.
- Yellow: hex ports formalization (function fields vs interfaces), CLI docs (demo), schema hardening across all artifacts, validator quality gates as hard policy.
- Red: executor runtime loop (rolling frontier), evidence redaction pipeline, telemetry/provenance ledger, resource arbitration at runtime.

---

## What’s Implemented vs Spec

### Deterministic Artifacts (spec: docs/formal-spec.md §§3, serialization rules)
- Canonical JSON: sorted map keys, arrays preserved, UTF‑8, LF, minimal‑decimal numbers, newline termination.
  - Implementation: `planner/internal/canonjson/canonjson.go:1`
- Hashing policy (preimage with `meta.artifact_hash` empty, then embed):
  - Implementation: `planner/internal/emitter/writer.go:1`
- Artifacts emitted: `features.json`, `tasks.json`, `dag.json`, `waves.json`, `coordinator.json`, plus `Plan.md` with a Hashes section.
  - Verified by running `tasksd plan` (see Repro appendix). Example plan: `planner/plans/Plan.md:1`.
- JSON Schemas embedded and validated (tightened requireds, allowed fields):
  - Schemas: `planner/internal/validate/schemas/*.schema.json`
  - Validation helpers: `planner/internal/validate/validate.go:1` (invoked by `tasksd validate`).

Assessment: ✅ Meets v8 requirements. Minimal-number normalization is implemented and exercised. Preimage hashing is documented and coded. Validation CLI is present.

### DAG Purity + Build (spec: DAG purity, acyclic + reduced, structural edges only)
- Structural-only edges in `dag.json`; resource edges are excluded from the graph and kept in `tasks.json` metadata.
  - Filter logic: `planner/internal/planner/dag/dag.go:1` (drops `resource` edges; confidence + hardness filters).
- Guarantees: cycle detection (DFS), Kahn layering (`depth`), critical path, transitive reduction, deterministic node/edge ordering.
  - Implementation: `planner/internal/planner/dag/dag.go:1`.
- Resource metadata kept adjacent to tasks (exclusive + limited):
  - Types: `planner/internal/model/task.go:1`, `planner/internal/model/resource_need.go:1`.

Assessment: ✅ DAG purity enforced; topo and reduction implemented and tested (`planner/internal/planner/dag/dag_test.go:1`).

### Waves Preview (spec: preview only; no feedback into DAG)
- Waves generator groups by Kahn depth and splits subwaves by exclusive resources; errors on missing task IDs; arrays stable.
  - Implementation: `planner/internal/planner/wavesim/wavesim.go:1`.
- Exposed via `plan.Service` and written to `waves.json`; DOT export produces `runtime.dot` alongside `dag.dot`.

Assessment: ✅ Preview semantics match spec; kept out of DAG and used as non‑authoritative preview.

### Validators Integration (spec: acceptance/evidence/interface; gating)
- Subprocess-based validator runner with caching, timeouts, and status normalization.
  - Code: `planner/internal/validators/*.go`
- `tasksd plan` accepts validator command flags and can operate in “strict” mode (convert failed reports to errors) or warning mode.
  - Flags wired in `planner/cmd/tasksd/main.go:220`.
- Reports are recorded under `tasks.json.meta.validator_reports` and included in the plan result.

Assessment: ✅ Integration present with cache + strict mode. ⚠️ Spec’s quality gates (e.g., “≥95% tasks/edges carry validated evidence”) are not hard‑enforced beyond strict mode; they’re policy rather than schema.

### Codebase Census (spec: repo analysis present in tasks metadata)
- Analyzer exists with typed counters and tests; results are embedded under `tasks.json.meta.codebase_analysis`.
  - Code: `planner/internal/analysis/census.go:1`.

Assessment: ✅ Implemented.

### Hexagonal Architecture Migration (spec: docs/hex-architecture-plan.md)
- Planning service exists and composes adapters (doc loader, census, deps, DAG, waves, validators, artifacts).
  - Service: `planner/internal/app/plan/service.go:1`
  - Factory: `planner/internal/app/plan/service_factory.go:1`
- Adapters provided:
  - Doc loader (markdown → tasks/specs): `planner/internal/app/plan/doc_loader.go:1`
  - DAG builder delegates to core DAG: `planner/internal/planner/dag/dag.go:1`
  - Waves builder delegates to core wavesim: `planner/internal/planner/wavesim/wavesim.go:1`
  - Validators runner + artifact writer wired via factory.
- CLI `tasksd` delegates orchestration to the service for `plan`, as intended.

Assessment: ✅ Hex migration largely in place on planner side. ⚠️ Ports are expressed as function fields on `Service` rather than explicit interfaces/types under a `ports/` package; still effective but differs from doc’s “ports module” recommendation.

### Executor (S.L.A.P.S.)
- CLI stub and service created; returns `ErrLoopNotImplemented`.
  - CLI: `planner/cmd/slapsd/main.go:1`
  - Service: `planner/internal/app/exec/service.go:1`, `service_factory.go:1`
- No rolling‑frontier runtime loop, resource arbitration, circuit breakers, or provenance ledger yet.

Assessment: ❌ Not implemented beyond scaffold.

### CLI Surface vs docs/cli-demo.md
- `tasksd` implements: `canonical`, `export-dot`, `plan`, `validate` (with node/edge labels and directory mode for DOT).
  - Source: `planner/cmd/tasksd/main.go:1`
- `docs/cli-demo.md` is currently a placeholder and should be backfilled with examples that reflect the current CLI.

Assessment: ⚠️ Functionality exists; documentation needs completion.

---

## Gaps and Deviations

1) Executor runtime (rolling frontier)
- Missing: task readiness frontier, global lock/queue ordering, quota/time‑window resources, retries/backoff, circuit breakers, hot‑patching, provenance ledger and JSONL execution logs.
- Impact: planner outputs are usable for visualization and audit; execution must be completed to realize end‑to‑end.

2) Evidence redaction pipeline
- Spec calls for redacting secrets in evidence excerpts before hashing. No redaction helpers found.
- Search: no `redact`/`redaction` implementation under `planner/`.

3) Schema hardening consistency
- Some schemas are tight (e.g., `dag.schema.json` nodes/edges with `additionalProperties: false`). Others (e.g., `tasks.schema.json`) allow additional properties at object levels to ease iteration.
- Action: adopt a unified policy per docs and progressively lock down.

4) Ports formalization
- Docs suggest an explicit `internal/ports` package. Current design uses function fields on `plan.Service`, which is pragmatic but less explicit for consumers.

5) CLI documentation
- `docs/cli-demo.md` is a placeholder; the implemented CLI offers multiple workflows not yet captured in docs.

6) Quality gates as enforceable policy
- Spec’s target gates (acyclic and transitive‑reduced DAG; ≥95% tasks/edges with validated evidence; ≥80% verb‑first titles; isolated nodes eliminated or justified) are not enforced as hard errors (beyond DAG validity and strict validator mode). A “lint/validate” pass could compute and gate these metrics.

---

## Risks and Mitigations

- Risk: Without executor, users may over‑interpret waves as a scheduler. Mitigation: keep “preview only” language prominent (already done) and document runtime responsibilities in README.
- Risk: Looser schemas on tasks/dependencies may let drift creep in. Mitigation: progressive hardening with explicit migrations and CI schema validation.
- Risk: Missing redaction may leak secrets in artifacts. Mitigation: add redaction helper and pre‑hash sanitizer and test with fixtures.

---

## Recommendations (Next Steps)

Short term (high leverage)
- Add a “policy-lint” in `tasksd validate`: compute and print quality gate metrics; support `--strict-policy` to fail under thresholds.
- Tighten JSON Schemas incrementally (align `additionalProperties` across artifacts; add patterns/enums for fields).
- Fill `docs/cli-demo.md` with working examples matching the current CLI.
- Implement evidence redaction helpers and invoke them before hashing.

Medium term
- Formalize port interfaces in a small `internal/ports` package while keeping function injection for tests; provide default adapters.
- Finish `slapsd` minimal rolling‑frontier loop: readiness queue, simple exclusive resource locks (global order), retry/backoff, and JSONL telemetry (append‑only).
- Wire CI to run `tasksd plan` on a tiny fixture + `tasksd validate` + DOT render to SVG as a smoke test.

Longer term
- Resource arbitration features (quotas/time windows/profiles) in executor; provenance ledger; hot‑patch injectors.
- Optional HTTP validator adapter and caching strategy pluggability.

---

## Repro Appendix (verified today)

- Run tests:
  - `cd planner && go test ./...`
- Create stub plan, validate, and export DOTs:
  - `cd planner && go run ./cmd/tasksd plan --out ./plans --repo ..`
  - `cd planner && go run ./cmd/tasksd validate --dir ./plans`
  - `cd planner && go run ./cmd/tasksd export-dot --dir ./plans`
- Canonicalize and hash demo JSON:
  - `cd planner && go run ./cmd/tasksd canonical ./test.json`

---

## Notes on Alignment with Design Docs

- formal-spec.md: Determinism, DAG purity, artifact set, and hash policy are implemented as specified. Validators are integrated; policy thresholds remain advisory.
- go-architecture.md: The planner is organized by internal packages with a service layer; CLI delegates to services; DOT export helpers exist. Executor scaffold exists; runtime behaviors remain to be implemented.
- hex-architecture-plan.md: Planner’s hexagonal service and adapters exist; an explicit ports package could close the last gap. Executor hex scaffold created with `ErrLoopNotImplemented` sentinel.
- cli-demo.md: Needs content; the implemented CLI provides the examples this doc intends to show.

