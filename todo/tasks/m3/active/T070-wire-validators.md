---
id: T070
milestone: m3
features:
  - F009
title: 'Wire acceptance, evidence, and interface validators into planner flow'
status: active
deps:
  - T030
  - T040
---

## User Story
As a compliance reviewer I need validator clients wired into `tasksd plan` so plans fail fast when acceptance, evidence, or interface requirements are missing.

## Summary
Integrate validator subprocesses or libraries, aggregate their diagnostics, and halt planning on gating failures while caching results for deterministic reruns.

## Scope
### In Scope
- add validator invocation hooks during plan generation
- capture diagnostics per task/edge and attach to Plan.md + tasks.json analysis blocks
- cache successful validator runs keyed by artifact hashes
### Out of Scope
- redaction logic (T170)
- executor runtime validation (future milestones)

## Execution Guidance
- respect validator exit codes and emit actionable remediation hints
- support `--offline` caching for repeated runs with unchanged inputs
- record validator fingerprints and timestamps for audit

## Acceptance (machine-verifiable)
```acceptance
[
  {"type":"command","cmd":"cd planner && go test ./cmd/tasksd","timeoutSeconds":240}
]
```

## Evidence & References
- docs/formal-spec.md §§6,10–11 (validator orchestration and gating)
- Feature F009 Validator Integration
