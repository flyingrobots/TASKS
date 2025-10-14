package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	analysis "github.com/james/tasks-planner/internal/analysis"
	"github.com/james/tasks-planner/internal/canonjson"
	"github.com/james/tasks-planner/internal/emitter"
	"github.com/james/tasks-planner/internal/export/dot"
	"github.com/james/tasks-planner/internal/hash"
	m "github.com/james/tasks-planner/internal/model"
	dagbuild "github.com/james/tasks-planner/internal/planner/dag"
	docp "github.com/james/tasks-planner/internal/planner/docparse"
	wavesim "github.com/james/tasks-planner/internal/planner/wavesim"
	"github.com/james/tasks-planner/internal/validate"
	validators "github.com/james/tasks-planner/internal/validators"
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
	if err := loadJSON(coordPath, &c); err != nil {
		return "", err
	}
	dotStr := dot.FromCoordinatorWithOptions(c, dot.Options{NodeLabel: nodeLabel, EdgeLabel: edgeLabel})
	if outPath == "" {
		return dotStr, nil
	}
	return "", os.WriteFile(outPath, []byte(dotStr), 0o644)
}

func writeDagDot(dagPath, tasksPath, outPath, nodeLabel, edgeLabel string) (string, error) {
	var tf m.TasksFile
	if err := loadJSON(tasksPath, &tf); err != nil {
		return "", err
	}
	titles := make(map[string]string)
	for _, t := range tf.Tasks {
		titles[t.ID] = t.Title
	}
	var df m.DagFile
	if err := loadJSON(dagPath, &df); err != nil {
		return "", err
	}
	dotStr := dot.FromDagWithOptions(df, titles, dot.Options{NodeLabel: nodeLabel, EdgeLabel: edgeLabel})
	if outPath == "" {
		return dotStr, nil
	}
	return "", os.WriteFile(outPath, []byte(dotStr), 0o644)
}

func emitDirArtifacts(dir, nodeLabel, edgeLabel string) error {
	// coordinator
	coord := join(dir, "coordinator.json")
	if exists(coord) {
		if _, err := writeCoordinatorDot(coord, join(dir, "runtime.dot"), nodeLabel, edgeLabel); err != nil {
			return err
		}
		fmt.Println("wrote", join(dir, "runtime.dot"))
	}
	// dag + tasks
	dag := join(dir, "dag.json")
	tasks := join(dir, "tasks.json")
	if exists(dag) && exists(tasks) {
		if _, err := writeDagDot(dag, tasks, join(dir, "dag.dot"), nodeLabel, edgeLabel); err != nil {
			return err
		}
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
		if *outPath == "" {
			fmt.Print(dotStr)
		}
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
	} else if *outPath == "" {
		fmt.Print(dotStr)
	}
}

func loadJSON(path string, v any) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
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
	acceptanceCmd := fs.String(
		"validators-acceptance",
		"",
		"Executable path or shell command for the acceptance validator; tasksd runs it via the system shell, streams canonical plan JSON on stdin, and expects a JSON report on stdout.",
	)
	evidenceCmd := fs.String(
		"validators-evidence",
		"",
		"Executable path or shell command for the evidence validator; supports arguments and runs via the system shell. Output must be JSON on stdout matching the validator report schema.",
	)
	interfaceCmd := fs.String(
		"validators-interface",
		"",
		"Executable path or shell command for the interface validator; executed through the system shell with JSON input on stdin and JSON report expected on stdout.",
	)
	validatorsCache := fs.String(
		"validators-cache",
		"",
		"Filesystem directory for cached validator reports (created if missing).",
	)
	validatorsTimeout := fs.Duration(
		"validators-timeout",
		30*time.Second,
		"Timeout applied to each validator command invocation.",
	)
	validatorsStrict := fs.Bool(
		"validators-strict",
		false,
		"Exit with an error when any validator fails (default logs warnings only).",
	)
	_ = fs.Parse(os.Args[2:])

	if err := os.MkdirAll(*out, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create output dir %s: %v\n", *out, err)
		os.Exit(1)
	}

	tasks, featuresList, docEdges, docProvided, err := buildTasksFromDoc(*doc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "plan: %v\n", err)
		os.Exit(1)
	}

	var analysisVal any = map[string]any{}
	if *repo != "" {
		if a, err := analysisPkg(*repo); err == nil {
			analysisVal = a
		} else {
			log.Printf("analysisPkg failed for repo %s: %v", *repo, err)
		}
	}

	tf := m.TasksFile{}
	tf.Meta.Version = "v8"
	tf.Meta.MinConfidence = 0.7
	tf.Meta.CodebaseAnalysis = analysisVal
	tf.Meta.Autonormalization.Split = []string{}
	tf.Meta.Autonormalization.Merged = []string{}
	tf.Tasks = tasks

	deps, resourceConflicts := inferDependencies(tasks, docEdges)
	tf.Dependencies = deps
	tf.ResourceConflicts = resourceConflicts

	if docProvided {
		for _, t := range tf.Tasks {
			if len(t.AcceptanceChecks) == 0 {
				fmt.Fprintf(os.Stderr, "task %s missing acceptance checks; add fenced ```acceptance``` block in spec\n", t.ID)
				os.Exit(1)
			}
		}
	}

	if err := validate.TasksFile(&tf); err != nil {
		fmt.Fprintf(os.Stderr, "tasks.json validation failed: %v\n", err)
		os.Exit(1)
	}

	df, err := dagbuild.Build(tasks, tf.Dependencies, tf.Meta.MinConfidence)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DAG build failed: %v\n", err)
		os.Exit(1)
	}
	if err := validate.DagFile(&df); err != nil {
		fmt.Fprintf(os.Stderr, "dag.json validation failed: %v\n", err)
		os.Exit(1)
	}

	coord := makeCoordinator(tasks, tf.Dependencies)

	waves, err := buildWaves(df, tasks)
	if err != nil {
		fmt.Fprintf(os.Stderr, "wavesim failed: %v\n", err)
		os.Exit(1)
	}

	features := makeFeaturesArtifact(featuresList, tasks)
	titles := taskTitles(tasks)

	validatorPayload := validators.Payload{
		Tasks:       &tf,
		Dag:         &df,
		Coordinator: &coord,
	}
	validatorCfg := validators.Config{
		AcceptanceCmd: *acceptanceCmd,
		EvidenceCmd:   *evidenceCmd,
		InterfaceCmd:  *interfaceCmd,
		CacheDir:      *validatorsCache,
		Timeout:       *validatorsTimeout,
	}
	var validatorReports []validators.Report
	if validatorCfg.AcceptanceCmd != "" || validatorCfg.EvidenceCmd != "" || validatorCfg.InterfaceCmd != "" {
		runner, err := validators.NewRunner(validatorCfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "validator runner: %v\n", err)
			os.Exit(1)
		}
		reports, runErr := runner.Run(context.Background(), validatorPayload)
		tf.Meta.ValidatorReports = convertValidatorReports(reports)
		validatorReports = reports
		if runErr != nil {
			if *validatorsStrict {
				fmt.Fprintf(os.Stderr, "validators reported issues: %v\n", runErr)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stderr, "validators warning: %v\n", runErr)
		}
	}

	if err := writeArtifacts(*out, &tf, &df, &coord, features, waves, titles, validatorReports); err != nil {
		fmt.Fprintf(os.Stderr, "write artifacts: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Plan stub written to", *out)
}

type featureSummary struct {
	ID    string
	Title string
}

func buildTasksFromDoc(docPath string) ([]m.Task, []featureSummary, []m.Edge, bool, error) {
	if docPath == "" || !exists(docPath) {
		tasks, features := stubTasksAndFeatures()
		return tasks, features, nil, false, nil
	}
	bs, err := os.ReadFile(docPath)
	if err != nil {
		return nil, nil, nil, false, fmt.Errorf("read --doc: %w", err)
	}
	feats, tks := docp.ParseMarkdown(string(bs))
	if len(feats) == 0 && len(tks) == 0 {
		fmt.Fprintln(os.Stderr, "Doc parsed empty; falling back to stub plan")
		tasks, features := stubTasksAndFeatures()
		return tasks, features, nil, false, nil
	}
	featuresList := make([]featureSummary, 0, len(feats))
	for _, f := range feats {
		featuresList = append(featuresList, featureSummary{ID: f.ID, Title: f.Title})
	}
	tasks := make([]m.Task, 0, len(tks))
	titleToID := map[string]string{}
	var parseErrors []string
	for i, spec := range tks {
		id := fmt.Sprintf("T%03d", i+1)
		if len(spec.Errors) > 0 {
			for _, e := range spec.Errors {
				parseErrors = append(parseErrors, fmt.Sprintf("%s: %s", spec.Title, e))
			}
		}
		task := m.Task{
			ID:        id,
			FeatureID: spec.FeatureID,
			Title:     spec.Title,
			Duration:  m.DurationPERT{Optimistic: 1, MostLikely: 2, Pessimistic: 3},
		}
		if spec.Hours > 0 {
			ml := spec.Hours
			task.Duration = m.DurationPERT{Optimistic: ml * 0.5, MostLikely: ml, Pessimistic: ml * 2}
		}
		if len(spec.Accept) > 0 {
			task.AcceptanceChecks = append(task.AcceptanceChecks, spec.Accept...)
		}
		applyTaskDefaults(&task)
		tasks = append(tasks, task)
		titleToID[normalizeKey(spec.Title)] = id
	}
	if len(parseErrors) > 0 {
		return nil, nil, nil, false, fmt.Errorf("doc parse errors: %s", strings.Join(parseErrors, "; "))
	}
	docEdges := []m.Edge{}
	for _, spec := range tks {
		toID := titleToID[normalizeKey(spec.Title)]
		if toID == "" {
			continue
		}
		for _, raw := range spec.After {
			fromID := resolveTaskID(raw, titleToID)
			if fromID == "" {
				continue
			}
			docEdges = append(docEdges, m.Edge{From: fromID, To: toID, Type: "sequential", IsHard: true, Confidence: 1.0})
		}
	}
	if len(featuresList) == 0 {
		featuresList = featuresFromTasks(tasks)
	}
	return tasks, featuresList, docEdges, true, nil
}

func stubTasksAndFeatures() ([]m.Task, []featureSummary) {
	base := []struct {
		id        string
		featureID string
		title     string
	}{
		{"T001", "F1", "Setup DB"},
		{"T002", "F1", "Migrate Schema"},
		{"T003", "F1", "API Handlers"},
	}
	tasks := make([]m.Task, 0, len(base))
	for _, spec := range base {
		task := m.Task{
			ID:        spec.id,
			FeatureID: spec.featureID,
			Title:     spec.title,
			Duration:  m.DurationPERT{Optimistic: 1, MostLikely: 2, Pessimistic: 3},
		}
		applyTaskDefaults(&task)
		tasks = append(tasks, task)
	}
	features := []featureSummary{{ID: "F1", Title: "Core DB + API"}}
	return tasks, features
}

func inferDependencies(tasks []m.Task, baseEdges []m.Edge) ([]m.Edge, map[string]any) {
	deps := append([]m.Edge{}, baseEdges...)
	if len(deps) == 0 && len(tasks) >= 2 {
		for i := 0; i < len(tasks)-1; i++ {
			typ := "sequential"
			if i > 0 {
				typ = "technical"
			}
			deps = append(deps, m.Edge{From: tasks[i].ID, To: tasks[i+1].ID, Type: typ, IsHard: true, Confidence: 1.0})
		}
	}
	resToTasks := map[string][]string{}
	for _, task := range tasks {
		for _, r := range task.Resources.Exclusive {
			resToTasks[r] = append(resToTasks[r], task.ID)
		}
	}
	resourceConflicts := map[string]any{}
	keys := make([]string, 0, len(resToTasks))
	for k := range resToTasks {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, r := range keys {
		ids := resToTasks[r]
		sort.Strings(ids)
		if len(ids) < 2 {
			continue
		}
		resourceConflicts[r] = map[string]any{"type": "exclusive", "tasks": ids}
		for i := 0; i < len(ids); i++ {
			for j := i + 1; j < len(ids); j++ {
				deps = append(deps, m.Edge{From: ids[i], To: ids[j], Type: "resource", Subtype: "mutual_exclusion", IsHard: true, Confidence: 1.0})
			}
		}
	}
	return deps, resourceConflicts
}

func makeCoordinator(tasks []m.Task, deps []m.Edge) m.Coordinator {
	coord := m.Coordinator{}
	coord.Version = "v8"
	coord.Graph.Nodes = tasks
	coord.Graph.Edges = deps
	coord.Config.Resources.Catalog = map[string]struct {
		Capacity  int    `json:"capacity"`
		Mode      string `json:"mode"`
		LockOrder int    `json:"lock_order"`
	}{}
	coord.Config.Resources.Profiles = map[string]map[string]int{"default": {}}
	coord.Config.Policies.ConcurrencyMax = 4
	coord.Config.Policies.LockOrdering = []string{}
	coord.Metrics.Estimates.P50TotalHours = 6
	coord.Metrics.Estimates.LongestPathLength = 3
	coord.Metrics.Estimates.WidthApprox = 1
	return coord
}

func buildWaves(df m.DagFile, tasks []m.Task) (map[string]any, error) {
	waves := map[string]any{"meta": map[string]any{"version": "v8", "planId": "", "artifact_hash": ""}}
	ids, err := wavesim.Generate(df, tasks)
	if err != nil {
		return nil, err
	}
	waves["waves"] = ids
	return waves, nil
}

func makeFeaturesArtifact(features []featureSummary, tasks []m.Task) map[string]any {
	if len(features) == 0 {
		features = featuresFromTasks(tasks)
	}
	sort.Slice(features, func(i, j int) bool { return features[i].ID < features[j].ID })
	entries := make([]any, 0, len(features))
	for _, f := range features {
		entries = append(entries, map[string]any{"id": f.ID, "title": f.Title})
	}
	return map[string]any{
		"meta":     map[string]any{"version": "v8", "artifact_hash": ""},
		"features": entries,
	}
}

func featuresFromTasks(tasks []m.Task) []featureSummary {
	seen := map[string]bool{}
	features := []featureSummary{}
	for _, t := range tasks {
		if t.FeatureID == "" {
			continue
		}
		if seen[t.FeatureID] {
			continue
		}
		features = append(features, featureSummary{ID: t.FeatureID, Title: t.FeatureID})
		seen[t.FeatureID] = true
	}
	sort.Slice(features, func(i, j int) bool { return features[i].ID < features[j].ID })
	return features
}

func taskTitles(tasks []m.Task) map[string]string {
	titles := make(map[string]string, len(tasks))
	for _, t := range tasks {
		titles[t.ID] = t.Title
	}
	return titles
}

const validatorDetailLimit = 2048

func convertValidatorReports(src []validators.Report) []m.ValidatorReport {
	if len(src) == 0 {
		return nil
	}
	// Compile-time compatibility check between validators.Report and model.ValidatorReport.
	var _ = func() m.ValidatorReport {
		var r validators.Report
		return m.ValidatorReport{
			Name:      r.Name,
			Status:    r.Status,
			Command:   r.Command,
			InputHash: r.InputHash,
			Cached:    r.Cached,
			Detail:    r.Detail,
			RawOutput: r.RawOutput,
		}
	}
	dst := make([]m.ValidatorReport, 0, len(src))
	for _, rep := range src {
		converted := m.ValidatorReport{
			Name:      rep.Name,
			Status:    rep.Status,
			Command:   rep.Command,
			InputHash: rep.InputHash,
			Cached:    rep.Cached,
			Detail:    truncateDetail(rep.Detail, validatorDetailLimit),
			RawOutput: rep.RawOutput,
		}
		dst = append(dst, converted)
	}
	return dst
}

func writeArtifacts(outDir string, tf *m.TasksFile, df *m.DagFile, coord *m.Coordinator, features map[string]any, waves map[string]any, titles map[string]string, validatorReports []validators.Report) error {
	hashes := map[string]string{}
	var errs []error
	writeWithHash := func(name string, value any, setHash func(string)) {
		hashValue, err := emitter.WriteWithArtifactHash(join(outDir, name), value, setHash)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", name, err))
			return
		}
		hashes[name] = hashValue
	}

	writeWithHash("tasks.json", tf, func(h string) { tf.Meta.ArtifactHash = h })
	if h, ok := hashes["tasks.json"]; ok {
		df.Meta.TasksHash = h
	}
	df.Meta.ArtifactHash = ""
	writeWithHash("dag.json", df, func(h string) { df.Meta.ArtifactHash = h })

	if meta, ok := waves["meta"].(map[string]any); ok {
		if h, ok := hashes["tasks.json"]; ok {
			meta["planId"] = h
		}
	}
	writeWithHash("waves.json", waves, func(h string) {
		if meta, ok := waves["meta"].(map[string]any); ok {
			meta["artifact_hash"] = h
		}
	})

	writeWithHash("features.json", features, func(h string) {
		if meta, ok := features["meta"].(map[string]any); ok {
			meta["artifact_hash"] = h
		}
	})

	writeWithHash("coordinator.json", coord, nil)

	if err := writePlanSummary(outDir, hashes, validatorReports); err != nil {
		errs = append(errs, err)
	}

	dagDot := dot.FromDagWithOptions(*df, titles, dot.Options{NodeLabel: "id-title", EdgeLabel: "type"})
	if err := os.WriteFile(join(outDir, "dag.dot"), []byte(dagDot), 0o644); err != nil {
		errs = append(errs, fmt.Errorf("write dag.dot: %w", err))
	}
	runtimeDot := dot.FromCoordinatorWithOptions(*coord, dot.Options{NodeLabel: "id-title", EdgeLabel: "type"})
	if err := os.WriteFile(join(outDir, "runtime.dot"), []byte(runtimeDot), 0o644); err != nil {
		errs = append(errs, fmt.Errorf("write runtime.dot: %w", err))
	}
	return errors.Join(errs...)
}

func applyTaskDefaults(task *m.Task) {
	if len(task.AcceptanceChecks) == 0 {
		task.AcceptanceChecks = []m.AcceptanceCheck{{Type: "command", Cmd: "echo ok", Timeout: 5}}
	}
	if task.DurationUnit == "" {
		task.DurationUnit = "hours"
	}
	task.ExecutionLogging.Format = "JSONL"
	if len(task.ExecutionLogging.RequiredFields) == 0 {
		task.ExecutionLogging.RequiredFields = []string{"timestamp", "task_id", "step", "status", "message"}
	}
	task.Compensation.Idempotent = true
}

func resolveTaskID(token string, titleToID map[string]string) string {
	trimmed := strings.TrimSpace(token)
	if len(trimmed) > 1 && (trimmed[0] == 'T' || trimmed[0] == 't') {
		isNumeric := true
		for _, r := range trimmed[1:] {
			if r < '0' || r > '9' {
				isNumeric = false
				break
			}
		}
		if isNumeric {
			return strings.ToUpper(trimmed)
		}
	}
	return titleToID[normalizeKey(trimmed)]
}

func normalizeKey(v string) string {
	return strings.ToLower(strings.TrimSpace(v))
}

func writePlanSummary(outDir string, hashes map[string]string, validatorReports []validators.Report) error {
	names := []string{"features.json", "tasks.json", "dag.json", "waves.json", "coordinator.json"}
	var md strings.Builder
	md.WriteString("# Plan (stub)\n\n")
	md.WriteString("## Hashes\n\n")
	for _, name := range names {
		hashValue := hashes[name]
		fmt.Fprintf(&md, "- %s: %s\n", name, hashValue)
	}
	if len(validatorReports) > 0 {
		md.WriteString("\n## Validators\n\n")
		for _, rep := range validatorReports {
			cached := ""
			if rep.Cached {
				cached = " (cached)"
			}
			detail := truncateDetail(rep.Detail, validatorDetailLimit)
			if detail == "" && len(rep.RawOutput) > 0 {
				detail = truncateDetail(string(rep.RawOutput), validatorDetailLimit)
			}
			if detail != "" {
				fmt.Fprintf(&md, "- %s: %s%s — %s\n", rep.Name, rep.Status, cached, detail)
			} else {
				fmt.Fprintf(&md, "- %s: %s%s\n", rep.Name, rep.Status, cached)
			}
		}
	}
	return os.WriteFile(join(outDir, "Plan.md"), []byte(md.String()), 0o644)
}

func truncateDetail(detail string, limit int) string {
	if limit <= 0 {
		return detail
	}
	runes := []rune(detail)
	if len(runes) <= limit {
		return detail
	}
	return string(runes[:limit]) + " … (truncated)"
}

// -----------------
// validate
// -----------------
func runValidate() {
	fs := flag.NewFlagSet("validate", flag.ExitOnError)
	dir := fs.String("dir", "", "Directory containing artifacts")
	_ = fs.Parse(os.Args[2:])
	if *dir == "" {
		fmt.Fprintln(os.Stderr, "Usage: tasksd validate --dir ./plans")
		os.Exit(1)
	}

	read := func(name string) ([]byte, bool) {
		p := join(*dir, name)
		b, err := os.ReadFile(p)
		if err != nil {
			return nil, false
		}
		return b, true
	}

	okAll := true

	// features.json
	if b, ok := read("features.json"); ok {
		if comp, stored, okHash, err := validate.CheckArtifactHash(b); err != nil || !okHash {
			okAll = false
			fmt.Fprintf(os.Stderr, "features.json hash mismatch: computed=%s stored=%s err=%v\n", comp, stored, err)
		}
		if err := validate.ValidateRaw("features.json", b); err != nil {
			okAll = false
			fmt.Fprintf(os.Stderr, "features.json schema: %v\n", err)
		} else {
			fmt.Println("OK features.json")
		}
	}

	// tasks.json
	if b, ok := read("tasks.json"); ok {
		if comp, stored, okHash, err := validate.CheckArtifactHash(b); err != nil || !okHash {
			okAll = false
			fmt.Fprintf(os.Stderr, "tasks.json hash mismatch: computed=%s stored=%s err=%v\n", comp, stored, err)
		}
		var tf m.TasksFile
		if err := json.Unmarshal(b, &tf); err != nil {
			okAll = false
			fmt.Fprintf(os.Stderr, "tasks.json parse: %v\n", err)
		} else if err := validate.TasksFile(&tf); err != nil {
			okAll = false
			fmt.Fprintf(os.Stderr, "tasks.json: %v\n", err)
		} else {
			fmt.Println("OK tasks.json")
		}
	}

	// dag.json
	if b, ok := read("dag.json"); ok {
		if comp, stored, okHash, err := validate.CheckArtifactHash(b); err != nil || !okHash {
			okAll = false
			fmt.Fprintf(os.Stderr, "dag.json hash mismatch: computed=%s stored=%s err=%v\n", comp, stored, err)
		}
		var df m.DagFile
		if err := json.Unmarshal(b, &df); err != nil {
			okAll = false
			fmt.Fprintf(os.Stderr, "dag.json parse: %v\n", err)
		} else if err := validate.DagFile(&df); err != nil {
			okAll = false
			fmt.Fprintf(os.Stderr, "dag.json: %v\n", err)
		} else {
			fmt.Println("OK dag.json")
		}
	}

	// waves.json
	if b, ok := read("waves.json"); ok {
		if comp, stored, okHash, err := validate.CheckArtifactHash(b); err != nil || !okHash {
			okAll = false
			fmt.Fprintf(os.Stderr, "waves.json hash mismatch: computed=%s stored=%s err=%v\n", comp, stored, err)
		}
		if err := validate.ValidateRaw("waves.json", b); err != nil {
			okAll = false
			fmt.Fprintf(os.Stderr, "waves.json schema: %v\n", err)
		} else {
			fmt.Println("OK waves.json")
		}
	}

	// coordinator.json
	if b, ok := read("coordinator.json"); ok {
		if err := validate.ValidateRaw("coordinator.json", b); err != nil {
			okAll = false
			fmt.Fprintf(os.Stderr, "coordinator.json schema: %v\n", err)
		} else {
			fmt.Println("OK coordinator.json")
		}
	}

	if !okAll {
		os.Exit(2)
	}
	fmt.Println("All artifacts valid.")
}

// analysisPkg wraps the codebase census but returns a compact map for embedding.
func analysisPkg(path string) (analysis.FileCensusCounts, error) {
	a, err := analysis.RunCensus(path)
	if err != nil {
		return analysis.FileCensusCounts{}, err
	}
	return analysis.FileCensusCounts{Files: len(a.Files), GoFiles: len(a.GoFiles)}, nil
}
