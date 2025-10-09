package dot

import (
    "fmt"
    "sort"
    "strings"

    m "github.com/james/tasks-planner/internal/model"
)

// Options controls DOT label styling.
type Options struct {
    // NodeLabel: "id" | "title" | "id-title"
    NodeLabel string
    // EdgeLabel: "none" | "type"
    EdgeLabel string
}

func (o Options) normalize() Options {
    out := o
    switch o.NodeLabel {
    case "id", "title", "id-title":
    default:
        out.NodeLabel = "id-title"
    }
    switch o.EdgeLabel {
    case "none", "type":
    default:
        out.EdgeLabel = "type"
    }
    return out
}

// FromDag renders a DOT/Graphviz string from a DagFile and an optional title map.
// titles[id] -> human label to append after the ID. If missing, only the ID is used.
// Only non-transitive edges are emitted to keep the view minimal.
func FromDag(d m.DagFile, titles map[string]string) string {
    return FromDagWithOptions(d, titles, Options{NodeLabel: "id-title", EdgeLabel: "type"})
}

// FromDagWithOptions renders DOT with label options.
func FromDagWithOptions(d m.DagFile, titles map[string]string, opts Options) string {
    opts = opts.normalize()
    var b strings.Builder
    // Graph header
    b.WriteString("digraph G {\n")
    b.WriteString("  rankdir=LR;\n")
    b.WriteString("  splines=true;\n")
    b.WriteString("  node [shape=box, style=\"rounded,filled\", fillcolor=white, fontname=Helvetica];\n\n")

    // Nodes (stable order by ID for deterministic output)
    ids := make([]string, 0, len(d.Nodes))
    // Map for quick lookup of critical path
    crit := make(map[string]bool, len(d.Nodes))
    for _, n := range d.Nodes {
        ids = append(ids, n.ID)
        if n.CriticalPath {
            crit[n.ID] = true
        }
    }
    sort.Strings(ids)
    for _, id := range ids {
        title := titles[id]
        label := nodeLabelFor(id, title, opts.NodeLabel)
        attr := ""
        if crit[id] {
            attr = ", color=red, penwidth=2"
        }
        fmt.Fprintf(&b, "  %s [label=\"%s\"%s];\n", safeID(id), escape(label), attr)
    }
    b.WriteString("\n")

    // Edges: only non-transitive
    type edge struct{ from, to, color, style, label string }
    edges := make([]edge, 0, len(d.Edges))
    for _, e := range d.Edges {
        if e.Transitive {
            continue
        }
        // Highlight critical path edges if both nodes are critical
        color := ""
        style := ""
        if crit[e.From] && crit[e.To] {
            color = "red"
        }
        // Edge label per options
        lbl := ""
        if opts.EdgeLabel == "type" {
            lbl = strings.TrimSpace(e.Type)
        }
        edges = append(edges, edge{from: e.From, to: e.To, color: color, style: style, label: lbl})
    }
    // Deterministic ordering
    sort.Slice(edges, func(i, j int) bool {
        if edges[i].from == edges[j].from {
            return edges[i].to < edges[j].to
        }
        return edges[i].from < edges[j].from
    })
    for _, e := range edges {
        attrs := kv(map[string]string{
            "color": e.color,
            "style": e.style,
            "label": e.label,
        })
        fmt.Fprintf(&b, "  %s -> %s%s;\n", safeID(e.from), safeID(e.to), attrs)
    }

    b.WriteString("}\n")
    return b.String()
}

func kv(m map[string]string) string {
    parts := make([]string, 0, len(m))
    for k, v := range m {
        if v == "" {
            continue
        }
        parts = append(parts, fmt.Sprintf("%s=\"%s\"", k, escape(v)))
    }
    if len(parts) == 0 {
        return ""
    }
    sort.Strings(parts)
    return " [" + strings.Join(parts, ", ") + "]"
}

func escape(s string) string {
    // minimal escaping for quotes and newlines
    s = strings.ReplaceAll(s, "\\", "\\\\")
    s = strings.ReplaceAll(s, "\"", "\\\"")
    s = strings.ReplaceAll(s, "\n", "\\n")
    return s
}

// safeID ensures the node identifiers are DOT-safe by quoting if needed.
// We keep simple: IDs produced by planner are expected alnum+dot/colon; return as-is.
func safeID(id string) string { return id }

// FromCoordinator renders a DOT/Graphviz string from a Coordinator (runtime view).
// Labels default to "ID: Title" if Title is non-empty, otherwise just ID.
func FromCoordinator(c m.Coordinator) string {
    return FromCoordinatorWithOptions(c, Options{NodeLabel: "id-title", EdgeLabel: "type"})
}

// FromCoordinatorWithOptions renders DOT from runtime view with label options.
func FromCoordinatorWithOptions(c m.Coordinator, opts Options) string {
    opts = opts.normalize()
    var b strings.Builder
    b.WriteString("digraph G {\n")
    b.WriteString("  rankdir=LR;\n")
    b.WriteString("  splines=true;\n")
    b.WriteString("  node [shape=box, style=\"rounded,filled\", fillcolor=white, fontname=Helvetica];\n\n")

    // Nodes
    ids := make([]string, 0, len(c.Graph.Nodes))
    titles := make(map[string]string, len(c.Graph.Nodes))
    for _, t := range c.Graph.Nodes {
        ids = append(ids, t.ID)
        titles[t.ID] = t.Title
    }
    sort.Strings(ids)
    for _, id := range ids {
        title := titles[id]
        label := nodeLabelFor(id, title, opts.NodeLabel)
        fmt.Fprintf(&b, "  %s [label=\"%s\"];\n", safeID(id), escape(label))
    }
    b.WriteString("\n")

    // Edges
    type edge struct{ from, to, label string }
    edges := make([]edge, 0, len(c.Graph.Edges))
    for _, e := range c.Graph.Edges {
        lbl := ""
        if opts.EdgeLabel == "type" {
            lbl = strings.TrimSpace(e.Type)
        }
        edges = append(edges, edge{from: e.From, to: e.To, label: lbl})
    }
    sort.Slice(edges, func(i, j int) bool {
        if edges[i].from == edges[j].from { return edges[i].to < edges[j].to }
        return edges[i].from < edges[j].from
    })
    for _, e := range edges {
        attrs := kv(map[string]string{"label": e.label})
        fmt.Fprintf(&b, "  %s -> %s%s;\n", safeID(e.from), safeID(e.to), attrs)
    }

    b.WriteString("}\n")
    return b.String()
}

func nodeLabelFor(id, title, mode string) string {
    switch mode {
    case "id":
        return id
    case "title":
        if title != "" { return title }
        return id
    case "id-title":
        if title != "" { return fmt.Sprintf("%s: %s", id, title) }
        return id
    default:
        if title != "" { return fmt.Sprintf("%s: %s", id, title) }
        return id
    }
}
