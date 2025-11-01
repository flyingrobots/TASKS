#!/usr/bin/env bash
set -euo pipefail

# Generate planner DAG SVGs from T.A.S.K.S. artifacts and write to docs/.
# Requires Go toolchain for tasksd. Uses dot if available, falls back to viz.js via npx.

for cmd in go jq; do
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "ERROR: required command '$cmd' not found" >&2
    exit 127
  fi
done

OUT_DIR=${1:-planner/plans}
DOCS_DIR=${2:-docs}

mkdir -p "$OUT_DIR" "$DOCS_DIR"

(cd planner && go run ./cmd/tasksd plan --out "${OUT_DIR}" --repo ..)
(cd planner && go run ./cmd/tasksd validate --dir "${OUT_DIR}")
(cd planner && go run ./cmd/tasksd export-dot --dir "${OUT_DIR}")

# Render DOT → SVG
plan_dot="${OUT_DIR}/dag.dot"
runtime_dot="${OUT_DIR}/runtime.dot"
plan_svg="${DOCS_DIR}/plan-dag.svg"
runtime_svg="${DOCS_DIR}/runtime.svg"

render(){
  local dotfile="$1" svgfile="$2"
  if ! test -f "$dotfile"; then return 0; fi
  if command -v dot >/dev/null 2>&1; then
    dot -Tsvg "$dotfile" -o "$svgfile"
  elif command -v npx >/dev/null 2>&1; then
    npx -y graphviz-cli -T svg -o "$svgfile" "$dotfile"
  else
    echo "WARN: no 'dot' or 'npx'; skipping render for $dotfile" >&2
    return 0
  fi
  # sanity header
  if ! (head -c 5 "$svgfile" | grep -q '<?xml' || head -n1 "$svgfile" | grep -q '^<svg'); then
    echo "ERROR: invalid SVG header in $svgfile" >&2
    return 1
  fi
}

render "$plan_dot" "$plan_svg"
render "$runtime_dot" "$runtime_svg"

echo "Wrote $plan_svg and $runtime_svg"
