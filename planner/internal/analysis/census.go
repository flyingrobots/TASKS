package analysis

import (
	"io/fs"
	"path/filepath"
	"strings"
)

// CodebaseAnalysis holds the results of a codebase census.
type CodebaseAnalysis struct {
	Files   []string // All files found
	GoFiles []string // Just the Go files found
}

// RunCensus performs a census of the codebase at the given path.
// It discovers files and categorizes them.
func RunCensus(path string) (*CodebaseAnalysis, error) {
	analysis := &CodebaseAnalysis{}

	err := filepath.WalkDir(path, func(currentPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			analysis.Files = append(analysis.Files, currentPath)
			if strings.HasSuffix(d.Name(), ".go") {
				analysis.GoFiles = append(analysis.GoFiles, currentPath)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return analysis, nil
}