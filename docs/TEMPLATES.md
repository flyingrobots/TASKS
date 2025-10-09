# Templates: Task, Feature, Milestone

Use these as starting points for new planning artifacts.

## Task Template (todo/tasks/mX/backlog/T###-your-title.md)

```yaml
---
id: T###
milestone: mX
features: [F###]
title: Verb-first, bounded title (2–8h)
status: backlog
deps: []
---
```

### Summary
- What is this task delivering?

### Acceptance (machine-verifiable)

```json
[
  {"type":"command","cmd":"go test ./...","timeoutSeconds":600}
]
```

### Test Plan
- Unit: list tests
- Integration: list tests
- CLI/UX: list checks

### Notes / Design
- How you plan to implement

### Links
- Feature: ../features/F###-your-feature.md

## Feature Template (todo/features/F###-your-feature.md)

```yaml
---
id: F###
title: Short feature name
depends: []
---
```

### Outcome
- What users can do when this ships

### Requirements
- Bulleted requirements, constraints

### Design
- High-level design and data contracts involved

### Acceptance
- Criteria that must be true to call this feature done

### Definition of Done
- Tests, docs, and release notes updated; examples run

### Test Plans
- Unit, integration, end-to-end

### Related Tasks
- T###, T###

## Milestone Template (todo/milestones/MX-your-milestone.md)

```markdown
# Milestone MX – Title

Scope: Intro to the scope.

Features: F###, F###

<!-- PROGRESS:START MX -->
(Progress will be updated by scripts)
<!-- PROGRESS:END MX -->
```

### Success Criteria / Demos
- “By now, you can …”

### Links
- Features: [F###](../features/F###-your-feature.md)
- Tasks: see `../tasks/mx/*`
