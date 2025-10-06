Alright James—no mysticism, just the formal guts. Here’s how SLAPS consumes coordinator.json from TASKS and what the executor actually does, end-to-end.

1) coordinator.json — canonical schema (planner → runtime contract)

{
  "version": "1.x",
  "graph": {
    "nodes": [
      {
        "id": "P1.T011",
        "name": "Deduplicate services",
        "kind": "task",
        "estimates": { "p50_h": 4, "p95_h": 8 },
        "priority": 0.62,                 // 0..1 (higher = sooner if ties)
        "acceptance": [{ "type": "cmd", "cmd": "pnpm test -- --grep=svc" }],
        "effects": [{ "type": "fs", "writes": ["apps/web/**", "packages/shared/**"]}],
        "resources": {                    // consumable/exclusive resources
          "locks": ["db:migrate"],        // exclusive lock names
          "quota": { "cpu": 1, "mem": "512Mi" }
        },
        "retry": { "max": 2, "backoff_s": [10, 30] },
        "checkpoint": true,               // snapshot after success
        "idempotent": true,               // can re-run safely
        "compensation": { "cmd": "git restore --source=@~1 -- ." }, // if not idempotent
        "tags": ["ui", "shared"]
      }
    ],
    "edges": [
      { "from": "P1.T009", "to": "P1.T011", "type": "hard" },
      { "from": "P1.T010", "to": "P1.T011", "type": "hard" },
      { "from": "P1.T020", "to": "P1.T026", "type": "soft" }       // prefer-after, not required
    ]
  },
  "policies": {
    "concurrency_max": 7,
    "fairness": "weighted",                 // weighted round-robin by priority bucket
    "circuit_breakers": {
      "repeat_error_threshold": 2,          // pattern trigger -> replan
      "thrash_window_s": 120,               // kill flappers
      "deadline_overrun_factor": 3.0        // P95 * factor -> abort
    },
    "lock_policy": "try-lock-wait-die",     // deadlock avoidance
    "update_mode": "hot",                   // allow dynamic edge/node inserts
    "telemetry": { "emit_interval_s": 2 }
  },
  "resources": {
    "locks": ["db:migrate", "vercel:deploy", "supabase:migrate"],
    "quotas": { "cpu_total": 8, "mem_total": "8Gi" }
  }
}

Notes
	•	Hard edges are must-happen-before; soft edges bias ordering but won’t block readiness.
	•	locks models external exclusivity (e.g., DB migration). This stays out of the DAG but constrains concurrency.
	•	idempotent + compensation is how SLAPS knows rollback semantics.
	•	checkpoint: true means capture a restore point after success (artifact hash, git ref, state snapshot).

2) SLAPS execution model (rolling frontier, event-driven)

Core idea: Maintain a ready frontier = tasks whose hard predecessors are done and whose required locks/quotas are currently available. Execute greedily subject to resource and policy constraints. On each task completion (or failure), recompute the frontier and keep flowing. No artificial “waves”.

Data structures
	•	predCount[id]: number of unfinished hard predecessors.
	•	succ[id]: adjacency list of outgoing edges.
	•	softSucc[id]: soft successors (used for tie-break).
	•	readyPQ: priority queue of ready tasks keyed by tuple
(-priorityBucket, depth, -parallel_opportunity, created_at)
where depth favors shallower critical path, and parallel_opportunity prefers tasks that unlock many successors.
	•	locksHeld: set of active lock names.
	•	quotaAvail: remaining CPU/MEM budget.
	•	errPatterns: rolling window of error fingerprints (for circuit breakers).

Scheduler loop (pseudocode)

initialize():
  compute predCount for all tasks using hard edges
  compute depth (longest-path from roots) for tie-breaking
  for each task with predCount==0 add to readyPQ if lock/quota feasible

while not all tasks completed:
  // 1) Launch as many as possible under constraints
  while readyPQ not empty and slots available and quotas/locks allow:
    t = readyPQ.pop()
    if acquireLocksAndQuotas(t):
      start(t)   // async
    else:
      defer t (requeue with small penalty)

  // 2) Wait for an event: completion/failure/update/timeout
  ev = waitForEvent()

  if ev.type == "task_succeeded":
    releaseLocksAndQuotas(ev.task)
    if ev.task.checkpoint: checkpoint(ev.task)
    for v in succ[ev.task]:
      predCount[v] -= 1
      if predCount[v]==0 and feasible(v): readyPQ.push(v)
    // soft successors get mild priority bump
    bumpSoftFollowers(ev.task)

  if ev.type == "task_failed":
    handleFailure(ev.task, ev.error)

  if ev.type == "hot_update":   // dynamic edge/node insert
    applyUpdate(ev.patch)       // adjust predCount/queues accordingly

  if ev.type == "timeout":
    abortOrRetry(ev.task)

Failure handling

handleFailure(t, error):
  recordErrorPattern(error)
  if circuitBreakerTriggered(error):
    replan()                 // see § hot updates below
    return

  if canRetry(t):
    scheduleRetry(t)         // backoff as per policy
  else:
    if t.idempotent == false and t.compensation:
      runCompensation(t)
    markAborted(t)
    propagateAbortIfRequired(t)

Circuit breaker triggers
	•	Repeated similar missing-dep errors (Cannot resolve @radix-ui/*, lucide-react, etc.) → batch-install task injection.
	•	Thrash: rapid fail-retry within thrash_window_s.
	•	Overrun: exceeding p95 * deadline_overrun_factor.

Deadlock avoidance
	•	Locks use try-lock with wait-die: older logical timestamp tasks may wait; younger ones fail fast and requeue (prevents circular waits).
	•	Global lock order by canonical name ensures consistent acquisition ordering when multiple locks are requested.

3) Hot updates (dynamic edges & re-plan)

When SLAPS detects a systemic pattern (e.g., a cascade of missing UI deps), it emits a planning patch and applies it live:

{
  "op": "inject_batch_task",
  "task": {
    "id": "AUTO.INSTALL.UI.DEPS",
    "name": "Install UI dependencies",
    "locks": ["pnpm:install"],
    "acceptance": [{"type":"cmd","cmd":"pnpm list class-variance-authority @radix-ui/react-avatar lucide-react"}],
    "idempotent": true,
    "checkpoint": false,
    "priority": 0.95
  },
  "edges": [
    {"from": "AUTO.INSTALL.UI.DEPS", "to": "P1.T013", "type": "soft"},
    {"from": "AUTO.INSTALL.UI.DEPS", "to": "P1.T014", "type": "soft"}
  ]
}

ApplyUpdate() behavior
	•	If a new node arrives with no hard predecessors, it can enter the frontier immediately (subject to locks/quotas).
	•	New hard edges increment predCount[target]. If the target is running, mark “must-finish” (don’t preempt); if not running and newly blocked, remove from readyPQ.
	•	Soft edges only affect priority ordering.

4) Correctness & guarantees
	•	Safety: No hard-edge violation (topological order maintained). No two tasks hold the same exclusive lock simultaneously. If non-idempotent, compensation (or manual escalation) restores invariants.
	•	Liveness: If at least one ready task is feasible under resource policy, the system eventually schedules it (fairness ensures no starvation).
	•	No deadlocks: Ensured by (a) acyclicity of hard edges; (b) global lock acquisition order + wait-die; (c) timeouts with requeue.
	•	Progress monotonicity: Completed hard predecessors never “un-complete.” Hot updates can only add constraints; they never invalidate recorded success.

5) Scheduling policy (greedy, bounded, resource-aware)

Objective (practical): Minimize makespan while avoiding thrash and respecting locks/quotas.

Heuristic key:
	•	Priority bucket: planner’s importance (SLO/critical path hint).
	•	Depth: prefer shallower nodes (finish “open fronts” quickly) or invert to drive critical path—configurable.
	•	Parallel opportunity: tasks that unlock many successors get a small boost.
	•	Soft-after bias: if A has a soft edge to B and both are ready, A first.

Complexity: Each event processes O(outdegree). With a binary heap, push/pop O(log n). Overall near-linear in edges for a typical run; hot updates add amortized O(ΔE log n).

6) State model per task

PENDING
  └─ready→ QUEUED
      └─dispatch→ RUNNING
          ├─success→ SUCCEEDED
          ├─fail & retries left→ BACKOFF (→ QUEUED)
          ├─fail & exhausted→ FAILED
          ├─timeout→ ABORTING → (compensate?) → FAILED
          └─cancel (global breaker)→ CANCELED

7) Rollback & checkpoints
	•	Checkpoint: (git commit or artifact hash + env snapshot). Stored after SUCCEEDED if checkpoint: true.
	•	Rollback policy:
	•	If an update requires reverting a cluster of tasks, choose latest common checkpoint dominating that cluster, run compensation for non-idempotent nodes between checkpoint and head, then re-enter scheduling.
	•	Minimal rollback set = nodes on any path feeding the failing node whose side-effects cannot be safely reapplied (derived from idempotent:false + effects).

8) Failure model & escalation
	•	Transient (network flake, npm registry, lock contention): auto-retry with jitter backoff.
	•	Systemic (repeated missing packages, schema drift): circuit breaker → inject diagnostic/repair task(s).
	•	Irreparable (broken invariant, missing human input): raise escalation block with precise questions (from TASKS acceptance gaps). Scheduler pauses only the affected subgraph; others proceed.

9) Metrics & telemetry (so you can see the “waves” light up)

Emit at telemetry.emit_interval_s:
	•	Per-task: state, start/stop, retries, last error digest.
	•	Frontier size over time; throughput (tasks/min).
	•	Lock utilization timeline; wait times.
	•	Makespan estimate (PERT on remaining subgraph).
	•	Circuit-breaker events & applied patches.

10) Minimal worked example (how the “rolling frontier” fixes your Wave-13 pain)
	1.	P1.T009, P1.T010 finish → P1.T011 becomes ready immediately (even if other Wave-5 tasks are still running).
	2.	P1.T011 starts; it fails with Cannot resolve @radix-ui/react-avatar.
	3.	Error digest matches 2 prior misses → breaker fires → inject AUTO.INSTALL.UI.DEPS as ready task with pnpm:install lock.
	4.	Installer runs once, succeeds; all UI-dependent tasks get a soft bias boost and re-enter the frontier.
	5.	No global pause. Other independent tasks keep flowing. Momentum preserved.

11) Why this works (short, blunt)
	•	We separate planning causality (hard edges) from execution scarcity (locks/quotas).
	•	We flow, don’t “wave”: tasks run the instant they’re truly unblocked.
	•	We adapt mid-run: hot updates patch the plan without trashing progress.
	•	We refuse to deadlock: consistent lock order + wait-die + timeouts.
	•	We own failure: idempotency/compensation rules make rollbacks sane.

If you want, I’ll turn this into:
	•	a coordinator.schema.json,
	•	a 150-line reference scheduler (TypeScript) with the exact loop above,
	•	and a tiny TUI that visualizes the frontier, the locks, and breaker events in real time.

Your move.