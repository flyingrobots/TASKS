package analysis

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// Helper function to create a mock codebase structure
func createMockCodebase(t *testing.T, tmpDir string, structure map[string]string) {
	for path, content := range structure {
		fullPath := filepath.Join(tmpDir, path)
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "census-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			createMockCodebase(t, tmpDir, tt.structure)

			analysisResult, err := RunCensus(tmpDir)
			if tt.expectError {
				if err == nil {
					t.Fatal("Expected an error, but got none")
				}
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

			if len(analysisResult.Files) != len(expectedAbsFiles) {
				t.Errorf("Expected %d total files, but found %d", len(expectedAbsFiles), len(analysisResult.Files))
			} else {
				for i, f := range analysisResult.Files {
					if f != expectedAbsFiles[i] {
						t.Errorf("File mismatch at index %d: expected %s, got %s", i, expectedAbsFiles[i], f)
					}
				}
			}

			if len(analysisResult.GoFiles) != len(expectedAbsGoFiles) {
				t.Errorf("Expected %d Go files, but found %d", len(expectedAbsGoFiles), len(analysisResult.GoFiles))
			} else {
				for i, f := range analysisResult.GoFiles {
					if f != expectedAbsGoFiles[i] {
						t.Errorf("Go file mismatch at index %d: expected %s, got %s", i, expectedAbsGoFiles[i], f)
					}
				}
			}
		})
	}
}