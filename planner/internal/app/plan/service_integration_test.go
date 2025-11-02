package plan

import (
    "context"
    "encoding/json"
    "os"
    "path/filepath"
    "testing"

    "github.com/james/tasks-planner/internal/export/dot"
    m "github.com/james/tasks-planner/internal/model"
    "github.com/james/tasks-planner/internal/validate"
)

// TestFullPlannerPipeline runs: Plan (stub) -> write artifacts -> schema+hash validate -> DOT export.
func TestFullPlannerPipeline(t *testing.T) {
    t.Parallel()
    out := t.TempDir()
    svc := NewDefaultService()
    res, err := svc.Plan(context.Background(), Request{DocPath: "", RepoPath: ".", OutDir: out})
    if err != nil {
        t.Fatalf("plan: %v", err)
    }
    if len(res.ArtifactHashes) == 0 {
        t.Fatalf("expected artifact hashes")
    }
    // Load and validate each artifact by schema + preimage hash policy
    must := func(name string) []byte {
        b, err := os.ReadFile(filepath.Join(out, name))
        if err != nil { t.Fatalf("read %s: %v", name, err) }
        return b
    }
    // tasks.json
    tb := must("tasks.json")
    if err := validate.ValidateRaw("tasks.json", tb); err != nil {
        t.Fatalf("tasks.json schema: %v", err)
    }
    if _, _, ok, err := validate.CheckArtifactHash(tb); err != nil || !ok {
        t.Fatalf("tasks.json hash mismatch: ok=%v err=%v", ok, err)
    }
    // dag.json
    db := must("dag.json")
    if err := validate.ValidateRaw("dag.json", db); err != nil {
        t.Fatalf("dag.json schema: %v", err)
    }
    if _, _, ok, err := validate.CheckArtifactHash(db); err != nil || !ok {
        t.Fatalf("dag.json hash mismatch: ok=%v err=%v", ok, err)
    }
    // coordinator.json
    cb := must("coordinator.json")
    if err := validate.ValidateRaw("coordinator.json", cb); err != nil {
        t.Fatalf("coordinator.json schema: %v", err)
    }
    if _, _, ok, err := validate.CheckArtifactHash(cb); err != nil || !ok {
        t.Fatalf("coordinator.json hash mismatch: ok=%v err=%v", ok, err)
    }
    // features.json
    fb := must("features.json")
    if err := validate.ValidateRaw("features.json", fb); err != nil {
        t.Fatalf("features.json schema: %v", err)
    }
    if _, _, ok, err := validate.CheckArtifactHash(fb); err != nil || !ok {
        t.Fatalf("features.json hash mismatch: ok=%v err=%v", ok, err)
    }
    // waves.json
    wb := must("waves.json")
    if err := validate.ValidateRaw("waves.json", wb); err != nil {
        t.Fatalf("waves.json schema: %v", err)
    }
    if _, _, ok, err := validate.CheckArtifactHash(wb); err != nil || !ok {
        t.Fatalf("waves.json hash mismatch: ok=%v err=%v", ok, err)
    }

    // DOT export (strings non-empty) by decoding and calling exporters directly
    // Coordinator
    var coord m.Coordinator
    if err := json.Unmarshal(cb, &coord); err != nil { t.Fatalf("decode coordinator: %v", err) }
    coordStr := dot.FromCoordinatorWithOptions(coord, dot.Options{NodeLabel: "id-title", EdgeLabel: "type"})
    if coordStr == "" { t.Fatalf("empty coordinator DOT") }
    // DAG
    var df m.DagFile
    if err := json.Unmarshal(db, &df); err != nil { t.Fatalf("decode dag: %v", err) }
    var tf m.TasksFile
    if err := json.Unmarshal(tb, &tf); err != nil { t.Fatalf("decode tasks: %v", err) }
    titles := map[string]string{}
    for _, tsk := range tf.Tasks { titles[tsk.ID] = tsk.Title }
    dagStr := dot.FromDagWithOptions(df, titles, dot.Options{NodeLabel: "id-title", EdgeLabel: "type"})
    if dagStr == "" { t.Fatalf("empty dag DOT") }
}
