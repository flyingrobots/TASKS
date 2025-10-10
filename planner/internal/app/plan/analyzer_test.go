package plan

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestCensusAnalyzerCountsFiles(t *testing.T) {
	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, "a.go"), []byte("package a"), 0o644); err != nil {
		t.Fatalf("write go file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "b.txt"), []byte("text"), 0o644); err != nil {
		t.Fatalf("write txt file: %v", err)
	}

	an := CensusAnalyzer{}
	counts, err := an.Analyze(context.Background(), tmp)
	if err != nil {
		t.Fatalf("analyze: %v", err)
	}
	if counts.Files != 2 || counts.GoFiles != 1 {
		t.Fatalf("unexpected counts: %+v", counts)
	}
}

func TestCensusAnalyzerPropagatesErrors(t *testing.T) {
	an := CensusAnalyzer{}
	_, err := an.Analyze(context.Background(), filepath.Join("/does", "not", "exist"))
	if err == nil {
		t.Fatalf("expected error for missing path")
	}
	var pathErr *os.PathError
	if !errors.As(err, &pathErr) {
		t.Fatalf("expected path error, got %v", err)
	}
}
