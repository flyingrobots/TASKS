package dag

import (
    "errors"
    "fmt"
    "sort"
    "strings"

    m "github.com/james/tasks-planner/internal/model"
)

// edgeRec is a small record for edge bookkeeping during build.
type edgeRec struct{ From, To, Type string }

// Build builds a minimized DAG from tasks and edges, applying confidence and hardness filters,
// detecting cycles, computing layering depths, longest path (critical path), and removing transitive edges.
func Build(tasks []m.Task, edges []m.Edge, minConfidence float64) (*m.DagFile, error) {
    df := &m.DagFile{
        Meta:    m.DagMeta{Version: m.SchemaVersion},
        Metrics: m.DagMetrics{KeptByType: map[string]int{}, DroppedByType: map[string]int{}},
    }

	// index tasks
	idx := map[string]int{}
	order := make([]string, 0, len(tasks))
	for i, t := range tasks {
		if prev, exists := idx[t.ID]; exists {
			df.Analysis.OK = false
			detail := fmt.Sprintf("duplicate task id %s (first index %d, duplicate index %d)", t.ID, prev, i)
			df.Analysis.Errors = append(df.Analysis.Errors, detail)
			return df, errors.New(detail)
		}
		idx[t.ID] = i
		order = append(order, t.ID)
	}

	// filter edges: structural only
	kept := make([]edgeRec, 0, len(edges))
	for _, e := range edges {
		typeKey := edgeTypeKey(e.Type)
		if !e.IsHard || e.Confidence < minConfidence {
			df.Metrics.DroppedByType[typeKey]++
			continue
		}
		if e.Type == "resource" {
			df.Metrics.DroppedByType[typeKey]++
			continue
		}
		if _, ok := idx[e.From]; !ok {
			df.Metrics.DroppedByType[typeKey]++
			continue
		}
		if _, ok := idx[e.To]; !ok {
			df.Metrics.DroppedByType[typeKey]++
			continue
		}
		kept = append(kept, edgeRec{From: e.From, To: e.To, Type: e.Type})
	}

	// build adjacency
	n := len(tasks)
	if n == 0 {
		df.Analysis.OK = false
		df.Analysis.Errors = append(df.Analysis.Errors, "no tasks to build DAG")
		return df, errors.New("no tasks to build DAG")
	}
	adj := make([][]int, n)
	indeg := make([]int, n)
	for i := range adj {
		adj[i] = []int{}
	}
	for _, e := range kept {
		u, v := idx[e.From], idx[e.To]
		adj[u] = append(adj[u], v)
		indeg[v]++
	}

	// cycle detection via DFS
	seen := make([]uint8, n) // 0=unseen,1=visiting,2=done
	var dfs func(int) bool
	dfs = func(u int) bool {
		seen[u] = 1
		for _, v := range adj[u] {
			if seen[v] == 1 {
				return true
			} // cycle
			if seen[v] == 0 && dfs(v) {
				return true
			}
		}
		seen[u] = 2
		return false
	}
	for i := 0; i < n; i++ {
		if seen[i] == 0 && dfs(i) {
			df.Analysis.OK = false
			df.Analysis.Errors = append(df.Analysis.Errors, "cycle detected in dependencies")
			return df, errors.New("cycle detected")
		}
	}

	// Kahn layering + topo order
	q := []int{}
	indeg2 := make([]int, n)
	copy(indeg2, indeg)
	for i := 0; i < n; i++ {
		if indeg2[i] == 0 {
			q = append(q, i)
		}
	}
	topo := []int{}
	depth := make([]int, n)
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		topo = append(topo, u)
		for _, v := range adj[u] {
			if depth[v] < depth[u]+1 {
				depth[v] = depth[u] + 1
			}
			indeg2[v]--
			if indeg2[v] == 0 {
				q = append(q, v)
			}
		}
	}
	if len(topo) != n {
		df.Analysis.OK = false
		df.Analysis.Errors = append(df.Analysis.Errors, "not all nodes reached in topo (unexpected)")
		return df, errors.New("invalid topo")
	}

	// Longest path (by edges count) and predecessor
	dist := make([]int, n)
	pred := make([]int, n)
	for i := range pred {
		pred[i] = -1
	}
	for _, u := range topo {
		for _, v := range adj[u] {
			if dist[v] < dist[u]+1 {
				dist[v] = dist[u] + 1
				pred[v] = u
			}
		}
	}
	// find sink on critical path
	sink := 0
	for i := 1; i < n; i++ {
		if dist[i] > dist[sink] {
			sink = i
		}
	}
	critPath := []string{}
	for x := sink; x != -1; x = pred[x] {
		critPath = append(critPath, tasks[x].ID)
	}
	// reverse critPath
	for i, j := 0, len(critPath)-1; i < j; i, j = i+1, j-1 {
		critPath[i], critPath[j] = critPath[j], critPath[i]
	}

	// Transitive reduction: remove (u->v) if there exists u->w->...->v
	// Compute reachability via DFS from each node.
	reach := make([]map[int]bool, n)
	for u := 0; u < n; u++ {
		vis := map[int]bool{}
		stack := []int{u}
		for len(stack) > 0 {
			x := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, y := range adj[x] {
				if !vis[y] {
					vis[y] = true
					stack = append(stack, y)
				}
			}
		}
		reach[u] = vis
	}
	keepEdge := make(map[[2]int]bool)
	for _, e := range kept {
		u, v := idx[e.From], idx[e.To]
		removable := false
		for _, w := range adj[u] {
			if w != v && reach[w][v] {
				removable = true
				break
			}
		}
		if !removable {
			keepEdge[[2]int{u, v}] = true
		}
	}

	// Fill nodes
	// Stable order by ID
	sortedIDs := append([]string(nil), order...)
	sort.Strings(sortedIDs)
    for _, id := range sortedIDs {
        i := idx[id]
        df.Nodes = append(df.Nodes, m.DagNode{ID: id, Depth: depth[i], CriticalPath: contains(critPath, id), ParallelOpportunity: 1})
    }

	// Fill edges (non-transitive only)
	// Render deterministically by from,to ordering
	type K struct {
		F, T string
		Ty   string
	}
	kept2 := []K{}
	for key := range keepEdge {
		u, v := key[0], key[1]
		kept2 = append(kept2, K{F: tasks[u].ID, T: tasks[v].ID, Ty: lookupEdgeType(kept, tasks[u].ID, tasks[v].ID)})
	}
	sort.Slice(kept2, func(i, j int) bool {
		if kept2[i].F == kept2[j].F {
			return kept2[i].T < kept2[j].T
		}
		return kept2[i].F < kept2[j].F
	})
	for _, e := range kept2 {
		df.Edges = append(df.Edges, struct {
			From       string `json:"from"`
			To         string `json:"to"`
			Type       string `json:"type"`
			Transitive bool   `json:"transitive"`
		}{From: e.F, To: e.T, Type: e.Ty, Transitive: false})
	}

	// Recompute kept edge counts after transitive reduction.
	counts := map[string]int{}
	for _, edge := range df.Edges {
		counts[edgeTypeKey(edge.Type)]++
	}
	df.Metrics.KeptByType = counts

	// Metrics
	df.Metrics.MinConfidenceApplied = minConfidence
	df.Metrics.Nodes = n
	df.Metrics.Edges = len(df.Edges)
	df.Metrics.LongestPathLength = dist[sink] + 1
	df.Metrics.CriticalPath = critPath
	// WidthApprox = max nodes per depth
	depthCount := map[int]int{}
	maxW := 0
	for i := 0; i < n; i++ {
		d := depth[i]
		depthCount[d]++
		if depthCount[d] > maxW {
			maxW = depthCount[d]
		}
	}
    df.Metrics.WidthApprox = maxW

    // Quality metrics derived from task metadata
    // Evidence coverage: fraction of tasks with >=1 evidence entries
    if len(tasks) > 0 {
        withEv := 0
        verbFirst := 0
        for _, t := range tasks {
            if len(t.Evidence) > 0 { withEv++ }
            if isVerbFirst(t.Title) { verbFirst++ }
        }
        df.Metrics.EvidenceCoverage = float64(withEv) / float64(len(tasks))
        df.Metrics.VerbFirstPct = float64(verbFirst) / float64(len(tasks))
    }

    df.Analysis.OK = true
    return df, nil
}

// isVerbFirst applies a lightweight heuristic: consider titles verb-first if
// the first token matches a small set of common imperative verbs (lowercased).
func isVerbFirst(title string) bool {
    if strings.TrimSpace(title) == "" { return false }
    s := strings.ToLower(strings.TrimSpace(title))
    // Normalize common punctuation
    s = strings.TrimLeft(s, "#- ")
    fields := strings.Fields(s)
    if len(fields) == 0 { return false }
    first := fields[0]
    verbs := map[string]struct{}{
        "add":{}, "build":{}, "change":{}, "create":{}, "deploy":{}, "document":{},
        "enable":{}, "ensure":{}, "extract":{}, "fix":{}, "implement":{}, "improve":{},
        "migrate":{}, "optimize":{}, "refactor":{}, "remove":{}, "rename":{}, "replace":{},
        "setup":{}, "set":{}, "test":{}, "update":{}, "upgrade":{}, "validate":{}, "verify":{},
        "wire":{}, "simulate":{}, "analyze":{}, "export":{}, "render":{},
    }
    if _, ok := verbs[first]; ok { return true }
    // Special-case "set up" as a two-word verb
    if first == "set" && len(fields) > 1 && fields[1] == "up" { return true }
    return false
}

func contains(a []string, s string) bool {
	for _, x := range a {
		if x == s {
			return true
		}
	}
	return false
}

func lookupEdgeType(es []edgeRec, f, t string) string { // edgeRec is defined above
	for _, e := range es {
		if e.From == f && e.To == t {
			return e.Type
		}
	}
	return ""
}

func edgeTypeKey(v string) string {
	if strings.TrimSpace(v) == "" {
		return "unknown"
	}
	return v
}
