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
	"strings"
	"sync"
	"time"

	"github.com/james/tasks-planner/internal/canonjson"
	"github.com/james/tasks-planner/internal/hash"
	m "github.com/james/tasks-planner/internal/model"
)

// Config drives validator execution.
type Config struct {
	AcceptanceCmd string
	EvidenceCmd   string
	InterfaceCmd  string
	CacheDir      string
	Timeout       time.Duration
}

// Report summarizes a validator run.
type Report struct {
	Name      string          `json:"name"`
	Command   string          `json:"command"`
	InputHash string          `json:"input_hash"`
	Status    string          `json:"status"`
	Detail    string          `json:"detail,omitempty"`
	RawOutput json.RawMessage `json:"raw_output,omitempty"`
	Cached    bool            `json:"cached"`
}

// Payload is serialized and sent to validators.
type Payload struct {
	Tasks       *m.TasksFile   `json:"tasks,omitempty"`
	Dag         *m.DagFile     `json:"dag,omitempty"`
	Coordinator *m.Coordinator `json:"coordinator,omitempty"`
}

// ExecFunc executes a command with the provided stdin.
type ExecFunc func(ctx context.Context, cmd string, stdin []byte) ([]byte, error)

// Runner orchestrates validator executions.
type Runner struct {
	cfg        Config
	execFn     ExecFunc
	cache      *cacheStore
	cacheMutex sync.Mutex
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
	var reports []Report
	entries := []struct {
		name string
		cmd  string
	}{
		{"acceptance", r.cfg.AcceptanceCmd},
		{"evidence", r.cfg.EvidenceCmd},
		{"interface", r.cfg.InterfaceCmd},
	}
	for _, entry := range entries {
		if entry.cmd == "" {
			continue
		}
		rep, err := r.runSingle(ctx, entry.name, entry.cmd, payload)
		if err != nil {
			return nil, err
		}
		reports = append(reports, rep)
	}
	return reports, nil
}

func (r *Runner) runSingle(ctx context.Context, name, cmd string, payload Payload) (Report, error) {
	inputBytes, inputHash, err := encodePayload(payload)
	if err != nil {
		return Report{}, err
	}
	r.cacheMutex.Lock()
	if cached, ok := r.cache.Load(name, inputHash); ok {
		r.cacheMutex.Unlock()
		cached.Cached = true
		return cached, nil
	}
	r.cacheMutex.Unlock()

	runCtx, cancel := context.WithTimeout(ctx, r.cfg.Timeout)
	defer cancel()
	stdout, err := r.execFn(runCtx, cmd, inputBytes)
	if err != nil {
		return Report{}, fmt.Errorf("validator %s failed: %w", name, err)
	}
	rep := Report{
		Name:      name,
		Command:   cmd,
		InputHash: inputHash,
		Status:    "ok",
		RawOutput: normalizeJSON(stdout),
	}
    if len(rep.RawOutput) > 0 {
        rep.Detail = string(rep.RawOutput)
    }
	r.cacheMutex.Lock()
	r.cache.Store(name, inputHash, rep)
	r.cacheMutex.Unlock()
	return rep, nil
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
	// treat as plain string payload
	quoted, _ := json.Marshal(string(trimmed))
	return json.RawMessage(quoted)
}

func defaultExec(ctx context.Context, cmd string, stdin []byte) ([]byte, error) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty command")
	}
	c := exec.CommandContext(ctx, parts[0], parts[1:]...)
	c.Stdin = bytes.NewReader(stdin)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	c.Stdout = &out
	c.Stderr = &errBuf
	if err := c.Run(); err != nil {
		if errBuf.Len() > 0 {
			return nil, fmt.Errorf("%w: %s", err, errBuf.String())
		}
		return nil, err
	}
	return out.Bytes(), nil
}
