package plan

import (
	"context"
	"errors"
	"testing"
	"time"

	analysis "github.com/james/tasks-planner/internal/analysis"
)

func TestCensusAnalyzerCancellation(t *testing.T) {
	an := CensusAnalyzer{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := an.Analyze(ctx, ".")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestCensusAnalyzerReportsCounts(t *testing.T) {
	an := CensusAnalyzer{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	tmp := t.TempDir()
	counts, err := an.Analyze(ctx, tmp)
	if err != nil {
		t.Fatalf("analyze: %v", err)
	}
	if counts.Files != 0 || counts.GoFiles != 0 {
		t.Fatalf("expected empty counts for empty dir, got %+v", counts)
	}
	// sanity: ensure context still valid
	if err := ctx.Err(); err != nil {
		t.Fatalf("context unexpectedly done: %v", err)
	}
	_ = analysis.CodebaseAnalysis{}
}
