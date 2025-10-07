#!/usr/bin/env node

import { WebSocketServer } from 'ws';
import { createServer } from 'http';
import { writeFileSync, readFileSync, existsSync } from 'fs';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const PORT = process.env.PORT || 8080;
const ANALYTICS_FILE = process.env.ANALYTICS_FILE || join(__dirname, 'agent-analytics.json');

// In-memory state
const taskStates = {};
const agentStats = {};
const taskHistory = [];

// WebSocket clients
const wsClients = new Set();

// Load existing analytics
function loadAnalytics() {
  if (existsSync(ANALYTICS_FILE)) {
    try {
      const data = JSON.parse(readFileSync(ANALYTICS_FILE, 'utf8'));
      Object.assign(agentStats, data.agentStats || {});
      taskHistory.push(...(data.taskHistory || []));
      console.log('Loaded analytics from', ANALYTICS_FILE);
    } catch (e) {
      console.error('Error loading analytics:', e);
    }
  }
}

// Save analytics periodically
function saveAnalytics() {
  try {
    writeFileSync(ANALYTICS_FILE, JSON.stringify({
      agentStats,
      taskHistory,
      savedAt: new Date().toISOString()
    }, null, 2));
  } catch (e) {
    console.error('Error saving analytics:', e);
  }
}

// Initialize agent if not exists
function initAgent(agentId) {
  if (!agentStats[agentId]) {
    agentStats[agentId] = {
      name: agentId,
      tasksCompleted: 0,
      tasksFailed: 0,
      tasksStarted: 0,
      totalTime: 0,
      completionTimes: [],
      path: [],
      currentTask: null,
      status: 'idle',
      firstSeen: Date.now(),
      lastSeen: Date.now()
    };
  }
}

// Calculate agent metrics
function calculateMetrics(agentId) {
  const stats = agentStats[agentId];
  if (!stats) return null;
  
  const times = stats.completionTimes;
  const medianTime = times.length > 0 
    ? times.sort((a,b) => a-b)[Math.floor(times.length / 2)]
    : 0;
  
  return {
    ...stats,
    averageTime: times.length > 0 ? stats.totalTime / times.length : 0,
    medianTime: medianTime,
    successRate: stats.tasksStarted > 0 
      ? (stats.tasksCompleted / stats.tasksStarted * 100).toFixed(1)
      : 0,
    efficiency: medianTime > 0 ? (1000000 / medianTime).toFixed(2) : 0
  };
}

// Get leaderboard
function getLeaderboard() {
  return Object.keys(agentStats)
    .map(agentId => calculateMetrics(agentId))
    .sort((a, b) => {
      // Sort by completed tasks, then by median time
      if (b.tasksCompleted !== a.tasksCompleted) {
        return b.tasksCompleted - a.tasksCompleted;
      }
      return a.medianTime - b.medianTime;
    });
}

// Process task update
function processUpdate(data) {
  const { agent, task, status, timestamp } = data;
  
  initAgent(agent);
  
  const agentData = agentStats[agent];
  agentData.lastSeen = Date.now();
  
  // Track task state
  const prevState = taskStates[task];
  taskStates[task] = {
    status,
    agent,
    timestamp,
    startTime: prevState?.startTime
  };
  
  // Update agent stats based on status
  switch(status) {
    case 'started':
      agentData.tasksStarted++;
      agentData.currentTask = task;
      agentData.status = 'working';
      agentData.path.push({ task, status, timestamp });
      taskStates[task].startTime = timestamp;
      break;
      
    case 'completed':
      if (taskStates[task].startTime) {
        const duration = timestamp - taskStates[task].startTime;
        agentData.completionTimes.push(duration);
        agentData.totalTime += duration;
        
        // Keep only last 100 completion times
        if (agentData.completionTimes.length > 100) {
          const removed = agentData.completionTimes.shift();
          agentData.totalTime -= removed;
        }
      }
      agentData.tasksCompleted++;
      agentData.currentTask = null;
      agentData.status = 'idle';
      agentData.path.push({ task, status, timestamp });
      break;
      
    case 'failed':
      agentData.tasksFailed++;
      agentData.currentTask = null;
      agentData.status = 'idle';
      agentData.path.push({ task, status, timestamp });
      break;
      
    case 'blocked':
      agentData.status = 'blocked';
      agentData.path.push({ task, status, timestamp });
      break;
      
    default:
      agentData.path.push({ task, status, timestamp });
  }
  
  // Keep path to last 50 entries
  if (agentData.path.length > 50) {
    agentData.path = agentData.path.slice(-50);
  }
  
  // Add to history
  taskHistory.push({
    agent,
    task,
    status,
    timestamp,
    date: new Date(timestamp * 1000).toISOString()
  });
  
  // Keep history to last 1000 entries
  if (taskHistory.length > 1000) {
    taskHistory.shift();
  }
}

// Broadcast to all WebSocket clients
function broadcast(type, data) {
  const message = JSON.stringify({ type, ...data });
  
  wsClients.forEach(client => {
    if (client.readyState === 1) {
      client.send(message);
    }
  });
}

// HTTP server
const server = createServer((req, res) => {
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Content-Type', 'application/json');
  
  if (req.method === 'POST' && req.url === '/update') {
    let body = '';
    
    req.on('data', chunk => {
      body += chunk.toString();
    });
    
    req.on('end', () => {
      try {
        const data = JSON.parse(body);
        
        // Process the update
        processUpdate(data);
        
        // Broadcast updates
        broadcast('task-update', {
          task: data.task,
          status: data.status,
          agent: data.agent
        });
        
        broadcast('agent-update', {
          agent: data.agent,
          stats: calculateMetrics(data.agent)
        });
        
        broadcast('leaderboard-update', {
          leaderboard: getLeaderboard()
        });
        
        // Respond
        res.writeHead(200);
        res.end(JSON.stringify({ ok: true }));
        
        console.log(`${data.agent} -> ${data.task}: ${data.status}`);
        
      } catch (error) {
        res.writeHead(400);
        res.end(JSON.stringify({ ok: false, error: error.message }));
      }
    });
    
  } else if (req.url === '/state') {
    res.writeHead(200);
    res.end(JSON.stringify({ taskStates, agentStats: Object.fromEntries(
      Object.entries(agentStats).map(([k, v]) => [k, calculateMetrics(k)])
    )}));
    
  } else if (req.url === '/leaderboard') {
    res.writeHead(200);
    res.end(JSON.stringify(getLeaderboard()));
    
  } else if (req.url === '/history') {
    res.writeHead(200);
    res.end(JSON.stringify(taskHistory.slice(-100)));
    
  } else if (req.url?.startsWith('/agent/')) {
    const agentId = req.url.split('/')[2];
    const metrics = calculateMetrics(agentId);
    if (metrics) {
      res.writeHead(200);
      res.end(JSON.stringify(metrics));
    } else {
      res.writeHead(404);
      res.end(JSON.stringify({ ok: false, error: 'Agent not found' }));
    }
    
  } else {
    res.writeHead(404);
    res.end(JSON.stringify({ ok: false, error: 'Not found' }));
  }
});

// WebSocket server
const wss = new WebSocketServer({ server });

wss.on('connection', (ws) => {
  console.log('Browser connected');
  wsClients.add(ws);
  
  // Send initial state
  ws.send(JSON.stringify({
    type: 'init',
    taskStates,
    agentStats: Object.fromEntries(
      Object.entries(agentStats).map(([k, v]) => [k, calculateMetrics(k)])
    ),
    leaderboard: getLeaderboard()
  }));
  
  ws.on('close', () => {
    wsClients.delete(ws);
  });
  
  ws.on('error', (err) => {
    console.error('WebSocket error:', err);
    wsClients.delete(ws);
  });
});

// Load analytics on startup
loadAnalytics();

// Save analytics every 30 seconds
setInterval(saveAnalytics, 30000);

// Graceful shutdown
process.on('SIGINT', () => {
  console.log('\nSaving analytics...');
  saveAnalytics();
  process.exit(0);
});

server.listen(PORT, () => {
  console.log(`
╔════════════════════════════════════════════════════════╗
║            DAG STATE SERVER WITH ANALYTICS            ║
╠════════════════════════════════════════════════════════╣
║  POST http://localhost:${PORT}/update                  ║
║  GET  http://localhost:${PORT}/state                   ║
║  GET  http://localhost:${PORT}/leaderboard             ║
║  GET  http://localhost:${PORT}/history                 ║
║  GET  http://localhost:${PORT}/agent/{agentId}         ║
║  WS   ws://localhost:${PORT}                           ║
╚════════════════════════════════════════════════════════╝

Example:
curl -X POST http://localhost:${PORT}/update \\
  -H "Content-Type: application/json" \\
  -d '{"agent":"mr_clean","task":"P1.T003","status":"started","timestamp":${Date.now()/1000}}'
  
curl http://localhost:${PORT}/leaderboard
  `);
});