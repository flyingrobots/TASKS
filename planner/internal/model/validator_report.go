package model

import "encoding/json"

// Validator report status values mirrored from the validators package.
const (
	ValidatorStatusPass  = "pass"
	ValidatorStatusFail  = "fail"
	ValidatorStatusError = "error"
	ValidatorStatusSkip  = "skip"
)

// ValidatorReport captures summary details for plan validators.
//
//   - Status uses the ValidatorStatus* constants (or "ok" for legacy validators).
//   - InputHash records the SHA-256 hash of the canonical validator payload.
//   - Cached indicates whether the report was reused from the local validator cache.
//   - Detail is a human-readable summary (potentially truncated).
//   - RawOutput stores the normalized JSON returned by the validator, when available.
type ValidatorReport struct {
	Name      string          `json:"name"`
	Status    string          `json:"status"`
	Command   string          `json:"command,omitempty"`
	InputHash string          `json:"input_hash"`
	Cached    bool            `json:"cached"`
	Detail    string          `json:"detail,omitempty"`
	RawOutput json.RawMessage `json:"raw_output,omitempty"`
}
