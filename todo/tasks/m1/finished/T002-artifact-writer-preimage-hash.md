---
id: T002
milestone: m1
features: [F001]
title: Enforce preimage hashing in artifact writer
status: finished
deps: [T001]
---

## User Story
As a release engineer I need the artifact writer to hash canonical bytes with the preimage rule so artifacts embed trustworthy `meta.artifact_hash` values.

## Summary
Update `emitter.WriteWithArtifactHash` to canonicalize JSON, compute SHA-256 with `meta.artifact_hash` cleared, set the hash via callback, and write newline-terminated canonical bytes with panic-safe handling.

## Scope
### In Scope
- adjust `planner/internal/emitter/writer.go`
- add tests for nil callbacks, panic recovery, and hash mismatches
- document hash ordering in AGENTS.md and code comments
### Out of Scope
- CLI surface output formatting (handled by T180)
- schema validation (handled by T003)

## Execution Guidance
- compute canonical bytes twice: once for hashing and once post-callback write
- wrap `setHash` invocations with panic recovery to return structured errors
- ensure file writes use `0o644` permissions and canonical newline termination

## Acceptance (machine-verifiable)
```acceptance
[
  {"type":"command","cmd":"cd planner && go test ./internal/emitter","timeoutSeconds":120},
  {"type":"command","cmd":"cd planner && go test ./...","timeoutSeconds":300}
]
```

## Evidence & References
- docs/formal-spec.md §§3,12 (artifact hashing workflow)
- Feature F001 Canonical JSON + Hashing
