#!/usr/bin/env node

import { WebSocketServer } from 'ws';
import { createServer } from 'http';
import { readFileSync } from 'fs';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const PORT = 3456;
const DAG_FILE = join(__dirname, 'public', 'dag.json');

// Load DAG data
const dagData = JSON.parse(readFileSync(DAG_FILE, 'utf8'));

// Task name mappings for realistic display
const taskNames = {
  'P1.T001': 'Initialize project structure',
  'P1.T002': 'Setup development environment',
  'P1.T003': 'Configure build tools',
  'P1.T004': 'Create base components',
  'P1.T005': 'Implement data models',
  'P1.T006': 'Setup authentication',
  'P1.T007': 'Create user service',
  'P1.T008': 'Setup database connections',
  'P1.T009': 'Implement API endpoints',
  'P1.T010': 'Add validation middleware',
  'P1.T011': 'Create admin dashboard',
  'P1.T012': 'Setup monitoring',
  'P1.T013': 'Configure logging',
  'P1.T014': 'Add error handling',
  'P1.T015': 'Write unit tests',
  'P1.T016': 'Setup CI/CD pipeline',
  'P1.T017': 'Configure deployment',
  'P1.T018': 'Add documentation',
  'P1.T019': 'Create API docs',
  'P1.T020': 'Performance optimization',
  'P1.T021': 'Security audit',
  'P1.T022': 'Load testing',
  'P1.T023': 'Code review',
  'P1.T024': 'Integration tests',
  'P1.T025': 'Deploy to production'
};

// Add task names to DAG data
if (dagData.topo_order) {
  dagData.task_metadata = {};
  dagData.topo_order.forEach(taskId => {
    dagData.task_metadata[taskId] = {
      name: taskNames[taskId] || taskId,
      description: `Task ${taskId}: ${taskNames[taskId] || 'Unknown task'}`,
      estimatedTime: Math.floor(Math.random() * 300) + 60 // 1-6 minutes
    };
  });
}

// Task states
const taskStates = {};
const agents = {};
const gitStats = {
  totalCommits: 0,
  totalLinesAdded: 0,
  totalLinesRemoved: 0,
  fileChangeFrequency: {}
};

// Agent names pool
const agentNames = [
  'agent-alpha',
  'agent-beta',
  'agent-gamma',
  'agent-delta',
  'agent-epsilon'
];

// Initialize agents
agentNames.forEach(name => {
  agents[name] = {
    name,
    tasksCompleted: 0,
    tasksFailed: 0,
    tasksStarted: 0,
    taskTimes: [],
    currentTask: null,
    totalTokens: 0,
    tokensPerTask: [],
    avgTokensPerTask: 0,
    tokenEfficiency: 0,
    commits: 0,
    linesAdded: 0,
    linesRemoved: 0,
    filesChanged: new Set(),
    commitSizes: [],
    lastCommit: null,
    tasksPerCommit: 0,
    linesPerTask: 0,
    commitFrequency: 0,
    tokensPerLine: 0,
    firstSeen: Date.now() / 1000,
    lastSeen: Date.now() / 1000,
    activeTime: 0,
    path: []
  };
});

// Initialize task states
dagData.topo_order.forEach(taskId => {
  taskStates[taskId] = {
    status: 'pending',
    agent: null,
    timestamp: null,
    startTime: null,
    tokensSpent: 0,
    startTokens: 0
  };
});

// Create HTTP server
const server = createServer((req, res) => {
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
const clients = new Set();

// Broadcast to all clients
function broadcast(message) {
  const data = JSON.stringify(message);
  clients.forEach(client => {
    if (client.readyState === 1) {
      client.send(data);
    }
  });
}

// Get current leaderboard
function getLeaderboard() {
  return Object.values(agents)
    .map(agent => ({
      ...agent,
      productivityScore: agent.tasksCompleted * 100 + agent.commits * 50,
      successRate: agent.tasksStarted > 0 
        ? `${((agent.tasksCompleted / agent.tasksStarted) * 100).toFixed(1)}%`
        : '0%',
      avgTaskTime: agent.taskTimes.length > 0
        ? agent.taskTimes.reduce((a, b) => a + b, 0) / agent.taskTimes.length
        : 0,
      medianTaskTime: agent.taskTimes.length > 0
        ? agent.taskTimes.sort((a, b) => a - b)[Math.floor(agent.taskTimes.length / 2)]
        : 0,
      totalTokensK: `${(agent.totalTokens / 1000).toFixed(1)}k`,
      tokenCost: `$${(agent.totalTokens * 0.00002).toFixed(2)}`,
      costPerTask: agent.tasksCompleted > 0
        ? `$${((agent.totalTokens * 0.00002) / agent.tasksCompleted).toFixed(3)}`
        : '$0',
      avgCommitSize: agent.commits > 0
        ? `+${Math.round((agent.linesAdded - agent.linesRemoved) / agent.commits)}`
        : '0',
      codeChurn: agent.linesAdded + agent.linesRemoved,
      netLines: agent.linesAdded - agent.linesRemoved,
      filesChangedCount: agent.filesChanged.size
    }))
    .sort((a, b) => b.productivityScore - a.productivityScore);
}

// Get git insights
function getGitInsights() {
  const recentCommits = [];
  Object.values(agents).forEach(agent => {
    if (agent.lastCommit) {
      recentCommits.push(agent.lastCommit);
    }
  });
  
  return {
    summary: {
      totalCommits: gitStats.totalCommits,
      totalLinesAdded: gitStats.totalLinesAdded,
      totalLinesRemoved: gitStats.totalLinesRemoved,
      netLines: gitStats.totalLinesAdded - gitStats.totalLinesRemoved,
      avgCommitSize: gitStats.totalCommits > 0
        ? `${Math.round((gitStats.totalLinesAdded + gitStats.totalLinesRemoved) / gitStats.totalCommits)} lines`
        : '0 lines'
    },
    hotFiles: Object.entries(gitStats.fileChangeFrequency)
      .map(([path, changes]) => ({ path, changes }))
      .sort((a, b) => b.changes - a.changes)
      .slice(0, 10),
    recentCommits: recentCommits.slice(-10),
    hourlyActivity: new Array(24).fill(0),
    topContributors: Object.values(agents)
      .filter(a => a.commits > 0)
      .map(a => ({ agent: a.name, commits: a.commits }))
      .sort((a, b) => b.commits - a.commits)
      .slice(0, 5),
    fileTypes: {}
  };
}

// Handle WebSocket connections
wss.on('connection', (ws) => {
  console.log('New client connected');
  clients.add(ws);
  
  // Send initial state with DAG data
  const initMessage = {
    type: 'init',
    dagData: dagData,
    tasks: taskStates,
    agents: agents,
    leaderboard: getLeaderboard(),
    gitInsights: getGitInsights(),
    recentEvents: []
  };
  ws.send(JSON.stringify(initMessage));
  console.log('Sent initial state with', Object.keys(taskStates).length, 'tasks');
  
  ws.on('close', () => {
    clients.delete(ws);
    console.log('Client disconnected');
  });
  
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
});

// Simulation logic
let currentTaskIndex = 0;
const maxConcurrentTasks = 3;
const activeTasks = new Map();

function getAvailableAgent() {
  const availableAgents = agentNames.filter(name => !agents[name].currentTask);
  if (availableAgents.length === 0) return null;
  return availableAgents[Math.floor(Math.random() * availableAgents.length)];
}

function getNextTask() {
  // Find next pending task that has no unmet dependencies
  for (let i = currentTaskIndex; i < dagData.topo_order.length; i++) {
    const taskId = dagData.topo_order[i];
    if (taskStates[taskId].status === 'pending') {
      // Check dependencies
      const hasDependencies = dagData.reduced_edges_sample?.some(
        ([source, target]) => target === taskId && taskStates[source].status !== 'completed'
      );
      if (!hasDependencies) {
        currentTaskIndex = i;
        return taskId;
      }
    }
  }
  return null;
}

function simulateTaskStart() {
  if (activeTasks.size >= maxConcurrentTasks) return;
  
  const taskId = getNextTask();
  if (!taskId) return;
  
  const agentName = getAvailableAgent();
  if (!agentName) return;
  
  const agent = agents[agentName];
  const timestamp = Date.now() / 1000;
  
  // Update task state
  taskStates[taskId] = {
    status: 'started',
    agent: agentName,
    timestamp,
    startTime: timestamp,
    startTokens: agent.totalTokens
  };
  
  // Update agent state
  agent.currentTask = taskId;
  agent.tasksStarted++;
  agent.lastSeen = timestamp;
  agent.path.push({ task: taskId, status: 'started', timestamp });
  
  activeTasks.set(taskId, {
    agent: agentName,
    startTime: timestamp,
    duration: 3000 + Math.random() * 7000 // 3-10 seconds
  });
  
  // Broadcast update
  broadcast({
    type: 'task-claimed',
    task: taskId,
    agent: agentName,
    timestamp
  });
  
  broadcast({
    type: 'event',
    eventType: 'task',
    agent: agentName,
    task: taskId,
    status: 'started',
    timestamp,
    message: `${agentName} started: ${taskNames[taskId] || taskId}`
  });
  
  console.log(`Task ${taskId} started by ${agentName}`);
}

function simulateTaskComplete(taskId) {
  const taskInfo = activeTasks.get(taskId);
  if (!taskInfo) return;
  
  const agent = agents[taskInfo.agent];
  const timestamp = Date.now() / 1000;
  const duration = timestamp - taskInfo.startTime;
  const tokensUsed = Math.floor(1000 + Math.random() * 5000);
  
  // Randomly fail some tasks (10% chance)
  const success = Math.random() > 0.1;
  const status = success ? 'completed' : 'failed';
  
  // Update task state
  taskStates[taskId].status = status;
  taskStates[taskId].timestamp = timestamp;
  taskStates[taskId].tokensSpent = tokensUsed;
  
  // Update agent state
  agent.currentTask = null;
  agent.lastSeen = timestamp;
  agent.totalTokens += tokensUsed;
  agent.tokensPerTask.push(tokensUsed);
  agent.path.push({ task: taskId, status, timestamp });
  
  if (success) {
    agent.tasksCompleted++;
    agent.taskTimes.push(duration);
    
    // Simulate git commit (30% chance)
    if (Math.random() < 0.3) {
      simulateGitCommit(taskInfo.agent, taskId);
    }
  } else {
    agent.tasksFailed++;
  }
  
  // Update agent averages
  if (agent.tokensPerTask.length > 0) {
    agent.avgTokensPerTask = agent.tokensPerTask.reduce((a, b) => a + b, 0) / agent.tokensPerTask.length;
  }
  if (agent.tasksCompleted > 0) {
    agent.tokenEfficiency = agent.totalTokens / agent.tasksCompleted;
  }
  
  activeTasks.delete(taskId);
  
  // Broadcast update
  broadcast({
    type: 'task-update',
    task: taskId,
    status,
    agent: taskInfo.agent,
    timestamp,
    tokensUsed
  });
  
  broadcast({
    type: 'event',
    eventType: 'task',
    agent: taskInfo.agent,
    task: taskId,
    status,
    timestamp,
    tokensUsed,
    duration,
    message: `${taskInfo.agent} ${status}: ${taskNames[taskId] || taskId} (${duration.toFixed(1)}s, ${tokensUsed} tokens)`
  });
  
  // Send stats update
  broadcast({
    type: 'stats-update',
    leaderboard: getLeaderboard(),
    gitInsights: getGitInsights()
  });
  
  console.log(`Task ${taskId} ${status} by ${taskInfo.agent}`);
}

function simulateGitCommit(agentName, taskId) {
  const agent = agents[agentName];
  const timestamp = Date.now() / 1000;
  const linesAdded = Math.floor(Math.random() * 200) + 10;
  const linesRemoved = Math.floor(Math.random() * 50);
  const filesChanged = Math.floor(Math.random() * 5) + 1;
  
  // Update agent git stats
  agent.commits++;
  agent.linesAdded += linesAdded;
  agent.linesRemoved += linesRemoved;
  agent.commitSizes.push(linesAdded + linesRemoved);
  
  // Add some fake files
  for (let i = 0; i < filesChanged; i++) {
    const fileName = `src/file${Math.floor(Math.random() * 100)}.ts`;
    agent.filesChanged.add(fileName);
    gitStats.fileChangeFrequency[fileName] = (gitStats.fileChangeFrequency[fileName] || 0) + 1;
  }
  
  // Update global git stats
  gitStats.totalCommits++;
  gitStats.totalLinesAdded += linesAdded;
  gitStats.totalLinesRemoved += linesRemoved;
  
  const commit = {
    agent: agentName,
    task: taskId,
    timestamp,
    message: `Complete: ${taskNames[taskId] || taskId}`,
    taskName: taskNames[taskId] || taskId,
    linesAdded,
    linesRemoved,
    filesChanged,
    hash: Math.random().toString(36).substring(2, 9)
  };
  
  agent.lastCommit = commit;
  
  // Calculate derived metrics
  if (agent.tasksCompleted > 0) {
    agent.tasksPerCommit = agent.tasksCompleted / agent.commits;
    agent.linesPerTask = (agent.linesAdded - agent.linesRemoved) / agent.tasksCompleted;
  }
  if (agent.totalTokens > 0 && agent.linesAdded > 0) {
    agent.tokensPerLine = agent.totalTokens / agent.linesAdded;
  }
  
  // Broadcast git event
  broadcast({
    type: 'git-commit',
    ...commit
  });
  
  broadcast({
    type: 'event',
    eventType: 'git',
    ...commit,
    message: `${agentName} committed: +${linesAdded} -${linesRemoved} (${filesChanged} files)`
  });
  
  console.log(`Git commit by ${agentName}: +${linesAdded} -${linesRemoved}`);
}

// Simulation loop
function runSimulation() {
  // Start new tasks
  simulateTaskStart();
  
  // Complete active tasks
  activeTasks.forEach((taskInfo, taskId) => {
    const elapsed = (Date.now() / 1000) - taskInfo.startTime;
    if (elapsed * 1000 >= taskInfo.duration) {
      simulateTaskComplete(taskId);
    }
  });
  
  // Check if simulation is complete
  const allTasksProcessed = Object.values(taskStates).every(
    state => state.status === 'completed' || state.status === 'failed'
  );
  
  if (allTasksProcessed) {
    console.log('Simulation complete!');
    broadcast({
      type: 'simulation-complete',
      timestamp: Date.now() / 1000
    });
  }
}

// Start the server
server.listen(PORT, () => {
  console.log(`WebSocket server listening on ws://localhost:${PORT}`);
  console.log(`Health check available at http://localhost:${PORT}/health`);
  console.log(`DAG contains ${dagData.topo_order.length} tasks`);
  console.log('\nStarting task simulation...\n');
  
  // Run simulation every 500ms
  setInterval(runSimulation, 500);
});

// Graceful shutdown
process.on('SIGINT', () => {
  console.log('\nShutting down...');
  wss.close();
  server.close();
  process.exit(0);
});