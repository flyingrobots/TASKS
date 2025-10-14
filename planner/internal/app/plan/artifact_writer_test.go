package plan

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	m "github.com/james/tasks-planner/internal/model"
)

func TestFileArtifactWriterWritesArtifacts(t *testing.T) {
	tmp := t.TempDir()
	writer := FileArtifactWriter{}

	tf := &m.TasksFile{}
	tf.Meta.Version = "v8"
	tf.Tasks = []m.Task{{ID: "T001", Title: "Do thing"}}
	df := &m.DagFile{}
	df.Meta.Version = "v8"
	df.Nodes = []struct {
		ID                  string `json:"id"`
		Depth               int    `json:"depth"`
		CriticalPath        bool   `json:"critical_path"`
		ParallelOpportunity int    `json:"parallel_opportunity"`
	}{{ID: "T001", Depth: 0}}

	waves := &m.WavesArtifact{Meta: m.WavesMeta{Version: "v8"}, Waves: [][]string{{"T001"}}}
	bundle := ArtifactBundle{
		TasksFile:   tf,
		DagFile:     df,
		Coordinator: &m.Coordinator{Version: "v8"},
		Features:    makeFeaturesArtifact([]FeatureSummary{{ID: "F001", Title: "Feature"}}),
		Waves:       waves,
		Titles:      taskTitles(tf.Tasks),
		ValidatorReports: []m.ValidatorReport{{
			Name:      "acceptance",
			Status:    m.ValidatorStatusPass,
			Command:   "accept",
			InputHash: strings.Repeat("a", 64),
			Detail:    strings.Repeat("x", 5000),
		}},
	}

	res, err := writer.Write(context.Background(), tmp, bundle)
	if err != nil {
		t.Fatalf("write: %v", err)
	}
	if len(res.Hashes) == 0 || res.Hashes["tasks.json"] == "" {
		t.Fatalf("expected hashes populated, got %+v", res.Hashes)
	}

	planBytes, err := os.ReadFile(filepath.Join(tmp, "Plan.md"))
	if err != nil {
		t.Fatalf("read Plan.md: %v", err)
	}
	content := string(planBytes)
	if !strings.Contains(content, "tasks.json") {
		t.Fatalf("plan summary missing hashes: %s", content)
	}
	if !strings.Contains(content, "… (truncated)") {
		t.Fatalf("expected truncated validator detail in Plan.md")
	}
	if _, err := os.Stat(filepath.Join(tmp, "dag.dot")); err != nil {
		t.Fatalf("dag.dot not written: %v", err)
	}
}

func TestFileArtifactWriterAggregatesErrors(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "file")
	if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	writer := FileArtifactWriter{}
	bundle := ArtifactBundle{
		TasksFile: &m.TasksFile{Meta: struct {
			Version           string  `json:"version"`
			MinConfidence     float64 `json:"min_confidence"`
			ArtifactHash      string  `json:"artifact_hash"`
			CodebaseAnalysis  any     `json:"codebase_analysis"`
			Autonormalization struct {
				Split  []string `json:"split"`
				Merged []string `json:"merged"`
			} `json:"autonormalization"`
			ValidatorReports []m.ValidatorReport `json:"validator_reports,omitempty"`
		}{Version: "v8"}},
		DagFile: &m.DagFile{Meta: struct {
			Version      string `json:"version"`
			ArtifactHash string `json:"artifact_hash"`
			TasksHash    string `json:"tasks_hash"`
		}{Version: "v8"}},
		Coordinator: &m.Coordinator{Version: "v8"},
		Features:    &m.FeaturesArtifact{Meta: m.ArtifactMeta{Version: "v8"}},
		Waves:       &m.WavesArtifact{Meta: m.WavesMeta{Version: "v8"}},
	}

	_, err := writer.Write(context.Background(), path, bundle)
	if err == nil {
		t.Fatalf("expected error")
	}
	var agg interface{ Unwrap() []error }
	if !errors.As(err, &agg) {
		t.Fatalf("expected aggregated error, got %T", err)
	}
}
