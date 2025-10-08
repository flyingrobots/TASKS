package emitter

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/james/tasks-planner/internal/canonjson"
    "github.com/james/tasks-planner/internal/hash"
)

// WriteWithArtifactHash writes v to path as canonical JSON, computing meta.artifact_hash over the
// canonical bytes with the hash field present but empty. The setHash callback should embed the computed
// hash value into v (e.g., set v.Meta.ArtifactHash = h) before the final write.
func WriteWithArtifactHash(path string, v any, setHash func(h string)) (string, error) {
    raw, err := json.Marshal(v)
    if err != nil { return "", fmt.Errorf("marshal %s: %w", path, err) }
    can, err := canonjson.ToCanonicalJSON(raw)
    if err != nil { return "", fmt.Errorf("canon %s: %w", path, err) }
    h := hash.HashCanonicalBytes(can)
    if setHash != nil { setHash(h) }
    raw2, err := json.Marshal(v)
    if err != nil { return "", fmt.Errorf("marshal2 %s: %w", path, err) }
    can2, err := canonjson.ToCanonicalJSON(raw2)
    if err != nil { return "", fmt.Errorf("canon2 %s: %w", path, err) }
    if err := os.WriteFile(path, can2, 0o644); err != nil { return "", fmt.Errorf("write %s: %w", path, err) }
    return h, nil
}

