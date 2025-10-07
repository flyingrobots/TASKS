#!/usr/bin/env node

import { WebSocketServer } from 'ws';
import { createServer } from 'http';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const PORT = process.env.PORT || 8080;

// In-memory state - just task -> status mapping
const taskStates = {};

// WebSocket clients
const wsClients = new Set();

// State colors
const statusColors = {
  pending: '#94a3b8',
  started: '#22c55e',    // GREEN for started
  in_progress: '#3b82f6',
  failed: '#ef4444',
  blocked: '#f97316',    // ORANGE for blocked
  completed: '#16a34a'
};

// Broadcast to all WebSocket clients
function broadcast(taskId, status) {
  const message = JSON.stringify({
    taskId,
    status,
    color: statusColors[status]
  });
  
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
        
        // Store state
        taskStates[data.task] = {
          status: data.status,
          agent: data.agent,
          timestamp: data.timestamp
        };
        
        // Broadcast to browsers
        broadcast(data.task, data.status);
        
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
    res.end(JSON.stringify(taskStates));
    
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
  
  // Send current state
  ws.send(JSON.stringify({
    type: 'init',
    states: taskStates
  }));
  
  ws.on('close', () => {
    wsClients.delete(ws);
  });
});

server.listen(PORT, () => {
  console.log(`
╔════════════════════════════════════════════════════════╗
║                  DAG STATE SERVER                      ║
╠════════════════════════════════════════════════════════╣
║  POST http://localhost:${PORT}/update                  ║
║  GET  http://localhost:${PORT}/state                   ║
║  WS   ws://localhost:${PORT}                           ║
╚════════════════════════════════════════════════════════╝

Example:
curl -X POST http://localhost:${PORT}/update \\
  -H "Content-Type: application/json" \\
  -d '{"agent":"mr_clean","task":"P1.T003","status":"started","timestamp":1755633601}'
  `);
});