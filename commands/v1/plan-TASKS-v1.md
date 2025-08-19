# T.A.S.K.S: Tasks Are Sequenced Key Steps – LLM Execution Spec (v1, tool-free)

*This command makes the LLM simulate feature extraction → tasklist → deps → DAG (with transitive reduction) → waves (Kahn antichains) → sync points, with evidence + confidence and basic timeline stats.*

## Mission

From a raw tech plan doc, produce a validated execution plan:

1. Extract features → break into tasks
2. Infer structural dependencies (no resource edges)
3. Build an acyclic, reduced DAG
4. Generate waves via Kahn layering and basic balancing
5. Propose sync points with quality gates
6. Output artifacts as JSON/Markdown

### Hard Rules

- **No chain-of-thought in outputs.** Do private reasoning; output only the artifacts.
- **Grounding required.** Each task and each dependency must cite at least one evidence snippet (quote or char offsets or section+lines).
- **Ban "resource" dependencies.** Edges represent technical/sequential/infra/knowledge only.
- **ID direction = prereq → dependent.**
- **Granularity target:** 2–8h typical (cap 16h). Auto-split/merge to hit it.
- **Confidence:** [0..1] for each edge; treat < 0.7 as speculative during reporting (keep them visible, exclude from DAG build unless explicitly included).
- **Deterministic output:** Use content hashing for stable IDs across runs.
- **Security hygiene:** Redact secrets/keys in evidence quotes with [REDACTED].
- **Units explicit:** All durations in hours, assume 5-day workweek, task independence for wave math.

### Inputs

- **PLAN_DOC** (raw text/markdown)
- **Optional:** MIN_CONFIDENCE (default 0.7), MAX_WAVE_SIZE (default 30)

### Required Outputs (exact file fences)

1. `---file:features.json`
2. `---file:tasks.json`
3. `---file:dag.json`
4. `---file:waves.json`
5. `---file:Plan.md`

## Schemas

### features.json

```json
{
  "generated": {
    "by": "T.A.S.K.S v1",
    "timestamp": "2025-08-13T00:00:00Z",
    "contentHash": "<sha256(normalizeJSON(file))>"
  },
  "features": [
    {
      "id": "F001",
      "title": "short name",
      "description": "1-3 sentences",
      "priority": "critical|high|medium|low",
      "source_evidence": [
        {
          "quote": "...",
          "loc": {
            "start": 123,
            "end": 178
          },
          "section": "H2: Authentication",
          "startLine": 123,
          "endLine": 131
        }
      ]
    }
  ]
}
```

### tasks.json

```json
{
  "meta": {
    "min_confidence": 0.7,
    "notes": "assumptions if any",
    "autonormalization": {
      "split": ["P1.T001->P1.T001a,P1.T001b (>16h)"],
      "merged": ["P1.T099,P1.T100 (<0.5h combined)"]
    }
  },
  "generated": {
    "by": "T.A.S.K.S v1",
    "timestamp": "2025-08-13T00:00:00Z",
    "contentHash": "<sha256(normalizeJSON(file))>"
  },
  "tasks": [
    {
      "id": "P1.T001",
      "feature_id": "F001",
      "title": "verb-first outcome",
      "description": "what artifact/contract/code is produced",
      "category": "foundation|implementation|integration|optimization",
      "skillsRequired": ["backend", "node"],
      "duration": {
        "optimistic": 2,
        "mostLikely": 4,
        "pessimistic": 8
      },
      "durationUnits": "hours",
      "interfaces_produced": ["AuthAPI:v1", "UserClaimsSchema:v1"],
      "interfaces_consumed": ["DB:UsersSchema:v1"],
      "acceptance_checks": [
        {
          "type": "command",
          "cmd": "npm test",
          "expect": {
            "passRateGte": 0.95
          }
        },
        {
          "type": "artifact",
          "path": "docs/api/v1.md",
          "expect": {
            "exists": true
          }
        },
        {
          "type": "metric",
          "name": "staging_p95_ms",
          "expect": {
            "lte": 200
          }
        }
      ],
      "source_evidence": [
        {
          "quote": "...",
          "loc": {
            "start": 200,
            "end": 240
          },
          "section": "H3: Implementation Details",
          "startLine": 45,
          "endLine": 52
        }
      ],
      "contentHash": "a1b2c3d4e5f6789a"
    }
  ],
  "dependencies": [
    {
      "from": "P1.T001",
      "to": "P1.T002",
      "type": "technical|sequential|infrastructure|knowledge",
      "reason": "why T002 requires T001",
      "evidence": [
        {
          "type": "doc|code|human|history",
          "reason": "...",
          "confidence": 0.9
        }
      ],
      "confidence": 0.85,
      "isHard": true
    }
  ]
}
```

### dag.json

```json
{
  "ok": true,
  "errors": [],
  "warnings": [],
  "generated": {
    "by": "T.A.S.K.S v1",
    "timestamp": "2025-08-13T00:00:00Z",
    "contentHash": "<sha256(normalizeJSON(file))>"
  },
  "metrics": {
    "minConfidenceApplied": 0.7,
    "keptByType": {
      "technical": 45,
      "sequential": 32,
      "infrastructure": 18,
      "knowledge": 26
    },
    "droppedByType": {
      "technical": 3,
      "sequential": 2,
      "infrastructure": 0,
      "knowledge": 1
    },
    "nodes": 84,
    "edges": 121,
    "edgeDensity": 0.017,
    "widthApprox": 18,
    "widthMethod": "kahn_layer_max",
    "longestPath": 7,
    "isolatedTasks": 0,
    "lowConfidenceEdgesExcluded": 6,
    "verbFirstPct": 0.92,
    "meceOverlapSuspects": 0
  },
  "topo_order": ["P1.T005", "P1.T001", "..."],
  "reduced_edges_sample": [
    ["P1.T001", "P1.T010"],
    ["P1.T010", "P1.T020"]
  ],
  "softDeps": [
    {
      "from": "P1.T010",
      "to": "P1.T020",
      "type": "knowledge",
      "reason": "Lessons from spike inform API design",
      "confidence": 0.6,
      "isHard": false
    }
  ],
  "lowConfidenceDeps": [
    {
      "from": "P1.T012",
      "to": "P1.T017",
      "type": "sequential",
      "confidence": 0.62,
      "reason": "UI mockups may inform backend data model"
    }
  ],
  "cycle_break_suggestions": []
}
```

### waves.json

```json
{
  "planId": "PLAN-a1b2c3d4",
  "generated": {
    "by": "T.A.S.K.S v1",
    "timestamp": "2025-08-13T00:00:00Z",
    "contentHash": "<sha256(normalizeJSON(file))>"
  },
  "config": {
    "maxWaveSize": 30,
    "barrier": {
      "kind": "quorum",
      "quorum": 0.95
    }
  },
  "waves": [
    {
      "waveNumber": 1,
      "tasks": ["P1.T005", "P1.T001", "..."],
      "estimates": {
        "units": "hours",
        "p50Hours": 10.5,
        "p80Hours": 12.0,
        "p95Hours": 14.2
      },
      "barrier": {
        "kind": "quorum",
        "quorum": 0.95,
        "timeoutMinutes": 120,
        "fallback": "deferOptional",
        "gateId": "W1→W2-q95"
      }
    }
  ]
}
```

### Plan.md (sections required)

- Header metrics (nodes, edges, widthApprox, longest path, edgeDensity, verbFirstPct)
- Wave list with barriers and P50/P80/P95 (units: hours)
- Sync Points & Quality Gates (checklists with timeouts and fallbacks)
- Low-Confidence Dependencies section (confidence < min_confidence edges)
- Parallelizable/Soft Dependencies section (isHard=false edges)
- Auto-normalization actions taken (if any)

## Algorithm

*The model must simulate these steps:*

### 1. Feature extraction

- Pull 5–25 features from PLAN_DOC: user-visible capabilities + enabling infrastructure
- Each: title, 1–3 sentence description, priority, at least one evidence snippet

### 2. Task breakdown (MECE)

- For each feature, create tasks (2–8h typical; cap 16h)
- Each task: skillsRequired[], PERT durations {a,m,b}, acceptance checks
- Ensure MECE ownership (no overlap). If overlap detected, refactor and note in tasks.json.meta.notes

### 3. Dependency discovery (structural only)

Classify edges as:
- **technical** (interfaces/artifacts/contracts required)
- **sequential** (information/order)
- **infrastructure** (env/tools)
- **knowledge** (research/learning)

Edge direction = prereq → dependent. Add reason, evidence[], and confidence ∈ [0,1].

**Hard vs Soft dependencies:**
- **isHard=true:** Critical path dependencies that block execution
- **isHard=false:** Soft dependencies that provide benefit but allow parallelization

Keep speculative edges (confidence < min_confidence) and soft edges (isHard=false) in tasks.json output, but they'll be excluded from the built DAG and listed separately for review.

### 4. DAG build + transitive reduction

- **Filter edges:** include only (confidence ≥ min_confidence AND isHard==true) for the graph. Keep others as "softDeps" and "lowConfidenceDeps" in dag.json
- **Cycle resolution:** attempt ≤ 2 automated passes. If still cyclic:
  - Set dag.json.ok=false
  - Populate dag.json.errors with cycle edges (e.g., "cycle detected: P1.T014→P1.T022→P1.T014")
  - Produce dag.json.cycle_break_suggestions[] with specific recommendations:
    - "Split P1.T022 into P1.T022a (API contract) and P1.T022b (implementation)"
    - "Insert interface: 'PaymentsAPI:v1' between T014 and T022"
  - Stop; do not emit waves.json; instruct caller to fix
- **Topological order:** Kahn's algorithm
- **Transitive reduction:** remove edge (u,v) if another path u → ... → v exists (keep minimal edges)
- **Compute metrics:** nodes, edges, edgeDensity = edges / (nodes*(nodes-1)), widthApprox (max antichain size via Kahn layering), longest path length, isolated tasks, confidence/type breakdowns, MECE overlap suspects
- **Heuristics:** edgeDensity < 0.05 → likely missing deps (warn). > 0.5 → over-constrained (warn)

### 5. Wave generation (Kahn antichains + balancing)

- **Layering:** Iteratively select zero-in-degree nodes → that set is Wave k; remove them; repeat
- **Balance:**
  - If a wave has > MAX_WAVE_SIZE, split by:
    - Keep tasks that produce/consume the **same interface** together (check interfaces_produced/interfaces_consumed arrays) unless the next barrier is kind=full
    - Prefer to move low-priority (shortest distance-to-sink) tasks to a new wave
  - If a wave has < 5 tasks and the next wave's tasks have no extra prereqs, merge waves
- **Wave duration estimates** (no Monte Carlo): For each task:
  - PERT mean μ = (a + 4m + b)/6, std σ = (b - a)/6
  - Wave p50 ≈ max_i μ_i
  - Wave p80 ≈ max_i (μ_i + 0.84·σ_i)
  - Wave p95 ≈ max_i (μ_i + 1.65·σ_i)
  - Units: hours, assumes 5-day workweek, task independence

*We're approximating the max of independent task durations—good enough for planning.*

### 6. Sync points & gates

- Between each wave k → k+1 define a barrier:
- **Default:** quorum with quorum = 0.95 (95% of wave tasks done)
- **Enhanced barrier config:** timeoutMinutes, fallback (proceed|halt|deferOptional), gateId for audit
- Use full barrier only if interfaces are tightly coupled
- **Quality gate checklist** (machine-checkable where possible). Example:
  - ✅ Unit tests >= 95% pass
  - ✅ API contract v1 published at /docs/api
  - ✅ Lint passes (0 errors)
  - ✅ p95 latency <= 200ms in staging
  - ✅ Security scan: 0 high severity

## Output Procedure

*Model must follow in one pass:*

1. **Auto-normalize tasks:** If mean > 16h, auto-split by artifact boundary; if < 0.5h, auto-merge to nearest parent feature task. Record actions in tasks.json.meta.autonormalization
2. **Emit features.json** with flexible evidence format and generated provenance block
3. **Emit tasks.json** (with all deps including low-confidence and soft deps, structured acceptance checks, interface inventories, content hashes, generated provenance)
4. **Build DAG in-model** using only edges with (confidence ≥ min_confidence AND isHard==true); attempt ≤2 cycle resolution passes; compute metrics; emit dag.json with lowConfidenceDeps and softDeps separated
5. **Generate waves** per algorithm (only if dag.json.ok==true); emit waves.json with planId derived from tasks.json hash
6. **Emit Plan.md:**
   - Summary metrics (using new field names, including verbFirstPct)
   - Waves with P50/P80/P95 and barrier details (kind, quorum, timeout, fallback, gateId)
   - Sync Points & Quality Gates (specific, machine-checkable checklists)
   - Low-Confidence Dependencies section (confidence < min_confidence edges)
   - Parallelizable/Soft Dependencies section (isHard=false edges)
   - Auto-normalization actions taken (if any)

## System Prompt

```
You are a planning engine that simulates a project planning CLI.
Produce ONLY the requested artifacts (JSON/Markdown) in file fences.
Do not reveal intermediate reasoning or chain-of-thought.
All tasks and edges must include evidence from the source document (accept either char offsets or section+lines).
Edges must be structural (technical, sequential, infrastructure, knowledge); resource constraints are disallowed.
Enforce dependency direction as prereq → dependent.
Target task granularity 2–8h (cap 16h); auto-split/merge as needed and record in meta.autonormalization.
Use min_confidence (default 0.7) and isHard flag to exclude low-confidence/soft edges from the built DAG, but list them for human review.
Generate waves using Kahn layering; apply interface-cohesion-aware balancing using interfaces_produced/interfaces_consumed arrays; compute wave P50/P80/P95 using PERT approximations.
Cap cycle resolution to 2 passes; fail with specific suggestions if still cyclic.
Ensure deterministic output via content hashing and normalization; derive planId from tasks.json hash.
Require at least one machine-verifiable acceptance check per task (command|artifact|metric types only).
Populate dag.json.lowConfidenceDeps[] for deps with confidence < min_confidence.
Compute dag.json.metrics.meceOverlapSuspects via title-similarity heuristic (>80% token overlap within a feature).
Before emitting evidence quotes, apply secret redaction heuristics; replace matches with [REDACTED].

Secret redaction rules - treat as secret and replace with "[REDACTED]" if a quote matches:
- 20+ char base64ish strings
- Hex strings ≥32 chars  
- Common key prefixes (AKIA|ASIA|ghp_|xoxb-|xoxp-|sk_live_|sk_test_)
- "-----BEGIN * PRIVATE KEY-----" ... "-----END * PRIVATE KEY-----"
Redact minimally (keep surrounding context).

Normalizing for hashing:
- Sort object keys lexicographically
- Remove null/undefined
- Trim strings
- Sort arrays of strings
- Round numbers to 3 decimals
```

## User Prompt Template

```
MIN_CONFIDENCE: {0.7 or custom}
MAX_WAVE_SIZE: {30 or custom}

INPUT DOCUMENT:
<<<
{PLAN_DOC}
>>>

Produce artifacts in this order and format:
---file:features.json
{...}
---file:tasks.json
{...}
---file:dag.json
{...}
---file:waves.json
{planId derived as "PLAN-" + sha256(normalizeJSON(tasks.json)).slice(0,8)}
---file:Plan.md
(markdown)
```

## Sanity Checklist

*The model should self-verify before emitting:*

- No cycles; dag.json.ok == true (if false, waves.json MUST NOT be emitted)
- No "resource" edges
- ≥70% tasks have durations; no task mean > 16h unless justified
- Wave 1 width > 1 unless problem is inherently serial
- Each wave has a barrier with timeout, fallback, and gateId; defaults to quorum 95%
- Low-Confidence Dependencies section present in Plan.md
- Parallelizable/Soft Dependencies section in Plan.md
- Every task has at least one structured machine-verifiable acceptance check (command|artifact|metric)
- Verb-first naming ≥80% compliance (and displayed in Plan.md header)
- Evidence quotes contain no secrets (redacted according to heuristics)
- Interface inventories populated for cohesion-aware wave splitting
- Generated provenance blocks present in all JSON files

## Scoring Rubric + Evaluator Spec

Use this to auto-check outputs and reject bad runs. It's strict where it should be, flexible where it matters.

### Simulated-Planctl Evaluator (v1)

**What it judges:**

Artifacts expected from the model:
- features.json
- tasks.json
- dag.json
- waves.json
- Plan.md

If any are missing → auto-reject.

### Hard fails (auto-reject, no debate)

1. **Cycles present:** dag.json.ok != true or errors contains "cycle"
2. **Resource edges used anywhere**
3. **Edge direction wrong** (any dep where from == dependent or obvious invert)
4. **< 50% tasks have PERT durations**
5. **No evidence on any task or any dependency**
6. **Wave violations:** a wave with > maxWaveSize and no split, or zero waves
7. **Missing machine-verifiable acceptance checks** on any task (must use command|artifact|metric types)
8. **Secrets in evidence quotes** (not redacted according to heuristics)
9. **Missing structured acceptance checks** (string format instead of {type, cmd/path/name, expect} objects)

If any of the above hits → grade REJECT with reasons.

### Scoring rubric (100 points)

| Category | Weight | How to score |
|----------|--------|--------------|
| A. Structural Validity | 20 | +20 if acyclic, edges all structural (tech/sequential/infra/knowledge), no resource edges, from→to correct. Deduct 5 for each violation (min 0). |
| B. Coverage & Granularity | 15 | % tasks with durations (≥70% = full 5, 50–69% = 2, <50% = 0); median PERT mean in 2–8h (5 points if true; 2 if slight drift); outliers resolved (≥90% within 0.5–16h = 5, else 0–3). Auto-normalization applied properly (2). |
| C. Evidence & Confidence | 15 | Tasks with evidence (≥95% = 5, 80–94% = 3, else 0); deps with evidence (same scale, 5); low-confidence table present in Plan.md and min_confidence applied in DAG (3); soft deps section present (2). |
| D. DAG Quality | 15 | EdgeDensity in [0.05, 0.5] (5 if yes; 2 if warned but justified; 0 if egregious). WidthApprox > 1 in Wave 1 (5 if yes; 0 if inherently serial but not explained; 3 if explained). Longest path vs wave count plausible (3). Confidence/type breakdowns present (2). |
| E. Wave Construction | 15 | Kahn layering respected (5), splits/merges sensible vs maxWaveSize with interface cohesion (5), per-wave P50/P80/P95 present and monotonic with units (5). |
| F. Sync Points & Gates | 10 | Each boundary has barrier kind with timeout/fallback/gateId (5), gate checklist is specific & machine-checkable (5). |
| G. MECE & Naming | 10 | No duplicate/overlapping tasks (3), verb-first titles ≥80% (4), machine-verifiable acceptance checks (3). |

### Passing bands

- **Excellent (90–100):** Ship it
- **Good (80–89):** Minor edits
- **Needs Work (70–79):** Return with specific fixes
- **Reject (<70):** Re-run generation

## Metrics the evaluator should compute

- **Durations coverage:** count(tasks with duration)/count(tasks)
- **PERT mean per task:** (a + 4m + b)/6; median and outlier counts (<0.5h or >16h)
- **Dependency sanity:**
  - Count by type; ensure none are resource
  - Evidence presence rate (tasks, deps)
  - Confidence distribution; count of < min_confidence
  - Hard vs soft dependency breakdown
- **Graph stats** (from dag.json or recomputed if needed):
  - nodes, edges, edgeDensity edges/(n*(n-1)), widthApprox (max antichain), longest path, isolated tasks
- **Wave checks:**
  - wave sizes vs config.maxWaveSize
  - P50 ≤ P80 ≤ P95 for every wave
  - Units explicitly stated
- **Quality metrics:**
  - % verb-first titles (regex check)
  - Machine-verifiable acceptance checks coverage
  - Barriers present with timeout/fallback/gateId
- **Security hygiene:** no unredacted secrets in evidence

## Evaluator output schema

```json
{
  "grade": "EXCELLENT | GOOD | NEEDS_WORK | REJECT",
  "score": 0,
  "reasons": ["short bullets of key issues"],
  "category_scores": {
    "A": 20,
    "B": 15,
    "C": 15,
    "D": 15,
    "E": 15,
    "F": 10,
    "G": 10
  },
  "metrics": {
    "tasks": 0,
    "durationsCoveragePct": 0,
    "pertMedian": 0,
    "outliersLow": 0,
    "outliersHigh": 0,
    "depsTotal": 0,
    "depsLowConfidence": 0,
    "depsSoft": 0,
    "evidenceCoverageTasksPct": 0,
    "evidenceCoverageDepsPct": 0,
    "verbFirstPct": 0,
    "machineVerifiableChecksPct": 0,
    "dag": {
      "nodes": 0,
      "edges": 0,
      "edgeDensity": 0,
      "widthApprox": 0,
      "longestPath": 0,
      "isolated": 0
    },
    "waves": [
      {
        "waveNumber": 1,
        "count": 0,
        "p50": 0,
        "p80": 0,
        "p95": 0,
        "barrier": "quorum",
        "quorum": 0.95,
        "hasTimeout": true,
        "hasFallback": true,
        "hasGateId": true
      }
    ]
  },
  "autoReject": false,
  "autoRejectReasons": []
}
```

## Evaluator steps (pseudo-logic)

1. **Presence check:** all 5 artifacts exist → else REJECT
2. **Parse JSON:** basic schema checks (ids present, arrays, types)
3. **Hard-fail scan:**
   - dag.json.ok !== true or "cycle" in errors → reject
   - Any dep type resource or "resource" in reason/type → reject
   - Durations coverage < 50% → reject
   - Missing evidence in any task or any dependency → reject
   - Wave count == 0, or any wave size > maxWaveSize without split → reject
4. **Compute metrics** (see above)
5. **Category scoring** using the rubric; clamp each to [0, weight]
6. **Total score:** assign band; include top 5 reasons (bullets)
7. **Return evaluation.json** per schema

## Quick heuristics the evaluator should use

- **Verb-first titles:** starts with a verb (regex-ish):
  ```regex
  /^(add|build|create|define|design|document|enable|export|generate|implement|integrate|migrate|publish|refactor|remove|rename|setup|test|upgrade)\b/i
  ```
- **Score partially** if ≥80% pass
- **Overlap sniff:** two tasks with >80% token overlap in titles → warn; if same feature_id and similar titles → deduct in G
- **Low-confidence deps section required:** Plan.md must contain "Low-Confidence Dependencies" section listing confidence < min_confidence edges. If absent but such edges exist in tasks.json → deduct in C
- **Soft deps section required:** Plan.md must contain "Parallelizable/Soft Dependencies" section listing isHard=false edges
- **EdgeDensity sanity:**
  - < 0.05 → warn "likely missing deps"
  - > 0.5 → warn "over-constrained"
- **Waves monotonicity:** for each wave, check p50 ≤ p80 ≤ p95; else deduct in E
- **Machine-verifiable checks:** acceptance_checks must be structured objects with type (command|artifact|metric), cmd/path/name, and expect fields with measurable criteria
- **Security hygiene:** evidence quotes must not contain secrets per redaction heuristics
- **Interface inventories:** tasks should have interfaces_produced/interfaces_consumed arrays populated for cohesion-aware splitting

## Minimal evaluator prompt (for an LLM judge)

### System

```
You are an exacting project plan evaluator. Read the five artifacts, compute the rubric scores, and output only evaluation.json per the schema. Do not include chain-of-thought. Be strict on hard fail rules.
```

### User

```
Artifacts:
---file:features.json
{...}
---file:tasks.json
{...}
---file:dag.json
{...}
---file:waves.json
{...}
---file:Plan.md
(markdown here)

Constraints:
- min_confidence default 0.7 unless meta overrides
- maxWaveSize default 30 unless config overrides

Now evaluate using the rubric and return a single JSON object:
(evaluation.json as per schema)
```

## Acceptance gates (use these in CI)

- **Merge gate:** grade ∈ {EXCELLENT, GOOD} and score ≥ 80
- **Soft gate (allow override):** NEEDS_WORK with score ∈ [70,79] and no Hard fails
- **Fail:** anything else

## Common fix suggestions (auto-generate from failures)

- **Cycles** → "Split large tasks by artifact boundary; insert interface contracts as separate tasks; ensure direction prereq→dependent."
- **Sparse edgeDensity** → "Add explicit interface or infra deps; avoid hand-wavy sequencing."
- **Low durations coverage** → "Add PERT durations to at least 70% of tasks; keep means 2–8h."
- **Evidence gaps** → "Attach quoted spans (or section+lines) for each task/edge; set confidence accordingly."
- **Waves too big** → "Split by category and by longest-path priority; enforce maxWaveSize while preserving interface cohesion."
- **Missing machine checks** → "Add measurable acceptance criteria with thresholds, commands, or artifact paths."
- **Security issues** → "Redact suspected secrets/keys/tokens in evidence quotes with [REDACTED]."

## Failure Handling & Auto-Remediation

### When to trigger

Run this protocol if the evaluator returns:
- grade: REJECT, or
- any Hard fail, or  
- score < 80 and CI policy requires ≥80

### Outputs (always produce)

1. `---file:repair_report.json` — What failed, what changed, why
2. Updated artifacts for any file you modify (same fences as spec)

### repair_report.json schema

```json
{
  "attempt": 1,
  "grade_before": "REJECT",
  "hard_fails": ["cycles_present", "missing_evidence"],
  "soft_issues": ["edgeDensity_low"],
  "actions": [
    {
      "type": "split_task",
      "target": "P1.T022",
      "reason": "mean>16h; cycle break"
    },
    {
      "type": "insert_interface",
      "name": "PaymentsAPI:v1",
      "between": ["P1.T014", "P1.T022"],
      "reason": "break cycle; make dep explicit"
    },
    {
      "type": "add_evidence",
      "target": "dep(P1.T005->P1.T009)",
      "source": "PLAN_DOC lines 142-148"
    }
  ],
  "results": {
    "dag_ok": true,
    "edgeDensity_before": 0.041,
    "edgeDensity_after": 0.067,
    "isolatedTasks_before": 3,
    "isolatedTasks_after": 0
  },
  "notes": "Deterministic IDs preserved where content unchanged."
}
```

### Auto-Repair Protocol (three passes max)

#### Pass L1: Fast, surgical fixes (no semantics drift)

1. **Cycles present**
   - Split any >16h task on the cycle boundary by artifact/interface
   - Insert interface contract task (e.g., XAPI:v1) to decouple
   - Rebuild DAG (2 passes max). If still cyclic → escalate to L2

2. **Resource edges detected**
   - Remove the edge. If structural intent exists, replace with technical|sequential|infrastructure|knowledge and add evidence
   - Add a note in tasks.json.meta.notes

3. **Edge direction wrong**
   - Swap to prereq → dependent and re-validate
   - If ambiguous, create an interface task as the prereq

4. **Durations coverage < 50%**
   - Infer PERT from neighbor tasks in same feature; set:
     - optimistic = 0.5×peer_median
     - mostLikely = peer_median  
     - pessimistic = 2×peer_median (cap at 16h)
   - Record in meta.notes: "inferred durations"

5. **Missing evidence (tasks or deps)**
   - Search PLAN_DOC for nearest section/lines; quote + confidence
   - If no support in doc: downgrade edge to isHard=false and set confidence=0.49; exclude from DAG; add to Plan.md Low-Confidence
   - If a task has no documentary basis: mark assumption in meta.notes, or merge into its parent feature's task

6. **Waves > maxWaveSize / zero waves**
   - Split by interface cohesion first, then by shortest distance-to-sink
   - Ensure each wave keeps related producer/consumer tasks together unless barrier=full

7. **Missing machine-verifiable acceptance checks**
   - For each task, add at least one {type: command|artifact|metric, expect: ...} object
   - Prefer npm test, pytest -q, artifact existence, or latency/error-rate metrics

8. **Secrets in evidence**
   - Re-emit quotes with the redaction heuristics; update Plan.md note: "Evidence redacted"

Run evaluator again. If REJECT persists → L2.

#### Pass L2: Structural tightening (small semantics allowed)

1. **edgeDensity < 0.05 (too sparse)**
   - Promote well-supported knowledge or infrastructure edges from soft → hard only if confidence ≥ MIN_CONFIDENCE and evidence exists
   - Add missing infra setup prereqs ("provision staging", "seed test data") where the doc clearly implies them

2. **edgeDensity > 0.5 (over-constrained)**
   - Remove transitive hard edges (keep them as soft rationale in dag.json.softDeps)
   - Prefer explicit interfaces to many-to-many hard constraints

3. **Isolated tasks > 0**
   - Either (a) move to Wave 1 foundation, or (b) connect via knowledge dep to the feature's first implementation task with evidence

4. **Verb-first < 80%**
   - Auto-rename non-verb titles using the verb list (setup|implement|publish…)

5. **Wave monotonicity broken (p50/p80/p95)**
   - Recompute PERT; correct σ = (b−a)/6
   - If outlier dominates max, split that task per artifact

Run evaluator again. If still REJECT or score < 80 → L3.

#### Pass L3: Escalate or ask for missing inputs

- **If evidence cannot be found for required edges/tasks:**
  - Do not hallucinate. Keep as soft, low-confidence; annotate in Plan.md "Needs Product/Arch decision"
  - Emit repair_report.json.escalation with precise questions:
    - "Confirm if API schema must exist before UI scaffolding (promote P1.T012→P1.T017 to technical hard?)"
    - "Provide acceptance metric for 'batch job success'; currently unspecified"

Stop after L3. Return updated artifacts + escalation questions.

### Triage Matrix (hard fail → action)

| Hard fail | Primary fix | Secondary fix |
|-----------|-------------|---------------|
| cycles_present | split task; insert interface | downgrade ambiguous edges to soft, re-plan |
| resource_edges | replace with structural dep+evidence | remove if truly scheduling-only |
| direction_wrong | swap to prereq→dependent | add interface task to clarify |
| durations_coverage_lt_50 | infer PERT from peers | mark assumption; cap pessimistic at 16h |
| missing_evidence | find quote; add confidence | downgrade to soft (exclude from DAG) |
| waves_violation | interface-cohesion split; cap size | merge undersized neighbor waves |
| missing_checks | add structured acceptance checks | tie to CI commands or artifacts |
| secrets_present | apply redaction heuristics | replace quote with line refs only |

### Auto-Repair System Prompt

```
You are a repair engine for T.A.S.K.S. Given the five artifacts and an evaluation failure, apply the Auto-Repair Protocol (L1→L3). 
- Make the smallest changes necessary.
- Preserve IDs and content hashes when content is unchanged.
- Do not fabricate evidence; demote to soft/low-confidence instead.
- Emit repair_report.json and updated artifacts (only those changed).
- Re-run internal validation before emitting.
```

### Auto-Repair User Prompt Template

```
EVALUATION:
<<<
{evaluation.json from prior run}
>>>

ARTIFACTS (original):
---file:features.json
{...}
---file:tasks.json
{...}
---file:dag.json
{...}
---file:waves.json
{...}
---file:Plan.md
{...}

Apply Auto-Repair Protocol. Return:
---file:repair_report.json
{...}
[updated files as needed in their fences]
```

### Guardrails

- **Determinism**: Do not churn IDs or planId unless affected content changes (recompute hashes only for changed files)
- **Minimal footprint**: Prefer adding interface tasks/edges over sweeping rewrites
- **Evidence-first**: If you can't cite, don't hard-link it
- **No lipstick on a DAG**: Don't "pad" density with bogus edges

## Implementation Notes

### Structured acceptance checks schema:
```json
{
  "type": "command|artifact|metric",
  "cmd": "npm test",                    // for command type
  "path": "docs/api/v1.md",             // for artifact type  
  "name": "staging_p95_ms",             // for metric type
  "expect": {
    "passRateGte": 0.95,               // for command
    "exists": true,                    // for artifact
    "lte": 200                         // for metric
  }
}
```

### PlanId derivation:
```javascript
planId = "PLAN-" + sha256(normalizeJSON(tasks.json)).slice(0,8)
```

### Secret redaction heuristics (applied before emitting quotes):
- 20+ char base64ish strings → `[REDACTED]`
- Hex strings ≥32 chars → `[REDACTED]`
- Key prefixes: AKIA|ASIA|ghp_|xoxb-|xoxp-|sk_live_|sk_test_ → `[REDACTED]`
- PEM blocks: `-----BEGIN * PRIVATE KEY-----` ... `-----END * PRIVATE KEY-----` → `[REDACTED]`

### Required structured fields:
- **tasks.json:** durationUnits, interfaces_produced/consumed, structured acceptance_checks, generated block
- **dag.json:** lowConfidenceDeps array, meceOverlapSuspects metric, enhanced cycle_break_suggestions
- **waves.json:** planId from tasks.json hash, generated block
- **Plan.md:** verbFirstPct in header, separate sections for low-confidence and soft deps

### Optional enhancements for future versions:
- **Duplicate detection:** cosine similarity on titles; auto-merge if >0.9 with same feature
- **Interface inventory:** enhanced interfaces with versioning for better cohesion analysis
- **Risk tagging:** riskLevel: low|med|high → show "critical path risk" in dag.json.metrics
