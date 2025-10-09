package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"

    "github.com/james/tasks-planner/internal/canonjson"
    analysis "github.com/james/tasks-planner/internal/analysis"
    "github.com/james/tasks-planner/internal/export/dot"
    "github.com/james/tasks-planner/internal/hash"
    m "github.com/james/tasks-planner/internal/model"
    docp "github.com/james/tasks-planner/internal/planner/docparse"
    dagbuild "github.com/james/tasks-planner/internal/planner/dag"
    wavesim "github.com/james/tasks-planner/internal/planner/wavesim"
    "github.com/james/tasks-planner/internal/validate"
    "github.com/james/tasks-planner/internal/emitter"
)

func usage() {
    fmt.Fprintf(os.Stderr, "tasksd commands:\n")
    fmt.Fprintf(os.Stderr, "  canonical <file.json>        Canonicalize JSON and print SHA-256.\n")
    fmt.Fprintf(os.Stderr, "  export-dot --dir DIR                  Emit dag.dot/runtime.dot from artifacts in DIR.\n")
    fmt.Fprintf(os.Stderr, "  export-dot --dag D --tasks T [--out O] Emit DOT from dag.json + tasks.json.\n")
    fmt.Fprintf(os.Stderr, "  export-dot --coordinator C [--out O]  Emit DOT from coordinator.json.\n")
    fmt.Fprintf(os.Stderr, "  plan [--doc FILE] [--repo DIR] [--out DIR]  Create stub artifacts and DOTs.\n")
    fmt.Fprintf(os.Stderr, "  validate --dir DIR                    Validate artifacts (hashes + schemas).\n")
}

func main() {
    if len(os.Args) < 2 {
        usage()
        os.Exit(1)
    }

    switch os.Args[1] {
    case "canonical":
        runCanonical()
    case "export-dot", "dot":
        runExportDot()
    case "plan":
        runPlan()
    case "validate":
        runValidate()
    default:
        // Back-compat: if a single path is provided, treat it as canonical
        if len(os.Args) == 2 {
            // shift to canonical path mode
            os.Args = []string{os.Args[0], "canonical", os.Args[1]}
            runCanonical()
            return
        }
        usage()
        os.Exit(1)
    }
}

func runCanonical() {
    if len(os.Args) < 3 {
        fmt.Fprintln(os.Stderr, "Usage: tasksd canonical <file.json>")
        os.Exit(1)
    }
    filePath := os.Args[2]
    file, err := os.Open(filePath)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: Failed to open file '%s': %v\n", filePath, err)
        os.Exit(1)
    }
    defer file.Close()

    inputBytes, err := io.ReadAll(file)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: Failed to read file '%s': %v\n", filePath, err)
        os.Exit(1)
    }
    canonicalBytes, err := canonjson.ToCanonicalJSON(inputBytes)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: Failed to process JSON: %v\n", err)
        os.Exit(1)
    }
    hashString := hash.HashCanonicalBytes(canonicalBytes)
    fmt.Println("--- Canonical JSON ---")
    fmt.Print(string(canonicalBytes))
    fmt.Println("--- SHA-256 Hash ---")
    fmt.Println(hashString)
}

func writeCoordinatorDot(coordPath, outPath, nodeLabel, edgeLabel string) (string, error) {
    var c m.Coordinator
    if err := loadJSON(coordPath, &c); err != nil { return "", err }
    dotStr := dot.FromCoordinatorWithOptions(c, dot.Options{NodeLabel: nodeLabel, EdgeLabel: edgeLabel})
    if outPath == "" { return dotStr, nil }
    return "", os.WriteFile(outPath, []byte(dotStr), 0o644)
}

func writeDagDot(dagPath, tasksPath, outPath, nodeLabel, edgeLabel string) (string, error) {
    var tf m.TasksFile
    if err := loadJSON(tasksPath, &tf); err != nil { return "", err }
    titles := make(map[string]string)
    for _, t := range tf.Tasks { titles[t.ID] = t.Title }
    var df m.DagFile
    if err := loadJSON(dagPath, &df); err != nil { return "", err }
    dotStr := dot.FromDagWithOptions(df, titles, dot.Options{NodeLabel: nodeLabel, EdgeLabel: edgeLabel})
    if outPath == "" { return dotStr, nil }
    return "", os.WriteFile(outPath, []byte(dotStr), 0o644)
}

func emitDirArtifacts(dir, nodeLabel, edgeLabel string) error {
    // coordinator
    coord := join(dir, "coordinator.json")
    if exists(coord) {
        if _, err := writeCoordinatorDot(coord, join(dir, "runtime.dot"), nodeLabel, edgeLabel); err != nil { return err }
        fmt.Println("wrote", join(dir, "runtime.dot"))
    }
    // dag + tasks
    dag := join(dir, "dag.json")
    tasks := join(dir, "tasks.json")
    if exists(dag) && exists(tasks) {
        if _, err := writeDagDot(dag, tasks, join(dir, "dag.dot"), nodeLabel, edgeLabel); err != nil { return err }
        fmt.Println("wrote", join(dir, "dag.dot"))
    }
    return nil
}

func runExportDot() {
    fs := flag.NewFlagSet("export-dot", flag.ExitOnError)
    dir := fs.String("dir", "", "Directory containing artifacts (emits defaults)")
    dagPath := fs.String("dag", "", "Path to dag.json")
    tasksPath := fs.String("tasks", "", "Path to tasks.json")
    coordPath := fs.String("coordinator", "", "Path to coordinator.json")
    outPath := fs.String("out", "", "Output file path (defaults to stdout)")
    nodeLabel := fs.String("node-label", "id-title", "Node label style: id|title|id-title")
    edgeLabel := fs.String("edge-label", "type", "Edge label style: none|type")
    _ = fs.Parse(os.Args[2:])

    // Directory mode: emit both if present
    if *dir != "" {
        if err := emitDirArtifacts(*dir, *nodeLabel, *edgeLabel); err != nil {
            fmt.Fprintf(os.Stderr, "export-dot: %v\n", err)
            os.Exit(1)
        }
        return
    }

    // Coordinator single-file mode
    if *coordPath != "" {
        dotStr, err := writeCoordinatorDot(*coordPath, *outPath, *nodeLabel, *edgeLabel)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Failed to write %s: %v\n", *outPath, err)
            os.Exit(1)
        }
        if *outPath == "" { fmt.Print(dotStr) }
        return
    }

    // dag + tasks mode
    if *dagPath == "" || *tasksPath == "" {
        fmt.Fprintln(os.Stderr, "Usage: tasksd export-dot --dir ./plans | --coordinator ./coordinator.json | --dag ./dag.json --tasks ./tasks.json [--out ./dag.dot]")
        os.Exit(1)
    }
    if dotStr, err := writeDagDot(*dagPath, *tasksPath, *outPath, *nodeLabel, *edgeLabel); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to export DOT: %v\n", err)
        os.Exit(1)
    } else if *outPath == "" { fmt.Print(dotStr) }
}

func loadJSON(path string, v any) error {
    f, err := os.Open(path)
    if err != nil { return err }
    defer f.Close()
    dec := json.NewDecoder(f)
    dec.DisallowUnknownFields()
    return dec.Decode(v)
}

func exists(path string) bool {
    _, err := os.Stat(path)
    return err == nil
}

func join(dir, name string) string { return filepath.Join(dir, name) }

// -----------------
// plan stub
// -----------------

func runPlan() {
    fs := flag.NewFlagSet("plan", flag.ExitOnError)
    doc := fs.String("doc", "", "Path to plan document (unused in stub)")
    repo := fs.String("repo", ".", "Path to codebase for census (optional)")
    out := fs.String("out", "./plans", "Output directory for artifacts")
    _ = fs.Parse(os.Args[2:])

    if err := os.MkdirAll(*out, 0o755); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create output dir %s: %v\n", *out, err)
        os.Exit(1)
    }

    // Build a tiny, valid plan (or parse from --doc)
    tasks := []m.Task{}
    featuresList := []struct{ ID, Title string }{}
    if *doc != "" && exists(*doc) {
        bs, err := os.ReadFile(*doc)
        if err != nil { fmt.Fprintf(os.Stderr, "Failed to read --doc: %v\n", err); os.Exit(1) }
        feats, tks := docp.ParseMarkdown(string(bs))
        if len(feats) == 0 && len(tks) == 0 {
            fmt.Fprintln(os.Stderr, "Doc parsed empty; falling back to stub plan")
        } else {
            // build features
            for _, f := range feats {
                featuresList = append(featuresList, struct{ ID, Title string }{ID: f.ID, Title: f.Title})
            }
            // build tasks with IDs T001, T002...
            idx := 0
            for _, ts := range tks {
                idx++
                id := fmt.Sprintf("T%03d", idx)
                t := m.Task{
                    ID:        id,
                    FeatureID: ts.FeatureID,
                    Title:     ts.Title,
                    Duration:  m.DurationPERT{Optimistic: 1, MostLikely: 2, Pessimistic: 3},
                    DurationUnit: "hours",
                    AcceptanceChecks: []m.AcceptanceCheck{{Type: "command", Cmd: "echo ok", Timeout: 5}},
                }
                if ts.Hours > 0 {
                    ml := ts.Hours
                    t.Duration = m.DurationPERT{Optimistic: ml * 0.5, MostLikely: ml, Pessimistic: ml * 2}
                }
                if len(ts.Errors) > 0 {
                    fmt.Fprintf(os.Stderr, "spec parse errors for task %s: %s\n", id, strings.Join(ts.Errors, "; "))
                    os.Exit(1)
                }
                if len(ts.Accept) > 0 {
                    t.AcceptanceChecks = ts.Accept
                }
                t.ExecutionLogging.Format = "JSONL"
                t.ExecutionLogging.RequiredFields = []string{"timestamp","task_id","step","status","message"}
                t.Compensation.Idempotent = true
                tasks = append(tasks, t)
            }
        }
    }
    if len(tasks) == 0 { // fallback stub
        tasks = []m.Task{
            func() m.Task { t := m.Task{
                ID:        "T001",
                FeatureID: "F1",
                Title:     "Setup DB",
                Duration:  m.DurationPERT{Optimistic: 1, MostLikely: 2, Pessimistic: 3},
                DurationUnit: "hours",
                AcceptanceChecks: []m.AcceptanceCheck{{Type: "command", Cmd: "echo ok", Timeout: 5}},
            }; t.ExecutionLogging.Format = "JSONL"; t.ExecutionLogging.RequiredFields = []string{"timestamp","task_id","step","status","message"}; t.Compensation.Idempotent = true; return t }(),
            func() m.Task { t := m.Task{
                ID:        "T002",
                FeatureID: "F1",
                Title:     "Migrate Schema",
                Duration:  m.DurationPERT{Optimistic: 1, MostLikely: 2, Pessimistic: 3},
                DurationUnit: "hours",
                AcceptanceChecks: []m.AcceptanceCheck{{Type: "command", Cmd: "echo ok", Timeout: 5}},
            }; t.ExecutionLogging.Format = "JSONL"; t.Compensation.Idempotent = true; return t }(),
            func() m.Task { t := m.Task{
                ID:        "T003",
                FeatureID: "F1",
                Title:     "API Handlers",
                Duration:  m.DurationPERT{Optimistic: 1, MostLikely: 2, Pessimistic: 3},
                DurationUnit: "hours",
                AcceptanceChecks: []m.AcceptanceCheck{{Type: "command", Cmd: "echo ok", Timeout: 5}},
            }; t.ExecutionLogging.Format = "JSONL"; t.Compensation.Idempotent = true; return t }(),
        }
        featuresList = []struct{ ID, Title string }{{ID: "F1", Title: "Core DB + API"}}
    }

    // optional census
    var analysis any = map[string]any{}
    if *repo != "" {
        if a, err := analysisPkg(*repo); err == nil {
            analysis = a
        }
    }

    tf := m.TasksFile{}
    tf.Meta.Version = "v8"
    tf.Meta.MinConfidence = 0.7
    tf.Meta.CodebaseAnalysis = analysis
    tf.Meta.Autonormalization.Split = []string{}
    tf.Meta.Autonormalization.Merged = []string{}
    tf.Tasks = tasks
    // Build dependencies: prefer explicit 'after:' from --doc; else linear fallback
    titleToID := map[string]string{}
    for _, t := range tasks { titleToID[strings.ToLower(strings.TrimSpace(t.Title))] = t.ID }
    depEdges := []m.Edge{}
    if *doc != "" && exists(*doc) {
        if bs, err := os.ReadFile(*doc); err == nil {
            _, tks := docp.ParseMarkdown(string(bs))
            for _, ts := range tks {
                if len(ts.After) == 0 { continue }
                toID := titleToID[strings.ToLower(strings.TrimSpace(ts.Title))]
                for _, dep := range ts.After {
                    d := strings.TrimSpace(dep)
                    dl := strings.ToLower(d)
                    fromID := ""
                    // Match TIDs like T001, T12, case-insensitive
                    if len(d) > 1 && (d[0] == 'T' || d[0] == 't') {
                        isNum := true
                        for _, r := range d[1:] { if r < '0' || r > '9' { isNum = false; break } }
                        if isNum { fromID = strings.ToUpper(d) }
                    }
                    if fromID == "" {
                        if v, ok := titleToID[dl]; ok { fromID = v } else { continue }
                    }
                    depEdges = append(depEdges, m.Edge{From: fromID, To: toID, Type: "sequential", IsHard: true, Confidence: 1.0})
                }
            }
        }
    }
    if len(depEdges) == 0 && len(tasks) >= 2 {
        for i := 0; i < len(tasks)-1; i++ {
            typ := "sequential"
            if i > 0 { typ = "technical" }
            depEdges = append(depEdges, m.Edge{From: tasks[i].ID, To: tasks[i+1].ID, Type: typ, IsHard: true, Confidence: 1.0})
        }
    }
    tf.Dependencies = depEdges

    // Enforce acceptance for doc-driven plans
    if *doc != "" {
        for _, t := range tf.Tasks {
            if len(t.AcceptanceChecks) == 0 {
                fmt.Fprintf(os.Stderr, "task %s missing acceptance checks; add fenced ```acceptance``` block in spec\n", t.ID)
                os.Exit(1)
            }
        }
    }

    // Validate tasks.json before DAG build
    if err := validate.TasksFile(&tf); err != nil {
        fmt.Fprintf(os.Stderr, "tasks.json validation failed: %v\n", err)
        os.Exit(1)
    }

    // features.json (minimal generic map)
    features := map[string]any{"meta": map[string]any{"version": "v8", "artifact_hash": ""}}
    featsArr := make([]any, 0, len(featuresList))
    for _, f := range featuresList { featsArr = append(featsArr, map[string]any{"id": f.ID, "title": f.Title}) }
    features["features"] = featsArr

    // Compute resource conflicts summary and resource edges (traceability; excluded from DAG)
    resToTasks := map[string][]string{}
    for _, t := range tf.Tasks {
        for _, r := range t.Resources.Exclusive { resToTasks[r] = append(resToTasks[r], t.ID) }
    }
    tf.ResourceConflicts = map[string]any{}
    for r, ids := range resToTasks {
        if len(ids) < 2 { continue }
        tf.ResourceConflicts[r] = map[string]any{"type":"exclusive","tasks": ids}
        for i := 0; i < len(ids); i++ {
            for j := i+1; j < len(ids); j++ {
                tf.Dependencies = append(tf.Dependencies, m.Edge{From: ids[i], To: ids[j], Type: "resource", Subtype: "mutual_exclusion", IsHard: true, Confidence: 1.0})
            }
        }
    }

    // Build DAG
    df, err := dagbuild.Build(tasks, tf.Dependencies, tf.Meta.MinConfidence)
    if err != nil {
        fmt.Fprintf(os.Stderr, "DAG build failed: %v\n", err)
        os.Exit(1)
    }
    if err := validate.DagFile(&df); err != nil {
        fmt.Fprintf(os.Stderr, "dag.json validation failed: %v\n", err)
        os.Exit(1)
    }

    // Coordinator (runtime view)
    coord := m.Coordinator{}
    coord.Version = "v8"
    coord.Graph.Nodes = tasks
    coord.Graph.Edges = tf.Dependencies
    coord.Config.Resources.Catalog = map[string]struct {
        Capacity int    `json:"capacity"`
        Mode     string `json:"mode"`
        LockOrder int   `json:"lock_order"`
    }{}
    coord.Config.Resources.Profiles = map[string]map[string]int{"default": {}}
    coord.Config.Policies.ConcurrencyMax = 4
    coord.Config.Policies.LockOrdering = []string{}
    coord.Metrics.Estimates.P50TotalHours = 6
    coord.Metrics.Estimates.LongestPathLength = 3
    coord.Metrics.Estimates.WidthApprox = 1

    // waves.json via wavesim (preview only; no feedback to DAG)
    w := map[string]any{"meta": map[string]any{"version":"v8","planId":"","artifact_hash":""}}
    w["waves"] = wavesim.Generate(df, tasks)

    // Write artifacts with canonical JSON + hashes
    hashes := map[string]string{}
    mustWriteJSONWithHash := func(name string, v any, setHash func(h string)) {
        h, err := emitter.WriteWithArtifactHash(join(*out, name), v, setHash)
        if err != nil { fmt.Fprintf(os.Stderr, "%v\n", err); os.Exit(1) }
        hashes[name] = h
    }

    // tasks.json
    mustWriteJSONWithHash("tasks.json", &tf, func(h string) { tf.Meta.ArtifactHash = h })
    // Fill dependent hashes
    df.Meta.ArtifactHash = "" // set by write
    df.Meta.TasksHash = hashes["tasks.json"]
    mustWriteJSONWithHash("dag.json", &df, func(h string) { df.Meta.ArtifactHash = h })
    if wm, ok := w["meta"].(map[string]any); ok { wm["planId"] = hashes["tasks.json"] }
    mustWriteJSONWithHash("waves.json", &w, func(h string) {
        if wm, ok := w["meta"].(map[string]any); ok { wm["artifact_hash"] = h }
    })
    mustWriteJSONWithHash("features.json", &features, func(h string) {
        if meta, ok := features["meta"].(map[string]any); ok {
            meta["artifact_hash"] = h
        }
    })
    mustWriteJSONWithHash("coordinator.json", &coord, func(h string) {})

    // Plan.md with hashes
    var md strings.Builder
    md.WriteString("# Plan (stub)\n\n")
    md.WriteString("## Hashes\n\n")
    for _, name := range []string{"features.json","tasks.json","dag.json","waves.json","coordinator.json"} {
        fmt.Fprintf(&md, "- %s: %s\n", name, hashes[name])
    }
    if err := os.WriteFile(join(*out, "Plan.md"), []byte(md.String()), 0o644); err != nil {
        fmt.Fprintf(os.Stderr, "write Plan.md: %v\n", err)
        os.Exit(1)
    }

    // Emit DOTs next to artifacts (directory mode)
    // dag.dot
    titles := map[string]string{}
    for _, t := range tasks { titles[t.ID] = t.Title }
    dagDot := dot.FromDagWithOptions(df, titles, dot.Options{NodeLabel: "id-title", EdgeLabel: "type"})
    if err := os.WriteFile(join(*out, "dag.dot"), []byte(dagDot), 0o644); err != nil { fmt.Fprintf(os.Stderr, "write dag.dot: %v\n", err); os.Exit(1) }
    // runtime.dot
    runtimeDot := dot.FromCoordinatorWithOptions(coord, dot.Options{NodeLabel: "id-title", EdgeLabel: "type"})
    if err := os.WriteFile(join(*out, "runtime.dot"), []byte(runtimeDot), 0o644); err != nil { fmt.Fprintf(os.Stderr, "write runtime.dot: %v\n", err); os.Exit(1) }

    fmt.Println("Plan stub written to", *out)
}

// -----------------
// validate
// -----------------
func runValidate() {
    fs := flag.NewFlagSet("validate", flag.ExitOnError)
    dir := fs.String("dir", "", "Directory containing artifacts")
    _ = fs.Parse(os.Args[2:])
    if *dir == "" { fmt.Fprintln(os.Stderr, "Usage: tasksd validate --dir ./plans"); os.Exit(1) }

    read := func(name string) ([]byte, bool) {
        p := join(*dir, name)
        b, err := os.ReadFile(p)
        if err != nil { return nil, false }
        return b, true
    }

    okAll := true

    // features.json
    if b, ok := read("features.json"); ok {
        if comp, stored, okHash, err := validate.CheckArtifactHash(b); err != nil || !okHash {
            okAll = false
            fmt.Fprintf(os.Stderr, "features.json hash mismatch: computed=%s stored=%s err=%v\n", comp, stored, err)
        }
        if err := validate.ValidateRaw("features.json", b); err != nil { okAll = false; fmt.Fprintf(os.Stderr, "features.json schema: %v\n", err) } else { fmt.Println("OK features.json") }
    }

    // tasks.json
    if b, ok := read("tasks.json"); ok {
        if comp, stored, okHash, err := validate.CheckArtifactHash(b); err != nil || !okHash {
            okAll = false
            fmt.Fprintf(os.Stderr, "tasks.json hash mismatch: computed=%s stored=%s err=%v\n", comp, stored, err)
        }
        var tf m.TasksFile
        if err := json.Unmarshal(b, &tf); err != nil { okAll = false; fmt.Fprintf(os.Stderr, "tasks.json parse: %v\n", err) } else if err := validate.TasksFile(&tf); err != nil { okAll = false; fmt.Fprintf(os.Stderr, "tasks.json: %v\n", err) } else { fmt.Println("OK tasks.json") }
    }

    // dag.json
    if b, ok := read("dag.json"); ok {
        if comp, stored, okHash, err := validate.CheckArtifactHash(b); err != nil || !okHash {
            okAll = false
            fmt.Fprintf(os.Stderr, "dag.json hash mismatch: computed=%s stored=%s err=%v\n", comp, stored, err)
        }
        var df m.DagFile
        if err := json.Unmarshal(b, &df); err != nil { okAll = false; fmt.Fprintf(os.Stderr, "dag.json parse: %v\n", err) } else if err := validate.DagFile(&df); err != nil { okAll = false; fmt.Fprintf(os.Stderr, "dag.json: %v\n", err) } else { fmt.Println("OK dag.json") }
    }

    // waves.json
    if b, ok := read("waves.json"); ok {
        if comp, stored, okHash, err := validate.CheckArtifactHash(b); err != nil || !okHash {
            okAll = false
            fmt.Fprintf(os.Stderr, "waves.json hash mismatch: computed=%s stored=%s err=%v\n", comp, stored, err)
        }
        if err := validate.ValidateRaw("waves.json", b); err != nil { okAll = false; fmt.Fprintf(os.Stderr, "waves.json schema: %v\n", err) } else { fmt.Println("OK waves.json") }
    }

    // coordinator.json
    if b, ok := read("coordinator.json"); ok {
        if err := validate.ValidateRaw("coordinator.json", b); err != nil { okAll = false; fmt.Fprintf(os.Stderr, "coordinator.json schema: %v\n", err) } else { fmt.Println("OK coordinator.json") }
    }

    if !okAll { os.Exit(2) }
    fmt.Println("All artifacts valid.")
}

// analysisPkg wraps the codebase census but returns a compact map for embedding.
func analysisPkg(path string) (any, error) {
    type fileCensus struct{ Files, GoFiles int }
    a, err := analysis.RunCensus(path)
    if err != nil { return nil, err }
    return fileCensus{Files: len(a.Files), GoFiles: len(a.GoFiles)}, nil
}
