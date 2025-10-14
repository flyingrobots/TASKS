package plan

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/james/tasks-planner/internal/emitter"
	"github.com/james/tasks-planner/internal/export/dot"
	m "github.com/james/tasks-planner/internal/model"
)

const validatorDetailLimit = 2048

// FileArtifactWriter writes planner artifacts to disk.
type FileArtifactWriter struct{}

// Write emits JSON artifacts and DOT previews, returning their hashes.
func (FileArtifactWriter) Write(ctx context.Context, out string, bundle ArtifactBundle) (ArtifactWriteResult, error) {
	hashes := map[string]string{}
	var errs artifactErrors
	writeWithHash := func(name string, value any, setHash func(string)) {
		hashValue, err := emitter.WriteWithArtifactHash(filepath.Join(out, name), value, setHash)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", name, err))
			return
		}
		hashes[name] = hashValue
	}

	writeWithHash("tasks.json", bundle.TasksFile, func(h string) { bundle.TasksFile.Meta.ArtifactHash = h })
	if h, ok := hashes["tasks.json"]; ok {
		bundle.DagFile.Meta.TasksHash = h
	}
	bundle.DagFile.Meta.ArtifactHash = ""
	writeWithHash("dag.json", bundle.DagFile, func(h string) { bundle.DagFile.Meta.ArtifactHash = h })

	if meta, ok := bundle.Waves["meta"].(map[string]any); ok {
		if h, ok := hashes["tasks.json"]; ok {
			meta["planId"] = h
		}
	}
	writeWithHash("waves.json", bundle.Waves, func(h string) {
		if meta, ok := bundle.Waves["meta"].(map[string]any); ok {
			meta["artifact_hash"] = h
		}
	})

	writeWithHash("features.json", bundle.Features, func(h string) {
		if meta, ok := bundle.Features["meta"].(map[string]any); ok {
			meta["artifact_hash"] = h
		}
	})

	writeWithHash("coordinator.json", bundle.Coordinator, func(string) {})

	if err := writePlanSummary(out, hashes, bundle.ValidatorReports); err != nil {
		errs = append(errs, err)
	}

	dagDot := dot.FromDagWithOptions(*bundle.DagFile, bundle.Titles, dot.Options{NodeLabel: "id-title", EdgeLabel: "type"})
	if err := os.WriteFile(filepath.Join(out, "dag.dot"), []byte(dagDot), 0o644); err != nil {
		errs = append(errs, fmt.Errorf("write dag.dot: %w", err))
	}
	runtimeDot := dot.FromCoordinatorWithOptions(*bundle.Coordinator, dot.Options{NodeLabel: "id-title", EdgeLabel: "type"})
	if err := os.WriteFile(filepath.Join(out, "runtime.dot"), []byte(runtimeDot), 0o644); err != nil {
		errs = append(errs, fmt.Errorf("write runtime.dot: %w", err))
	}

	if len(errs) > 0 {
		return ArtifactWriteResult{}, errs
	}
	return ArtifactWriteResult{Hashes: hashes}, nil
}

type artifactErrors []error

func (a artifactErrors) Error() string {
	parts := make([]string, len(a))
	for i, err := range a {
		parts[i] = err.Error()
	}
	return strings.Join(parts, "; ")
}

func (a artifactErrors) Unwrap() []error { return []error(a) }

func writePlanSummary(out string, hashes map[string]string, validatorReports []m.ValidatorReport) error {
	names := []string{"features.json", "tasks.json", "dag.json", "waves.json", "coordinator.json"}
	var sb strings.Builder
	sb.WriteString("# Plan (stub)\n\n")
	sb.WriteString("## Hashes\n\n")
	for _, name := range names {
		hashValue := hashes[name]
		sb.WriteString(fmt.Sprintf("- %s: %s\n", name, hashValue))
	}
	if len(validatorReports) > 0 {
		sb.WriteString("\n## Validators\n\n")
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
				sb.WriteString(fmt.Sprintf("- %s: %s%s — %s\n", rep.Name, rep.Status, cached, detail))
			} else {
				sb.WriteString(fmt.Sprintf("- %s: %s%s\n", rep.Name, rep.Status, cached))
			}
		}
	}
	return os.WriteFile(filepath.Join(out, "Plan.md"), []byte(sb.String()), 0o644)
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

func taskTitles(tasks []m.Task) map[string]string {
	titles := make(map[string]string, len(tasks))
	for _, t := range tasks {
		titles[t.ID] = t.Title
	}
	return titles
}
