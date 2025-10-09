---
id: F019
title: Docs & Samples
depends: []
priority: P1
status: backlog
---

## User Story
As a new contributor I need curated docs and runnable samples so I can understand the planner/runtime contracts and reproduce an end-to-end plan quickly.

## Outcome
Author onboarding guides, reference docs, and sample plans that walk through canonicalization, planning, validation, and execution flows aligned with the normative spec.

## Scope & Boundaries
- refresh README, AGENTS.md, and tutorial docs with current CLI usage
- provide sample plans plus scripts to regenerate hashes and DOT exports
- align docs with validator + executor expectations called out in the spec
- exclude long-form blog content (out of scope)

## Acceptance Criteria
- sample plan artifacts regenerate cleanly using documented commands
- docs cite the normative spec and clarify contract boundaries for planner/executor
- doc tests or automated checks ensure snippets remain in sync with code

## Evidence & References
- docs/formal-spec.md §§1–3,18–19 (governing principles and scaffolding)

## Linked Tasks
- T200-docs-and-samples
