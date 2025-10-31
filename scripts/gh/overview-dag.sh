#!/usr/bin/env bash
set -euo pipefail

# Overview: Generates a project-wide DAG of GitHub issues by reading:
# - Epic → Child links from epic issue checklists ("- [ ] #<num>")
# - Blocked-by from body lines containing "Blocked by" followed by issue refs (e.g., "#23, #24")
#
# Outputs:
# - docs/overview-dag.dot
# - docs/overview-dag.svg (requires graphviz `dot`)
#
# Requirements:
# - gh CLI authenticated (GITHUB_TOKEN in CI works)
# - jq, sed, awk, dot

REPO=${REPO:-${GITHUB_REPOSITORY:-}}
if [[ -z "${REPO}" ]]; then
  # Try from git remote
  remote=$(git remote get-url origin 2>/dev/null || true)
  case "$remote" in
    git@github.com:*) REPO=${remote#git@github.com:}; REPO=${REPO%.git} ;;
    https://github.com/*) REPO=${remote#https://github.com/} ; REPO=${REPO%.git} ;;
    http://github.com/*) REPO=${remote#http://github.com/} ; REPO=${REPO%.git} ;;
  esac
fi
if [[ -z "${REPO}" ]]; then
  echo "REPO not set and could not detect from git remote" >&2; exit 1
fi

mkdir -p docs
dot_out=${1:-docs/overview-dag.dot}
svg_out=${2:-docs/overview-dag.svg}

# Fetch open issues (increase limit if needed)
tmpdir=$(mktemp -d)
trap 'rm -rf "$tmpdir"' EXIT

echo "Fetching issues from $REPO ..." >&2
gh issue list --repo "$REPO" --state open -L 300 --json number,title,labels,url > "$tmpdir/issues.json"

# Get epic checklist bodies (we use these to infer epic→child edges)
epics=$(jq -r '.[] | select(.title|test("^Epic:")) | .number' "$tmpdir/issues.json")
>"$tmpdir/edges.txt"

for n in $epics; do
  body=$(gh issue view --repo "$REPO" "$n" --json body --jq .body 2>/dev/null || echo "")
  # Parse "- [ ] #<num> ..." lines
  while IFS= read -r line; do
    if [[ "$line" =~ ^-\ \[.\]\ \#([0-9]+) ]]; then
      child="${BASH_REMATCH[1]}"
      echo "epic $n -> $child" >> "$tmpdir/edges.txt"
    fi
  done < <(printf '%s
' "$body" | sed -n '1,400p')
done

# Parse "Blocked by" lines in bodies (best-effort)
all_nums=$(jq -r '.[].number' "$tmpdir/issues.json")
for n in $all_nums; do
  body=$(gh issue view --repo "$REPO" "$n" --json body --jq .body 2>/dev/null || echo "")
  mapfile -t lines < <(printf '%s
' "$body" | grep -i "Blocked by" || true)
  for l in "${lines[@]}"; do
    # Extract all #<num> tokens on that line
    while [[ "$l" =~ \#([0-9]+) ]]; do
      blk="${BASH_REMATCH[1]}"
      echo "block $n <- $blk" >> "$tmpdir/edges.txt"
      l="${l#*#${blk}}" # advance
    done
  done
done

# Build DOT
{
  echo 'digraph G {'
  echo '  rankdir=LR;'
  echo '  node [shape=box, style="rounded,filled", fontname=Helvetica, fontsize=10, fillcolor=white];'

  # Nodes with labels
  jq -r '.[] | [.number, (.title|gsub("\""; "\\\""))] | @tsv' "$tmpdir/issues.json" |
  while IFS=$'\t' read -r num title; do
    printf '  I%s [label="#%s: %s"];
' "$num" "$num" "$title"
  done

  # Edges
  while read -r kind rest; do
    [[ -z "$kind" ]] && continue
    if [[ "$kind" == "epic" ]]; then
      # epic parent -> child
      parent=${rest%% *}; child=${rest##* }
      printf '  I%s -> I%s [color="gray50", style=dashed, label="parent"];
' "$parent" "$child"
    elif [[ "$kind" == "block" ]]; then
      # block N <- B means edge B -> N
      target=${rest%% *}; blocker=${rest##* }
      printf '  I%s -> I%s [color="red", label="blocks"];
' "$blocker" "$target"
    fi
  done < "$tmpdir/edges.txt"

  echo '}'
} > "$dot_out"

echo "Wrote $dot_out"

if command -v dot >/dev/null 2>&1; then
  dot -Tsvg "$dot_out" -o "$svg_out"
  echo "Wrote $svg_out"
else
  echo "Graphviz 'dot' not found; skipped SVG generation" >&2
fi
