package model

import "encoding/json"

// ValidatorReport captures summary details for plan validators.
type ValidatorReport struct {
	Name      string          `json:"name"`
	Status    string          `json:"status"`
	Command   string          `json:"command,omitempty"`
	InputHash string          `json:"input_hash"`
	Cached    bool            `json:"cached"`
	Detail    string          `json:"detail,omitempty"`
	RawOutput json.RawMessage `json:"raw_output,omitempty"`
}
