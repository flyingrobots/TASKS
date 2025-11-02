package plan

import (
    "context"
    "encoding/hex"
    "encoding/json"
    "errors"
    "fmt"
    "math"
    "math/rand"
    "sort"
    "strings"

    analysis "github.com/james/tasks-planner/internal/analysis"
    "github.com/james/tasks-planner/internal/canonjson"
    "github.com/james/tasks-planner/internal/hash"
    m "github.com/james/tasks-planner/internal/model"
    "github.com/james/tasks-planner/internal/validators"
)

const (
    defaultMinConfidence = 0.7
    schemaVersion        = "v9" // breaking contract changes: camelCase tags across artifacts
)

// FeatureSummary represents a lightweight feature descriptor produced by the spec loader.
type FeatureSummary struct {
	ID    string
	Title string
}

// TasksResult captures outcomes from the spec/doc loader.
type TasksResult struct {
	Tasks             []m.Task
	Features          []FeatureSummary
	Dependencies      []m.Edge
	ResourceConflicts map[string]any
	DocProvided       bool
}

// ArtifactBundle bundles planner artifacts for writing via the artifact port.
type ArtifactBundle struct {
	TasksFile        *m.TasksFile
	DagFile          *m.DagFile
	Coordinator      *m.Coordinator
	Features         *m.FeaturesArtifact
	Waves            *m.WavesArtifact
	Titles           *m.TitlesArtifact
	ValidatorReports []m.ValidatorReport
}

// ArtifactWriteResult summarizes the outcome of the artifact writer.
type ArtifactWriteResult struct {
	Hashes map[string]string
}

// ValidatorRunner mirrors the validator adapter contract required by the service.
type ValidatorRunner interface {
	Run(ctx context.Context, payload validators.Payload) ([]validators.Report, error)
}

// Service orchestrates the planning workflow via injected adapters/ports.
type Service struct {
	BuildTasks         func(ctx context.Context, docPath string) (TasksResult, error)
	AnalyzeRepo        func(ctx context.Context, repo string) (analysis.FileCensusCounts, error)
	ResolveDeps        func(tasks []m.Task, docEdges []m.Edge) ([]m.Edge, map[string]any)
	BuildDAG           func(ctx context.Context, tasks []m.Task, deps []m.Edge, minConfidence float64) (*m.DagFile, error)
	BuildCoordinator   func(tasks []m.Task, deps []m.Edge) m.Coordinator
	ValidateTasks      func(tf *m.TasksFile) error
	ValidateDAG        func(df *m.DagFile) error
	BuildWaves         func(ctx context.Context, df *m.DagFile, tasks []m.Task) (*m.WavesArtifact, error)
	WriteArtifacts     func(ctx context.Context, out string, bundle ArtifactBundle) (ArtifactWriteResult, error)
	NewValidatorRunner func(cfg validators.Config) (ValidatorRunner, error)
}

// Request describes a plan invocation.
type Request struct {
	DocPath          string
	RepoPath         string
	OutDir           string
	MinConfidence    *float64
	ValidatorConfig  validators.Config
	StrictValidators bool
}

// Result summarizes planner execution output.
type Result struct {
	ArtifactHashes   map[string]string
	ValidatorReports []m.ValidatorReport
	Warnings         []string
}

// Plan executes the planning workflow.
func (s Service) Plan(ctx context.Context, req Request) (Result, error) {
	if s.BuildTasks == nil || s.AnalyzeRepo == nil || s.BuildDAG == nil || s.ValidateTasks == nil || s.ValidateDAG == nil || s.BuildWaves == nil || s.WriteArtifacts == nil {
		return Result{}, errors.New("plan service: missing required adapters")
	}

	tasksRes, err := s.BuildTasks(ctx, req.DocPath)
	if err != nil {
		return Result{}, fmt.Errorf("load tasks: %w", err)
	}

	tf := &m.TasksFile{}
	tf.Meta.Version = schemaVersion
	switch {
	case req.MinConfidence == nil:
		tf.Meta.MinConfidence = defaultMinConfidence
	case *req.MinConfidence < 0 || *req.MinConfidence > 1:
		return Result{}, fmt.Errorf("minConfidence out of range [0,1]: %v", *req.MinConfidence)
	default:
		tf.Meta.MinConfidence = *req.MinConfidence
	}
	census, err := s.AnalyzeRepo(ctx, req.RepoPath)
	if err != nil {
		return Result{}, fmt.Errorf("analysis: %w", err)
	}
	tf.Meta.CodebaseAnalysis = census
	tf.Meta.Autonormalization.Split = []string{}
	tf.Meta.Autonormalization.Merged = []string{}
	tf.Tasks = tasksRes.Tasks
	if s.ResolveDeps != nil {
		deps, conflicts := s.ResolveDeps(tasksRes.Tasks, tasksRes.Dependencies)
		tf.Dependencies = deps
		tf.ResourceConflicts = conflicts
	} else {
		tf.Dependencies = tasksRes.Dependencies
		tf.ResourceConflicts = tasksRes.ResourceConflicts
	}

	if tasksRes.DocProvided {
		var missing []string
		for _, task := range tf.Tasks {
			if len(task.AcceptanceChecks) == 0 {
				missing = append(missing, task.ID)
			}
		}
		if len(missing) > 0 {
			return Result{}, fmt.Errorf("missing acceptance checks for tasks: %s", strings.Join(missing, ", "))
		}
	}

	if err := s.ValidateTasks(tf); err != nil {
		return Result{}, fmt.Errorf("validate tasks: %w", err)
	}

	dagFile, err := s.BuildDAG(ctx, tf.Tasks, tf.Dependencies, tf.Meta.MinConfidence)
	if err != nil {
		return Result{}, fmt.Errorf("build dag: %w", err)
	}
	if err := s.ValidateDAG(dagFile); err != nil {
		return Result{}, fmt.Errorf("validate dag: %w", err)
	}

	waves, err := s.BuildWaves(ctx, dagFile, tf.Tasks)
	if err != nil {
		return Result{}, fmt.Errorf("build waves: %w", err)
	}

	features := tasksRes.Features
	if len(features) == 0 {
		features = featuresFromTasks(tf.Tasks)
	}

	titles := taskTitles(tf.Tasks)

	var coord m.Coordinator
	if s.BuildCoordinator != nil {
		coord = s.BuildCoordinator(tf.Tasks, tf.Dependencies)
	} else {
		coord = makeCoordinator(tf.Tasks, tf.Dependencies)
	}

    // Populate coordinator metrics from the DAG metrics, with Monte Carlo p50 for total hours.
    if dagFile != nil {
        coord.Metrics.Estimates.LongestPathLength = dagFile.Metrics.LongestPathLength
        coord.Metrics.Estimates.WidthApprox = dagFile.Metrics.WidthApprox

        // Deterministic seed derived from the canonical tasks preimage
        // Ensure we don't mutate tf while hashing
        preimage, err := json.Marshal(tf)
        if err != nil {
            return Result{}, fmt.Errorf("marshal tasks preimage: %w", err)
        }
        can, err := canonjson.ToCanonicalJSON(preimage)
        if err != nil {
            return Result{}, fmt.Errorf("canonicalize tasks preimage: %w", err)
        }
        h := hash.HashCanonicalBytes(can) // hex string
        // Use first 8 bytes for a 64-bit seed
        var seed int64
        if b, err := hex.DecodeString(h[:16]); err == nil && len(b) == 8 {
            for i := 0; i < 8; i++ { seed = (seed << 8) | int64(b[i]) }
        } else {
            // Fallback deterministic seed
            seed = 42
        }

        // Compute Monte Carlo median for plan makespan (hours)
        // Samples scaled modestly by DAG size but capped for CI friendliness
        n := len(tf.Tasks)
        samples := 2000
        if n > 100 { samples = 1000 } // basic cap for very large DAGs
        p50 := mcP50MakespanHours(tf.Tasks, edgesFromDag(dagFile.Edges), samples, seed)
        if !math.IsNaN(p50) && !math.IsInf(p50, 0) {
            // Round to 6 decimal places to avoid excessive float noise in artifacts
            p50 = math.Round(p50*1e6) / 1e6
            coord.Metrics.Estimates.P50TotalHours = p50
        }
    }

	var validatorReports []m.ValidatorReport
	var warnings []string
	if s.NewValidatorRunner != nil && validatorConfigured(req.ValidatorConfig) {
		runner, err := s.NewValidatorRunner(req.ValidatorConfig)
		if err != nil {
			return Result{}, fmt.Errorf("validator runner: %w", err)
		}
		payload := validators.Payload{Tasks: tf, Dag: dagFile, Coordinator: &coord}
		reports, runErr := runner.Run(ctx, payload)
		modelReports := convertValidatorReports(reports)
		validatorReports = modelReports
		if len(modelReports) > 0 {
			tf.Meta.ValidatorReports = modelReports
		}

		var failedReports []m.ValidatorReport
		for _, rep := range modelReports {
			status := strings.ToLower(strings.TrimSpace(rep.Status))
			if status == m.ValidatorStatusFail || status == m.ValidatorStatusError || status == "failed" {
				failedReports = append(failedReports, rep)
			}
		}
		if req.StrictValidators {
			if len(failedReports) > 0 {
				names := make([]string, 0, len(failedReports))
				for _, rep := range failedReports {
					names = append(names, rep.Name)
				}
				return Result{}, fmt.Errorf("validators failed: %s", strings.Join(names, ", "))
			}
			if runErr != nil {
				return Result{}, fmt.Errorf("validators: %w", runErr)
			}
		} else {
			if runErr != nil {
				warnings = append(warnings, runErr.Error())
			}
			for _, rep := range failedReports {
				msg := fmt.Sprintf("validator %s reported %s", rep.Name, rep.Status)
				if detail := strings.TrimSpace(rep.Detail); detail != "" {
					msg += fmt.Sprintf(" — %s", detail)
				}
				warnings = append(warnings, msg)
			}
		}
	}

	artifactBundle := ArtifactBundle{
		TasksFile:        tf,
		DagFile:          dagFile,
		Coordinator:      &coord,
		Features:         makeFeaturesArtifact(features),
		Waves:            waves,
		Titles:           titles,
		ValidatorReports: validatorReports,
	}

	writeResult, err := s.WriteArtifacts(ctx, req.OutDir, artifactBundle)
	if err != nil {
		return Result{}, fmt.Errorf("write artifacts: %w", err)
	}

	result := Result{
		ArtifactHashes:   writeResult.Hashes,
		ValidatorReports: tf.Meta.ValidatorReports,
		Warnings:         warnings,
	}
	return result, nil
}

func validatorConfigured(cfg validators.Config) bool {
	return cfg.AcceptanceCmd != "" || cfg.EvidenceCmd != "" || cfg.InterfaceCmd != ""
}

func featuresFromTasks(tasks []m.Task) []FeatureSummary {
	seen := map[string]bool{}
	summaries := []FeatureSummary{}
	for _, t := range tasks {
		if t.FeatureID == "" {
			continue
		}
		if seen[t.FeatureID] {
			continue
		}
		summaries = append(summaries, FeatureSummary{ID: t.FeatureID, Title: t.FeatureID})
		seen[t.FeatureID] = true
	}
	sort.Slice(summaries, func(i, j int) bool { return summaries[i].ID < summaries[j].ID })
	return summaries
}

func makeFeaturesArtifact(features []FeatureSummary) *m.FeaturesArtifact {
	entries := make([]m.FeatureEntry, 0, len(features))
	for _, f := range features {
		entries = append(entries, m.FeatureEntry{ID: f.ID, Title: f.Title})
	}
	return &m.FeaturesArtifact{
		Meta:     m.ArtifactMeta{Version: schemaVersion, ArtifactHash: ""},
		Features: entries,
	}
}

func makeCoordinator(tasks []m.Task, deps []m.Edge) m.Coordinator {
    coord := m.Coordinator{}
    setCoordinatorVersion(&coord, schemaVersion)
    coord.Graph.Nodes = tasks
    coord.Graph.Edges = deps
    coord.Config.Resources.Catalog = map[string]m.ResourceSpec{}
    coord.Config.Resources.Profiles = map[string]map[string]int{"default": {}}
    coord.Config.Policies.LockOrdering = []string{}
    return coord
}

func convertValidatorReports(src []validators.Report) []m.ValidatorReport {
	if len(src) == 0 {
		return nil
	}
	out := make([]m.ValidatorReport, 0, len(src))
	for _, rep := range src {
		detail := strings.TrimSpace(rep.Detail)
		out = append(out, m.ValidatorReport{
			Name:      rep.Name,
			Status:    rep.Status,
			Command:   rep.Command,
			InputHash: rep.InputHash,
			Cached:    rep.Cached,
			Detail:    detail,
			RawOutput: rep.RawOutput,
		})
	}
	return out
}

// mcP50MakespanHours computes a Monte Carlo estimate of the plan makespan (median in hours)
// respecting precedence edges only. It samples each task's duration from a triangular
// distribution parameterized by (Optimistic, MostLikely, Pessimistic). The result is
// deterministic for a given seed and input graph.
func mcP50MakespanHours(tasks []m.Task, edges []m.Edge, samples int, seed int64) float64 {
    if samples <= 0 || len(tasks) == 0 {
        return 0
    }
    // Map task IDs to indices for arrays
    idx := make(map[string]int, len(tasks))
    for i, t := range tasks { idx[t.ID] = i }
    // Build predecessors and topological order (Kahn)
    pred := make([][]int, len(tasks))
    indeg := make([]int, len(tasks))
    for _, e := range edges {
        u, okU := idx[e.From]
        v, okV := idx[e.To]
        if !okU || !okV { continue }
        pred[v] = append(pred[v], u)
        indeg[v]++
    }
    order := make([]int, 0, len(tasks))
    q := make([]int, 0, len(tasks))
    // init q with zero indegree nodes
    for i := range tasks { if indeg[i] == 0 { q = append(q, i) } }
    for len(q) > 0 {
        v := q[0]; q = q[1:]
        order = append(order, v)
        // decrease indegree of successors
        for w := range tasks {
            // scan pred list of w for v (cheap for small DAGs; OK here)
            for _, p := range pred[w] { if p == v { indeg[w]--; if indeg[w] == 0 { q = append(q, w) } ; break } }
        }
    }
    if len(order) != len(tasks) {
        // cycle or missing edges mapping; fall back to simple sum of ML on critical path-like guess
        var sum float64
        for _, t := range tasks { sum += t.Duration.MostLikely }
        return sum
    }

    rng := rand.New(rand.NewSource(seed))
    draws := make([]float64, samples)
    dur := make([]float64, len(tasks))
    ef := make([]float64, len(tasks))
    for s := 0; s < samples; s++ {
        // sample durations
        for i, t := range tasks {
            a := t.Duration.Optimistic
            m := t.Duration.MostLikely
            b := t.Duration.Pessimistic
            if b < a { a, b = b, a }
            if m < a { m = a } else if m > b { m = b }
            dur[i] = triangular(a, m, b, rng)
        }
        // forward pass for earliest finish
        for _, v := range order {
            maxPred := 0.0
            for _, p := range pred[v] { if ef[p] > maxPred { maxPred = ef[p] } }
            ef[v] = maxPred + dur[v]
        }
        // makespan is max ef
        maxEF := 0.0
        for _, v := range order { if ef[v] > maxEF { maxEF = ef[v] } }
        draws[s] = maxEF
    }
    sort.Float64s(draws)
    mid := samples / 2
    if samples%2 == 1 { return draws[mid] }
    return 0.5 * (draws[mid-1] + draws[mid])
}

func triangular(a, m, b float64, rng *rand.Rand) float64 {
    if a == b { return a }
    u := rng.Float64()
    c := 0.0
    if b > a { c = (m - a) / (b - a) }
    if u < c {
        return a + math.Sqrt(u*(b-a)*(m-a))
    }
    return b - math.Sqrt((1-u)*(b-a)*(b-m))
}

func edgesFromDag(in []m.DagEdge) []m.Edge {
    if len(in) == 0 { return nil }
    out := make([]m.Edge, 0, len(in))
    for _, e := range in {
        out = append(out, m.Edge{From: e.From, To: e.To, Type: e.Type, IsHard: true, Confidence: 1})
    }
    return out
}
