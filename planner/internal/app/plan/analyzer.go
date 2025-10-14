package plan

import (
	"context"

	analysis "github.com/james/tasks-planner/internal/analysis"
)

// Analyzer abstracts repository census discovery.
type Analyzer interface {
	Analyze(ctx context.Context, path string) (analysis.FileCensusCounts, error)
}

// CensusAnalyzer implements Analyzer using the analysis package.
type CensusAnalyzer struct{}

func (CensusAnalyzer) Analyze(ctx context.Context, path string) (analysis.FileCensusCounts, error) {
	if err := ctx.Err(); err != nil {
		return analysis.FileCensusCounts{}, err
	}

	type result struct {
		report *analysis.CodebaseAnalysis
		err    error
	}
	done := make(chan result, 1)
	go func() {
		report, err := analysis.RunCensus(path)
		done <- result{report: report, err: err}
	}()

	select {
	case <-ctx.Done():
		return analysis.FileCensusCounts{}, ctx.Err()
	case res := <-done:
		if res.err != nil {
			return analysis.FileCensusCounts{}, res.err
		}
		return analysis.FileCensusCounts{
			Files:   len(res.report.Files),
			GoFiles: len(res.report.GoFiles),
		}, nil
	}
}
