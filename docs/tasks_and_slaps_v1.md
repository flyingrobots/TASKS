TASKS + SLAPS: Formal Declaration (v1)

Dream it. Plan it. Build it. Ship it.
The planner with proofs; the runtime that finishes.

⸻

0) Scope (Normative)

This document defines the contract between T.A.S.K.S. (the planner) and S.L.A.P.S. (the executor). It specifies:
	•	the required artifacts T.A.S.K.S. must output,
	•	the canonical JSON Schemas for those artifacts,
	•	the coordinator contract consumed by S.L.A.P.S.,
	•	the execution algorithm, including rollback & replanning rules,
	•	the logging contract and deterministic hashing requirements.

Artifacts: features.json, tasks.json, dag.json, waves.json, Plan.md, Decisions.md (exact set & naming are required).  ￼

⸻

1) Principles (Authoritative)
	1.	Evidence-first planning. Every task and dependency cites source evidence and carries a confidence score.  ￼
	2.	Codebase-aware by default. Before planning, inventory the repo (polyglot ast-grep patterns included).  ￼  ￼
	3.	Boundaries and receipts. Each task defines complexity, done-criteria, scope, and logging guidance.  ￼
	4.	Resources are real. Structural deps remain king, but infrastructure resource edges and capacities are allowed (no human deps).  ￼
	5.	Deterministic outputs. Canonicalize and hash artifacts; record the hash in meta.artifact_hash.  ￼  ￼
	6.	Quality metrics are explicit. Edge density warning band 0.01–0.5; naming & evidence targets apply.  ￼  ￼

⸻

2) Artifact Set (Produced by T.A.S.K.S.)
	•	features.json – capabilities & outcomes
	•	tasks.json – atomic, bounded tasks + boundaries, evidence, resources
	•	dag.json – nodes/edges + metrics (edge density, critical path set)
	•	waves.json – optional pre-layering (SLAPS may ignore if doing rolling-frontier)
	•	Plan.md – human narrative with codebase analysis & links to evidence
	•	Decisions.md – alternatives & rationale

Generation order is fixed and MUST be respected for hashing/provenance.  ￼

⸻

3) Logging Contract (For SLAPS agents & evaluators)

All task executors MUST emit JSON Lines (one JSON object per line) with:
timestamp (ISO-8601), task_id, step, status (start|progress|done|error), message, optional percent (0–100), optional data (object). On status=error, include data.error_code and data.stack if available.  ￼

Validation rule: “Execution guidance” in tasks MUST include these logging instructions verbatim.  ￼

Example (JSONL):

{"timestamp":"2025-08-23T08:00:01Z","task_id":"P1.T001","step":"generate-migrations","status":"start","message":"begin"}
{"timestamp":"2025-08-23T08:00:12Z","task_id":"P1.T001","step":"generate-migrations","status":"progress","percent":60,"message":"users table created"}
{"timestamp":"2025-08-23T08:00:44Z","task_id":"P1.T001","step":"generate-migrations","status":"done","message":"both migrations written"}

￼

⸻

4) Deterministic Output & Hashing (Reproducibility)

Canonicalize JSON (lexicographic keys at all depths, minimal numeric form, UTF-8, newline-terminated), LF line endings, 2-space indent, no trailing spaces; compute SHA-256 and write into meta.artifact_hash for each JSON artifact; list all hashes in Plan.md/Decisions.md.  ￼  ￼

⸻

5) JSON Schemas (Draft 2020-12)

$id values are illustrative. All objects below MUST include meta.artifact_hash post-canonicalization (see §4).

5.1 features.json

{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://spec.tasks-slaps.io/schemas/features.json",
  "type": "object",
  "required": ["meta", "features"],
  "properties": {
    "meta": {
      "type": "object",
      "required": ["generated_by", "generated_at", "artifact_hash"],
      "properties": {
        "generated_by": { "type": "string" },
        "generated_at": { "type": "string", "format": "date-time" },
        "artifact_hash": { "type": "string", "pattern": "^[a-f0-9]{64}$" }
      },
      "additionalProperties": false
    },
    "features": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["id", "name", "outcome", "priority"],
        "properties": {
          "id": { "type": "string" },
          "name": { "type": "string" },
          "outcome": { "type": "string" },
          "priority": { "type": "string", "enum": ["P0", "P1", "P2", "P3"] },
          "evidence": { "$ref": "#/$defs/evidenceArray" }
        },
        "additionalProperties": false
      }
    }
  },
  "$defs": {
    "evidence": {
      "type": "object",
      "required": ["type", "source", "confidence", "rationale"],
      "properties": {
        "type": { "type": "string", "enum": ["plan","code","commit","doc"] },
        "source": { "type": "string" },
        "excerpt": { "type": "string" },
        "confidence": { "type": "number", "minimum": 0, "maximum": 1 },
        "rationale": { "type": "string" }
      },
      "additionalProperties": false
    },
    "evidenceArray": { "type": "array", "items": { "$ref": "#/$defs/evidence" }, "minItems": 1 }
  }
}

Evidence schema matches the explicit requirement.  ￼

5.2 tasks.json

{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://spec.tasks-slaps.io/schemas/tasks.json",
  "type": "object",
  "required": ["meta", "tasks", "dependencies"],
  "properties": {
    "meta": {
      "type": "object",
      "required": ["min_confidence", "artifact_hash", "codebase_analysis"],
      "properties": {
        "min_confidence": { "type": "number", "minimum": 0, "maximum": 1 },
        "artifact_hash": { "type": "string", "pattern": "^[a-f0-9]{64}$" },
        "codebase_analysis": {
          "type": "object",
          "properties": {
            "existing_apis": { "type": "array", "items": { "type": "string" } },
            "reused_components": { "type": "array", "items": { "type": "string" } },
            "extension_points": { "type": "array", "items": { "type": "string" } },
            "shared_resources": {
              "type": "object",
              "additionalProperties": {
                "type": "object",
                "oneOf": [
                  {
                    "properties": {
                      "type": { "const": "exclusive" },
                      "location": { "type": "string" },
                      "constraint": { "type": "string" },
                      "reason": { "type": "string" }
                    },
                    "required": ["type","location","constraint","reason"],
                    "additionalProperties": false
                  },
                  {
                    "properties": {
                      "type": { "const": "shared_limited" },
                      "location": { "type": "string" },
                      "capacity": { "type": "integer", "minimum": 1 },
                      "reason": { "type": "string" }
                    },
                    "required": ["type","location","capacity","reason"],
                    "additionalProperties": false
                  }
                ]
              }
            }
          },
          "additionalProperties": false
        }
      },
      "additionalProperties": false
    },
    "tasks": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["id","feature_id","title","boundaries","execution_guidance","evidence"],
        "properties": {
          "id": { "type": "string" },
          "feature_id": { "type": "string" },
          "title": { "type": "string" },
          "shared_resources": {
            "type": "object",
            "properties": {
              "exclusive": { "type": "array", "items": { "type": "string" } },
              "shared_limited": { "type": "array", "items": { "type": "string" } },
              "creates": { "type": "array", "items": { "type": "string" } },
              "modifies": { "type": "array", "items": { "type": "string" } }
            },
            "additionalProperties": false
          },
          "boundaries": {
            "type": "object",
            "required": ["expected_complexity","definition_of_done","scope"],
            "properties": {
              "expected_complexity": {
                "type": "object",
                "required": ["value"],
                "properties": {
                  "value": { "type": "string" },
                  "breakdown": { "type": "string" }
                },
                "additionalProperties": false
              },
              "definition_of_done": {
                "type": "object",
                "required": ["criteria"],
                "properties": {
                  "criteria": { "type": "array", "items": { "type": "string" }, "minItems": 1 },
                  "stop_when": { "type": "string" }
                },
                "additionalProperties": false
              },
              "scope": {
                "type": "object",
                "properties": {
                  "include": { "type": "array", "items": { "type": "string" } },
                  "exclude": { "type": "array", "items": { "type": "string" } }
                },
                "additionalProperties": false
              }
            },
            "additionalProperties": false
          },
          "execution_guidance": {
            "type": "object",
            "required": ["logging"],
            "properties": {
              "logging": {
                "type": "object",
                "required": ["format","fields","on_error"],
                "properties": {
                  "format": { "const": "jsonl" },
                  "fields": {
                    "type": "array",
                    "items": { "type": "string" },
                    "contains": { "enum": ["timestamp","task_id","step","status","message"] }
                  },
                  "on_error": { "type": "array", "items": { "type": "string", "enum": ["data.error_code","data.stack"] } }
                },
                "additionalProperties": false
              },
              "checkpoints": {
                "type": "array",
                "items": { "type": "string" }
              }
            },
            "additionalProperties": false
          },
          "evidence": { "$ref": "https://spec.tasks-slaps.io/schemas/features.json#/$defs/evidenceArray" }
        },
        "additionalProperties": false
      }
    },
    "dependencies": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["from","to","type","confidence","isHard"],
        "properties": {
          "from": { "type": "string" },
          "to": { "type": "string" },
          "type": {
            "type": "string",
            "enum": ["technical","sequential","interface","mutual_exclusion","resource_limited","knowledge"]
          },
          "shared_resource": { "type": "string" },
          "capacity": { "type": "integer", "minimum": 1 },
          "reason": { "type": "string" },
          "evidence": { "$ref": "https://spec.tasks-slaps.io/schemas/features.json#/$defs/evidenceArray" },
          "confidence": { "type": "number", "minimum": 0, "maximum": 1 },
          "isHard": { "type": "boolean" }
        },
        "additionalProperties": false
      }
    }
  },
  "additionalProperties": false
}

Resource types & examples reflect the shared resource model and mutual exclusion edges.  ￼

5.3 dag.json

{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://spec.tasks-slaps.io/schemas/dag.json",
  "type": "object",
  "required": ["meta","nodes","edges","metrics"],
  "properties": {
    "meta": {
      "type": "object",
      "required": ["artifact_hash"],
      "properties": { "artifact_hash": { "type": "string", "pattern": "^[a-f0-9]{64}$" } },
      "additionalProperties": false
    },
    "nodes": {
      "type": "array",
      "items": { "type": "object", "required": ["id"], "properties": { "id": { "type": "string" } }, "additionalProperties": true }
    },
    "edges": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["from","to","type"],
        "properties": {
          "from": { "type": "string" },
          "to": { "type": "string" },
          "type": { "type": "string", "enum": ["technical","sequential","interface","mutual_exclusion","resource_limited","knowledge"] }
        },
        "additionalProperties": false
      }
    },
    "metrics": {
      "type": "object",
      "required": ["edge_density","verb_first_ratio","evidence_coverage"],
      "properties": {
        "edge_density": { "type": "number", "minimum": 0 },
        "verb_first_ratio": { "type": "number", "minimum": 0, "maximum": 1 },
        "evidence_coverage": { "type": "number", "minimum": 0, "maximum": 1 },
        "critical_path": { "type": "array", "items": { "type": "string" } }
      },
      "additionalProperties": false
    }
  },
  "additionalProperties": false
}

edge_density warning band 0.01–0.5 is enforced by evaluator/policy (see §8).  ￼

5.4 waves.json

{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://spec.tasks-slaps.io/schemas/waves.json",
  "type": "object",
  "required": ["meta","waves"],
  "properties": {
    "meta": {
      "type": "object",
      "required": ["artifact_hash","max_wave_size"],
      "properties": {
        "artifact_hash": { "type": "string", "pattern": "^[a-f0-9]{64}$" },
        "max_wave_size": { "type": "integer", "minimum": 1 }
      },
      "additionalProperties": false
    },
    "waves": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["id","tasks"],
        "properties": {
          "id": { "type": "string" },
          "tasks": { "type": "array", "items": { "type": "string" } },
          "barriers": { "type": "array", "items": { "type": "string", "enum": ["gate:tests","gate:review","gate:deploy"] } }
        },
        "additionalProperties": false
      }
    }
  },
  "additionalProperties": false
}


⸻

6) Coordinator Contract (coordinator.json → consumed by SLAPS)

{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://spec.tasks-slaps.io/schemas/coordinator.json",
  "type": "object",
  "required": ["meta","mode","limits","policies","resources"],
  "properties": {
    "meta": {
      "type": "object",
      "required": ["artifact_hash","tasks_hash","dag_hash","waves_hash"],
      "properties": {
        "artifact_hash": { "type": "string", "pattern": "^[a-f0-9]{64}$" },
        "tasks_hash": { "type": "string", "pattern": "^[a-f0-9]{64}$" },
        "dag_hash": { "type": "string", "pattern": "^[a-f0-9]{64}$" },
        "waves_hash": { "type": "string", "pattern": "^[a-f0-9]{64}$" }
      },
      "additionalProperties": false
    },
    "mode": {
      "type": "string",
      "enum": ["rolling_frontier","layered_waves"],
      "default": "rolling_frontier"
    },
    "limits": {
      "type": "object",
      "properties": {
        "max_parallel": { "type": "integer", "minimum": 1 },
        "per_resource_caps": {
          "type": "object",
          "additionalProperties": { "type": "integer", "minimum": 1 }
        }
      },
      "additionalProperties": false
    },
    "resources": {
      "type": "object",
      "description": "Mirror of tasks.meta.codebase_analysis.shared_resources (names + capacities)"
    },
    "policies": {
      "type": "object",
      "required": ["retry","backoff","checkpoint","circuit_breakers","rollback","replan"],
      "properties": {
        "retry": {
          "type": "object",
          "properties": { "max_attempts": { "type": "integer", "minimum": 0 } },
          "additionalProperties": false
        },
        "backoff": {
          "type": "object",
          "properties": {
            "type": { "type": "string", "enum": ["exponential","linear","none"] },
            "base_ms": { "type": "integer", "minimum": 0 },
            "max_ms": { "type": "integer", "minimum": 0 }
          },
          "additionalProperties": false
        },
        "checkpoint": {
          "type": "object",
          "properties": {
            "enabled": { "type": "boolean", "default": true },
            "frequency": { "type": "string", "enum": ["task_start","task_end","custom"] }
          },
          "additionalProperties": false
        },
        "circuit_breakers": {
          "type": "object",
          "properties": {
            "repeated_errors": {
              "type": "object",
              "properties": {
                "threshold": { "type": "integer", "minimum": 1 },
                "window_s": { "type": "integer", "minimum": 1 },
                "action": { "type": "string", "enum": ["pause_task","rollback_and_replan"] }
              },
              "additionalProperties": false
            },
            "thrash_detect": {
              "type": "object",
              "properties": {
                "min_loops": { "type": "integer", "minimum": 1 },
                "action": { "type": "string", "enum": ["kill_and_escalate","rollback_and_replan"] }
              },
              "additionalProperties": false
            }
          },
          "additionalProperties": false
        },
        "rollback": {
          "type": "object",
          "properties": {
            "strategy": { "type": "string", "enum": ["per_task","wave_scope","frontier_scope"] },
            "max_depth": { "type": "integer", "minimum": 0 }
          },
          "additionalProperties": false
        },
        "replan": {
          "type": "object",
          "properties": {
            "enabled": { "type": "boolean", "default": true },
            "triggers": {
              "type": "array",
              "items": { "type": "string",
                "enum": ["dependency_violation","missing_resource_capacity","pattern_of_missing_deps","task_timeout","external_constraint"]
              }
            },
            "actions": {
              "type": "array",
              "items": { "type": "string",
                "enum": ["inject_edges","insert_tasks","batch_install","recompute_frontier","shrink_parallelism"]
              }
            }
          },
          "additionalProperties": false
        }
      },
      "additionalProperties": false
    }
  },
  "additionalProperties": false
}


⸻

7) SLAPS Execution Algorithm (Rolling-Frontier)

Mode: rolling_frontier (default). Tasks launch as soon as all structural deps are satisfied and resource locks can be acquired; no artificial wave syncs.

Pseudocode (normative-ish):

ready = { t | indegree(t)==0 }
while ready or running:
  for each t in ready sorted by priority:
    if acquire_all_resource_locks(t):
       launch(t)
       running.add(t)
       ready.remove(t)

  on task_event(e):
    log(JSONL e)
    if e.status == done:
       release_locks(e.t)
       mark_complete(e.t)
       newly_unblocked = { u | (forall p in preds(u): complete(p)) }
       ready |= filter(resource_feasible, newly_unblocked)

    if e.status == error:
       apply_retry_or_backoff(e.t)
       if breaker_tripped(e.t) or error_non_retryable(e.t):
          execute_rollback(policy, scope_for(e.t))
          planned_changes = consider_replan(triggers_from_logs)
          if planned_changes:
             mutate_frontier(planned_changes)   # inject edges/tasks, shrink parallelism, etc.
          else:
             escalate(e.t)

Key mechanics
	•	Locks: exclusive / shared-limited (capacity), aligned with tasks.json.meta.codebase_analysis.shared_resources.  ￼
	•	Admission control: per-resource capacities + global max_parallel.
	•	Frontier mutation (replan): allowed actions include edge injection, task insertion, batch install, parallelism shrink.
	•	Telemetry: all decisions are backed by log-derived signals (see §3).

⸻

8) Rollback & Replan: Decision Rules (Operational)

SLAPS MUST evaluate the following in order on task failure:
	1.	Retry if error is transient and attempts < policies.retry.max_attempts.
	2.	Backoff per policies.backoff.
	3.	Circuit break if repeated errors or thrash pattern exceed thresholds → pause/kill, then rollback.
	4.	Rollback scope (choose smallest that heals consistency):
	•	per_task: revert task-changes; keep frontier.
	•	frontier_scope: revert this task and any dependent tasks that started after its last good checkpoint.
	•	wave_scope: only if running in layered_waves.
	5.	Replan if triggers fire:
	•	pattern_of_missing_deps ⇒ run repo scan, batch-install missing deps, inject edges, restart.
	•	missing_resource_capacity ⇒ inject resource_limited edges, reduce parallelism.
	•	dependency_violation ⇒ inject interface edge or split task (planner-style fix).
	•	task_timeout ⇒ shrink parallelism or re-queue with higher priority.
	•	external_constraint ⇒ add mutual_exclusion/capacity as needed.

These behaviors operationalize the earlier “thrash shields”, “hot edge injection”, “deadlock immunity”, and “self-healing swarm orchestration”. (See also logging & evidence mandates.)  ￼

⸻

9) Evaluator & Quality Bars
	•	Hard fails: missing task boundaries; missing logging guidance; absent codebase analysis; missing Decisions.md.  ￼
	•	Scoring (120 pts) includes Task Boundaries and Codebase Awareness dimensions.  ￼
	•	Passing bands: Excellent 108–120; Good 96–107; Needs Work 84–95; Reject <84.  ￼

⸻

10) Implementation Notes (Incorporated Feedback)
	•	Standard JSONL logs (agents & evaluator).  ￼
	•	Polyglot ast-grep cookbook (Prisma/Drizzle/Knex/Sequelize/Rails + deploy/test/env).  ￼  ￼
	•	Infrastructure resource edges allowed; human-resource edges disallowed.  ￼
	•	Edge-density warning set to 0.01–0.5.  ￼
	•	Explicit Evidence schema for both tasks and edges.  ￼
	•	Deterministic hashing incl. placement of artifact_hash.  ￼

⸻

11) Minimal End-to-End: What “Done” Looks Like
	1.	T.A.S.K.S. consumes PLAN_DOC, scans the repo with ast-grep, outputs artifacts in the exact order above (with canonical JSON + hashes).  ￼  ￼
	2.	S.L.A.P.S. loads tasks.json, dag.json, waves.json, and coordinator.json and executes in rolling_frontier unless told otherwise, honoring resource locks and capacities.
	3.	Agents stream JSONL logs; evaluator parses them against the contract.  ￼
	4.	On failures, rollback → replan per §8, mutating the frontier and/or injecting constraints.
	5.	Finish when all nodes complete; write final provenance (artifact hashes, log digests, decision trail).

⸻

12) One-Line Summary (for the sticker)

TASKS writes the math. SLAPS makes it true.

⸻

￼