package model

// AcceptanceCheck defines a machine-verifiable criterion for task completion.
type AcceptanceCheck struct {
	Type    string         `json:"type"`
	Cmd     string         `json:"cmd,omitempty"`
	Path    string         `json:"path,omitempty"`
	Expect  map[string]any `json:"expect,omitempty"`
	Timeout int            `json:"timeoutSeconds,omitempty"`
}
