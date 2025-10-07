#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get config file path (can be overridden by passing it as first argument)
CONFIG_FILE="${1:-./config.json}"

echo -e "${YELLOW}Starting DAG State Server...${NC}"

# Function to kill process on port
kill_port() {
    local port=$1
    local pid=$(lsof -ti:$port)
    if [ ! -z "$pid" ]; then
        echo -e "${YELLOW}Killing process on port $port (PID: $pid)${NC}"
        kill -9 $pid 2>/dev/null
        sleep 1
    fi
}

# Kill any existing process on port 8080
echo -e "${GREEN}Clearing port 8080...${NC}"
kill_port 8080

# Start the DAG state server with config file
echo -e "${GREEN}Starting DAG state server using config: $CONFIG_FILE${NC}"
node ../dag-state-server-full.js "$CONFIG_FILE"