package emitter

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/james/tasks-planner/internal/canonjson"
    "github.com/james/tasks-planner/internal/hash"
)

// WriteWithArtifactHash writes v to path as canonical JSON.
// Hash policy: we compute meta.artifact_hash over the canonical bytes with the hash field present
// but set to an empty string (the "preimage"), then embed that hex into v via setHash, and finally
// write the canonical JSON to disk. The stored hash therefore reflects the preimage, not the final
// file bytes that include the hash, by design (see AGENTS.md: Hash preimage).
//
// Concurrency: v must not be mutated concurrently while this function executes. The setHash callback
// is invoked synchronously; callers must not share v across goroutines during this call.
func WriteWithArtifactHash(path string, v any, setHash func(h string)) (string, error) {
    // Marshal + canonicalize preimage (with empty hash field pre-set by caller)
    raw, err := json.Marshal(v)
    if err != nil { return "", fmt.Errorf("marshal %s: %w", path, err) }
    can, err := canonjson.ToCanonicalJSON(raw)
    if err != nil { return "", fmt.Errorf("canon %s: %w", path, err) }
    h := hash.HashCanonicalBytes(can)
    if setHash == nil {
        // No mutation; write the preimage as-is
        if err := os.WriteFile(path, can, 0o644); err != nil { return "", fmt.Errorf("write %s: %w", path, err) }
        return h, nil
    }
    // Set hash with panic safety, then re-marshal and canonicalize for final write
    var panicErr error
    func() {
        defer func() { if r := recover(); r != nil { panicErr = fmt.Errorf("setHash panicked: %v", r) } }()
        setHash(h)
    }()
    if panicErr != nil { return "", panicErr }
    raw2, err := json.Marshal(v)
    if err != nil { return "", fmt.Errorf("marshal after setHash %s: %w", path, err) }
    can2, err := canonjson.ToCanonicalJSON(raw2)
    if err != nil { return "", fmt.Errorf("canon after setHash %s: %w", path, err) }
    if err := os.WriteFile(path, can2, 0o644); err != nil { return "", fmt.Errorf("write %s: %w", path, err) }
    return h, nil
}
