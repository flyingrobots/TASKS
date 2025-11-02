Version v9: Contract changes and migration

Summary
- Artifact field names use camelCase across all planner artifacts (tasks.json, dag.json, waves.json, coordinator.json, features.json).
- Meta.version is now "v9". The validator enforces the preimage hash policy by blanking meta.artifactHash in-memory when recomputing digests.

What changed
- Former snake_case names (e.g., acceptance_checks, feature_id, artifact_hash) are now acceptanceChecks, featureID, artifactHash.
- No functional changes to the data shape beyond field names.

Backwards compatibility
- The validator accepts artifacts that contain either camelCase artifactHash or snake_case artifact_hash for one release window while repositories migrate.
- For hashing, the validator computes the digest over canonical bytes with the hash field set to an empty string, regardless of the field name.

How to migrate
1) Re-run `tasksd plan` to regenerate artifacts with v9.
2) If you maintain custom producers/consumers, update JSON field names to camelCase.
3) Update any JSON Path, jq, or integration queries to the new names.

Notes
- Arrays remain in computed order; object keys are lexicographically sorted; numbers use minimal decimal rendering; UTF-8; LF line endings; newline-terminated.

