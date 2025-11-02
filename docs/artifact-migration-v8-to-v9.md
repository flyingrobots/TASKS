Artifact Migration: v8 → v9

Summary
- v9 standardizes field names to camelCase across all artifacts and clarifies the preimage hash policy (artifactHash present with computed digest; validator blanks the field in-memory when recomputing).
- meta.version now emits "v9". JSON Schemas validate shape, not the version literal, to avoid pinning repos.

Field rename map (representative)

| v8 (snake_case)     | v9 (camelCase)      | Applies to                  |
|---------------------|---------------------|-----------------------------|
| artifact_hash       | artifactHash        | meta in all artifacts       |
| feature_id          | featureID           | tasks.json, coordinator     |
| acceptance_checks   | acceptanceChecks    | tasks.json, coordinator     |
| duration_units      | durationUnits       | tasks.json, coordinator     |
| min_confidence      | minConfidence       | tasks.json meta             |
| tasks_hash          | tasksHash           | dag.json meta               |
| execution_logging   | executionLogging    | tasks.json, coordinator     |
| required_fields     | requiredFields      | executionLogging            |

Notes:
- Arrays are unchanged; object keys sort lexicographically; numbers use minimal decimal formatting; UTF‑8 + LF; newline‑terminated.

Migration options

1) Regenerate artifacts (recommended)
   - Run: `cd planner && go run ./cmd/tasksd plan --out ./plans --repo ..`
   - Validate: `cd planner && go run ./cmd/tasksd validate --dir ./plans`

2) jq transform (best‑effort for simple cases)
   Example for tasks.json (rename meta.artifact_hash → artifactHash and nested keys):
   ```sh
   jq '
     .meta |= (.artifactHash = (.artifact_hash // .artifactHash) | del(.artifact_hash))
     | .meta |= (.minConfidence = (.min_confidence // .minConfidence) | del(.min_confidence))
     | .tasks[] |= (.featureID = (.feature_id // .featureID) | del(.feature_id))
     | .tasks[] |= (.acceptanceChecks = (.acceptance_checks // .acceptanceChecks) | del(.acceptance_checks))
     | .tasks[] |= (.durationUnits = (.duration_units // .durationUnits) | del(.duration_units))
   ' tasks.json > tasks.v9.json
   ```
   Validate with `tasksd validate` and keep a backup of the original.

Deprecation / compatibility
- Validator accepts both artifactHash and artifact_hash during the transition window, but emits a deprecation warning in logs for artifact_hash.
- v8 artifacts are not rejected by schema on version alone, but downstream tools may assume v9 naming. Prefer regeneration.

