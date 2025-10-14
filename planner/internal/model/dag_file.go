package model

// DagMeta captures metadata for dag.json.
type DagMeta struct {
	Version      string `json:"version"`
	ArtifactHash string `json:"artifact_hash"`
	TasksHash    string `json:"tasks_hash"`
}

// DagNode represents a node entry in dag.json.
type DagNode struct {
	ID                  string `json:"id"`
	Depth               int    `json:"depth"`
	CriticalPath        bool   `json:"critical_path"`
	ParallelOpportunity int    `json:"parallel_opportunity"`
}

// DagEdge represents an edge entry in dag.json.
type DagEdge struct {
	From       string `json:"from"`
	To         string `json:"to"`
	Type       string `json:"type"`
	Transitive bool   `json:"transitive"`
}

// DagMetrics records aggregate measurements for the DAG.
type DagMetrics struct {
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
}

// DagAnalysis carries validation results for the DAG.
type DagAnalysis struct {
	OK       bool     `json:"ok"`
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
	SoftDeps []Edge   `json:"soft_deps"`
}

// DagFile represents the canonical dag.json artifact.
type DagFile struct {
	Meta     DagMeta     `json:"meta"`
	Nodes    []DagNode   `json:"nodes"`
	Edges    []DagEdge   `json:"edges"`
	Metrics  DagMetrics  `json:"metrics"`
	Analysis DagAnalysis `json:"analysis"`
}
