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
func createMockCodebase(t *testing.T, tmpDir string, structure map[string]string, chmodDir string) {
	// First, create all files and directories
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

	// Then, apply chmod if specified
	if chmodDir != "" {
		unreadableDirPath := filepath.Join(tmpDir, chmodDir)
		// Ensure the directory exists before chmodding
		if err := os.MkdirAll(unreadableDirPath, 0755); err != nil {
			t.Fatalf("Failed to create unreadable dir %s: %v", unreadableDirPath, err)
		}
		// Defer permission restore BEFORE the test function's defer os.RemoveAll
		defer func() {
			if err := os.Chmod(unreadableDirPath, 0755); err != nil {
				t.Logf("Failed to restore permissions for %s: %v", unreadableDirPath, err)
			}
		}()
		if err := os.Chmod(unreadableDirPath, 0000); err != nil { // chmod 000
			t.Fatalf("Failed to chmod dir %s: %v", unreadableDirPath, err)
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
		chmodDir      string // Directory to chmod 000 for permission denied test
		chmodRoot     bool // If true, chmod tmpDir itself to 000
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
		{
			name: "hidden files and dot-directories",
			structure: map[string]string{
				".git/config": "git config",
				".env":        "ENV_VAR=value",
				"main.go":     "package main",
				".vscode/settings.json": "{}",
			},
			expectedFiles: []string{`.env`, `.git/config`, `.vscode/settings.json`, `main.go`},
			expectedGoFiles: []string{"main.go"},
		},
		{
			name: "filenames with spaces and unicode",
			structure: map[string]string{
				`My File.go`:      "package main",
				`résumé.pdf`:      "pdf content",
				`你好世界.txt`:      "hello world",
				`another file.txt`: "content",
			},
			expectedFiles: []string{`My File.go`, `another file.txt`, `résumé.pdf`, `你好世界.txt`},
			expectedGoFiles: []string{`My File.go`},
		},
		{
			name: "permission denied root directory", // New test case
			structure: map[string]string{
				"file.go": "package main",
			},
			expectedFiles: []string{`file.go`},
			expectedGoFiles: []string{`file.go`},
			expectError: true,
			chmodRoot: true, // Chmod tmpDir itself
		},
		{
			name: "permission denied subdirectory", // Original test case, renamed
			structure: map[string]string{
				"readable/file.go": "package main",
				"unreadable/file.go": "package main", // This file won't be found
			},
			expectedFiles: []string{"readable/file.go"},
			expectedGoFiles: []string{"readable/file.go"},
			expectError: true, // Expect error from WalkDir
			chmodDir: "unreadable", // Directory to make unreadable
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
				createMockCodebase(t, tmpDir, tt.structure, tt.chmodDir) // Pass chmodDir here

				if tt.chmodRoot { // Chmod the root directory itself
					defer func() {
						if err := os.Chmod(tmpDir, 0755); err != nil {
							t.Logf("Failed to restore permissions for root dir %s: %v", tmpDir, err)
						}
					}()
					if err := os.Chmod(tmpDir, 0000); err != nil {
						t.Fatalf("Failed to chmod root dir %s: %v", tmpDir, err)
					}
				}
			} else {
				// For error cases, use a path that definitely doesn't exist
				tmpDir = filepath.Join(os.TempDir(), "non-existent-path-for-census-test")
			}

			analysisResult, err := RunCensus(tmpDir)
			t.Logf("DEBUG: RunCensus returned err: %v", err) // Replaced with t.Logf

			if tt.expectError {
				if err == nil {
					t.Fatal("Expected an error, but got none")
				}
				
				// Check if the error is a permission denied error for chmodRoot case
				if tt.chmodRoot && !os.IsPermission(err) {
					t.Errorf("Expected permission denied error, got %v", err)
				}
				
				t.Logf("DEBUG: Test case '%s' correctly returned error: %v", tt.name, err) // Replaced with t.Logf
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