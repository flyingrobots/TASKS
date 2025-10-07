package model

// DagFile represents the canonical dag.json artifact.
type DagFile struct {
	Meta struct {
		Version      string `json:"version"`
		ArtifactHash string `json:"artifact_hash"`
		TasksHash    string `json:"tasks_hash"`
	} `json:"meta"`
	Nodes []struct {
		ID                  string `json:"id"`
		Depth               int    `json:"depth"`
		CriticalPath        bool   `json:"critical_path"`
		ParallelOpportunity int    `json:"parallel_opportunity"`
	} `json:"nodes"`
	Edges []struct {
		From       string `json:"from"`
		To         string `json:"to"`
		Type       string `json:"type"`
		Transitive bool   `json:"transitive"`
	} `json:"edges"`
	Metrics struct {
		MinConfidenceApplied float64        `json:"min_confidence_applied"`
		KeptByType           map[string]int `json:"kept_by_type"`
		DroppedByType        map[string]int `json:"dropped_by_type"`
		Nodes                int            `json:"nodes"`
		Edges                int            `json:"edges"`
		EdgeDensity          float64        `json:"edge_density"`
		WidthApprox          int            `json:"width_approx"`
		LongestPathLength    int            `json:"longest_path_length"`
		CriticalPath         []string       `json:"critical_path"`
		IsolatedTasks        int            `json:"isolated_tasks"`
		VerbFirstPct         float64        `json:"verb_first_pct"`
		EvidenceCoverage     float64        `json:"evidence_coverage"`
	} `json:"metrics"`
	Analysis struct {
		OK       bool     `json:"ok"`
		Errors   []string `json:"errors"`
		Warnings []string `json:"warnings"`
		SoftDeps []Edge   `json:"soft_deps"`
	} `json:"analysis"`
}
