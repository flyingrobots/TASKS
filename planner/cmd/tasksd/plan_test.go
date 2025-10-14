package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	m "github.com/james/tasks-planner/internal/model"
)

func TestPlanCommandWithValidators(t *testing.T) {
	tmp := t.TempDir()
	outDir := filepath.Join(tmp, "out")
	validatorBin := buildMockValidatorBinary(t, tmp)
	repoRoot := repoRoot(t)
	cmd := exec.Command("go", "run", "./cmd/tasksd",
		"plan",
		"--out", outDir,
		"--validators-acceptance", validatorBin,
	)
	cmd.Dir = repoRoot
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("tasksd plan: %v\nstdout: %s\nstderr: %s", err, stdout.String(), stderr.String())
	}
	contents, err := os.ReadFile(filepath.Join(outDir, "tasks.json"))
	if err != nil {
		t.Fatalf("read tasks.json: %v", err)
	}
	var tf m.TasksFile
	if err := json.Unmarshal(contents, &tf); err != nil {
		t.Fatalf("unmarshal tasks.json: %v", err)
	}
	if len(tf.Meta.ValidatorReports) == 0 {
		t.Fatalf("expected validator reports recorded")
	}
	report := tf.Meta.ValidatorReports[0]
	if report.Name == "" || report.Status == "" || report.InputHash == "" {
		t.Fatalf("validator report incomplete: %+v", report)
	}
}

func buildMockValidatorBinary(t *testing.T, dir string) string {
	t.Helper()
	if _, err := exec.LookPath("go"); err != nil {
		t.Skipf("skipping validator build; go binary not found: %v", err)
	}
	source := `package main
import (
  "encoding/json"
  "os"
)
func main() {
  var payload map[string]any
  _ = json.NewDecoder(os.Stdin).Decode(&payload)
  json.NewEncoder(os.Stdout).Encode(map[string]string{"status":"pass","detail":"ok"})
}
`
	src := filepath.Join(dir, "validator.go")
	if err := os.WriteFile(src, []byte(source), 0o644); err != nil {
		t.Fatalf("write validator: %v", err)
	}
	bin := filepath.Join(dir, "validator")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	cmd := exec.Command("go", "build", "-o", bin, src)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("build validator: %v\n%s", err, out.String())
	}
	return bin
}

func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for {
		candidate := filepath.Join(dir, "go.mod")
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("repo root: could not locate go.mod starting from %s", dir)
		}
		dir = parent
	}
}
