# Markdown-First Todo Workflow

This repo embeds a Git-native planning system under `todo/`.

- Milestones: `todo/milestones/` (specs + progress bars)
- Features: `todo/features/` (requirements, acceptance, links)
- Tasks: `todo/tasks/{milestone}/{backlog|active|finished|merged}/T###-*.md`

## Scripts

- Set active: `npm run todo:task:set-active -- T001`
- Set finished: `npm run todo:task:set-finished -- T001`
- Set merged: `npm run todo:task:set-merged -- T001`
- Update progress: `npm run todo:update`

These move the task markdown file between status folders, update its frontmatter `status`, and refresh progress sections inside `todo/README.md` and milestone docs.

## Frontmatter

Each task file has a minimal YAML-like header:

```
---
id: T001
milestone: m1
features: [F001]
title: Implement minimal-number canonical JSON
status: backlog
deps: [T002]
---
```

## Example Session

1) Start T001:
```
npm run todo:task:set-active -- T001
```
2) Finish T001:
```
npm run todo:task:set-finished -- T001
```
3) Merge T001 (after PR merged):
```
npm run todo:task:set-merged -- T001
```
4) Recompute progress bars:
```
npm run todo:update
```

## Later: Hubless Integration

This layout mirrors Hubless concepts. When Hubless CLI is ready, we can:
- Export markdown tasks to `@hubless/**` JSON
- Drive the same state transitions through a Go CLI/TUI
