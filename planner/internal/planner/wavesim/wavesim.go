package wavesim

import (
    "sort"

    m "github.com/james/tasks-planner/internal/model"
)

// Generate returns waves as a slice of task-id slices. It respects Kahn layering via df.Nodes[].Depth
// and splits each layer into subwaves so tasks that contend on the same exclusive resource do not
// share a subwave. This is a planning-time preview only; it does not feed back into the DAG.
func Generate(df m.DagFile, tasks []m.Task) [][]string {
    // Build maps
    idToTask := make(map[string]m.Task, len(tasks))
    for _, t := range tasks { idToTask[t.ID] = t }
    // Group nodes by depth
    depthToIDs := map[int][]string{}
    depths := []int{}
    for _, n := range df.Nodes {
        depthToIDs[n.Depth] = append(depthToIDs[n.Depth], n.ID)
    }
    for d := range depthToIDs { depths = append(depths, d) }
    sort.Ints(depths)
    // For each depth, pack tasks into subwaves by exclusive resources
    var waves [][]string
    for _, d := range depths {
        ids := depthToIDs[d]
        sort.Strings(ids)
        // greedy packing
        type set map[string]struct{}
        subwaves := []struct{ ids []string; used set }{}
        for _, id := range ids {
            t := idToTask[id]
            exclusives := make(set)
            for _, r := range t.Resources.Exclusive { exclusives[r] = struct{}{} }
            placed := false
            for i := range subwaves {
                conflict := false
                for r := range exclusives {
                    if _, ok := subwaves[i].used[r]; ok { conflict = true; break }
                }
                if !conflict {
                    subwaves[i].ids = append(subwaves[i].ids, id)
                    for r := range exclusives { subwaves[i].used[r] = struct{}{} }
                    placed = true
                    break
                }
            }
            if !placed {
                used := make(set)
                for r := range exclusives { used[r] = struct{}{} }
                subwaves = append(subwaves, struct{ ids []string; used set }{ids: []string{id}, used: used})
            }
        }
        for _, sw := range subwaves {
            waves = append(waves, sw.ids)
        }
    }
    return waves
}

