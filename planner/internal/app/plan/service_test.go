package plan_test

import (
    "context"
    "errors"
    "testing"

    analysis "github.com/james/tasks-planner/internal/analysis"
    "github.com/james/tasks-planner/internal/app/plan"
    m "github.com/james/tasks-planner/internal/model"
    "github.com/james/tasks-planner/internal/validators"
)

type stubRunner struct {
    reports []validators.Report
    err     error
    ran     bool
}

func (s *stubRunner) Run(ctx context.Context, payload validators.Payload) ([]validators.Report, error) {
    s.ran = true
    return s.reports, s.err
}

func TestServicePlanSuccess(t *testing.T) {
    tasks := []m.Task{{ID: "T001", Title: "Do thing", FeatureID: "F1", AcceptanceChecks: []m.AcceptanceCheck{{Type: "command", Cmd: "echo ok"}}, Duration: m.DurationPERT{Optimistic: 1, MostLikely: 2, Pessimistic: 3}, DurationUnit: "hours"}}
    features := []plan.FeatureSummary{{ID: "F1", Title: "Feature One"}}
    deps := []m.Edge{}

    runner := &stubRunner{reports: []validators.Report{{Name: "acceptance", Status: m.ValidatorStatusPass, Command: "echo pass", InputHash: "abcd", Detail: "ok"}}}

    svc := plan.Service{
        BuildTasks: func(ctx context.Context, docPath string) (plan.TasksResult, error) {
            if docPath != "spec.md" {
                t.Fatalf("expected doc path spec.md, got %s", docPath)
            }
            return plan.TasksResult{Tasks: tasks, Features: features, Dependencies: deps, DocProvided: true}, nil
        },
        AnalyzeRepo: func(ctx context.Context, repo string) (analysis.FileCensusCounts, error) {
            if repo != "./repo" {
                t.Fatalf("unexpected repo path: %s", repo)
            }
            return analysis.FileCensusCounts{Files: 10, GoFiles: 5}, nil
        },
        BuildDAG: func(ctx context.Context, tasks []m.Task, deps []m.Edge, minConfidence float64) (m.DagFile, error) {
            df := m.DagFile{}
            df.Meta.Version = "v8"
            df.Nodes = []struct {
                ID                  string `json:"id"`
                Depth               int    `json:"depth"`
                CriticalPath        bool   `json:"critical_path"`
                ParallelOpportunity int    `json:"parallel_opportunity"`
            }{{ID: "T001", Depth: 0}}
            return df, nil
        },
        ValidateTasks: func(tf *m.TasksFile) error { return nil },
        ValidateDag:   func(df *m.DagFile) error { return nil },
        BuildWaves: func(ctx context.Context, df m.DagFile, tasks []m.Task) (map[string]any, error) {
            return map[string]any{"meta": map[string]any{"version": "v8"}, "waves": [][]string{{"T001"}}}, nil
        },
        WriteArtifacts: func(ctx context.Context, out string, bundle plan.ArtifactBundle) (plan.ArtifactWriteResult, error) {
            if out != "./plans" {
                t.Fatalf("unexpected out dir: %s", out)
            }
            if len(bundle.ValidatorReports) != 1 {
                t.Fatalf("expected validator report in bundle")
            }
            return plan.ArtifactWriteResult{Hashes: map[string]string{"tasks.json": "hash1"}}, nil
        },
        NewValidatorRunner: func(cfg validators.Config) (plan.ValidatorRunner, error) {
            if cfg.AcceptanceCmd != "accept" {
                t.Fatalf("unexpected validator config: %+v", cfg)
            }
            return runner, nil
        },
    }

    req := plan.Request{
        DocPath:          "spec.md",
        RepoPath:         "./repo",
        OutDir:           "./plans",
        MinConfidence:    0.5,
        ValidatorConfig:  validators.Config{AcceptanceCmd: "accept"},
        StrictValidators: false,
    }

    res, err := svc.Plan(context.Background(), req)
    if err != nil {
        t.Fatalf("plan: %v", err)
    }
    if !runner.ran {
        t.Fatalf("expected validator runner to execute")
    }
    if res.ArtifactHashes["tasks.json"] != "hash1" {
        t.Fatalf("unexpected artifact hashes: %+v", res.ArtifactHashes)
    }
    if len(res.ValidatorReports) != 1 || res.ValidatorReports[0].Status != m.ValidatorStatusPass {
        t.Fatalf("unexpected validator reports: %+v", res.ValidatorReports)
    }
    if len(res.Warnings) != 0 {
        t.Fatalf("expected no warnings, got %+v", res.Warnings)
    }
}

func TestServicePlanValidatorWarnings(t *testing.T) {
    tasks := []m.Task{{ID: "T001", Title: "Do thing", AcceptanceChecks: []m.AcceptanceCheck{{Type: "command", Cmd: "echo ok"}}}}

    runner := &stubRunner{
        reports: []validators.Report{{Name: "acceptance", Status: m.ValidatorStatusFail, Command: "accept", Detail: "broken"}},
        err:     errors.New("validator failed"),
    }

    svc := plan.Service{
        BuildTasks: func(ctx context.Context, docPath string) (plan.TasksResult, error) {
            return plan.TasksResult{Tasks: tasks, DocProvided: false}, nil
        },
        AnalyzeRepo: func(ctx context.Context, repo string) (analysis.FileCensusCounts, error) {
            return analysis.FileCensusCounts{}, nil
        },
        BuildDAG: func(ctx context.Context, tasks []m.Task, deps []m.Edge, minConfidence float64) (m.DagFile, error) {
            return m.DagFile{}, nil
        },
        ValidateTasks: func(tf *m.TasksFile) error { return nil },
        ValidateDag:   func(df *m.DagFile) error { return nil },
        BuildWaves: func(ctx context.Context, df m.DagFile, tasks []m.Task) (map[string]any, error) {
            return map[string]any{"meta": map[string]any{"version": "v8"}}, nil
        },
        WriteArtifacts: func(ctx context.Context, out string, bundle plan.ArtifactBundle) (plan.ArtifactWriteResult, error) {
            return plan.ArtifactWriteResult{Hashes: map[string]string{"tasks.json": "hash1"}}, nil
        },
        NewValidatorRunner: func(cfg validators.Config) (plan.ValidatorRunner, error) {
            return runner, nil
        },
    }

    req := plan.Request{
        DocPath:          "",
        RepoPath:         "",
        OutDir:           "./plans",
        ValidatorConfig:  validators.Config{AcceptanceCmd: "accept"},
        StrictValidators: false,
    }

    res, err := svc.Plan(context.Background(), req)
    if err != nil {
        t.Fatalf("plan: %v", err)
    }
    if len(res.ValidatorReports) != 1 || res.ValidatorReports[0].Status != m.ValidatorStatusFail {
        t.Fatalf("unexpected validator reports: %+v", res.ValidatorReports)
    }
    if len(res.Warnings) != 1 {
        t.Fatalf("expected warning, got %+v", res.Warnings)
    }
}

func TestServicePlanValidatorStrictFailure(t *testing.T) {
    svc := plan.Service{
        BuildTasks: func(ctx context.Context, docPath string) (plan.TasksResult, error) {
            return plan.TasksResult{Tasks: []m.Task{{ID: "T001", Title: "Do", AcceptanceChecks: []m.AcceptanceCheck{{Type: "command", Cmd: "echo ok"}}}}, DocProvided: false}, nil
        },
        AnalyzeRepo: func(context.Context, string) (analysis.FileCensusCounts, error) { return analysis.FileCensusCounts{}, nil },
        BuildDAG: func(context.Context, []m.Task, []m.Edge, float64) (m.DagFile, error) { return m.DagFile{}, nil },
        ValidateTasks: func(*m.TasksFile) error { return nil },
        ValidateDag:   func(*m.DagFile) error { return nil },
        BuildWaves: func(ctx context.Context, df m.DagFile, tasks []m.Task) (map[string]any, error) {
            return map[string]any{"meta": map[string]any{"version": "v8"}}, nil
        },
        WriteArtifacts: func(context.Context, string, plan.ArtifactBundle) (plan.ArtifactWriteResult, error) {
            return plan.ArtifactWriteResult{Hashes: map[string]string{"tasks.json": "hash"}}, nil
        },
        NewValidatorRunner: func(validators.Config) (plan.ValidatorRunner, error) {
            return &stubRunner{err: errors.New("boom")}, nil
        },
    }

    _, err := svc.Plan(context.Background(), plan.Request{
        OutDir:           "./plans",
        ValidatorConfig:  validators.Config{AcceptanceCmd: "accept"},
        StrictValidators: true,
    })
    if err == nil {
        t.Fatalf("expected error with strict validators")
    }
}
