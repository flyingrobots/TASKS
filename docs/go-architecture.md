A Go implementation that takes your v8 spec and nails it to the wall: DAG purity in the planner, resource realism in the runtime, validators as deterministic guardrails, and a clean JSON contract between the two. Where this plan aligns with the earlier TASKS/SLAPS artifacts and tool specs, I cite them inline.

---

## **The shape of the system (Go-first, production-ready)**

  

Think of the system as three cooperating binaries, plus a shared model package:

1. **Planner** (binary: tasksd): parses inputs, runs codebase analysis, compiles the tasks list, infers and minimizes the precedence DAG, simulates resource-aware waves, and emits the canonical artifacts (features.json, tasks.json, dag.json, waves.json, Plan.md, Decisions.md). This phase is “LLM + tools,” but the outputs are deterministic and validated by hard tools, not vibes. That division of responsibilities follows your v2 orchestration philosophy: let the “brain” decompose; let the tools prove and constrain.        
2. **Executor** (binary: slapsd): consumes the planner’s **coordinator.json** (the contract), maintains a rolling frontier, arbitrates resources (locks, quotas) outside the graph, applies circuit breakers and hot updates, executes acceptance checks via workers, and records provenance. This replaces brittle “waves” with a high-throughput loop, as argued in your FlowOps notes. 
3. **Tooling daemons** (binaries: acceptance-validator, evidence-validator, interface-resolver): deterministic, stateless services with JSON I/O; they’re invoked by the planner during compilation and by the executor at run time when checks must be executed or re-validated. The acceptance/evidence/interface contracts are already fully specified; we keep them out-of-process for reproducibility and to avoid coupling planning/runtime to a specific library build.   

The **shared model** (package model) holds canonical Go types for tasks, edges, resources, evidence, acceptance checks, graph metrics, waves, and the planner→executor contract. Both tasksd and slapsd import model, which eliminates drift.

This is the most boring possible architecture on purpose. Boring ships.

---

## **Repository layout that survives scale**

Use a single Go module with clear boundaries and zero import cycles:

```
/cmd
  /tasksd            # planner daemon/CLI (plan, validate, emit artifacts)
  /slapsd            # executor daemon/CLI (run, inspect, patch)
  /acceptance-validator
  /evidence-validator
  /interface-resolver
/internal
  /model             # canonical structs, JSON schema helpers, hashing
  /canonjson         # canonical JSON serializer (sorted keys, minimal numbers)
  /hash              # SHA-256 helpers and artifact manifest typing
  /analysis          # codebase scanners (wrappers over ast-grep, etc.)
  /planner
    /feature         # extraction from PLAN_DOC
    /taskgen         # granular tasks with boundaries, acceptance, evidence
    /deps            # edge inference (structural vs resource), confidence policy
    /dag             # cycle detection, transitive reduction, metrics
    /wavesim         # resource-aware wave simulation preview
    /validators      # thin RPC/CLI clients for external tools
  /executor
    /frontier        # priority queue; topo-ready filter; starvation aging
    /locks           # global-order lock manager; readers-writer semantics
    /quota           # semaphores; backpressure; profiles
    /dispatch        # worker selection by capability; retry policy
    /breaker         # error fingerprinting; hot-update trigger
    /patch           # add_task/add_edge/modify_resource apply engine
    /provenance      # append-only JSONL ledger
    /logging         # JSONL task log ingestion; structured event bus
  /schemas           # JSON Schema for features/tasks/dag/waves/coordinator
  /redaction         # secret scrubbing for evidence excerpts
  /evaluator         # optional plan scoring as per your rubric
/pkg
  /cli               # shared cobra-style command scaffolding
  /httpapi           # optional admin HTTP endpoints for slapsd/tasksd
```

Everything under /internal is explicitly non-public. Tools are batteries-included binaries so CI can run them without linking to your code.

---

## **The canonical model (Go types that match your spec)**

We define the types once and encode deterministically everywhere. A representative subset (edited for brevity):

```
// internal/model/types.go
package model

type Evidence struct {
  Type       string  `json:"type"`   // "plan" | "code" | "commit" | "doc"
  Source     string  `json:"source"`
  Excerpt    string  `json:"excerpt,omitempty"` // redacted before hashing
  Confidence float64 `json:"confidence"`        // 0..1
  Rationale  string  `json:"rationale"`
}

type AcceptanceCheck struct {
  Type    string         `json:"type"` // "command"|"artifact"|"metric"|...
  Cmd     string         `json:"cmd,omitempty"`
  Path    string         `json:"path,omitempty"`
  Expect  map[string]any `json:"expect,omitempty"`
  Timeout int            `json:"timeoutSeconds,omitempty"`
}

type DurationPERT struct {
  Optimistic  float64 `json:"optimistic"`
  MostLikely  float64 `json:"mostLikely"`
  Pessimistic float64 `json:"pessimistic"`
}

type InterfaceProduced struct {
  Name    string `json:"name"`
  Version string `json:"version"` // "v1" | "v1.2" | ...
  Type    string `json:"type"`    // "api"|"database_schema"|...
}

type InterfaceConsumed struct {
  Name               string `json:"name"`
  VersionRequirement string `json:"version_requirement"`
  Type               string `json:"type"`
  Required           bool   `json:"required"`
}

type ResourceNeed struct {
  Name   string `json:"name"`            // e.g., "cpu", "db.migrate"
  Units  int    `json:"units,omitempty"` // semaphores
  Access string `json:"access,omitempty"`// "read"|"write" for RW locks
}

type Task struct {
  ID           string              `json:"id"`         // "P1.T001"
  FeatureID    string              `json:"feature_id"` // "F001"
  Title        string              `json:"title"`
  Description  string              `json:"description,omitempty"`
  Category     string              `json:"category,omitempty"`
  Duration     DurationPERT        `json:"duration"`
  DurationUnit string              `json:"durationUnits"` // "hours"
  InterfacesProduced []InterfaceProduced `json:"interfaces_produced,omitempty"`
  InterfacesConsumed []InterfaceConsumed `json:"interfaces_consumed,omitempty"`
  AcceptanceChecks   []AcceptanceCheck   `json:"acceptance_checks"`
  Evidence           []Evidence          `json:"source_evidence"`
  Resources struct {
    Exclusive []string        `json:"exclusive,omitempty"` // names in catalog
    Limited   []ResourceNeed  `json:"limited,omitempty"`
  } `json:"resources"`
  ExecutionLogging struct {
    Format         string   `json:"format"` // "JSONL"
    RequiredFields []string `json:"required_fields"`
  } `json:"execution_logging"`
  Compensation struct {
    Idempotent   bool   `json:"idempotent"`
    RollbackCmd  string `json:"rollback_cmd,omitempty"`
  } `json:"compensation"`
}

type Edge struct {
  From       string      `json:"from"`
  To         string      `json:"to"`
  Type       string      `json:"type"` // "technical"|"sequential"|"infrastructure"|"knowledge"|"resource"
  Subtype    string      `json:"subtype,omitempty"` // "mutual_exclusion"|"resource_limited"
  IsHard     bool        `json:"isHard"`
  Confidence float64     `json:"confidence"`
  Evidence   []Evidence  `json:"evidence,omitempty"`
}

type TasksFile struct {
  Meta struct {
    Version         string  `json:"version"`
    MinConfidence   float64 `json:"min_confidence"`
    ArtifactHash    string  `json:"artifact_hash"`
    CodebaseAnalysis any    `json:"codebase_analysis"`
    Autonormalization struct {
      Split  []string `json:"split"`
      Merged []string `json:"merged"`
    } `json:"autonormalization"`
  } `json:"meta"`
  Tasks        []Task `json:"tasks"`
  Dependencies []Edge `json:"dependencies"`
  ResourceConflicts map[string]any `json:"resource_conflicts,omitempty"`
}

type DagFile struct {
  Meta struct {
    Version     string `json:"version"`
    ArtifactHash string `json:"artifact_hash"`
    TasksHash    string `json:"tasks_hash"`
  } `json:"meta"`
  Nodes   []struct {
    ID                  string `json:"id"`
    Depth               int    `json:"depth"`
    CriticalPath        bool   `json:"critical_path"`
    ParallelOpportunity int    `json:"parallel_opportunity"`
  } `json:"nodes"`
  Edges   []struct {
    From string `json:"from"`
    To   string `json:"to"`
    Type string `json:"type"`
    Transitive bool `json:"transitive"`
  } `json:"edges"`
  Metrics struct {
    MinConfidenceApplied float64         `json:"min_confidence_applied"`
    KeptByType           map[string]int  `json:"kept_by_type"`
    DroppedByType        map[string]int  `json:"dropped_by_type"`
    Nodes                int             `json:"nodes"`
    Edges                int             `json:"edges"`
    EdgeDensity          float64         `json:"edge_density"`
    WidthApprox          int             `json:"width_approx"`
    LongestPathLength    int             `json:"longest_path_length"`
    CriticalPath         []string        `json:"critical_path"`
    IsolatedTasks        int             `json:"isolated_tasks"`
    VerbFirstPct         float64         `json:"verb_first_pct"`
    EvidenceCoverage     float64         `json:"evidence_coverage"`
  } `json:"metrics"`
  Analysis struct {
    OK        bool     `json:"ok"`
    Errors    []string `json:"errors"`
    Warnings  []string `json:"warnings"`
    SoftDeps  []Edge   `json:"soft_deps"`
  } `json:"analysis"`
}

type Coordinator struct {
  Version string `json:"version"`
  Graph   struct {
    Nodes []Task `json:"nodes"` // runtime-ready
    Edges []Edge `json:"edges"` // hard structural deps only
  } `json:"graph"`
  Config struct {
    Resources struct {
      Catalog  map[string]struct {
        Capacity  int    `json:"capacity"`
        Mode      string `json:"mode"`       // "exclusive"|"semaphore"
        LockOrder int    `json:"lock_order"` // global ordering
      } `json:"catalog"`
      Profiles map[string]map[string]int `json:"profiles"`
    } `json:"resources"`
    Policies struct {
      ConcurrencyMax         int               `json:"concurrency_max"`
      LockOrdering           []string          `json:"lock_ordering"`
      CircuitBreakerThresholds map[string]any  `json:"circuit_breaker_thresholds"`
    } `json:"policies"`
  } `json:"config"`
  Metrics struct {
    Estimates struct {
      P50TotalHours     float64 `json:"p50_total_hours"`
      LongestPathLength int     `json:"longest_path_length"`
      WidthApprox       int     `json:"width_approx"`
    } `json:"estimates"`
  } `json:"metrics"`
}
```

These types encode exactly the contracts you’ve been iterating on: **structural edges feed the DAG; resource edges are recorded for traceability but excluded from topological sort**. That “don’t pollute the graph with resource constraints” rule is non-negotiable; the earlier specs concur, while also making infrastructure resource edges first-class evidence-backed records for schedulers to respect.   

---

## **Deterministic artifacts and hashing in Go**

  

We ship a tiny package (`/internal/canonjson`) that emits canonical JSON:

- Maps are materialized as ordered key arrays and serialized lexicographically.
- Numbers are rendered in minimal decimal form.
- LF endings, two-space indentation, newline-terminated.
- Hashes are `sha256(lowercase hex)` over the canonical bytes and written into `meta.artifact_hash` for every JSON artifact; the markdown files carry a “Hashes” section for cross-checks.

That is directly aligned with the deterministic-output guarantees you called out; we integrate it at the model boundary so every binary uses the same serialization. 

---

## **The planner (tasksd) in Go, end to end**

The planner does six things in order, as prose—not magic.

> Implementation note: `cmd/tasksd` delegates orchestration to `internal/app/plan.Service`. The default wiring lives in `plan.NewDefaultService()`, which assembles the Markdown doc loader, census analyzer, dependency resolver, wave simulator, coordinator builder, artifact writer, and validator runner behind small interfaces. The CLI only overrides behavior (e.g., logging analysis failures) while the service coordinates task loading, DAG build, validation, wave preview, and artifact emission. This keeps the domain workflow pure and lets tests exercise the planner with in-memory adapters.

**First**, it ingests the PLAN_DOC and runs a codebase census. In practice, we shell out to ast-grep for polyglot patterns (Prisma, Drizzle, Knex, Sequelize, Rails), deployment YAMLs, test fixtures, and environment access hotspots; we then normalize findings into tasks.json.meta.codebase_analysis.shared_resources with names, locations, constraints, and rationales. This is straight from your enhanced codebase-first method. 

**Second**, it extracts features (5–25), each with priority and evidence, from the spec; the feature extractor is deliberately simple: section headings + requirement phrases → features, with evidence tracked for later validation and confidence. That matches the v1/v2 approach to grounded features.   

**Third**, it decomposes tasks to 2–8 hour units, auto-splitting >16h and auto-merging <0.5h tasks, and it **requires** every task to declare (a) expected complexity, (b) definition of done, (c) explicit scope boundaries, (d) at least one machine-verifiable acceptance check, and (e) execution logging instructions that conform to the JSONL log contract (timestamp, task_id, step, status, message; error_code/stack on errors). The acceptance checks are validated immediately by the acceptance validator, which returns an executable test suite. That validator is fully specified; we just wrap it.   

**Fourth**, it infers dependencies. Structural edges—technical, sequential, infrastructure, and knowledge—enter the precedence DAG. Resource edges are recorded with evidence and confidence (often 1.0) but explicitly excluded from the DAG. Evidence for every task and edge is schema-validated: type, source anchor, excerpt (redacted), confidence, and rationale. Evidence quotes are validated against sources using exact, fuzzy, and semantic strategies; the validator’s output is persisted for audit. 

**Fifth**, it builds and optimizes the DAG: filter on min_confidence and isHard, detect and report cycles (≤2 automated resolution passes by inserting/splitting interface tasks), and perform transitive reduction (Gries-Herman). Compute edge density, width, critical path, and related metrics. This is exactly the v1/v2 pipeline, formalized.   

**Sixth**, it simulates waves purely for planning-time preview using Kahn layering plus resource feasibility checks per wave, recording any capacity-induced separations **without** injecting them back into the precedence graph. The executor is free to exploit real-time availability later. This keeps the plan minimal and the runtime honest. 

All along, the planner talks to the **interface resolver** to ensure that produced and consumed interfaces are version-compatible and acyclic at the interface level, with automatic upgrades when allowed by policy. Its spec is done; we implement a thin client and fold results into the DAG and tasks.

Finally, tasksd emits artifacts with canonical JSON and embedded SHA-256 hashes, and a derived planId that seeds waves.json and coordinator.json. Any hard fails (cycles, missing evidence, missing machine checks) abort emission of coordinator.json, by design. 

---

## **The executor (slapsd) in Go, exactly how it runs**

At startup, slapsd loads **coordinator.json**, registers resource capacities per profile (local/ci/prod), and builds a pure structural graph (edges: hard structural only). It spawns the main loop:

1. **Frontier formation**: all tasks whose hard predecessors are complete move to the ready set.
2. **Resource arbitration**: a task is runnable if it can acquire all required locks and quotas in a single shot. Locks obey a canonical global order; writes block other writers and readers sharing the same resource; reads may share. Quotas are token semaphores (CPU, test DB slots, runners).
3. **Dispatch**: tasks are scored by a priority heuristic (shallower depth, higher downstream parallel opportunity, higher confidence, lower rollback cost, plus aging) and sent to capability-matched workers.
4. **Event-driven progress**: slapsd consumes task JSONL logs to update progress, detect errors, and compute error fingerprints.
5. **Circuit breakers and hot updates**: repeated fingerprints within a window trip breakers that synthesize patches—injecting remedial tasks (e.g., batch install missing UI deps via import census), adding soft ordering hints (add_edge), or adjusting capacities. The patcher applies these changes to the running plan safely: additions only, never removing existing hard edges. This is FlowOps’ rolling frontier and breaker playbook—just formalized in Go. 

**Workers** are just goroutines pulling assignments; the executor spawns external processes to do the work so we can capture exact JSONL output and keep the runtime crash-only. Acceptance checks are run at task completion via the acceptance-validator’s generated test suite, so “done” is never subjective. 

**Provenance** is an append-only JSONL ledger: one line per task execution with timestamps, exit codes, and artifact deltas. If a non-idempotent task fails, the executor runs the declared compensation and rewinds safely along a reverse-topological “rollback frontier.”

All of this is deliberately boring Go: channels, container/heap for the priority queue, sync for concurrency control, context.Context for cancellation, log/slog with a JSON handler to mirror the JSONL format.

---

## **The interprocess contracts (and why they’re JSON, not protobuf)**

We standardize on content-addressed **canonical JSON** for artifacts and for tool RPCs. For long-lived admin endpoints (introspection, metrics), slapsd exposes a small HTTP API:

- GET /v1/frontier → JSON snapshot of ready/running/completed tasks.
- POST /v1/patch → apply hot updates (add_task, add_edge, modify_resource).
- GET /v1/resources → current capacities and usage.
- GET /v1/provenance?task_id=... → last N ledger entries.

We prefer JSON here because (a) artifacts are already JSON; (b) the tooling daemons use JSON; (c) determinism and auditability matter more than a couple of microseconds. If you later want streaming and stricter types, add gRPC on top—nothing here prevents it.

---

## **Validators: how Go calls them and how failures gate the run**

Each validator has two adapters: a subprocess client (exec the tool with JSON stdin/stdout) and an optional HTTP client. The planner never trusts an unvalidated plan:

- **Acceptance**: every task must have at least one machine-checkable criterion, and the validator returns a runnable test suite; the executor will later run that suite as the last step of the task. 
- **Evidence**: every task and edge must carry at least one evidence object; the validator verifies that the excerpt exists (exact/fuzzy/semantic) and returns per-claim details; secrets are redacted **before** hashing and emission. 
- **Interface**: produced/consumed compatibility is checked, missing producers are flagged, and version conflicts are auto-resolved when allowed by policy.

Validator subprocesses follow a predictable contract:

- Inputs are canonical JSON payloads (`tasks`, `dag`, `coordinator`) hashed with SHA-256 prior to invocation.
- Each payload hash is used as a cache key so repeated `tasksd plan` runs with unchanged inputs reuse validator output deterministically.
- Cached validator reports (status, command, raw JSON) are persisted alongside artifacts and surfaced in `tasks.json.meta.validator_reports` and Plan.md for auditing.

If any of those fails under “strict,” tasksd stops and emits an actionable error report instead of a half-baked plan. That’s consistent with your earlier v1/v2 “no cycles, no speculation” posture.   

---

## **Resource modeling with teeth (and no graph pollution)**

Resources live **next to** the graph, not inside it. The planner records exclusive and limited resources discovered during analysis (e.g., migrations, deployment pipelines, test DB pools). Those show up as resource edges for traceability only; the executor enforces them by locks and semaphores. This aligns with the “constraints outside the DAG” rule in your updated specs and with the rolling frontier design’s emphasis on throughput rather than artificial serialization.   

The lock manager enforces a **global order** to eliminate deadlocks; the quota manager applies **backpressure** to avoid thrash. Readers–writer semantics allow multiple readers; a single writer excludes all. Profiles (local/ci/prod) adjust capacities without mutating the artifacts.

---

## **Logging, telemetry, and provenance (no guesswork)**

Every worker writes **JSON Lines** with timestamp, task_id, step, status (start|progress|done|error), and message; on error, data.error_code and data.stack are required. The executor ingests those lines in real time, updates the frontier, and fingerprints errors for the breakers. The format is stable, simple, and already documented in your drafts; we enforce it in the Go dispatcher with a validator shim so we never silently drop malformed lines. 

Provenance is separate from logs: an append-only file provenance/yyyymmddTHHMMSS.jsonl with {task_id, started, finished, exit_code, worker_id, artifacts_changed, logs_path}. That’s your “repo-native history that tells the story” principle—formalized. 

---

## **Algorithmic commitments (and the exact places they live)**

- Cycle detection: DFS + recursion stack in /internal/planner/dag.
- Topological sort and antichains: Kahn layering in /internal/planner/dag.
- Transitive reduction: Gries–Herman for sparse DAGs, with unit tests on edge preservation.
- Wave simulation: /internal/planner/wavesim applies capacity checks over each antichain to produce waves.json, without ever feeding those constraints back into the precedence DAG.
- Priority queue: container/heap implementation in /internal/executor/frontier with a strict comparator (depth asc, unblock potential desc, rollback cost asc, confidence desc, wait time asc).
    These match the v1/v2 algorithms and the FlowOps runtime model.     

---

## **Concrete CLI, so engineers don’t guess**

There are only a few commands worth remembering.

1. Planning
    tasksd plan --doc spec.md --repo . --out ./plans/projectA
    Writes features.json, tasks.json, dag.json, waves.json, Plan.md, Decisions.md with canonical hashes. Fails hard on cycles, missing evidence, or non-executable acceptance.
1. Validation-only (CI preflight)
    tasksd validate --in ./plans/projectA
    Recomputes hashes, re-runs validators, and refuses to proceed if artifacts drift.
2. Execution
    slapsd run --coordinator ./plans/projectA/coordinator.json --profile ci
    Starts the rolling frontier; exposes /v1/* admin endpoints; writes provenance.
3. Hot patch
    slapsd patch --file ./patches/2025-10-05.json
    Adds tasks/edges/capacity tweaks safely, as allowed by the patch schema.

The defaults align with the earlier “stateless JSON tools” philosophy: simple CLIs that can be composed in any CI pipeline without a bespoke server. 

---

## **Security and redaction (applied early, applied always)**

Evidence excerpts are redacted before hashing and emission; the redactor runs the same heuristics you already enumerated (API keys, JWTs, PEM blocks, long base64/hex). The planner fails if the evidence validator finds an unredacted secret, and the executor refuses to run with unvalidated artifacts. That’s exactly how you keep “audit trail” without “data spill.” 

---

## **Where this explicitly aligns with earlier artifacts (so no one re‑debates it)**

- **LLM + tools orchestration** and quality gates: we inherit the v2 agent/pipeline pattern, but in Go, with validators as subprocesses or HTTP services. 
- **Grounded features → bounded tasks → validated DAG → Kahn waves**: that v1 spine remains intact; we add canonical hashing and resource-aware wave _simulation_. 
- **Shared resources recorded, not embedded**: preserve DAG purity, enforce resources at runtime. Your v2.1 resource modeling carries forward cleanly. 
- **Executor = rolling frontier with circuit breakers and playbooks**: FlowOps is codified here as SLAPS’ event loop, not a slide. 
- **Validators are normative**: acceptance (machine-checkable only), evidence (quote verification), interface (compatibility and versioning). We don’t water that down.   

---

## **Implementation checklist (engineer-facing, finite)**

1. **Model & canon**: implement internal/model and internal/canonjson; write golden tests that re-serialize JSON fixtures and verify the same SHA-256 across platforms. 
2. **Validators**: wrap acceptance/evidence/interface tools with subprocess adapters; add retry and structured error mapping.   
3. **Planner**: codebase analysis wrappers over ast-grep; feature extraction; task generation with boundary enforcement; dependency inference; DAG build + reduction; wave simulation; emit artifacts with hashes; stop on any hard fail.   
4. **Executor**: frontier heap; lock/quota managers (global order + RW semantics; semaphore quotas); dispatcher; JSONL ingestion; circuit breakers and hot-patch engine; provenance ledger; admin HTTP. 
5. **Security**: evidence redaction pre-hash; refuse to run on unvalidated artifacts. 
6. **E2E tests**: seed repo with a sample plan; run tasksd plan → slapsd run --profile local; verify acceptance suites pass; prove that circuit breaker playbook injects a batch-install task under synthetic “module not found” failures. 

---

## **One last opinionated call**

Do not collapse planner and executor into a single service. The contract boundary (coordinator.json) is your guarantee of reproducibility, auditability, and safe parallel development. You already lived this separation in the earlier designs; implementing it in Go this way means your team can refactor the planner’s fancy decomposition logic without risking the runtime’s safety, and vice versa. Keep the DAG pure, keep resources out of it, and let SLAPS police the lanes. That’s how you get seven workers chewing a repo like a freight train and not a funeral procession.   
