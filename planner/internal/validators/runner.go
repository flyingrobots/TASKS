package validators

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/james/tasks-planner/internal/canonjson"
	"github.com/james/tasks-planner/internal/hash"
	m "github.com/james/tasks-planner/internal/model"
)

// Report is an alias to the planner model's validator report, ensuring shared semantics.
type Report = m.ValidatorReport

// Config drives validator execution.
type Config struct {
	AcceptanceCmd string
	EvidenceCmd   string
	InterfaceCmd  string
	CacheDir      string
	Timeout       time.Duration
}

// Payload is serialized and sent to validators.
type Payload struct {
	Tasks       *m.TasksFile   `json:"tasks,omitempty"`
	Dag         *m.DagFile     `json:"dag,omitempty"`
	Coordinator *m.Coordinator `json:"coordinator,omitempty"`
}

// ExecFunc executes a command with the provided stdin and returns stdout, stderr.
type ExecFunc func(ctx context.Context, command string, stdin []byte) ([]byte, []byte, error)

// Runner orchestrates validator executions.
type Runner struct {
	cfg    Config
	execFn ExecFunc
	cache  *cacheStore
}

// NewRunner instantiates a runner with sane defaults.
func NewRunner(cfg Config) (*Runner, error) {
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	cacheDir := cfg.CacheDir
	if cacheDir == "" {
		home, _ := os.UserHomeDir()
		if home == "" {
			return nil, errors.New("validators: resolve cache dir (set --validators-cache)")
		}
		cacheDir = filepath.Join(home, ".tasksd", "validator-cache")
	}
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return nil, fmt.Errorf("validators: mkdir cache: %w", err)
	}
	store, err := newCacheStore(cacheDir)
	if err != nil {
		return nil, err
	}
	return &Runner{cfg: cfg, execFn: defaultExec, cache: store}, nil
}

// SetExecFunc overrides the executor (useful for tests).
func (r *Runner) SetExecFunc(fn ExecFunc) {
	r.execFn = fn
}

// Run executes configured validators with the payload.
func (r *Runner) Run(ctx context.Context, payload Payload) ([]Report, error) {
	type entry struct {
		name string
		cmd  string
	}
	entries := []entry{
		{"acceptance", r.cfg.AcceptanceCmd},
		{"evidence", r.cfg.EvidenceCmd},
		{"interface", r.cfg.InterfaceCmd},
	}
	reports := make([]Report, 0, len(entries))
	var errs multiError
	for _, entry := range entries {
		if entry.cmd == "" {
			continue
		}
		rep, err := r.runSingle(ctx, entry.name, entry.cmd, payload)
		reports = append(reports, rep)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return reports, errs
	}
	return reports, nil
}

func (r *Runner) runSingle(ctx context.Context, name, cmd string, payload Payload) (Report, error) {
	inputBytes, inputHash, err := encodePayload(payload)
	if err != nil {
		return Report{}, err
	}
	rep, ok, err := r.cache.Load(name, inputHash)
	if err != nil {
		return Report{}, err
	}
	if ok {
		rep.Cached = true
		return rep, nil
	}

	runCtx, cancel := context.WithTimeout(ctx, r.cfg.Timeout)
	defer cancel()
	stdout, stderr, execErr := r.execFn(runCtx, cmd, inputBytes)
	raw := normalizeJSON(stdout)
	report := Report{
		Name:      name,
		Command:   cmd,
		InputHash: inputHash,
		RawOutput: raw,
	}
	status, detail := interpretOutput(raw)
	if detail != "" {
		report.Detail = detail
	}
	if status != "" {
		report.Status = status
	}
	if execErr != nil {
		if runCtx.Err() != nil {
			execErr = runCtx.Err()
		}
		if report.Detail == "" && len(stderr) > 0 {
			report.Detail = strings.TrimSpace(string(stderr))
		}
		if report.Detail == "" {
			report.Detail = execErr.Error()
		}
		if report.Status == "" {
			report.Status = m.ValidatorStatusError
		}
		return report, fmt.Errorf("validator %s: %w", name, execErr)
	}
	if report.Status == "" {
		report.Status = m.ValidatorStatusPass
	}
	if report.Detail == "" && len(stderr) > 0 {
		report.Detail = strings.TrimSpace(string(stderr))
	}
	if err := r.cache.Store(name, inputHash, report); err != nil {
		return report, err
	}
	return report, nil
}

func encodePayload(payload Payload) ([]byte, string, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, "", fmt.Errorf("validator payload marshal: %w", err)
	}
	can, err := canonjson.ToCanonicalJSON(raw)
	if err != nil {
		return nil, "", fmt.Errorf("validator canonicalize: %w", err)
	}
	return can, hash.HashCanonicalBytes(can), nil
}

func normalizeJSON(raw []byte) json.RawMessage {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return nil
	}
	if trimmed[0] == '{' || trimmed[0] == '[' {
		return json.RawMessage(trimmed)
	}
	quoted, _ := json.Marshal(string(trimmed))
	return json.RawMessage(quoted)
}

func interpretOutput(raw json.RawMessage) (string, string) {
	if len(raw) == 0 {
		return "", ""
	}
	var parsed struct {
		Status string `json:"status"`
		Detail string `json:"detail"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", string(raw)
	}
	return normalizeStatus(parsed.Status), parsed.Detail
}

func normalizeStatus(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case m.ValidatorStatusPass:
		return m.ValidatorStatusPass
	case m.ValidatorStatusFail:
		return m.ValidatorStatusFail
	case m.ValidatorStatusError:
		return m.ValidatorStatusError
	case m.ValidatorStatusSkip:
		return m.ValidatorStatusSkip
	case "ok":
		return "ok"
	default:
		return ""
	}
}

func defaultExec(ctx context.Context, command string, stdin []byte) ([]byte, []byte, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "cmd", "/c", command)
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", command)
	}
	cmd.Stdin = bytes.NewReader(stdin)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		if ctx.Err() != nil {
			return out.Bytes(), errBuf.Bytes(), ctx.Err()
		}
		return out.Bytes(), errBuf.Bytes(), err
	}
	return out.Bytes(), errBuf.Bytes(), nil
}

type multiError []error

func (m multiError) Error() string {
	if len(m) == 0 {
		return ""
	}
	parts := make([]string, len(m))
	for i, err := range m {
		parts[i] = err.Error()
	}
	return strings.Join(parts, "; ")
}

func (m multiError) Unwrap() []error {
	if len(m) == 0 {
		return nil
	}
	return []error(m)
}
