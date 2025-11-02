#!/usr/bin/env bash
set -euo pipefail

# Dependency checks (dot is optional and handled later)
for cmd in gh jq sed awk; do
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "ERROR: required command '$cmd' not found in PATH" >&2
    exit 127
  fi
done

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

LIMIT=${GH_LIST_LIMIT:-300}
echo "Fetching issues from $REPO (limit $LIMIT via GH_LIST_LIMIT) ..." >&2
if ! gh issue list --repo "$REPO" --state open -L "$LIMIT" --json number,title,labels,url > "$tmpdir/issues.json"; then
  echo "ERROR: failed to list issues for $REPO" >&2
  exit 1
fi

# Initialize edges file (truncate)
: >"$tmpdir/edges.txt"

# Batch fetch open issue bodies via GraphQL to avoid N+1 API calls
owner="${REPO%%/*}"; name="${REPO##*/}"
# Paginate through open issues (max 100 per page)
echo "[]" > "$tmpdir/all_nodes.json"
after=""
while :; do
  if [[ -z "$after" ]]; then
    query='query($owner:String!,$name:String!){ repository(owner:$owner,name:$name){ issues(first:100, states:OPEN, orderBy:{field:CREATED_AT, direction:DESC}){ pageInfo{ hasNextPage endCursor } nodes { number title body } } } }'
    res=$(gh api graphql -f owner="$owner" -f name="$name" -f query="$query")
  else
    query='query($owner:String!,$name:String!,$after:String!){ repository(owner:$owner,name:$name){ issues(first:100, states:OPEN, orderBy:{field:CREATED_AT, direction:DESC}, after:$after){ pageInfo{ hasNextPage endCursor } nodes { number title body } } } }'
    res=$(gh api graphql -f owner="$owner" -f name="$name" -f after="$after" -f query="$query")
  fi
  if [[ -z "$res" ]]; then
    echo "ERROR: failed to fetch issue page via GraphQL" >&2
    exit 1
  fi
  jq '.data.repository.issues.nodes' <<<"$res" > "$tmpdir/nodes.json"
  jq -s 'add' "$tmpdir/all_nodes.json" "$tmpdir/nodes.json" > "$tmpdir/all_nodes.tmp" && mv "$tmpdir/all_nodes.tmp" "$tmpdir/all_nodes.json"
  hasNext=$(jq -r '.data.repository.issues.pageInfo.hasNextPage' <<<"$res")
  after=$(jq -r '.data.repository.issues.pageInfo.endCursor' <<<"$res")
  [[ "$hasNext" == "true" && -n "$after" ]] || break
done

# Derive edges from titles and bodies
jq -r '.[] | @base64' "$tmpdir/all_nodes.json" | while read -r row; do
  _jq(){ echo "$row" | base64 --decode | jq -r "$1"; }
  num=$(_jq '.number')
  title=$(_jq '.title // ""')
  body=$(_jq '.body // ""')
  # Epic children from checklist lines
  if [[ "$title" =~ ^Epic: ]]; then
    printf '%s\n' "$body" | sed -n '1,800p' | awk '/^- \[[ x]\] #[0-9]+/ { for(i=1;i<=NF;i++) if ($i ~ /^#[0-9]+$/){ gsub("#","",$i); printf("epic %s -> %s\n", num, $i) }}' num="$num" >> "$tmpdir/edges.txt"
  fi
  # Dependency entries: lines that mention common dependency phrases and issue refs
  while read -r l; do
    if ! echo "$l" | grep -Eiq "blocked[- ]?by|blocks|depends on|requires|after:"; then
      continue
    fi
    while [[ "$l" =~ \#([0-9]+) ]]; do
      blk="${BASH_REMATCH[1]}"
      echo "block $num <- $blk" >> "$tmpdir/edges.txt"
      l="${l#*#$blk}" # advance using $blk (no extra braces)
    done
  done < <(printf '%s\n' "$body" | grep -Ei "blocked[- ]?by|blocks|depends on|requires|after:" || true)
done

# Build DOT
{
  echo 'digraph G {'
  echo '  rankdir=LR;'
  echo '  node [shape=box, style="rounded,filled", fontname=Helvetica, fontsize=10, fillcolor=white];'

  # Nodes with labels
  jq -r '.[] | [.number, (.title|gsub("\""; "\\\""))] | @tsv' "$tmpdir/issues.json" |
  while IFS=$'\t' read -r num title; do
    printf '  I%s [label="#%s: %s"];\n' "$num" "$num" "$title"
  done

  # Edges
  while read -r kind rest; do
    [[ -z "$kind" ]] && continue
    if [[ "$kind" == "epic" ]]; then
      # epic parent -> child
      parent=${rest%% *}; child=${rest##* }
      printf '  I%s -> I%s [color="gray50", style=dashed, label="parent"];\n' "$parent" "$child"
    elif [[ "$kind" == "block" ]]; then
      # block N <- B means edge B -> N
      target=${rest%% *}; blocker=${rest##* }
      printf '  I%s -> I%s [color="red", label="blocks"];\n' "$blocker" "$target"
    fi
  done < "$tmpdir/edges.txt"

  echo '}'
} > "$dot_out"

echo "Wrote $dot_out"

if command -v dot >/dev/null 2>&1; then
  dot -Tsvg "$dot_out" -o "$svg_out"
  echo "Wrote $svg_out"
elif command -v npx >/dev/null 2>&1; then
  # JS Graphviz fallback. Prefer stdout redirection form.
  if npx -y graphviz-cli -T svg "$dot_out" > "$svg_out" \
     || npx -y graphviz-cli dot -Tsvg "$dot_out" > "$svg_out" \
     || npx -y @viz-js/viz -T svg -o "$svg_out" "$dot_out"; then
    echo "Wrote $svg_out (via JS Graphviz)"
  else
    echo "graphviz JS fallback failed; SVG not generated" >&2
  fi
else
  echo "Neither 'dot' nor 'npx' available; skipped SVG generation" >&2
fi
