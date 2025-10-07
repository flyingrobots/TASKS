package model

// Edge represents a dependency between two tasks.
type Edge struct {
	From       string     `json:"from"`
	To         string     `json:"to"`
	Type       string     `json:"type"`
	Subtype    string     `json:"subtype,omitempty"`
	IsHard     bool       `json:"isHard"`
	Confidence float64    `json:"confidence"`
	Evidence   []Evidence `json:"evidence,omitempty"`
}
