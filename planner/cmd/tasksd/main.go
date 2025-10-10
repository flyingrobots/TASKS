package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	analysis "github.com/james/tasks-planner/internal/analysis"
	"github.com/james/tasks-planner/internal/app/plan"
	"github.com/james/tasks-planner/internal/canonjson"
	"github.com/james/tasks-planner/internal/export/dot"
	"github.com/james/tasks-planner/internal/hash"
	m "github.com/james/tasks-planner/internal/model"
	dagbuild "github.com/james/tasks-planner/internal/planner/dag"
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

	docLoader := plan.NewMarkdownDocLoader()
	analyzer := plan.CensusAnalyzer{}

	artifactWriter := plan.FileArtifactWriter{}

	svc := plan.Service{
		BuildTasks: func(ctx context.Context, docPath string) (plan.TasksResult, error) {
			res, err := docLoader.Load(ctx, docPath)
			if err != nil {
				return plan.TasksResult{}, err
			}
			deps, conflicts := inferDependencies(res.Tasks, res.Dependencies)
			res.Dependencies = deps
			res.ResourceConflicts = conflicts
			return res, nil
		},
		AnalyzeRepo: func(ctx context.Context, repo string) (analysis.FileCensusCounts, error) {
			if repo == "" {
				return analysis.FileCensusCounts{}, nil
			}
			counts, err := analyzer.Analyze(ctx, repo)
			if err != nil {
				log.Printf("analysis failed for repo %s: %v", repo, err)
				return analysis.FileCensusCounts{}, nil
			}
			return counts, nil
		},
		BuildDAG: func(ctx context.Context, tasks []m.Task, deps []m.Edge, minConfidence float64) (m.DagFile, error) {
			return dagbuild.Build(tasks, deps, minConfidence)
		},
		ValidateTasks: validate.TasksFile,
		ValidateDag:   validate.DagFile,
		BuildWaves: func(ctx context.Context, df m.DagFile, tasks []m.Task) (map[string]any, error) {
			return buildWaves(df, tasks)
		},
		WriteArtifacts: artifactWriter.Write,
		NewValidatorRunner: func(cfg validators.Config) (plan.ValidatorRunner, error) {
			return validators.NewRunner(cfg)
		},
	}

	req := plan.Request{
		DocPath:       *doc,
		RepoPath:      *repo,
		OutDir:        *out,
		MinConfidence: 0.7,
		ValidatorConfig: validators.Config{
			AcceptanceCmd: *acceptanceCmd,
			EvidenceCmd:   *evidenceCmd,
			InterfaceCmd:  *interfaceCmd,
			CacheDir:      *validatorsCache,
			Timeout:       *validatorsTimeout,
		},
		StrictValidators: *validatorsStrict,
	}

	res, err := svc.Plan(context.Background(), req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	for _, warn := range res.Warnings {
		fmt.Fprintf(os.Stderr, "validators warning: %s\n", warn)
	}

	fmt.Println("Plan stub written to", *out)
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

const validatorDetailLimit = 2048

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
