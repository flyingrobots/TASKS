---
id: F001
title: Canonical JSON + Hashing
depends: []
priority: P0
status: shaping
---

## User Story
As a planner engineer I need canonical serialization and hashing so that every artifact we emit is deterministic and auditable across runs, letting S.L.A.P.S. trust the contract without re-validating raw specs.

## Outcome
We provide a reusable Go package that normalizes JSON (sorted keys, minimal decimal rendering, newline termination) and a companion hashing helper that computes `meta.artifact_hash` over the canonical bytes with the preimage convention described in the spec.

## Scope & Boundaries
- normalize objects, arrays, and numeric encodings for planner artifacts
- emit SHA-256 hashes that incorporate the canonical newline
- expose a CLI entry point for canonicalization and hashing used by CI
- exclude higher-level planner logic (feature parsing, DAG build)

## Acceptance Criteria
- canonicalizer test suite covers key-sorting, number minimization, newline termination
- hash helper validates preimage policy for all artifact types
- CLI command `tasksd canonical <file>` outputs canonical bytes and hash in deterministic format
- documentation explains how downstream writers integrate the helpers

## Evidence & References
- docs/formal-spec.md §§3,12 (deterministic artifacts and hashing contract)

## Linked Tasks
- T001-canonical-minimal-numbers
- T002-artifact-writer-preimage-hash
