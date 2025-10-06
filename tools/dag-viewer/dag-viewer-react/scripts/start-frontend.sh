#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Starting React Dev Server...${NC}"

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

# Kill any existing Vite dev server
echo -e "${GREEN}Clearing Vite ports...${NC}"
kill_port 5173
kill_port 5174

# Start the React dev server
echo -e "${GREEN}Starting React dev server...${NC}"
echo -e "${GREEN}React app will be available at: http://localhost:5173${NC}"
echo -e "${YELLOW}Make sure the DAG server is running on port 8080 for live updates${NC}"
npm run dev