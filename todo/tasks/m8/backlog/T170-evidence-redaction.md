---
id: T170
milestone: m8
features: [F016]
title: Implement evidence redaction + fail-closed policy
status: backlog
deps: [T070]
---

Redact secrets in evidence excerpts (JWTs, keys, PEMs); validators accept redacted spans; block on leaks.
