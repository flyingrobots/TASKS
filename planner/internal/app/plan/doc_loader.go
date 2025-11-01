package plan

import (
	"context"
	"fmt"
	"os"
	"strings"

	m "github.com/james/tasks-planner/internal/model"
	docp "github.com/james/tasks-planner/internal/planner/docparse"
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
		{"T001", "F001", "Setup DB"},
		{"T002", "F001", "Migrate Schema"},
		{"T003", "F001", "API Handlers"},
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
	return tasks, []FeatureSummary{{ID: "F001", Title: "Core DB + API"}}
}

func applyTaskDefaults(task *m.Task) {
    if len(task.AcceptanceChecks) == 0 {
        switch strings.ToLower(task.Title) {
        case "setup db":
            task.AcceptanceChecks = []m.AcceptanceCheck{{Type: "command", Cmd: `psql "$DB_DSN" -c "\\conninfo" >/dev/null 2>&1`, Timeout: 10}}
        case "migrate schema":
            task.AcceptanceChecks = []m.AcceptanceCheck{{Type: "command", Cmd: `db/migrate status | grep -q Applied`, Timeout: 15}}
        case "api handlers":
            task.AcceptanceChecks = []m.AcceptanceCheck{{Type: "command", Cmd: `curl -fsS ${API_BASE:-http://localhost:8080}/healthz >/dev/null`, Timeout: 10}}
        default:
            // No safe generic check — force authors to provide a real acceptance
            // by using a failing placeholder so pipelines catch missing checks.
            task.AcceptanceChecks = []m.AcceptanceCheck{{Type: "command", Cmd: `sh -c 'echo "missing acceptance checks" >&2; exit 1'`, Timeout: 5}}
        }
    }
    if task.DurationUnit == "" {
        task.DurationUnit = "hours"
    }
    // Execution logging defaults + light variation to improve coverage realism
    wasEmpty := task.ExecutionLogging.Format == ""
    baseFields := []string{"timestamp", "task_id", "step", "status", "message"}
    switch strings.ToLower(task.Title) {
    case "setup db":
        task.ExecutionLogging.RequiredFields = append(baseFields, "db_response_time")
    case "migrate schema":
        task.ExecutionLogging.RequiredFields = append(baseFields, "migration_version")
    case "api handlers":
        if wasEmpty {
            task.ExecutionLogging.Format = "JSON"
        } else {
            // leave user-provided non-empty format untouched
        }
        task.ExecutionLogging.RequiredFields = append(baseFields, "service_version")
    default:
        if len(task.ExecutionLogging.RequiredFields) == 0 {
            task.ExecutionLogging.RequiredFields = baseFields
        }
    }
    if wasEmpty && task.ExecutionLogging.Format == "" {
        task.ExecutionLogging.Format = "JSONL"
    }
    task.Compensation.Idempotent = true
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
