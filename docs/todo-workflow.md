# Markdown-First Todo Workflow

See `todo/README.md` for the directory structure, scripts, and everyday usage. This document focuses on helper commands, formatting standards, and Hubless integration.

## Scripts (helpers)

### Helper Commands

- Create a task branch off origin/main:
  - `npm run todo:branch -- <feature-slug> <TASK_ID>`
  - Example: `npm run todo:branch -- validators T070`
- Open a PR for the current branch (via GitHub CLI):
  - `npm run todo:pr:create`
  - Fills title/body from commits; base = main; head = current branch

The todo CLI moves the task markdown file between status folders, updates its frontmatter `status`, and refreshes progress sections inside `todo/README.md` and milestone docs.

## Frontmatter

Each task file has a minimal YAML-like header:

```yaml
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
```bash
npm run todo:task:set-active -- T001
```
2) Finish T001:
```bash
npm run todo:task:set-finished -- T001
```
3) Merge T001 (after PR merged):
```bash
npm run todo:task:set-merged -- T001
```
4) Recompute progress bars:
```bash
npm run todo:update
```

5) Open a PR for your branch:
```bash
npm run todo:pr:create
```

## Later: Hubless Integration

This layout mirrors Hubless concepts. When Hubless CLI is ready, we can:
- Export markdown tasks to `@hubless/**` JSON
- Drive the same state transitions through a Go CLI/TUI
