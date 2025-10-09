---
id: F004
title: Feature Extraction
depends: [F003]
priority: P0
status: backlog
---

## User Story
As a planner I need to translate a markdown charter into a prioritized capability list so future decomposition is bounded, auditable, and traceable to the source narrative.

## Outcome
Parse the primary spec, identify 5–25 outcome-oriented capabilities, attach priorities and evidence excerpts, and serialize them as canonical `features.json` records.

## Scope & Boundaries
- extend `docparse` to recognize feature headings, priorities, and evidence anchors
- normalize verb-first feature titles and add rationale text blocks
- annotate each feature with references back to the source document and census signals
- exclude downstream task decomposition (covered by F005)

## Acceptance Criteria
- generated `features.json` contains ≥5 features with priorities and evidence anchors
- missing evidence or malformed sections cause actionable validation errors
- documentation clarifies authoring guidance for future feature specs

## Evidence & References
- docs/formal-spec.md §5 (feature extraction guidance)

## Linked Tasks
- T020-parse-markdown-features-tasks
