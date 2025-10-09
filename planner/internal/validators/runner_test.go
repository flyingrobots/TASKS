package validators

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/james/tasks-planner/internal/model/testutil"
)

func TestRunnerExecutesAndCaches(t *testing.T) {
	tmp := t.TempDir()
	bin := buildMockValidator(t, tmp, successValidatorSource)
	runner := newRunnerForTest(t, Config{AcceptanceCmd: bin, CacheDir: tmp, Timeout: 2 * time.Second})
	payload := Payload{Tasks: testutil.StubTasksFile()}
	reports, err := runner.Run(context.Background(), payload)
	if err != nil {
		t.Fatalf("first run: %v", err)
	}
	if len(reports) != 1 || reports[0].Status != StatusPass {
		t.Fatalf("unexpected reports: %+v", reports)
	}
	if reports[0].Cached {
		t.Fatalf("expected fresh run, got cached")
	}
	reports2, err := runner.Run(context.Background(), payload)
	if err != nil {
		t.Fatalf("second run: %v", err)
	}
	if !reports2[0].Cached {
		t.Fatalf("expected cached run")
	}
}

func TestRunnerErrorHandling(t *testing.T) {
	tmp := t.TempDir()
	bin := buildMockValidator(t, tmp, errorValidatorSource)
	runner := newRunnerForTest(t, Config{AcceptanceCmd: bin, CacheDir: tmp, Timeout: time.Second})
	payload := Payload{Tasks: testutil.StubTasksFile()}
	reports, err := runner.Run(context.Background(), payload)
	if err == nil {
		t.Fatalf("expected error")
	}
	if len(reports) != 1 || reports[0].Status != StatusError {
		t.Fatalf("expected error status, got %+v", reports)
	}
	if reports[0].Cached {
		t.Fatalf("error report must not be cached")
	}
}

func TestRunnerTimeout(t *testing.T) {
	tmp := t.TempDir()
	bin := buildMockValidator(t, tmp, slowValidatorSource)
	runner := newRunnerForTest(t, Config{AcceptanceCmd: bin, CacheDir: tmp, Timeout: 200 * time.Millisecond})
	payload := Payload{Tasks: testutil.StubTasksFile()}
	_, err := runner.Run(context.Background(), payload)
	if err == nil {
		t.Fatalf("expected timeout error")
	}
}

func TestRunnerMultipleValidators(t *testing.T) {
	tmp := t.TempDir()
	bin := buildMockValidator(t, tmp, successValidatorSource)
	runner := newRunnerForTest(t, Config{
		AcceptanceCmd: bin,
		EvidenceCmd:   bin,
		CacheDir:      tmp,
		Timeout:       time.Second,
	})
	payload := Payload{Tasks: testutil.StubTasksFile()}
	reports, err := runner.Run(context.Background(), payload)
	if err != nil {
		t.Fatalf("multi run: %v", err)
	}
	if len(reports) != 2 {
		t.Fatalf("expected two reports, got %d", len(reports))
	}
}

func TestRunnerCacheInvalidation(t *testing.T) {
	tmp := t.TempDir()
	bin := buildMockValidator(t, tmp, successValidatorSource)
	runner := newRunnerForTest(t, Config{AcceptanceCmd: bin, CacheDir: tmp, Timeout: time.Second})
	payload := Payload{Tasks: testutil.StubTasksFile()}
	reports, err := runner.Run(context.Background(), payload)
	if err != nil {
		t.Fatalf("first run: %v", err)
	}
	if reports[0].Cached {
		t.Fatalf("expected uncached")
	}
	payload.Tasks.Tasks[0].Title = "Different"
	reports2, err := runner.Run(context.Background(), payload)
	if err != nil {
		t.Fatalf("second run: %v", err)
	}
	if reports2[0].Cached {
		t.Fatalf("payload changed; cache should miss")
	}
}

func TestRunnerConcurrentCacheAccess(t *testing.T) {
	tmp := t.TempDir()
	bin := buildMockValidator(t, tmp, successValidatorSource)
	runner := newRunnerForTest(t, Config{AcceptanceCmd: bin, CacheDir: tmp, Timeout: time.Second})
	payload := Payload{Tasks: testutil.StubTasksFile()}
	if _, err := runner.Run(context.Background(), payload); err != nil {
		t.Fatalf("priming run: %v", err)
	}
	var wg sync.WaitGroup
	errCh := make(chan error, 8)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			reps, err := runner.Run(context.Background(), payload)
			if err != nil {
				errCh <- err
				return
			}
			if len(reps) != 1 || !reps[0].Cached {
				errCh <- fmt.Errorf("expected cached report, got %+v", reps)
			}
		}()
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			t.Fatalf("concurrent run: %v", err)
		}
	}
}

const successValidatorSource = `package main
import (
  "encoding/json"
  "os"
)
type payload struct {
  Tasks any ` + "`json:\"tasks\"`" + `
}
func main() {
  var p payload
  _ = json.NewDecoder(os.Stdin).Decode(&p)
  json.NewEncoder(os.Stdout).Encode(map[string]string{"status":"pass","detail":"ok"})
}
`

const errorValidatorSource = `package main
import (
  "fmt"
  "os"
)
func main() {
  fmt.Fprintln(os.Stderr, "boom")
  os.Exit(2)
}
`

const slowValidatorSource = `package main
import (
  "time"
)
func main() {
  time.Sleep(2 * time.Second)
}
`

func newRunnerForTest(t *testing.T, cfg Config) *Runner {
	t.Helper()
	runner, err := NewRunner(cfg)
	if err != nil {
		t.Fatalf("new runner: %v", err)
	}
	return runner
}

func buildMockValidator(t *testing.T, dir, source string) string {
	t.Helper()
	src := filepath.Join(dir, "main.go")
	if err := os.WriteFile(src, []byte(source), 0o644); err != nil {
		t.Fatalf("write validator: %v", err)
	}
	bin := filepath.Join(dir, "validator")
	cmd := exec.Command("go", "build", "-o", bin, src)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("build validator: %v", err)
	}
	if runtime.GOOS == "windows" {
		return bin + ".exe"
	}
	return bin
}
