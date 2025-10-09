#!/usr/bin/env bash
set -euo pipefail

# From repo root or this directory
ROOT_DIR="$(cd "$(dirname "$0")"/../.. && pwd)"
EX_DIR="$ROOT_DIR/examples/dot-export"

echo "Generating DOT from dag.json + tasks.json..."
pushd "$ROOT_DIR/planner" >/dev/null
go run ./cmd/tasksd export-dot --dir "$EX_DIR"
popd >/dev/null

if command -v dot >/dev/null 2>&1; then
  echo "Rendering SVG via Graphviz..."
  dot -Tsvg "$EX_DIR/dag.dot" -o "$EX_DIR/dag.svg"
  echo "Done: $EX_DIR/dag.svg"
else
  echo "Graphviz 'dot' not found; generated $EX_DIR/dag.dot only." >&2
fi

echo "OK"
