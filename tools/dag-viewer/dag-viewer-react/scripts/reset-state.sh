#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Resetting DAG state...${NC}"

# Get config file path (can be overridden by passing it as first argument)
CONFIG_FILE="${1:-./config.json}"

# Extract state file path from config
STATE_FILE=$(node -e "
  const config = require('$CONFIG_FILE');
  console.log(config.paths.stateFile);
" 2>/dev/null)

if [ -z "$STATE_FILE" ]; then
  STATE_FILE="/Users/james/git/TASKS/tools/dag-viewer/dag-state-full.json"
  echo -e "${YELLOW}Could not read config, using default state file: $STATE_FILE${NC}"
fi

echo -e "${GREEN}Clearing state file: $STATE_FILE${NC}"

# Create empty state
cat > "$STATE_FILE" << 'EOF'
{
  "tasks": {},
  "agents": {},
  "commits": [],
  "dag": null,
  "gitStats": {
    "totalCommits": 0,
    "totalLinesAdded": 0,
    "totalLinesRemoved": 0,
    "fileChangeFrequency": {},
    "agentCommits": {},
    "hourlyActivity": [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
    "dailyActivity": {},
    "hotFiles": {},
    "commitSizes": []
  },
  "events": []
}
EOF

echo -e "${GREEN}State reset successfully!${NC}"
echo -e "${YELLOW}You may need to restart the server for changes to take effect.${NC}"