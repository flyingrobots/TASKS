package validate

import (
    "bytes"
    "embed"
    "encoding/json"
    "fmt"

    m "github.com/james/tasks-planner/internal/model"
    jsonschema "github.com/santhosh-tekuri/jsonschema/v5"
    "github.com/james/tasks-planner/internal/canonjson"
    "github.com/james/tasks-planner/internal/hash"
)

//go:embed schemas/*.json
var schemaFS embed.FS

var (
    compiled = map[string]*jsonschema.Schema{}
)

func compileSchemas() error {
    if len(compiled) > 0 { return nil }
    c := jsonschema.NewCompiler()
    // load tasks
    tb, err := schemaFS.ReadFile("schemas/tasks.schema.json")
    if err != nil { return err }
    if err := c.AddResource("mem://tasks.schema.json", bytes.NewReader(tb)); err != nil { return err }
    ts, err := c.Compile("mem://tasks.schema.json")
    if err != nil { return err }
    compiled["tasks.json"] = ts
    // load dag
    db, err := schemaFS.ReadFile("schemas/dag.schema.json")
    if err != nil { return err }
    if err := c.AddResource("mem://dag.schema.json", bytes.NewReader(db)); err != nil { return err }
    ds, err := c.Compile("mem://dag.schema.json")
    if err != nil { return err }
    compiled["dag.json"] = ds

    // features
    fb, err := schemaFS.ReadFile("schemas/features.schema.json")
    if err != nil { return err }
    if err := c.AddResource("mem://features.schema.json", bytes.NewReader(fb)); err != nil { return err }
    fs, err := c.Compile("mem://features.schema.json")
    if err != nil { return err }
    compiled["features.json"] = fs

    // waves
    wb, err := schemaFS.ReadFile("schemas/waves.schema.json")
    if err != nil { return err }
    if err := c.AddResource("mem://waves.schema.json", bytes.NewReader(wb)); err != nil { return err }
    ws, err := c.Compile("mem://waves.schema.json")
    if err != nil { return err }
    compiled["waves.json"] = ws

    // coordinator
    cb, err := schemaFS.ReadFile("schemas/coordinator.schema.json")
    if err != nil { return err }
    if err := c.AddResource("mem://coordinator.schema.json", bytes.NewReader(cb)); err != nil { return err }
    cs, err := c.Compile("mem://coordinator.schema.json")
    if err != nil { return err }
    compiled["coordinator.json"] = cs
    return nil
}

// ValidateRaw validates raw JSON bytes against a named schema key.
func ValidateRaw(schemaKey string, raw []byte) error {
    if err := compileSchemas(); err != nil { return err }
    var inst any
    if err := json.Unmarshal(raw, &inst); err != nil { return err }
    s, ok := compiled[schemaKey]
    if !ok { return fmt.Errorf("no schema for %s", schemaKey) }
    if err := s.Validate(inst); err != nil { return err }
    return nil
}

// CheckArtifactHash recomputes the artifact hash over canonical JSON with meta.artifact_hash blank.
// Returns (computedHash, storedHash, ok, error)
func CheckArtifactHash(raw []byte) (string, string, bool, error) {
    var v any
    if err := json.Unmarshal(raw, &v); err != nil { return "", "", false, err }
    mobj, ok := v.(map[string]any)
    if !ok { return "", "", false, fmt.Errorf("root not object") }
    meta, ok := mobj["meta"].(map[string]any)
    if !ok { return "", "", true, nil } // no meta; nothing to check
    stored, _ := meta["artifact_hash"].(string)
    meta["artifact_hash"] = ""
    raw2, err := json.Marshal(mobj)
    if err != nil { return "", stored, false, err }
    can, err := canonjson.ToCanonicalJSON(raw2)
    if err != nil { return "", stored, false, err }
    computed := hash.HashCanonicalBytes(can)
    return computed, stored, computed == stored, nil
}

func TasksFile(tf *m.TasksFile) error {
    if tf.Meta.Version == "" { return fmt.Errorf("meta.version required") }
    if len(tf.Tasks) == 0 { return fmt.Errorf("at least one task required") }
    for i, t := range tf.Tasks {
        if t.ID == "" { return fmt.Errorf("task[%d].id required", i) }
        if t.Title == "" { return fmt.Errorf("task[%d].title required", i) }
        if len(t.AcceptanceChecks) == 0 { return fmt.Errorf("task[%s] acceptance_checks required", t.ID) }
    }
    // JSON Schema validation
    if err := compileSchemas(); err != nil { return err }
    b, _ := json.Marshal(tf)
    var inst any
    _ = json.Unmarshal(b, &inst)
    if err := compiled["tasks.json"].Validate(inst); err != nil { return fmt.Errorf("tasks.json schema: %w", err) }
    return nil
}

func DagFile(df *m.DagFile) error {
    if df.Meta.Version == "" { return fmt.Errorf("meta.version required") }
    if len(df.Nodes) == 0 { return fmt.Errorf("nodes required") }
    // JSON Schema validation
    if err := compileSchemas(); err != nil { return err }
    b, _ := json.Marshal(df)
    var inst any
    _ = json.Unmarshal(b, &inst)
    if err := compiled["dag.json"].Validate(inst); err != nil { return fmt.Errorf("dag.json schema: %w", err) }
    return nil
}
