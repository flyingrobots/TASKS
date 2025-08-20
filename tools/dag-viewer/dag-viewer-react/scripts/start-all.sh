#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Starting DAG Viewer Services...${NC}"

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

# Kill processes on relevant ports
echo -e "${GREEN}Clearing ports...${NC}"
kill_port 3456  # WebSocket server port
kill_port 8080  # Alternative WebSocket server port
kill_port 5173  # Vite default port
kill_port 5174  # Vite alternate port

# Extra check: forcefully kill anything on 8080
echo -e "${YELLOW}Ensuring port 8080 is free...${NC}"
lsof -ti:8080 | xargs -r kill -9 2>/dev/null || true

# Start the DAG state server
echo -e "${GREEN}Starting DAG state server on port 8080...${NC}"
PORT=8080 DAG_FILE=/Users/james/git/pf3/docs/plans/professional-quality-code/dag.json node ../dag-state-server-full.js 2>&1 > ./logs.txt &
SERVER_PID=$!
echo -e "${GREEN}DAG state server started (PID: $SERVER_PID)${NC}"

# Give the server a moment to start
sleep 2

# Start the React dev server
echo -e "${GREEN}Starting React dev server...${NC}"
npm run dev &
DEV_PID=$!
echo -e "${GREEN}React dev server started (PID: $DEV_PID)${NC}"

# Function to handle cleanup on exit
cleanup() {
    echo -e "\n${YELLOW}Shutting down services...${NC}"
    kill $SERVER_PID 2>/dev/null
    kill $DEV_PID 2>/dev/null
    kill_port 3456
    kill_port 8080
    kill_port 5173
    kill_port 5174
    echo -e "${GREEN}Services stopped${NC}"
    exit 0
}

# Set up trap to cleanup on Ctrl+C
trap cleanup INT TERM

# Wait for both processes
echo -e "${GREEN}Services are running. Press Ctrl+C to stop all services.${NC}"
echo -e "${GREEN}WebSocket server: ws://localhost:8080${NC}"
echo -e "${GREEN}React app: http://localhost:5173${NC}"

# Keep script running
wait