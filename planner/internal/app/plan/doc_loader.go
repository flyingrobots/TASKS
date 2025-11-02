package plan

import (
    "context"
    "fmt"
    "os"
    "strings"

    m "github.com/james/tasks-planner/internal/model"
    docp "github.com/james/tasks-planner/internal/planner/docparse"
)

// Canonical task titles used by stub plans and defaults. Keep verb-first.
const (
    TaskTitleSetupDB              = "Setup DB"
    TaskTitleMigrateSchema        = "Migrate Schema"
    TaskTitleImplementAPIHandlers = "Implement API Handlers"
)

// DocLoader loads tasks/features from a specification document.
type DocLoader interface {
	Load(ctx context.Context, docPath string) (TasksResult, error)
}

// MarkdownDocLoader loads markdown specs from the filesystem.
type MarkdownDocLoader struct {
	ReadFile func(string) ([]byte, error)
}

// NewMarkdownDocLoader creates a loader using os.ReadFile.
func NewMarkdownDocLoader() MarkdownDocLoader {
	return MarkdownDocLoader{ReadFile: os.ReadFile}
}

func (l MarkdownDocLoader) Load(ctx context.Context, docPath string) (TasksResult, error) {
	if err := ctx.Err(); err != nil {
		return TasksResult{}, err
	}
	if docPath == "" {
		tasks, features := stubPlan()
		return TasksResult{Tasks: tasks, Features: features, DocProvided: false}, nil
	}
	if err := ctx.Err(); err != nil {
		return TasksResult{}, err
	}
	info, err := os.Stat(docPath)
	if err != nil {
		if os.IsNotExist(err) {
			tasks, features := stubPlan()
			return TasksResult{Tasks: tasks, Features: features, DocProvided: false}, nil
		}
		return TasksResult{}, fmt.Errorf("stat --doc: %w", err)
	}
	if info.IsDir() {
		return TasksResult{}, fmt.Errorf("--doc points to a directory: %s", docPath)
	}
	if err := ctx.Err(); err != nil {
		return TasksResult{}, err
	}
	raw, err := l.read(ctx, docPath)
	if err != nil {
		return TasksResult{}, fmt.Errorf("read --doc: %w", err)
	}
	feats, specs := docp.ParseMarkdown(string(raw))
	if len(feats) == 0 && len(specs) == 0 {
		tasks, features := stubPlan()
		return TasksResult{Tasks: tasks, Features: features, DocProvided: false}, nil
	}
	features := make([]FeatureSummary, 0, len(feats))
	for _, f := range feats {
		features = append(features, FeatureSummary{ID: f.ID, Title: f.Title})
	}
	tasks := make([]m.Task, 0, len(specs))
	titleToID := map[string]string{}
	var parseErrors []string
	for i, spec := range specs {
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
		key := normalizeKey(spec.Title)
		if prev, ok := titleToID[key]; ok && prev != id {
			parseErrors = append(parseErrors, fmt.Sprintf("duplicate task title %q (IDs %s and %s) — use explicit IDs in 'after:'", spec.Title, prev, id))
		} else {
			titleToID[key] = id
		}
	}
	if len(parseErrors) > 0 {
		return TasksResult{}, fmt.Errorf("doc parse errors: %s", strings.Join(parseErrors, "; "))
	}

	edges := []m.Edge{}
	for _, spec := range specs {
		toID := titleToID[normalizeKey(spec.Title)]
		if toID == "" {
			continue
		}
		for _, raw := range spec.After {
			fromID := resolveTaskID(raw, titleToID)
			if fromID == "" {
				continue
			}
			edges = append(edges, m.Edge{From: fromID, To: toID, Type: "sequential", IsHard: true, Confidence: 1})
		}
	}
	if len(features) == 0 {
		features = featuresFromTasks(tasks)
	}
	return TasksResult{Tasks: tasks, Features: features, Dependencies: edges, DocProvided: true}, nil
}

func (l MarkdownDocLoader) read(ctx context.Context, path string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var (
		data []byte
		err  error
	)
	if l.ReadFile != nil {
		data, err = l.ReadFile(path)
	} else {
		data, err = os.ReadFile(path)
	}
	if err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return data, nil
}

func stubPlan() ([]m.Task, []FeatureSummary) {
    base := []struct {
        id        string
        featureID string
        title     string
    }{
        {"T001", "F001", TaskTitleSetupDB},
        {"T002", "F001", TaskTitleMigrateSchema},
        {"T003", "F001", TaskTitleImplementAPIHandlers},
    }
	tasks := make([]m.Task, 0, len(base))
    for _, spec := range base {
        task := m.Task{ID: spec.id, FeatureID: spec.featureID, Title: spec.title}
        switch strings.ToLower(spec.title) {
        case strings.ToLower(TaskTitleSetupDB):
            task.Duration = m.DurationPERT{Optimistic: 1, MostLikely: 2.5, Pessimistic: 4}
        case strings.ToLower(TaskTitleMigrateSchema):
            task.Duration = m.DurationPERT{Optimistic: 1, MostLikely: 3.5, Pessimistic: 8}
        case strings.ToLower(TaskTitleImplementAPIHandlers):
            task.Duration = m.DurationPERT{Optimistic: 2, MostLikely: 6, Pessimistic: 12}
        default:
            task.Duration = m.DurationPERT{Optimistic: 1, MostLikely: 2, Pessimistic: 3}
        }
        applyTaskDefaults(&task)
        tasks = append(tasks, task)
    }
    return tasks, []FeatureSummary{{ID: "F001", Title: "Core DB + API"}}
}

func applyTaskDefaults(task *m.Task) {
    applyAcceptanceDefaults(task)
    if task.DurationUnit == "" { task.DurationUnit = "hours" }
    applyLoggingDefaults(task)
    seedEvidenceDefaults(task)
}

func applyAcceptanceDefaults(task *m.Task) {
    if len(task.AcceptanceChecks) != 0 { return }
    lower := strings.ToLower(task.Title)
    switch {
    case strings.Contains(lower, "setup db"):
        task.AcceptanceChecks = []m.AcceptanceCheck{{
            Type:    "command",
            Cmd:     `sh -c 'test -n "$DB_DSN" && psql "$DB_DSN" -c "\\conninfo" >/dev/null 2>&1'`,
            Timeout: 10,
        }}
    case strings.Contains(lower, "migrate schema"):
        task.AcceptanceChecks = []m.AcceptanceCheck{{
            Type:    "command",
            Cmd:     `sh -c 'command -v db/migrate >/dev/null 2>&1 && db/migrate status | grep -q Applied'`,
            Timeout: 15,
        }}
    case strings.Contains(lower, "api handler"):
        task.AcceptanceChecks = []m.AcceptanceCheck{{
            Type:    "command",
            Cmd:     `sh -c 'curl -fsS "${API_BASE-http://localhost:8080}/healthz" >/dev/null'`,
            Timeout: 10,
        }}
    default:
        // Force authors to provide an acceptance by failing fast.
        task.AcceptanceChecks = []m.AcceptanceCheck{{Type: "command", Cmd: `sh -c 'echo "missing acceptance checks" >&2; exit 1'`, Timeout: 5}}
    }
}

func applyLoggingDefaults(task *m.Task) {
    wasEmpty := task.ExecutionLogging.Format == ""
    baseFields := []string{"timestamp", "task_id", "step", "status", "message", "error", "error_type", "error_details"}
    if len(task.ExecutionLogging.RequiredFields) == 0 {
        task.ExecutionLogging.RequiredFields = append([]string{}, baseFields...)
    }
    ensureField := func(field string) {
        for _, f := range task.ExecutionLogging.RequiredFields { if f == field { return } }
        task.ExecutionLogging.RequiredFields = append(task.ExecutionLogging.RequiredFields, field)
    }
    lower := strings.ToLower(task.Title)
    switch {
    case strings.Contains(lower, "setup db"):
        ensureField("db_response_time")
    case strings.Contains(lower, "migrate schema"):
        ensureField("migration_version")
    case strings.Contains(lower, "api handler"):
        if wasEmpty { task.ExecutionLogging.Format = "JSON" }
        ensureField("service_version")
    }
    if wasEmpty && task.ExecutionLogging.Format == "" {
        task.ExecutionLogging.Format = "JSONL"
    }
}

func seedEvidenceDefaults(task *m.Task) {
    if len(task.Evidence) != 0 { return }
    lower := strings.ToLower(task.Title)
    switch {
    case strings.Contains(lower, "setup db"):
        task.Evidence = append(task.Evidence, m.Evidence{
            Type:       "code_analysis",
            Source:     "planner/internal/app/plan/doc_loader.go#applyTaskDefaults",
            Excerpt:    "DB setup defaults: acceptance check via psql; logging requires db_response_time",
            Confidence: 0.9,
            Rationale:  "stub default grounded in code defaults",
        })
    case strings.Contains(lower, "migrate schema"):
        task.Evidence = append(task.Evidence, m.Evidence{
            Type:       "docs",
            Source:     "docs/go-architecture.md",
            Excerpt:    "Migrations modeled as technical edges with ordering",
            Confidence: 0.85,
            Rationale:  "ordering requirement documented in architecture",
        })
    case strings.Contains(lower, "api handler"):
        task.Evidence = append(task.Evidence, m.Evidence{
            Type:       "docs",
            Source:     "README.md",
            Excerpt:    "CLI/demo health endpoint used for acceptance (/healthz)",
            Confidence: 0.8,
            Rationale:  "acceptance derived from documented demo",
        })
    default:
        task.Evidence = append(task.Evidence, m.Evidence{
            Type:       "docs",
            Source:     "docs/v8/v8.md",
            Excerpt:    "Tasks require machine-verifiable acceptance and traceable evidence",
            Confidence: 0.75,
            Rationale:  "fallback evidence for stub plans",
        })
    }
}

func resolveTaskID(token string, titleToID map[string]string) string {
	trimmed := strings.TrimSpace(token)
	if trimmed == "" {
		return ""
	}
	if len(trimmed) > 1 && (trimmed[0] == 'T' || trimmed[0] == 't') {
		for _, r := range trimmed[1:] {
			if r < '0' || r > '9' {
				return ""
			}
		}
		return strings.ToUpper(trimmed)
	}
	return titleToID[normalizeKey(trimmed)]
}

func normalizeKey(v string) string {
	return strings.ToLower(strings.TrimSpace(v))
}
