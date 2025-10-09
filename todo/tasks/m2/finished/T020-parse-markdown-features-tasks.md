---
id: T020
milestone: m2
features: [F004]
title: Parse markdown charter into features and raw tasks
status: finished
deps: [T010]
---

## User Story
As a planning analyst I need the doc parser to extract prioritized features and provisional tasks from the charter so later stages can normalize them.

## Summary
Extend `docparse` to detect feature headings, priority hints, and task bullets (including durations/after clauses) and emit structured `Feature` and `TaskSpec` objects with evidence pointers.

## Scope
### In Scope
- handle fenced `acceptance` blocks and convert JSON/YAML payloads to structs
- parse `after:` dependency text and optimistic duration hints
- capture parser errors per task for downstream surfacing
### Out of Scope
- final task normalization (T030)
- dependency inference (T040)

## Execution Guidance
- add streaming scanner support for large documents with configurable buffer sizes
- produce detailed parser diagnostics for missing fences or malformed metadata
- keep parsing deterministic and order-preserving

## Acceptance (machine-verifiable)
```acceptance
[
  {"type":"command","cmd":"cd planner && go test ./internal/planner/docparse","timeoutSeconds":180}
]
```

## Evidence & References
- docs/formal-spec.md §5 (feature extraction and doc parsing guidance)
- Feature F004 Feature Extraction
