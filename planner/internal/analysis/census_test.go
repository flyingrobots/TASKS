package analysis

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp" // Import go-cmp
)

// Helper function to create a mock codebase structure
func createMockCodebase(t *testing.T, tmpDir string, structure map[string]string) {
	for path, content := range structure {
		// Validate path: reject absolute paths, normalize, ensure it doesn't escape tmpDir
		if filepath.IsAbs(path) {
			t.Fatalf("Test setup error: Absolute path '%s' not allowed in mock structure", path)
		}
		cleanedPath := filepath.Clean(path)
		// Check if the cleaned path attempts to go up the directory tree
		if strings.HasPrefix(cleanedPath, "..") {
			t.Fatalf("Test setup error: Path '%s' attempts to escape temp directory", path)
		}
		
		fullPath := filepath.Join(tmpDir, cleanedPath)
		
		// Ensure the path doesn't escape tmpDir after joining
		rel, err := filepath.Rel(tmpDir, fullPath)
		if err != nil {
			t.Fatalf("Test setup error: Failed to get relative path for '%s': %v", fullPath, err)
		}
		if strings.HasPrefix(rel, "..") {
			t.Fatalf("Test setup error: Path '%s' escapes temp directory after join", path)
		}

		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create dir %s: %v", dir, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write mock file %s: %v", fullPath, err)
		}
	}
}

func TestRunCensus(t *testing.T) {
	tests := []struct {
		name          string
		structure     map[string]string // path -> content
		expectedFiles []string
		expectedGoFiles []string
		expectError   bool
		setupError    bool // For testing non-existent path
	}{
		{
			name: "empty directory",
			structure: map[string]string{},
			expectedFiles: []string{},
			expectedGoFiles: []string{},
		},
		{
			name: "flat directory with mixed files",
			structure: map[string]string{
				"main.go":    "package main",
				"go.mod":     "module test",
				"README.md":  "# Test",
				"util.go":    "package util",
			},
			expectedFiles: []string{"go.mod", "main.go", "README.md", "util.go"},
			expectedGoFiles: []string{"main.go", "util.go"},
		},
		{
			name: "nested directory with Go files",
			structure: map[string]string{
				"main.go":             "package main",
				"internal/db/db.go":   "package db",
				"internal/api/api.go": "package api",
				"internal/config.yaml": "config: true",
			},
			expectedFiles: []string{"internal/api/api.go", "internal/config.yaml", "internal/db/db.go", "main.go"},
			expectedGoFiles: []string{"internal/api/api.go", "internal/db/db.go", "main.go"},
		},
		{
			name: "directory with only non-Go files",
			structure: map[string]string{
				"config.json": "{}",
				"template.html": "<html>",
			},
			expectedFiles: []string{"config.json", "template.html"},
			expectedGoFiles: []string{},
		},
		{
			name: "non-existent directory (error case)",
			structure: map[string]string{},
			expectError: true,
			setupError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tmpDir string
			var err error

			if !tt.setupError { // Only create temp dir if not testing a non-existent path
				tmpDir, err = os.MkdirTemp("", "census-test-*")
				if err != nil {
					t.Fatalf("Failed to create temp dir: %v", err)
				}
				defer os.RemoveAll(tmpDir)
				createMockCodebase(t, tmpDir, tt.structure)
			} else {
				// For error cases, use a path that definitely doesn't exist
				tmpDir = filepath.Join(os.TempDir(), "non-existent-path-for-census-test")
			}

			analysisResult, err := RunCensus(tmpDir)

			if tt.expectError {
				if err == nil {
					t.Fatal("Expected an error, but got none")
				}
				// For now, we just check for an error. More specific error type checks can be added later.
				return
			}
			if err != nil {
				t.Fatalf("RunCensus failed: %v", err)
			}

			if analysisResult == nil {
				t.Fatal("Expected a non-nil analysis result")
			}

			// Adjust expected paths to be absolute for comparison
			expectedAbsFiles := make([]string, len(tt.expectedFiles))
			for i, f := range tt.expectedFiles {
				expectedAbsFiles[i] = filepath.Join(tmpDir, f)
			}
			sort.Strings(expectedAbsFiles) // Sort to ensure deterministic comparison

			expectedAbsGoFiles := make([]string, len(tt.expectedGoFiles))
			for i, f := range tt.expectedGoFiles {
				expectedAbsGoFiles[i] = filepath.Join(tmpDir, f)
			}
			sort.Strings(expectedAbsGoFiles) // Sort to ensure deterministic comparison

			// Sort actual results for deterministic comparison
			sort.Strings(analysisResult.Files)
			sort.Strings(analysisResult.GoFiles)

			// Use cmp.Diff for idiomatic slice comparison
			if diff := cmp.Diff(expectedAbsFiles, analysisResult.Files); diff != "" {
				t.Errorf("analysisResult.Files mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(expectedAbsGoFiles, analysisResult.GoFiles); diff != "" {
				t.Errorf("analysisResult.GoFiles mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
