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

# Resolve OUT_DIR to an absolute path before cd'ing into planner to avoid
# accidentally creating planner/planner/plans when OUT_DIR starts with 'planner/'.
mkdir -p "$OUT_DIR" "$DOCS_DIR"
OUT_ABS=$(cd "$(dirname "$OUT_DIR")" && pwd)/$(basename "$OUT_DIR")

(cd planner && go run ./cmd/tasksd plan --out "${OUT_ABS}" --repo ..)
(cd planner && go run ./cmd/tasksd validate --dir "${OUT_ABS}")
(cd planner && go run ./cmd/tasksd export-dot --dir "${OUT_ABS}")

# Render DOT → SVG
plan_dot="${OUT_ABS}/dag.dot"
runtime_dot="${OUT_ABS}/runtime.dot"
plan_svg="${DOCS_DIR}/plan-dag.svg"
runtime_svg="${DOCS_DIR}/runtime.svg"

render(){
  local dotfile="$1" svgfile="$2"
  if ! test -f "$dotfile"; then return 0; fi
  if command -v dot >/dev/null 2>&1; then
    dot -Tsvg "$dotfile" -o "$svgfile"
  elif command -v npx >/dev/null 2>&1; then
    if npx -y @hpcc-js/wasm-graphviz-cli -K dot -T svg "$dotfile" > "$svgfile" \
       || npx -y viz.js-cli -Kdot -Tsvg "$dotfile" -o "$svgfile" \
       || npx -y graphviz-cli -K dot -T svg "$dotfile" > "$svgfile"; then
      :
    else
      echo "graphviz JS fallback failed; skipping $dotfile" >&2
      return 0
    fi
  else
    echo "WARN: no 'dot' or 'npx'; skipping render for $dotfile" >&2
    return 0
  fi
  # Sanity header (warn-only)
  if [ -s "$svgfile" ] && ! (head -c 5 "$svgfile" | grep -q '<?xml' || head -n1 "$svgfile" | grep -q '^<svg'); then
    echo "WARN: invalid SVG header in $svgfile" >&2
  fi
}

render "$plan_dot" "$plan_svg"
render "$runtime_dot" "$runtime_svg"

echo "Wrote $plan_svg and $runtime_svg"
