package model

// DagMeta captures metadata for dag.json.
type DagMeta struct {
	Version      string `json:"version"`
	ArtifactHash string `json:"artifactHash"`
	TasksHash    string `json:"tasksHash"`
}

// DagNode represents a node entry in dag.json.
type DagNode struct {
	ID                  string `json:"id"`
	Depth               int    `json:"depth"`
	CriticalPath        bool   `json:"criticalPath"`
	ParallelOpportunity int    `json:"parallelOpportunity"`
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
	MinConfidenceApplied float64        `json:"minConfidenceApplied"`
	KeptByType           map[string]int `json:"keptByType"`
	DroppedByType        map[string]int `json:"droppedByType"`
	Nodes                int            `json:"nodes"`
	Edges                int            `json:"edges"`
	EdgeDensity          float64        `json:"edgeDensity"`
	WidthApprox          int            `json:"widthApprox"`
	LongestPathLength    int            `json:"longestPathLength"`
	CriticalPath         []string       `json:"criticalPath"`
	IsolatedTasks        int            `json:"isolatedTasks"`
	VerbFirstPct         float64        `json:"verbFirstPct"`
	EvidenceCoverage     float64        `json:"evidenceCoverage"`
}

// DagAnalysis carries validation results for the DAG.
type DagAnalysis struct {
	OK       bool     `json:"ok"`
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
	SoftDeps []Edge   `json:"softDeps"`
}

// DagFile represents the canonical dag.json artifact.
type DagFile struct {
	Meta     DagMeta     `json:"meta"`
	Nodes    []DagNode   `json:"nodes"`
	Edges    []DagEdge   `json:"edges"`
	Metrics  DagMetrics  `json:"metrics"`
	Analysis DagAnalysis `json:"analysis"`
}
