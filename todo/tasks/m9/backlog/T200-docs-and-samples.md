---
id: T200
milestone: m9
features: [F019]
title: Refresh docs and runnable samples for planner + executor
status: backlog
deps: [T001, T002, T020, T050, T080]
---

## User Story
As a new contributor I need complete docs and runnable samples so I can reproduce plan generation, validation, and export flows without digging through source.

## Summary
Update README, AGENTS.md, and docs under `docs/` with the latest CLI guidance, add runnable sample scripts that regenerate artifacts, and ensure references cite the normative spec.

## Scope
### In Scope
- rewrite onboarding and workflow docs with task/feature templates
- provide sample script(s) to regenerate `plans/sample` and DOT exports
- add troubleshooting and FAQ sections for validator failures
### Out of Scope
- blog/tutorial content outside repo
- external site deployment

## Execution Guidance
- keep docs aligned with templates defined in `docs/TEMPLATES.md`
- cite docs/formal-spec.md in every artifact-facing guide
- add doc tests or CI checks to ensure commands stay current

## Acceptance (machine-verifiable)
```acceptance
[
  {"type":"command","cmd":"npm run docs:lint","timeoutSeconds":420},
  {"type":"command","cmd":"cd planner && go run ./cmd/tasksd plan --out ./plans/sample && cd planner && go run ./cmd/tasksd validate --dir ./plans/sample","timeoutSeconds":480}
]
```

## Evidence & References
- docs/formal-spec.md §§1–3,18–19 (scope, artifacts, scaffolding)
- Feature F019 Docs & Samples
