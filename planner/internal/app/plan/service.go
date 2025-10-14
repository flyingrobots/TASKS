package plan

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	analysis "github.com/james/tasks-planner/internal/analysis"
	m "github.com/james/tasks-planner/internal/model"
	"github.com/james/tasks-planner/internal/validators"
)

const (
	defaultMinConfidence = 0.7
	schemaVersion        = "v8"
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
	coord.Version = schemaVersion
	coord.Graph.Nodes = tasks
	coord.Graph.Edges = deps
	coord.Config.Resources.Catalog = map[string]struct {
		Capacity  int    `json:"capacity"`
		Mode      string `json:"mode"`
		LockOrder int    `json:"lock_order"`
	}{}
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
