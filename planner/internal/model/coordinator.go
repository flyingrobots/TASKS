package model

// Coordinator represents the coordinator.json contract passed from the planner to the executor.
type Coordinator struct {
	Version string `json:"version"`
	Graph   struct {
		Nodes []Task `json:"nodes"`
		Edges []Edge `json:"edges"`
	} `json:"graph"`
	Config struct {
		Resources struct {
			Catalog map[string]struct {
				Capacity  int    `json:"capacity"`
				Mode      string `json:"mode"`
				LockOrder int    `json:"lock_order"`
			} `json:"catalog"`
			Profiles map[string]map[string]int `json:"profiles"`
		} `json:"resources"`
		Policies struct {
			ConcurrencyMax           int              `json:"concurrency_max"`
			LockOrdering             []string         `json:"lock_ordering"`
			CircuitBreakerThresholds map[string]any `json:"circuit_breaker_thresholds"`
		} `json:"policies"`
	} `json:"config"`
	Metrics struct {
		Estimates struct {
			P50TotalHours     float64 `json:"p50_total_hours"`
			LongestPathLength int     `json:"longest_path_length"`
			WidthApprox       int     `json:"width_approx"`
		} `json:"estimates"`
	} `json:"metrics"`
}
