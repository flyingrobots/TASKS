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
		titleToID[normalizeKey(spec.Title)] = id
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
	return tasks, []FeatureSummary{{ID: "F1", Title: "Core DB + API"}}
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
	if trimmed == "" {
		return ""
	}
	if len(trimmed) > 1 && (trimmed[0] == 'T' || trimmed[0] == 't') {
		for _, r := range trimmed[1:] {
			if r < '0' || r > '9' {
				return ""
			}
		}
		if len(trimmed) == 4 {
			return strings.ToUpper(trimmed)
		}
	}
	return titleToID[normalizeKey(trimmed)]
}

func normalizeKey(v string) string {
	return strings.ToLower(strings.TrimSpace(v))
}
