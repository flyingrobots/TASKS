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
	report, err := analysis.RunCensus(path)
	if err != nil {
		return analysis.FileCensusCounts{}, err
	}
	return analysis.FileCensusCounts{
		Files:   len(report.Files),
		GoFiles: len(report.GoFiles),
	}, nil
}
