---
id: T180
milestone: m8
features: [F017]
title: Export deterministic DOT views for DAG and runtime
status: finished
deps: [T050, T080]
---

## User Story
As a reviewer I need faithful DOT exports for both the precedence DAG and runtime graph so I can inspect plans visually during reviews.

## Summary
Refactor `tasksd export-dot` into reusable helpers that load coordinator or dag/tasks pairs, render deterministic DOT output with configurable labels, and document usage examples.

## Scope
### In Scope
- share dot-generation helpers across directory, single DAG, and coordinator modes
- support node label (`id`, `title`, `id-title`) and edge label toggles
- emit DOT outputs into plan directories alongside canonical artifacts
### Out of Scope
- standalone web viewer (moved to external repo)
- runtime scheduler instrumentation (other milestones)

## Execution Guidance
- canonicalize JSON inputs before DOT rendering for consistent nodes
- ensure CLI exits with clear errors when required artifacts are missing
- update docs and examples illustrating DOT → SVG workflows

## Acceptance (machine-verifiable)
```acceptance
[
  {"type":"command","cmd":"cd planner && go test ./cmd/tasksd","timeoutSeconds":300},
  {"type":"command","cmd":"cd examples/dot-export && ./render.sh","timeoutSeconds":240}
]
```

## Evidence & References
- docs/formal-spec.md §§7–8,17 (visualization and evaluation guidance)
- Feature F017 Visualization & Export
