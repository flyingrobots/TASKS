package model

// Evidence represents a piece of evidence supporting a task or dependency.
type Evidence struct {
	Type       string  `json:"type"`
	Source     string  `json:"source"`
	Excerpt    string  `json:"excerpt,omitempty"`
	Confidence float64 `json:"confidence"`
	Rationale  string  `json:"rationale"`
}