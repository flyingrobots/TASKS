# PR Review Checklist (Grouped by File)

## .github/workflows/overview-dag.yml
- [x] [Trivial] YAML linter flags spacing in the branch list; keep the existing "on:" trigger but remove the extra spaces inside the brackets so the branches line reads with no spaces inside the brackets (changed "[ main ]" to "[main]") to satisfy yamllint formatting.
- [x] [Trivial] fix branches array spacing ([ main ] -> [main])
- [x] [Trivial] remove redundant self-copy in commit step (cp -f docs/overview-dag.svg docs/overview-dag.svg)
  - Implemented: deleted the no-op cp line; the generated file is added directly.

## .github/workflows/overview-dag.yml around lines 22-23
- [x] [Major] the workflow calls bash scripts/gh/overview-dag.sh without any validation or error handling; update the step to export required env (GITHUB_TOKEN, REPO), add "set -e" to fail fast, verify the script exists (if [ ! -f scripts/gh/overview-dag.sh ] exit with an error message), verify required tools like "dot" exist (command -v dot || exit with an error message), then run the script and ensure you check its exit status so failures produce clear logs and cause the job to fail.
  - Implemented: added a “Preflight checks” step (verifies script + dot), and guarded the generation step with set -e.
  - [x] (follow-up): add a step to echo Graphviz version and file sizes of generated artifacts for easier debugging.

## .github/workflows/overview-dag.yml around lines 24 to 30
- [x] [Major] the workflow uploads DAG artifacts even when generation may have failed, leading to missing or partial files being stored; update the upload step to run only on overall job success by adding an if: success() condition to the upload step, and add a preceding verification step that checks the two files exist (use shell test -f docs/overview-dag.dot and test -f docs/overview-dag.svg and exit non‑zero with a clear message if either is missing) so the job fails before upload when artifacts are absent.
  - Implemented: added a “Verify DAG outputs” step, and `if: success()` to artifact upload.
  - [x] (follow-up): add a small summary step printing the node/edge counts parsed from the dot file.

## Follow-ups processed on 2025-11-02
- [x] overview-dag.sh: derive node declarations from the same paginated GraphQL dataset used for edges; remove GH_LIST_LIMIT mismatch risk and keep `gh issue list` only as a light metadata fetch.
- [x] Rendering fallbacks: replace brittle `@viz-js/viz` invocation with `@hpcc-js/wasm-graphviz-cli` primary fallback, keep `viz.js-cli` and `graphviz-cli` as secondary fallbacks; use stdout redirection form `-T svg input.dot > output.svg` to avoid "Missing output output format"/ENOENT errors.
- [x] Workflow temp files: write PR comment markdown to `${{ runner.temp }}/dag-comment.md` and read from that path in the action.
- [x] Quoting improvement: ensure Graphviz version echo uses a safe pipeline and strips CRs `$(dot -V 2>&1 | head -n1 | sed 's/\r$//')`.
- [x] Exec service tests: rename ambiguous tests, add sentinel-based assertions and a table test for missing adapters; verify error precedence and context cancellation.
- [x] doc_loader defaults: make execution logging `requiredFields` additive/idempotent; stop clobbering `Compensation.Idempotent`; harden acceptance commands for portability and guard missing env/commands.

## docs/state-of-the-repo.md
- [x] [Trivial] add a single blank line before each heading like "Assessment:" so the heading is separated from the prior paragraph; then in the "Gaps and Deviations" block normalize the ordered list to use consistent "1." prefixes for every item (or alternatively update the markdownlint config to allow sequential numbering) so the list complies with the linter; run markdownlint to verify and commit the cleaned file.
  - Implemented: inserted blank lines and normalized numbering; ran local checks.
  - [ ] (follow-up): add a simple markdownlint config to the repo to keep style consistent in future edits.

## planner/internal/app/exec/service_factory_internal_test.go
- [x] [Major] replace inline/fallback fixtures with realistic coordinator JSON (nodes, edges, resources, policies, metrics) and strengthen assertions beyond Version.
- [x] [Major] add negative cases: invalid JSON, missing fields (zero-values), wrong types, and empty payload.
  - [ ] (follow-up): add schema validation of coordinator and corresponding unit tests once executor adds validation.

## planner/internal/app/exec/service_test.go
- [x] [Minor] current tests only cover InitRuntime and RunLoop single-error propagation; add unit tests for the missing edge cases: (1) LoadCoordinator returns a zero-value or invalid coordinator with no error — assert Run returns an appropriate validation error and InitRuntime is not called; (2) context cancellation — add tests that cancel the context before calling Run and that simulate cancellation during InitRuntime and during RunLoop (InitRuntime should return ctx.Err() and Run should surface context.Canceled without entering subsequent phases); (3) simultaneous/multiple errors — add a test where InitRuntime and RunLoop both return errors and assert which error Run returns (define expected precedence in the test), and ensure mocks track whether subsequent phases were invoked; implement these tests using small inline function mocks, context.WithCancel/WithTimeout to simulate cancellation, and table-driven subtests for clarity.
  - Implemented: added validation for empty version, context cancellation tests (before/during), and precedence test asserting init error dominates and RunLoop not invoked.
  - [ ] (follow-up): extend validation scope once executor gains richer coordinator checks.
- [x] [Minor] test only checks that an error occurred but does not assert the error message; update the assertion to verify the error content (e.g., confirm err is non-nil and that err.Error() contains or equals "loop failed", or use errors.Is/errors.As if the error may be wrapped) so the test fails if the specific loop failure message is not propagated.

## planner/internal/app/plan/analyzer_test.go
- [ ] [Major] there is a blank identifier assignment used solely to keep analysis.CodebaseAnalysis imported; remove the dead import and the `_ = analysis.CodebaseAnalysis{}` line if the type is not needed, or if the type should be referenced for documentation/compile-time guarantees, replace the hack with a real usage in the test (e.g., declare a named variable of that type or assert its zero value in a test helper) so the import is legitimately used.

## planner/internal/app/plan/artifact_writer_test.go
- [ ] [Minor] add explicit length assertion for truncated validator detail (e.g., assert count of 'x' <= 2000)
- [ ] [Minor] verify truncation actually limits output length (assert count of 'x' <= 2000)

## planner/internal/app/plan/coordinator_adapter.go
- [ ] [Critical] implement metrics in Build() (p50_total_hours, longest_path_length, width_approx) and populate coordinator.json

## planner/internal/app/plan/dependency_adapter.go
- [ ] [Critical] Fallback edge creation doesn't guard against empty task IDs (add empty-ID guard to fallback edges or validate upstream)
- [ ] [Critical] code guards against empty task IDs before processing resources but does not apply the same check when creating fallback edges; update the fallback edge creation code to either skip creating edges when task.ID == "" or validate/return an error earlier so no fallback uses an empty ID, ensuring every use of task.ID is consistently guarded; modify the fallback edge path to check task.ID non-empty (or propagate the guard) before building or appending edges.

## planner/plans/Plan.md
- [ ] [Critical] coordinator.json currently has a null meta object and therefore lacks meta.artifact_hash; regenerate or update coordinator.json to include a meta object with "artifact_hash": "" (i.e. {"meta":{"artifact_hash":""},...}) and ensure config/graph remain intact, then compute the canonical JSON byte representation (use the project's canonical serializer: artifact_hash field left as empty string for preimage), take the SHA-256 of those bytes, write that resulting hex into coordinator.json.meta.artifact_hash, and finally replace the coordinator.json hash entry in Plan.md with the newly computed hex value so the artifact and Plan.md match.

## planner/plans/coordinator.json
- [ ] [Critical] change one or two acceptance_checks (different command, timeout or add a script/file check), add resource constraints to at least one node (cpu, memory, disk or network tags), and alter one compensation to non-idempotent to cover rollback behavior; ensure execution_logging required_fields remain valid or add a node-specific extra field so tests exercise differing logging schemas.
- [ ] [Critical] lines 1-142: the top-level "meta.artifact_hash" field required by the planner→executor contract is missing; add a top-level "meta" object with "artifact_hash" set to the SHA-256 hex of the canonical bytes computed with meta.artifact_hash present but set to an empty string (preimage hash policy). To fix: insert "meta":{"artifact_hash":""} into the JSON, then compute the canonical byte representation (deterministic JSON: stable key ordering, no insignificant whitespace, UTF-8) of the entire document with artifact_hash set to "" and compute SHA-256 over those bytes, then replace artifact_hash with the resulting hex digest; ensure the final file is saved using the same canonicalization rules so the hash validates.
- [ ] [Critical] lines ~13-21) never sets Metrics; fix by computing and populating Metrics in the Build method: call the existing DAG functions to obtain longest_path_length and width_approx, implement and call a new P50TotalHours computation (e.g., aggregate task duration estimates, derive the 50th-percentile total plan duration using existing task duration distributions or conservative deterministic durations if distributions are unavailable), and set metrics.estimates.{longest_path_length,p50_total_hours,width_approx} on the returned coordinator object before returning; ensure the dag.go metric functions are exposed/used or moved into a shared package and unit-tested so coordinator.json no longer contains zero stubs.
- [ ] [Trivial] reduce copy-paste — vary durations, acceptance_checks, resources; avoid identical node configs

## planner/plans/dag.json
- [ ] [Critical] resolve quality gate failures — raise evidence_coverage ≥ 0.95 and fix verb_first_pct to reflect titles; update analyzer or artifacts accordingly

## planner/plans/features.json
- [ ] [Critical] meta.artifact_hash value is incorrect; replace the existing hash string with the correct SHA-256 value "eaf4b5f4c07120d59c768bab50694249b3e0ee975588c01588f8931b73ea00ad" and commit the change; after updating, verify by computing the SHA-256 over the canonical JSON bytes with meta.artifact_hash set to an empty string and ensure the result equals the provided hash.

## planner/plans/tasks.json
- [ ] [Critical] "echo ok") which provides no real validation; replace each task's stub check with concrete command-based checks appropriate to the task (for T001 ensure DB connectivity and that the expected schema/table exists, for T002 verify latest migration version is present in schema_migrations, for T003 hit the service HTTP health endpoint and assert expected JSON status), set sensible timeoutSeconds (e.g. 5–10s) and keep type "command", and ensure the commands exit non-zero on failure so the acceptance check actually fails when the task is not completed.
- [ ] [Critical] around line 62 the "source_evidence" field is null for every task which causes evidence_coverage to be 0; update each task to reference concrete evidence (GitHub issues, docs, code analysis or test results), populate source_evidence objects with fields like type (e.g. "github_issue" or "code_analysis"), url or files, an excerpt/summary, and a confidence score (0.0–1.0); ensure references point to real existing issues/docs/code, update all tasks so ≥95% have non-null validated evidence, then re-generate the DAG and verify evidence_coverage ≥ 0.95 before committing.
- [ ] [Major] every task compensation is marked "idempotent": true without justification; update each task's compensation to reflect its real idempotency and strategy: for each task (T001, T002, T003) set "idempotent" correctly (true only if operations are safe to rerun), add a "strategy" field describing how idempotency is achieved (e.g., "check_before_create", "versioned_migrations", "blue_green_deploy"), and include rollback metadata like "rollback_cmd" or "requires_manual_cleanup": true for non-idempotent operations; ensure comments or brief rationale are added for each to document why the chosen value is correct.
- [ ] [Minor] PERT duration values are identical for all tasks which indicates copy‑paste placeholders; update each task with realistic, task‑specific PERT estimates (optimistic/mostLikely/pessimistic) reflecting actual complexity (e.g., T001 database setup: optimistic 1, mostLikely 2–3, pessimistic 4; T002 schema migration: optimistic 1, mostLikely 3–6, pessimistic 8; T003 API handlers: optimistic 2, mostLikely 4–8, pessimistic 12) or, if you haven't estimated yet, mark these entries explicitly as placeholders/TBD and add a short comment/field like "estimateSource":"TBD" so reviewers know they're not final.

## planner/plans/waves.json
- [ ] [Critical] meta.artifact_hash is incorrect; replace the current value with the correct preimage-policy hash 7df36cbe62f64107b9c42117851cad2bfe28eb5b20460ead0f6784a8d133d3bf. To do this, set meta.artifact_hash to that exact string, and if you need to re-verify, recompute the canonical JSON hash using the preimage policy (set meta.artifact_hash = "" for the preimage, sort keys, ensure newline-terminated canonical JSON, then compute the hash) and confirm it matches the provided value before committing.

## z_tmp_excl_37772.go
- [ ] [Critical] delete file and add gitignore rule 'z_tmp_*'; ensure all Go packages live under planner/internal/*
- [ ] [Critical] this is a temporary stray file that should not be committed; delete z_tmp_excl_37772.go from the branch/PR, remove it from the index (git rm --cached or git rm) or delete and commit the deletion, and push the commit; add a gitignore entry to ignore files matching z_tmp_* (or the project-wide temp pattern) so future temp files are not tracked; ensure no references to this package exist and verify all legitimate planner packages remain under planner/internal/* per guidelines before finalizing the PR.
