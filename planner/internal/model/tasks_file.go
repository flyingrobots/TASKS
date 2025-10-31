package model

// TasksFile represents the canonical tasks.json artifact.
type TasksFile struct {
	Meta struct {
		Version           string  `json:"version"`
		MinConfidence     float64 `json:"minConfidence"`
		ArtifactHash      string  `json:"artifactHash"`
		CodebaseAnalysis  any     `json:"codebaseAnalysis"`
		Autonormalization struct {
			Split  []string `json:"split"`
			Merged []string `json:"merged"`
		} `json:"autonormalization"`
		ValidatorReports []ValidatorReport `json:"validatorReports,omitempty"`
	} `json:"meta"`
	Tasks             []Task         `json:"tasks"`
	Dependencies      []Edge         `json:"dependencies"`
	ResourceConflicts map[string]any `json:"resourceConflicts,omitempty"`
}
