#!/bin/bash
pkill -f "node ../dag-state-server-full.js" 2>/dev/null
pkill -f vite 2>/dev/null
lsof -ti:8080 | xargs kill -9 2>/dev/null || true
lsof -ti:5173 | xargs kill -9 2>/dev/null || true

PORT=8080 DAG_FILE=/Users/james/git/pf3/docs/plans/professional-quality-code/dag.json node ../dag-state-server-full.js &
BACKEND=$!
sleep 2
npm run dev &
FRONTEND=$!

echo "Backend PID: $BACKEND"
echo "Frontend PID: $FRONTEND"
echo "Servers running at:"
echo "  Backend: http://localhost:8080"
echo "  Frontend: http://localhost:5173"

wait