---
id: F001
title: Canonical JSON + Hashing
depends: []
---

Outcome: deterministic serializer (sorted keys, minimal numbers, LF, newline-terminated) and SHA‑256 artifact hashing with correct preimage policy.

Acceptance
- Serializer test suite passes map-key sorting and minimal-number cases
- Artifact writer computes meta.artifact_hash over canonical JSON with empty hash field; hashes listed in Plan.md

Tasks: T001, T002
