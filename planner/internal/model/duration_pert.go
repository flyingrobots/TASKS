package model

// DurationPERT provides optimistic, most likely, and pessimistic estimates for a task's duration.
type DurationPERT struct {
	Optimistic  float64 `json:"optimistic"`
	MostLikely  float64 `json:"mostLikely"`
	Pessimistic float64 `json:"pessimistic"`
}
