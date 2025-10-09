# DOT Export Example

Minimal, runnable example that exports a DOT graph from `dag.json` and `tasks.json` and renders it to SVG.

## Files
- `tasks.json` — three tasks with titles
- `dag.json` — structural edges (one transitive edge included to show it’s suppressed)
- `render.sh` — convenience script to export DOT and render via Graphviz

## Quick Start

From repo root:

```bash
cd examples/dot-export
./render.sh
```

Outputs:
- `dag.dot` — Graphviz DOT
- `dag.svg` — rendered SVG (requires Graphviz `dot`)

If `dot` isn’t installed, you can still inspect the DOT:

```bash
cd examples/dot-export
cd .. && go run ../planner/cmd/tasksd export-dot --dag ./dot-export/dag.json --tasks ./dot-export/tasks.json
```
