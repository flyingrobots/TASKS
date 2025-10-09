# **T.A.S.K.S. + S.L.A.P.S. v8.0**

**Document status:** normative specification

**Compatibility:** supersedes prior drafts; aligns with the v2 planning pipeline and the FlowOps/rolling‑frontier execution guidance; integrates acceptance, evidence, and interface validation tool specifications

**Contract of the system:** the dependency graph describes precedence only; resource constraints are modeled, recorded, and enforced by the runtime outside the graph.
  
## **1. Scope and governing principle**

This specification defines two cooperating subsystems. **T.A.S.K.S.** is the planning engine that turns a specification and a codebase into an optimized, audit‑ready plan: a clean, acyclic precedence DAG plus deterministic artifacts. **S.L.A.P.S.** is the execution runtime that consumes that plan and executes it with a rolling frontier scheduler, resource arbitration, self‑healing behaviors, and strict provenance.

The governing principle is **DAG purity**. Precedence edges mean “must happen before,” nothing more. Resource constraints—exclusivity, capacities, time windows—are cataloged and enforced **adjacent to** the graph by the runtime. Keeping the graph minimal and structural maximizes parallelism and keeps the plan mathematically valid; letting the executor arbitrate resources keeps reality in the loop without polluting the topology. This separation was implicit in v2; here it is explicit and non‑negotiable.   

## **2. Responsibilities and division of labor**

T.A.S.K.S. is responsible for understanding, decomposition, proof, and minimality. It reads the plan document and the codebase; inventories the interfaces and extension points already present; extracts a small set of top‑level capabilities; decomposes those into bounded tasks with crisp boundaries; infers structural dependencies and attaches evidence; and compiles a minimal DAG through cycle detection and transitive reduction. It produces canonical JSON artifacts with SHA‑256 hashes, plus human‑readable rationales that a reviewer can audit in one pass. This aligns with the v2 “LLM as orchestrator, tools as validators” architecture and its deterministic artifact set. 

  

S.L.A.P.S. is responsible for throughput, safety, and adaptation. It maintains a rolling frontier of ready tasks; acquires locks and quotas in a global order to prevent deadlocks; dispatches to capability‑matched workers; monitors telemetry; trips circuit breakers on repeated failure patterns; injects hot patches into the running plan when necessary; and records an append‑only provenance ledger for audit. It is the direct successor to the FlowOps guidance that replaced rigid wave barriers with an event‑driven scheduler and practical circuit‑breaker playbooks. 

  

The arbitration between creative decomposition and deterministic validation follows the v2 orchestration model: the planner may be opinionated, but tools make the gates hard. Where this document refers to validators and orchestration patterns, it is referencing the v2 agent and tool design directly. 

  

## **3. Inputs, outputs, and determinism**

  

The planning engine accepts (a) a primary plan document in text or markdown and (b) access to the codebase. The outputs are five canonical artifacts—features.json, tasks.json, dag.json, waves.json, and Plan.md—and, for execution, a coordinator.json contract passed to S.L.A.P.S. Each JSON artifact is serialized canonically and newline‑terminated, then hashed with SHA‑256 and records that hash in meta.artifact_hash. Plan.md includes a “Hashes” section listing the exact SHA‑256 of every artifact.

  

Serialization is deterministic. Object keys are sorted lexicographically at all depths; arrays remain in computed order; numbers are rendered in minimal decimal form; encoding is UTF‑8; line endings are LF; indentation is two spaces; trailing whitespace is not permitted. Hashes are computed over the canonical bytes plus a final newline and emitted as lowercase hex. This practice originated in v1/v2 and is promoted here to a strict requirement. 

  

## **4. Codebase‑first planning**

  

The planner begins with a census of the existing system. It uses static queries (for example, ast-grep) across the languages actually present to locate interfaces and extension points, migration frameworks and schema touchpoints, deployment and CI plumbing, test scaffolds, and environment configuration. In modern JS/TS stacks that often means Prisma, Drizzle, Knex, or Sequelize, with Rails kept for polyglot repositories. The findings populate tasks.json.meta.codebase_analysis, including a catalog of shared resources whose exclusivity or capacity will matter at execution time. This approach preserves the “reuse first” ethic of v2—extend interfaces you have before inventing new ones—and is the practical groundwork for realistic plans. 

  

## **5. Feature extraction and task definition**

  

The planner extracts five to twenty‑five capabilities from the source document, each with a short, outcome‑oriented description, a priority, and at least one evidence citation to the plan text or code. It then decomposes those capabilities into tasks that fit within a target of two to eight hours of work, with a hard maximum of sixteen hours. Any task that exceeds that cap is split along artifact or interface boundaries; any task that would take less than half an hour is merged with a sibling to keep the surface area manageable. The record of these normalizations appears as a structured note in tasks.json.meta.autonormalization. This sizing discipline and its audit trail were explicit in early T.A.S.K.S. drafts and remain required. 

  

Every task has boundaries that leave nothing to interpretation: a quantifiable complexity declaration; a “definition of done” that stops scope creep cold; scope constraints that name the files and subsystems to touch or to avoid; execution guidance that tells a worker how to log progress; explicit resource requirements; a declaration of idempotency and a compensation action if it is not idempotent; and at least one **machine‑verifiable** acceptance check. The acceptance checks are not prose—they are structured, executable criteria of standard types and are validated by a dedicated tool. The acceptance tool is normative and its implementation is already specified; this planner/runtime contract depends on it. 

  

## **6. Evidence and redaction**

  

Every task and every dependency carries one or more evidence objects. An evidence object states a type (“plan,” “code,” “commit,” or “doc”), a source location, an excerpt with secrets redacted, a confidence score between zero and one, and a short rationale that states why the evidence supports the claim. Evidence is validated by a deterministic tool that verifies the quotes exist in the sources with exact, fuzzy, and semantic matching, respects section context, and reports per‑claim validation status. Its configuration and behaviors are normative and already defined. Before hashing or publication, sensitive substrings such as API keys, JWTs, and PEM blocks must be replaced with explicit [REDACTED_*] markers. 

  

## **7. Dependency modeling and the minimal graph**

  

Dependencies are inferred only where the plan or code supports them. Structural edges—the only edges that enter the precedence DAG—include technical dependencies where one task produces an interface that another consumes, sequential dependencies where information flow dictates order, infrastructure prerequisites such as environments or tools that must exist, and knowledge edges where a research spike or context discovery materially de‑risks an implementation step. Resource edges are **not** structural; they are recorded for traceability but are excluded from topological sorting. Human‑resource edges are prohibited. This resolves an early contradiction in the project history by allowing infrastructure resource edges for audit while keeping them out of the graph, which was the practical guidance that emerged in v2.1. 

  

Edges carry isHard and confidence. When the DAG is compiled, only edges with isHard == true and a confidence at or above the plan’s minimum (default 0.7) are included. Cycles are detected with a DFS/stack pass. The planner may attempt up to two surgical, semantics‑preserving interventions—splitting a composite task at an interface boundary or inserting a contract task—to break a cycle; if a cycle remains, the planner must stop and emit a diagnostic with explicit suggestions rather than continuing into wave generation. Following validation, the DAG is reduced transitively using a sparse‑graph algorithm to remove redundant constraints and keep the graph minimal. These algorithms and their quality metrics were formalized in v1/v2 and are unchanged here except for the stricter hard‑fail policy. 

  

The compiled dag.json reports the applied confidence threshold; a count by dependency type of edges kept and dropped; the counts of nodes and edges; edge density; the approximated maximum width (a proxy for parallel opportunity) via Kahn layering; the length of the critical path; the percentage of verb‑first task names; and a list of soft and low‑confidence dependencies for human review. The recommended edge density band is **0.01–0.50**. Values below this range often imply missing constraints; values above it typically imply over‑constraint or failure to reduce transitively.

  

## **8. Waves as simulation, not shackles**

  

Wave planning is a planning‑time preview of parallelism, not a runtime barrier. The planner layers the DAG into antichains using Kahn’s algorithm, then simulates resource feasibility inside each candidate wave to forecast contention. Where exclusive resources or finite quotas would force separations, the planner records the decision and the reasoning in waves.json. For each wave, it estimates P50, P80, and P95 completion times using standard PERT approximations from per‑task optimistic, likely, and pessimistic durations, under the conservative assumption that waves are paced by their longest member. This is exactly the v2 wave generator’s approach; here it is explicitly framed as a forecast rather than a harness. 

  

## **9. Standard telemetry**

  

Every worker process—human or machine—emits JSON Lines logs to standard output or an agreed log file. Each line contains a timestamp, a task identifier, a step name, a status enumerated as “start,” “progress,” “done,” or “error,” and a short message; optional fields include a progress percentage and a structured data object. On error, the data object must include an error code and, where available, a diagnostic stack. This logging contract is mandatory so that the runtime can reason about progress and circuit‑breaker conditions, and so that the provenance ledger can link to primary telemetry without scraping human prose. The v2.1 planning notes emphasized explicit logging instructions inside each task; this elevates the format to a shared contract that both planner and executor depend upon. 

  

## **10. The execution runtime**

  

S.L.A.P.S. executes the plan with a **rolling frontier** scheduler. As soon as a task’s hard predecessors complete, the task becomes eligible. The runtime checks that the task’s required resources can be acquired under current profiles and policies; if so, it dispatches the task to a worker with the required capabilities. If not, it leaves the task in the frontier and pulls another eligible task that can run within resource constraints. The priority heuristic prefers tasks at shallower depth (to maximize downstream unblocking), tasks whose completion unlocks wider parallelism, tasks with higher evidence confidence, and tasks with lower rollback cost; an aging term prevents starvation. This model is the direct continuation of the FlowOps “replace waves with frontier” guidance; the value is simple: no idle agents because a calendar barrier hasn’t tripped. 

  

Resources are enforced by two managers. The lock manager governs exclusive resources such as schema migrations or repository‑wide file operations. It uses a canonical global lock ordering (numeric or lexicographic) and an all‑or‑nothing acquire to prevent partial holds. The quota manager governs finite pools such as CPU cores, CI runners, or test databases and applies back‑pressure when a dispatch would exceed capacity. Read/write semantics are available to permit concurrent reads to shared resources while serializing writes. Profiles such as local, ci, and prod alter capacities and windows without altering the plan artifacts.

  

The runtime watches telemetry for failure fingerprints. When patterns cross configured thresholds, circuit breakers trip and a hot‑update patch is generated. Patches add tasks, add edges, or adjust resource policies; they never remove hard structural edges. Common playbooks include scanning for missing front‑end dependencies and batch‑installing them based on import census, quarantining flapping tasks and gating their dependents, or reducing concurrency temporarily when resource exhaustion is detected. These “pattern playbooks,” and the practice of late‑bound soft edges inferred from file‑touch sets, are direct lifts from the FlowOps field notes and are now first‑class runtime behaviors. 

  

A task that is non‑idempotent must define a compensation action. On failure, the runtime executes compensations in reverse topological order for any already‑applied dependents, marks the failed task incomplete, and returns it to the frontier for retry in accordance with policy. Whether or not a task is idempotent is declared in tasks.json and affects scheduling decisions through the “rollback cost” priority factor.

  

Every dispatch and completion is recorded in a provenance ledger as an append‑only JSONL file. Each record contains the task id, start and finish timestamps, exit code, worker id, a list of artifacts changed with before/after hashes, a pointer to the task’s JSONL telemetry, and a checkpoint identifier. This practice was recommended alongside “commit per task, tag per frontier depth” and remains useful for forensic reconstruction. 

  

## **11. Validators and quality gates**

  

Three validators are normative.

  

The **acceptance validator** confirms that each acceptance criterion is machine‑checkable and translates into executable commands with crisp pass/fail semantics. It rejects subjective criteria; it can synthesize a minimal test suite from validated criteria and report functional and risk coverage. The tool’s API, input/output schemas, and algorithm are already specified and binding. 

  

The **evidence validator** verifies that every evidence excerpt appears in the cited source, with exact, fuzzy, and semantic matching under defined thresholds and with section‑aware context checks. It emits a complete validation record per claim and a summary with validation rate and confidence distribution. Its configuration and error modes are already defined and binding. 

  

The **interface resolver** validates that produced and consumed interfaces are compatible under semantic versioning rules, detects missing producers and circularities at the interface layer, proposes or applies safe auto‑resolutions, and builds a registry showing producers, consumers, and dependency chains. Its contract is likewise normative. 

  

The orchestration of these tools—when they run, how their diagnostics gate the pipeline, and how their results are reflected back into artifacts—follows the v2 agent and pipeline definition. The planner does the semantics; the tools make the math and the gates. 

  

## **12. Data contracts**

  

The JSON artifacts are contracts; their shape and meaning are not negotiable.

  

features.json records the feature set extracted from the plan, with generated provenance and a JSON‑schema‑compliant evidence array for each feature. tasks.json holds the planning universe: a meta block with minimum confidence and the codebase analysis; an audit of automatic normalizations; and a tasks array where each task has identifiers, boundaries, execution guidance including the standard telemetry requirements, resource requirements, duration estimates, interfaces produced and consumed, acceptance checks, evidence, reuse notes, idempotency and compensation, and capability tags for worker matching. It also includes a dependencies array where edge objects carry from, to, type, isHard, confidence, and evidence. Resource edges are recorded for traceability but are intentionally excluded when building the precedence DAG.

  

dag.json reports whether the graph is valid; lists any errors and warnings; enumerates metrics including density, width, and the critical path; lists soft and low‑confidence edges; and—if the graph is invalid—spells out cycle break suggestions. waves.json is keyed by a plan id derived from the hash of tasks.json, declares the default barrier configuration, and records the simulated waves with per‑wave estimates and resource usage summaries. These shapes and their intent were introduced in v1/v2; their semantics here are the same but the requirements for evidence, logging, and resource separation are stricter. 

  

coordinator.json is the execution contract. It embeds the graph’s **hard structural** nodes and edges, the resource catalog with capacities, access modes, lock ordering, and environment profiles, plus runtime policies such as maximum concurrency, canonical lock order, and circuit‑breaker thresholds. It includes plan‑level estimates such as P50 total hours, approximated width, and critical path length. S.L.A.P.S. consumes this contract verbatim at startup.

  

## **13. Quality criteria**

  

A plan that passes must be acyclic, minimal, and evidenced. The graph has zero cycles; the density is in the recommended band or the reason to diverge is recorded; at least ninety‑five percent of tasks and edges carry validated evidence objects; at least one machine‑verifiable acceptance check exists for every task; at least eighty percent of task titles lead with a verb; isolated nodes are eliminated or justified; soft and low‑confidence dependencies are explicitly separated for human review. These are the same axes the v1/v2 evaluator scored; the thresholds here are firm. 

  

## **14. Security and redaction**

  

Before emission or hashing, evidence quotes must be scrubbed for secrets. Replace matched substrings—API keys, long hex or base64 tokens, JWTs, PEM blocks—with explicit [REDACTED_*] markers, keeping surrounding context intact so evidence still proves the point without leaking credentials. The rules for redaction and the validator’s tolerance for redacted spans are inherited from the evidence tool’s specification. 

  

## **15. Human‑in‑the‑loop**

  

Automation does the counting; people do the deciding. There are four mandatory checkpoints. First, a document quality gate confirms that the source specification is complete and internally consistent. Second, a sanity pass over features and tasks checks that decomposition matches domain reality. Third, if cycles remain after two automated passes, an architect resolves them explicitly by splitting tasks or inserting interface contracts. Fourth, a final approval signs off the plan before execution. These checkpoints mirror the “automate the tedious, validate with humans, decide with humans” guidance from v2 and are retained because they prevent expensive tail risks. 

  

## **16. Algorithms in plain language**

  

The planner’s graph pass works as follows. It filters edges to those that are both hard and above the confidence threshold; it constructs adjacency lists; it runs a DFS with a recursion stack to detect cycles and halts if they remain after two surgical corrections; it runs transitive reduction to delete any edge implied by a longer path; it then computes width and the critical path using Kahn’s layering, which also yields an initial wave partition. None of these choices are exotic; they are the same choices that made v1/v2 predictable and audit‑friendly. 

  

The runtime’s scheduling loop is equally plain. It maintains a set of tasks whose predecessors are all complete. For each such task, in priority order, it attempts to acquire required resources using a global lock order and all‑or‑nothing semantics; if successful, it dispatches to a worker that advertises the necessary capabilities; otherwise, it leaves the task in the frontier and continues. It blocks for the next event—completion, failure, or timeout—handles it, and repeats. If error fingerprints exceed configured thresholds, it applies a playbook that may inject a new task (for example, batch‑install dependencies discovered by import census), add a soft ordering edge, or modify resource capacities for a profile. That is FlowOps in practice, now formalized. 

  

## **17. Testing and evaluation**

  

A complete plan includes its own tests. The acceptance validator can assemble an executable test suite from task‑level criteria; the evidence validator can produce an evidence validation report; the interface resolver can generate an interface registry and conflict report. These tool outputs should be integrated into CI, with strict gates that mirror the categories used by the v1/v2 evaluator: structural validity, coverage and granularity, evidence and confidence, DAG quality, wave construction, gates, and naming discipline. The acceptance tool’s test generation and coverage analysis are especially important; they catch “unrunnable” criteria before they reach runtime.     

  

## **18. Implementation checklist**

  

A team implementing this system can proceed in this order.

1. Establish canonical JSON and hashing utilities in the planner repository and add a minimal tool that prints the canonical bytes and hash for a given object.
    
2. Implement the codebase census for your languages and frameworks, and persist the results in tasks.json.meta.codebase_analysis.
    
3. Implement feature extraction, task decomposition, boundary enforcement, and auto‑normalization; wire in the standard logging contract and resource annotations at the task level.
    
4. Integrate the acceptance validator and reject any task lacking at least one machine‑verifiable check; integrate the evidence validator and reject any task or edge whose evidence fails validation; integrate the interface resolver and halt on conflicts unless an explicit override is recorded after human review.
    
5. Implement dependency inference and the DAG compiler with cycle detection, transitive reduction, and metrics; emit dag.json with soft and low‑confidence edges separated and a clear ok flag.
    
6. Implement wave simulation as Kahn layering with resource feasibility checks; emit waves.json with per‑wave estimates and resource usage summaries.
    
7. Emit coordinator.json with resource catalog, profiles, and runtime policies; run the rolling frontier scheduler against this contract with lock and quota managers; write the provenance ledger and verify task telemetry.
    
8. Implement circuit‑breaker pattern matching and hot‑update patching; start with the import‑census/batch‑install playbook that saved real teams from “missing UI deps” pinball.
    
9. Wire all validators and the scheduler into CI; reject on cycles, missing evidence, missing acceptance checks, invalid interfaces, or security leaks, and insist on deterministic hashes across a short battery of reproducibility tests.
    

  

This ordering follows the v2 pipeline logic, FlowOps’ execution pragmatics, and the tool specifications’ zero‑questions APIs.         

  

## **19. Repository layout and scaffolding**

  

A minimal implementation places planning artifacts under plans/<project>/, validator specs and stubs under tools/, a simple bin/ with a CLI that runs the planner and validators in order, and a provenance/ directory that the runtime writes. The orchestration prompt and tool wiring can be adapted from the v2 agent definition, with the same “LLM for decomposition; tools for validation” stance. 

---

### **Closing note**

  

Keep the DAG clean, keep the evidence honest, and keep resources where they belong: **next to** the graph, not inside it. Plans that follow these rules run fast, heal themselves, and leave a paper trail that leadership and auditors actually trust. That is the point—and that is what this Editor’s Edition locks in.   

  

If you want me to emit an initial project skeleton—canonical JSON/hash helpers, logging shims, validator harness, and a runnable rolling‑frontier loop wired to coordinator.json—I can generate those files exactly to this spec now, with sample tasks and a working provenance ledger.