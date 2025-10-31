#!/usr/bin/env bash
set -euo pipefail

# Dependency checks
for cmd in gh jq; do
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "ERROR: required command '$cmd' not found in PATH" >&2
    exit 127
  fi
done

# Seed rough Depth values for the Project v2 titled "TASKS + SLAPS OVERVIEW".
# Usage:
#   GITHUB_TOKEN=... [OWNER=owner] [REPO_NAME=repo] ./scripts/gh/seed-overview-depths.sh
#
# Notes:
# - OWNER and REPO_NAME default to values derived from the current git remote.

owner=${OWNER:-$(git remote get-url origin | sed -E 's#.*github.com[:/]([^/]+)/.*#\1#')}
repo=${REPO_NAME:-$(git remote get-url origin | sed -E 's#.*/([^/]+?)(\.git)?$#\1#' | sed 's/.git$//')}

proj_id=$(gh api graphql -f query='query($login:String!){user(login:$login){projectsV2(first:50){nodes{id title}}}}' -F login="$owner" --jq '.data.user.projectsV2.nodes[] | select(.title=="TASKS + SLAPS OVERVIEW") | .id' 2>/dev/null || true)
if [[ -z "$proj_id" ]]; then
  proj_id=$(gh api graphql -f query='query($login:String!){organization(login:$login){projectsV2(first:50){nodes{id title}}}}' -F login="$owner" --jq '.data.organization.projectsV2.nodes[] | select(.title=="TASKS + SLAPS OVERVIEW") | .id')
fi

fields=$(gh api graphql -f query='query($id:ID!){node(id:$id){... on ProjectV2{fields(first:50){nodes{__typename ... on ProjectV2Field{ id name dataType }}}}}}' -F id="$proj_id" --jq '.data.node.fields.nodes')
depth_id=$(echo "$fields" | jq -r '.[] | select(.name=="Depth") | .id')
if [[ -z "$depth_id" || "$depth_id" == null ]]; then
  echo "Depth field is missing on the overview project; create it in the UI or via GraphQL first." >&2
  exit 1
fi

add_or_get_item(){ local num="$1"; local cid=$(gh issue view --repo "$owner/$repo" "$num" --json id --jq .id); gh api graphql -f query='mutation($proj:ID!,$content:ID!){ addProjectV2ItemById(input:{projectId:$proj, contentId:$content}){ item { id } } }' -F proj="$proj_id" -F content="$cid" --jq '.data.addProjectV2ItemById.item.id' 2>/dev/null || true; }
set_depth(){ local num="$1" depth="$2"; local item=$(add_or_get_item "$num"); if [[ -z "$item" || "$item" == null ]]; then echo "Skip #$num"; return; fi; gh api graphql -f query='mutation($proj:ID!,$item:ID!,$field:ID!,$num:Float!){ updateProjectV2ItemFieldValue(input:{projectId:$proj, itemId:$item, fieldId:$field, value:{number:$num}}){ projectV2Item{ id } } }' -F proj="$proj_id" -F item="$item" -F field="$depth_id" -F num="$depth" >/dev/null && echo "Depth $depth set for #$num" || echo "Failed #$num"; }

# Known depths (rough), adjust over time
pairs=(
  15:0 16:0 17:0 18:0 19:1 20:2 21:2 22:0 23:1 24:1 25:2 26:3 27:0 28:0 29:0 30:0 31:0
  40:0 41:0 42:0 43:0 44:0 45:0 46:1 47:1 48:2 49:2 50:1 51:3 52:2 53:2 54:3 55:3 56:2 57:3 58:3 59:2 60:3 61:3
)
for kv in "${pairs[@]}"; do num=${kv%%:*}; d=${kv#*:}; set_depth "$num" "$d"; done
