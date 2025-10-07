package analysis

import (
	"io/fs"
	"path/filepath"
	"strings"
)

// CodebaseAnalysis holds the results of a codebase census.
// It includes a list of all files found and a specific list of Go files.
type CodebaseAnalysis struct {
	Files   []string // All files found within the scanned path.
	GoFiles []string // Paths to Go files (.go extension) found within the scanned path.
}

// RunCensus performs a census of the codebase at the given path.
// It recursively walks the directory, discovers all files, and categorizes them.
// Currently, it identifies all files and specifically Go files.
//
// Parameters:
//   path: The root directory path of the codebase to analyze.
//
// Returns:
//   A pointer to a CodebaseAnalysis struct containing the discovered files,
//   or an error if the directory cannot be walked.
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
