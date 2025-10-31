package model

// ResourceSpec defines a resource catalog entry for coordinator.json.
type ResourceSpec struct {
	Capacity  int    `json:"capacity"`
	Mode      string `json:"mode"`
	LockOrder int    `json:"lockOrder"`
}

// Coordinator represents the coordinator.json contract passed from the planner to the executor.
type Coordinator struct {
	Version string `json:"version"`
	Graph   struct {
		Nodes []Task `json:"nodes"`
		Edges []Edge `json:"edges"`
	} `json:"graph"`
	Config struct {
		Resources struct {
			Catalog  map[string]ResourceSpec   `json:"catalog"`
			Profiles map[string]map[string]int `json:"profiles"`
		} `json:"resources"`
		Policies struct {
			ConcurrencyMax           int            `json:"concurrencyMax"`
			LockOrdering             []string       `json:"lockOrdering"`
			CircuitBreakerThresholds map[string]any `json:"circuitBreakerThresholds"`
		} `json:"policies"`
	} `json:"config"`
	Metrics struct {
		Estimates struct {
			P50TotalHours     float64 `json:"p50TotalHours"`
			LongestPathLength int     `json:"longestPathLength"`
			WidthApprox       int     `json:"widthApprox"`
		} `json:"estimates"`
	} `json:"metrics"`
}
