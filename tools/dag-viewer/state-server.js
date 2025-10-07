#!/usr/bin/env node

import { WebSocketServer } from 'ws';
import { readFileSync, watchFile, existsSync } from 'fs';
import { createServer } from 'http';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

// Configuration
const PORT = process.env.PORT || 8080;
const STATE_FILE = process.env.STATE_FILE || join(__dirname, 'task-states.json');

// Create HTTP server for serving the HTML
const server = createServer((req, res) => {
  // Enable CORS
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Access-Control-Allow-Methods', 'GET');
  
  if (req.url === '/health') {
    res.writeHead(200, { 'Content-Type': 'text/plain' });
    res.end('OK');
  } else {
    res.writeHead(404);
    res.end('Not found');
  }
});

// Create WebSocket server
const wss = new WebSocketServer({ server });

// Store active connections
const clients = new Set();

// Read state file
function readStateFile() {
  try {
    if (!existsSync(STATE_FILE)) {
      console.log(`State file not found: ${STATE_FILE}, creating empty state file`);
      return {};
    }
    const data = readFileSync(STATE_FILE, 'utf8');
    return JSON.parse(data);
  } catch (error) {
    console.error('Error reading state file:', error);
    return {};
  }
}

// Broadcast state to all connected clients
function broadcastState() {
  const state = readStateFile();
  const message = JSON.stringify({
    type: 'state-update',
    data: state,
    timestamp: new Date().toISOString()
  });
  
  clients.forEach(client => {
    if (client.readyState === 1) { // WebSocket.OPEN
      client.send(message);
    }
  });
  
  console.log(`Broadcasted state update to ${clients.size} clients`);
}

// Watch the state file for changes
watchFile(STATE_FILE, { interval: 500 }, (curr, prev) => {
  if (curr.mtime !== prev.mtime) {
    console.log('State file changed, broadcasting update...');
    broadcastState();
  }
});

// Handle WebSocket connections
wss.on('connection', (ws) => {
  console.log('New client connected');
  clients.add(ws);
  
  // Send initial state
  const state = readStateFile();
  ws.send(JSON.stringify({
    type: 'state-update',
    data: state,
    timestamp: new Date().toISOString()
  }));
  
  // Handle client messages
  ws.on('message', (message) => {
    try {
      const data = JSON.parse(message);
      if (data.type === 'ping') {
        ws.send(JSON.stringify({ type: 'pong' }));
      }
    } catch (error) {
      console.error('Error parsing message:', error);
    }
  });
  
  // Handle disconnection
  ws.on('close', () => {
    console.log('Client disconnected');
    clients.delete(ws);
  });
  
  ws.on('error', (error) => {
    console.error('WebSocket error:', error);
    clients.delete(ws);
  });
});

// Start server
server.listen(PORT, () => {
  console.log(`
╔═══════════════════════════════════════════════════════╗
║         DAG State WebSocket Server Running            ║
╠═══════════════════════════════════════════════════════╣
║  WebSocket Port: ${PORT}                              ║
║  State File: ${STATE_FILE}                            ║
║                                                       ║
║  Clients can connect to ws://localhost:${PORT}        ║
╚═══════════════════════════════════════════════════════╝
  `);
  
  // Create example state file if it doesn't exist
  if (!existsSync(STATE_FILE)) {
    const exampleState = {
      "task1": "completed",
      "task2": "in_progress",
      "task3": "pending",
      "task4": "failed",
      "task5": "blocked"
    };
    require('fs').writeFileSync(STATE_FILE, JSON.stringify(exampleState, null, 2));
    console.log(`Created example state file: ${STATE_FILE}`);
  }
});