package analysis

import (
	"fmt"
	"io/fs"
	"os"
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
	analysis := &CodebaseAnalysis{
		Files:   []string{}, // Initialize to empty slice
		GoFiles: []string{}, // Initialize to empty slice
	}

	// Check if the path exists and is a directory
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to access path '%s': %w", path, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("path '%s' is not a directory", path)
	}

	err = filepath.WalkDir(path, func(currentPath string, d fs.DirEntry, err error) error {
		if err != nil {
			// fmt.Fprintf(os.Stderr, "DEBUG: WalkDir callback error at %s: %v\n", currentPath, err) // DEBUG PRINT
			return err
		}
		// Explicitly check for permission errors on directories
		if d.IsDir() {
			// Try to read the directory to force a permission error if it's unreadable
			// This is a workaround for filepath.WalkDir's subtle behavior on some OSes
			// where it might not return an error for unreadable subdirectories.
			_, statErr := os.ReadDir(currentPath)
			if statErr != nil && os.IsPermission(statErr) {
				return statErr // Propagate permission error
			}
		} else { // It's a file
			analysis.Files = append(analysis.Files, currentPath)
			if strings.HasSuffix(d.Name(), ".go") {
				analysis.GoFiles = append(analysis.GoFiles, currentPath)
			}
		}
		return nil
	})

	if err != nil {
		// fmt.Fprintf(os.Stderr, "DEBUG: RunCensus returning error from WalkDir: %v\n", err) // DEBUG PRINT
		return nil, fmt.Errorf("failed to walk directory '%s': %w", path, err)
	}

	return analysis, nil
}
