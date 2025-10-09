# DOT Exporter

Tiny helper to render planner artifacts as Graphviz DOT for quick visualization.

## What it does
- Converts `dag.json` + `tasks.json` into a DOT graph.
- Emits only non‑transitive, structural edges to preserve DAG minimality.
- Highlights critical‑path nodes/edges in red.
- Labels nodes as `ID: Title` when task titles are available.

## Usage
From repo root:

```
cd planner
go run ./cmd/tasksd export-dot --dag ./plans/dag.json --tasks ./plans/tasks.json --out ./plans/dag.dot
# Render with Graphviz (install graphviz):
dot -Tsvg ./plans/dag.dot -o ./plans/dag.svg
```

Directory mode (writes `dag.dot` and `runtime.dot` when artifacts exist):

```
cd planner
go run ./cmd/tasksd export-dot --dir ./plans
```

Coordinator view directly:

```
cd planner
go run ./cmd/tasksd export-dot --coordinator ./plans/coordinator.json --out ./plans/runtime.dot
```

Label customization:

```
# Node labels: id | title | id-title (default)
# Edge labels: none | type (default)
go run ./cmd/tasksd export-dot --dir ./plans --node-label title --edge-label none
```

## Notes
- Output is deterministic (stable node/edge ordering).
- If you don’t pass `--out`, DOT is printed to stdout.
- Critical path information is read from `dag.json.nodes[].critical_path`.
