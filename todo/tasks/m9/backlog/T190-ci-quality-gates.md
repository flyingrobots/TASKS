---
id: T190
milestone: m9
features: [F018]
title: Enforce CI quality gates for planner artifacts and hashes
status: backlog
deps: [T001, T002, T003, T050, T070]
---

## User Story
As a release captain I need CI pipelines that fail on schema/hash drift or validator regressions so broken plans never merge to main.

## Summary
Set up CI workflows that run canonicalization, hashing, schema validation, DAG checks, and validator suites; publish status badges and accelerate reproducibility jobs.

## Scope
### In Scope
- GitHub Actions (or equivalent) workflow definitions for build/test/validate
- nightly reproducibility job that reruns `tasksd plan` and compares hashes
- status badges + documentation updates showing gate coverage
### Out of Scope
- release packaging (T200)
- executor deployment automation (future work)

## Execution Guidance
- leverage caching where deterministic; rerun canonicalization before hashing comparisons
- fail fast with actionable log messages linking to docs
- store CI artifacts (canonical JSON, DOTs) for post-failure inspection

## Acceptance (machine-verifiable)
```acceptance
[
  {"type":"command","cmd":"npm run lint && npm test","timeoutSeconds":600},
  {"type":"command","cmd":"./.github/workflows/ci-local.sh","timeoutSeconds":900}
]
```

## Evidence & References
- docs/formal-spec.md §§3,12,17 (determinism and evaluation gates)
- Feature F018 CI + Quality Gates
