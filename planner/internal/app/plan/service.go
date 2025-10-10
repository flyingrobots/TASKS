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

const defaultMinConfidence = 0.7

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
	Features         map[string]any
	Waves            map[string]any
	Titles           map[string]string
	ValidatorReports []validators.Report
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
	BuildDAG           func(ctx context.Context, tasks []m.Task, deps []m.Edge, minConfidence float64) (m.DagFile, error)
	ValidateTasks      func(tf *m.TasksFile) error
	ValidateDag        func(df *m.DagFile) error
	BuildWaves         func(ctx context.Context, df m.DagFile, tasks []m.Task) (map[string]any, error)
	WriteArtifacts     func(ctx context.Context, out string, bundle ArtifactBundle) (ArtifactWriteResult, error)
	NewValidatorRunner func(cfg validators.Config) (ValidatorRunner, error)
}

// Request describes a plan invocation.
type Request struct {
	DocPath          string
	RepoPath         string
	OutDir           string
	MinConfidence    float64
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
	if s.BuildTasks == nil || s.AnalyzeRepo == nil || s.BuildDAG == nil || s.ValidateTasks == nil || s.ValidateDag == nil || s.BuildWaves == nil || s.WriteArtifacts == nil {
		return Result{}, errors.New("plan service: missing required adapters")
	}

	tasksRes, err := s.BuildTasks(ctx, req.DocPath)
	if err != nil {
		return Result{}, fmt.Errorf("load tasks: %w", err)
	}

	tf := &m.TasksFile{}
	tf.Meta.Version = "v8"
	tf.Meta.MinConfidence = req.MinConfidence
	if tf.Meta.MinConfidence == 0 {
		tf.Meta.MinConfidence = defaultMinConfidence
	}
	census, err := s.AnalyzeRepo(ctx, req.RepoPath)
	if err != nil {
		return Result{}, fmt.Errorf("analysis: %w", err)
	}
	tf.Meta.CodebaseAnalysis = census
	tf.Meta.Autonormalization.Split = []string{}
	tf.Meta.Autonormalization.Merged = []string{}
	tf.Tasks = tasksRes.Tasks
	tf.Dependencies = tasksRes.Dependencies
	tf.ResourceConflicts = tasksRes.ResourceConflicts

	if tasksRes.DocProvided {
		for _, task := range tf.Tasks {
			if len(task.AcceptanceChecks) == 0 {
				return Result{}, fmt.Errorf("task %s missing acceptance checks", task.ID)
			}
		}
	}

	if err := s.ValidateTasks(tf); err != nil {
		return Result{}, fmt.Errorf("validate tasks: %w", err)
	}

	dagFile, err := s.BuildDAG(ctx, tf.Tasks, tf.Dependencies, tf.Meta.MinConfidence)
	if err != nil {
		return Result{}, fmt.Errorf("build dag: %w", err)
	}
	if err := s.ValidateDag(&dagFile); err != nil {
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

	coord := makeCoordinator(tf.Tasks, tf.Dependencies)

	var validatorReports []validators.Report
	var warnings []string
	if s.NewValidatorRunner != nil && validatorConfigured(req.ValidatorConfig) {
		runner, err := s.NewValidatorRunner(req.ValidatorConfig)
		if err != nil {
			return Result{}, fmt.Errorf("validator runner: %w", err)
		}
		payload := validators.Payload{Tasks: tf, Dag: &dagFile, Coordinator: &coord}
		reports, runErr := runner.Run(ctx, payload)
		validatorReports = reports
		modelReports := convertValidatorReports(reports)
		if len(modelReports) > 0 {
			tf.Meta.ValidatorReports = modelReports
		}
		if runErr != nil {
			if req.StrictValidators {
				return Result{}, fmt.Errorf("validators: %w", runErr)
			}
			warnings = append(warnings, runErr.Error())
		}
	}

	artifactBundle := ArtifactBundle{
		TasksFile:        tf,
		DagFile:          &dagFile,
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

func makeFeaturesArtifact(features []FeatureSummary) map[string]any {
	entries := make([]any, 0, len(features))
	for _, f := range features {
		entries = append(entries, map[string]any{"id": f.ID, "title": f.Title})
	}
	return map[string]any{
		"meta":     map[string]any{"version": "v8", "artifact_hash": ""},
		"features": entries,
	}
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
