package model

// TasksFile represents the canonical tasks.json artifact.
type TasksFile struct {
	Meta struct {
		Version           string  `json:"version"`
		MinConfidence     float64 `json:"min_confidence"`
		ArtifactHash      string  `json:"artifact_hash"`
		CodebaseAnalysis  any     `json:"codebase_analysis"`
		Autonormalization struct {
			Split  []string `json:"split"`
			Merged []string `json:"merged"`
		} `json:"autonormalization"`
	} `json:"meta"`
	Tasks        []Task         `json:"tasks"`
	Dependencies []Edge         `json:"dependencies"`
	ResourceConflicts map[string]any `json:"resource_conflicts,omitempty"`
}
