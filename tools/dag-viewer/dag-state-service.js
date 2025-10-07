#!/usr/bin/env node

import { WebSocketServer } from 'ws';
import { createServer as createHttpServer } from 'http';
import { createServer as createUnixServer } from 'net';
import { unlink, existsSync, writeFileSync, readFileSync } from 'fs';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

// Configuration
const WS_PORT = process.env.WS_PORT || 8080;
const UNIX_SOCKET = process.env.UNIX_SOCKET || '/tmp/dag-state.sock';
const STATE_FILE = process.env.STATE_FILE || join(__dirname, 'dag-state.json');
const PERSIST_INTERVAL = 5000; // Save state to disk every 5 seconds

// In-memory state
let dagState = {};
let stateVersion = 0;
let isDirty = false;

// WebSocket clients
const wsClients = new Set();

// Load initial state from disk
function loadState() {
  try {
    if (existsSync(STATE_FILE)) {
      const data = readFileSync(STATE_FILE, 'utf8');
      dagState = JSON.parse(data);
      console.log(`Loaded state from ${STATE_FILE}`);
    }
  } catch (error) {
    console.error('Error loading state:', error);
    dagState = {};
  }
}

// Save state to disk
function saveState() {
  if (!isDirty) return;
  
  try {
    writeFileSync(STATE_FILE, JSON.stringify(dagState, null, 2));
    isDirty = false;
    console.log(`State persisted to ${STATE_FILE}`);
  } catch (error) {
    console.error('Error saving state:', error);
  }
}

// Broadcast state to all WebSocket clients
function broadcastState(taskId = null) {
  const message = JSON.stringify({
    type: 'state-update',
    version: stateVersion,
    timestamp: new Date().toISOString(),
    taskId: taskId,
    fullState: taskId ? null : dagState,
    taskState: taskId ? { [taskId]: dagState[taskId] } : null
  });
  
  wsClients.forEach(client => {
    if (client.readyState === 1) { // WebSocket.OPEN
      client.send(message);
    }
  });
}

// Update task state
function updateTaskState(taskId, state, metadata = {}) {
  const validStates = ['pending', 'started', 'in_progress', 'failed', 'blocked', 'completed'];
  
  if (!validStates.includes(state)) {
    throw new Error(`Invalid state: ${state}. Must be one of: ${validStates.join(', ')}`);
  }
  
  dagState[taskId] = {
    state: state,
    lastUpdated: new Date().toISOString(),
    ...metadata
  };
  
  stateVersion++;
  isDirty = true;
  
  console.log(`Task ${taskId} -> ${state}`);
  broadcastState(taskId);
}

// Batch update multiple tasks
function batchUpdateTasks(updates) {
  for (const [taskId, data] of Object.entries(updates)) {
    if (typeof data === 'string') {
      dagState[taskId] = {
        state: data,
        lastUpdated: new Date().toISOString()
      };
    } else {
      dagState[taskId] = {
        ...data,
        lastUpdated: new Date().toISOString()
      };
    }
  }
  
  stateVersion++;
  isDirty = true;
  
  console.log(`Batch updated ${Object.keys(updates).length} tasks`);
  broadcastState();
}

// Create Unix domain socket server
function createUnixSocketServer() {
  // Remove existing socket file
  if (existsSync(UNIX_SOCKET)) {
    unlink(UNIX_SOCKET, (err) => {
      if (err) console.error('Error removing socket:', err);
    });
  }
  
  const server = createUnixServer((client) => {
    console.log('Unix socket client connected');
    
    let buffer = '';
    
    client.on('data', (data) => {
      buffer += data.toString();
      
      // Process complete JSON messages (newline delimited)
      const lines = buffer.split('\n');
      buffer = lines.pop(); // Keep incomplete line in buffer
      
      lines.forEach(line => {
        if (!line.trim()) return;
        
        try {
          const message = JSON.parse(line);
          
          switch (message.type) {
            case 'update':
              updateTaskState(message.taskId, message.state, message.metadata);
              client.write(JSON.stringify({ success: true, version: stateVersion }) + '\n');
              break;
              
            case 'batch':
              batchUpdateTasks(message.updates);
              client.write(JSON.stringify({ success: true, version: stateVersion }) + '\n');
              break;
              
            case 'query':
              const response = message.taskId 
                ? { state: dagState[message.taskId] }
                : { state: dagState };
              client.write(JSON.stringify(response) + '\n');
              break;
              
            case 'clear':
              dagState = {};
              stateVersion++;
              isDirty = true;
              broadcastState();
              client.write(JSON.stringify({ success: true }) + '\n');
              break;
              
            default:
              client.write(JSON.stringify({ error: 'Unknown message type' }) + '\n');
          }
        } catch (error) {
          console.error('Error processing message:', error);
          client.write(JSON.stringify({ error: error.message }) + '\n');
        }
      });
    });
    
    client.on('end', () => {
      console.log('Unix socket client disconnected');
    });
    
    client.on('error', (err) => {
      console.error('Unix socket error:', err);
    });
  });
  
  server.listen(UNIX_SOCKET, () => {
    console.log(`Unix socket listening at: ${UNIX_SOCKET}`);
  });
  
  return server;
}

// Create HTTP server for WebSocket
const httpServer = createHttpServer((req, res) => {
  res.setHeader('Access-Control-Allow-Origin', '*');
  
  if (req.url === '/health') {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ 
      status: 'healthy',
      version: stateVersion,
      taskCount: Object.keys(dagState).length,
      wsClients: wsClients.size
    }));
  } else if (req.url === '/state') {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify(dagState));
  } else if (req.url.startsWith('/update/') && req.method === 'POST') {
    // Handle POST /update/:taskId/:state
    const parts = req.url.split('/');
    const taskId = parts[2];
    const state = parts[3];
    
    if (taskId && state) {
      try {
        updateTaskState(taskId, state);
        res.writeHead(200, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ success: true, version: stateVersion }));
      } catch (error) {
        res.writeHead(400, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: error.message }));
      }
    } else {
      res.writeHead(400, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ error: 'Invalid request format' }));
    }
  } else {
    res.writeHead(404);
    res.end('Not found');
  }
});

// Create WebSocket server
const wss = new WebSocketServer({ server: httpServer });

wss.on('connection', (ws) => {
  console.log('WebSocket client connected');
  wsClients.add(ws);
  
  // Send initial state
  ws.send(JSON.stringify({
    type: 'state-update',
    version: stateVersion,
    timestamp: new Date().toISOString(),
    fullState: dagState
  }));
  
  ws.on('message', (message) => {
    try {
      const data = JSON.parse(message);
      if (data.type === 'ping') {
        ws.send(JSON.stringify({ type: 'pong' }));
      }
    } catch (error) {
      console.error('WebSocket message error:', error);
    }
  });
  
  ws.on('close', () => {
    console.log('WebSocket client disconnected');
    wsClients.delete(ws);
  });
  
  ws.on('error', (error) => {
    console.error('WebSocket error:', error);
    wsClients.delete(ws);
  });
});

// Periodic state persistence
setInterval(saveState, PERSIST_INTERVAL);

// Graceful shutdown
process.on('SIGINT', () => {
  console.log('\nShutting down gracefully...');
  saveState();
  
  // Close WebSocket connections
  wsClients.forEach(client => client.close());
  
  // Remove Unix socket
  if (existsSync(UNIX_SOCKET)) {
    unlink(UNIX_SOCKET, () => {
      process.exit(0);
    });
  } else {
    process.exit(0);
  }
});

// Initialize
loadState();
createUnixSocketServer();

httpServer.listen(WS_PORT, () => {
  console.log(`
╔════════════════════════════════════════════════════════╗
║           DAG State Management Service                 ║
╠════════════════════════════════════════════════════════╣
║  Unix Socket: ${UNIX_SOCKET.padEnd(40)} ║
║  WebSocket:   ws://localhost:${String(WS_PORT).padEnd(25)} ║
║  HTTP API:    http://localhost:${String(WS_PORT).padEnd(23)} ║
║  State File:  ${STATE_FILE.padEnd(40)} ║
║                                                        ║
║  Endpoints:                                            ║
║    /health - Service health check                     ║
║    /state  - Current DAG state (GET)                  ║
╚════════════════════════════════════════════════════════╝
  `);
});