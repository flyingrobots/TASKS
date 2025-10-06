#!/usr/bin/env node

import { WebSocketServer } from 'ws';
import { createServer } from 'http';
import { writeFileSync, readFileSync, existsSync } from 'fs';
import { fileURLToPath } from 'url';
import { dirname, join, resolve } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

// Load configuration
const CONFIG_PATH = process.argv[2] || join(__dirname, 'dag-viewer-react', 'config.json');
let config = {
  server: { port: 8080, host: 'localhost' },
  paths: {
    dagFile: join(__dirname, '../../../dag.json'),
    stateFile: join(__dirname, 'dag-state-full.json')
  }
};

if (existsSync(CONFIG_PATH)) {
  try {
    config = JSON.parse(readFileSync(CONFIG_PATH, 'utf8'));
    console.log('Loaded configuration from', CONFIG_PATH);
  } catch (e) {
    console.error('Error loading config:', e);
    console.log('Using default configuration');
  }
} else {
  console.log('Config file not found at', CONFIG_PATH, '- using defaults');
}

const PORT = config.server.port;
const STATE_FILE = resolve(config.paths.stateFile);

// In-memory state
const state = {
  tasks: {},           // taskId -> { status, agent, timestamp }
  agents: {},          // agentId -> { stats }
  commits: [],         // array of commit events
  dag: null,           // DAG structure loaded from file
  gitStats: {          // aggregate git statistics
    totalCommits: 0,
    totalLinesAdded: 0,
    totalLinesRemoved: 0,
    fileChangeFrequency: {},
    agentCommits: {},
    hourlyActivity: new Array(24).fill(0),
    dailyActivity: {},
    hotFiles: {},      // files changed most frequently
    commitSizes: []    // distribution of commit sizes
  },
  events: []          // all events in order
};

// WebSocket clients
const wsClients = new Set();

// Load state
function loadState() {
  if (existsSync(STATE_FILE)) {
    try {
      const data = JSON.parse(readFileSync(STATE_FILE, 'utf8'));
      Object.assign(state, data);
      console.log('Loaded state from', STATE_FILE);
    } catch (e) {
      console.error('Error loading state:', e);
    }
  }
}

// Load DAG structure and related files
function loadDAG() {
  const dagPath = resolve(config.paths.dagFile);
  try {
    if (existsSync(dagPath)) {
      const rawDag = JSON.parse(readFileSync(dagPath, 'utf8'));
      
      // Convert DAG format if needed
      if (rawDag.topo_order && rawDag.reduced_edges) {
        // Convert from topo_order/reduced_edges format to nodes/edges format
        const nodes = rawDag.topo_order.map(id => ({ id }));
        const edges = rawDag.reduced_edges.map(edge => ({
          source: edge[0],
          target: edge[1]
        }));
        
        state.dag = {
          ...rawDag,
          nodes,
          edges
        };
        console.log(`Loaded and converted DAG from ${dagPath}: ${nodes.length} nodes, ${edges.length} edges`);
      } else if (rawDag.nodes && rawDag.edges) {
        // Already in the correct format
        state.dag = rawDag;
        console.log(`Loaded DAG from ${dagPath}: ${rawDag.nodes.length} nodes, ${rawDag.edges.length} edges`);
      } else {
        state.dag = rawDag;
        console.log('Loaded DAG from', dagPath);
      }
    } else {
      console.warn('DAG file not found at', dagPath);
    }
  } catch (e) {
    console.error('Error loading DAG:', e);
  }
  
  // Load features if available
  if (config.paths.featuresFile) {
    const featuresPath = resolve(config.paths.featuresFile);
    try {
      if (existsSync(featuresPath)) {
        state.features = JSON.parse(readFileSync(featuresPath, 'utf8'));
        console.log('Loaded features from', featuresPath);
      }
    } catch (e) {
      console.error('Error loading features:', e);
    }
  }
  
  // Load waves if available
  if (config.paths.wavesFile) {
    const wavesPath = resolve(config.paths.wavesFile);
    try {
      if (existsSync(wavesPath)) {
        state.waves = JSON.parse(readFileSync(wavesPath, 'utf8'));
        console.log('Loaded waves from', wavesPath);
      }
    } catch (e) {
      console.error('Error loading waves:', e);
    }
  }
  
  // Load tasks if available
  if (config.paths.tasksFile) {
    const tasksPath = resolve(config.paths.tasksFile);
    try {
      if (existsSync(tasksPath)) {
        state.tasksMetadata = JSON.parse(readFileSync(tasksPath, 'utf8'));
        console.log('Loaded tasks metadata from', tasksPath);
      }
    } catch (e) {
      console.error('Error loading tasks metadata:', e);
    }
  }
}

// Get ready tasks (all dependencies completed)
function getReadyTasks() {
  if (!state.dag || (!state.dag.edges && !state.dag.reduced_edges_sample)) {
    return { ready: [], blocked: [], error: 'DAG not loaded' };
  }
  
  const ready = [];
  const blocked = [];
  const inProgress = [];
  
  // Get all tasks from DAG
  const allTasks = new Set();
  const edges = state.dag.edges || state.dag.reduced_edges_sample || [];
  
  if (state.dag.edges) {
    // Handle standard edges format
    edges.forEach(edge => {
      allTasks.add(edge.source);
      allTasks.add(edge.target);
    });
  } else if (state.dag.reduced_edges_sample) {
    // Handle reduced_edges_sample format [source, target]
    edges.forEach(edge => {
      allTasks.add(edge[0]);
      allTasks.add(edge[1]);
    });
  }
  
  // Add tasks from topo_order if available
  if (state.dag.topo_order) {
    state.dag.topo_order.forEach(task => allTasks.add(task));
  }
  
  // Build dependency map
  const dependencies = {};
  allTasks.forEach(task => {
    dependencies[task] = [];
  });
  
  if (state.dag.edges) {
    // Handle standard edges format
    state.dag.edges.forEach(edge => {
      dependencies[edge.target].push(edge.source);
    });
  } else if (state.dag.reduced_edges_sample) {
    // Handle reduced_edges_sample format [source, target]
    state.dag.reduced_edges_sample.forEach(edge => {
      dependencies[edge[1]].push(edge[0]);
    });
  }
  
  // Check each task
  allTasks.forEach(taskId => {
    const taskState = state.tasks[taskId];
    const status = taskState?.status || 'pending';
    
    // Skip already completed tasks
    if (status === 'completed') {
      return;
    }
    
    // Track in-progress tasks
    if (status === 'started' || status === 'in_progress') {
      inProgress.push({
        task: taskId,
        agent: taskState.agent,
        status
      });
      return;
    }
    
    // Check if all dependencies are completed
    const deps = dependencies[taskId] || [];
    const allDepsCompleted = deps.every(dep => {
      const depState = state.tasks[dep];
      return depState && depState.status === 'completed';
    });
    
    if (allDepsCompleted) {
      // Check if any dependency is blocking
      const hasBlockedDep = deps.some(dep => {
        const depState = state.tasks[dep];
        return depState && depState.status === 'blocked';
      });
      
      if (hasBlockedDep) {
        blocked.push({
          task: taskId,
          reason: 'dependency_blocked',
          dependencies: deps
        });
      } else if (deps.length === 0) {
        // No dependencies - ready to go
        ready.push({
          task: taskId,
          priority: 'high', // Entry points have high priority
          dependencies: []
        });
      } else {
        // All deps completed - ready to go
        ready.push({
          task: taskId,
          priority: 'normal',
          dependencies: deps
        });
      }
    } else {
      // Some dependencies not completed
      const incompleteDeps = deps.filter(dep => {
        const depState = state.tasks[dep];
        return !depState || depState.status !== 'completed';
      });
      
      blocked.push({
        task: taskId,
        reason: 'dependencies_incomplete',
        dependencies: deps,
        incomplete: incompleteDeps
      });
    }
  });
  
  // Sort ready tasks by priority and wave if available
  ready.sort((a, b) => {
    if (a.priority !== b.priority) {
      return a.priority === 'high' ? -1 : 1;
    }
    return 0;
  });
  
  return {
    ready,
    blocked,
    inProgress,
    summary: {
      readyCount: ready.length,
      blockedCount: blocked.length,
      inProgressCount: inProgress.length,
      completedCount: Array.from(allTasks).filter(t => 
        state.tasks[t]?.status === 'completed'
      ).length,
      totalTasks: allTasks.size
    }
  };
}

// Save state
function saveState() {
  try {
    writeFileSync(STATE_FILE, JSON.stringify(state, null, 2));
  } catch (e) {
    console.error('Error saving state:', e);
  }
}

// Initialize agent
function initAgent(agentId) {
  if (!state.agents[agentId]) {
    state.agents[agentId] = {
      name: agentId,
      // Task metrics
      tasksCompleted: 0,
      tasksFailed: 0,
      tasksStarted: 0,
      taskTimes: [],
      currentTask: null,
      // Token metrics
      totalTokens: 0,
      tokensPerTask: [],
      avgTokensPerTask: 0,
      tokenEfficiency: 0, // tokens per completed task
      // Git metrics
      commits: 0,
      linesAdded: 0,
      linesRemoved: 0,
      filesChanged: new Set(),
      commitSizes: [],
      lastCommit: null,
      // Productivity metrics
      tasksPerCommit: 0,
      linesPerTask: 0,
      commitFrequency: 0,
      tokensPerLine: 0,
      // Time tracking
      firstSeen: Date.now(),
      lastSeen: Date.now(),
      activeTime: 0,
      path: []
    };
  }
}

// Process task update
function processTaskUpdate(data) {
  const { agent, task, status, timestamp, total_tokens } = data;
  
  initAgent(agent);
  
  const agentData = state.agents[agent];
  agentData.lastSeen = Date.now();
  
  // Update task state
  const prevState = state.tasks[task];
  const prevAgentTokens = agentData.lastKnownTokens || 0;
  
  // Calculate tokens spent on this specific update
  let tokensForThisUpdate = 0;
  if (total_tokens) {
    tokensForThisUpdate = total_tokens - prevAgentTokens;
    agentData.lastKnownTokens = total_tokens;
    agentData.totalTokens = total_tokens;
  }
  
  state.tasks[task] = { 
    status, 
    agent, 
    timestamp, 
    startTime: prevState?.startTime,
    tokensSpent: (prevState?.tokensSpent || 0) + tokensForThisUpdate,
    startTokens: prevState?.startTokens
  };
  
  // Update agent stats
  switch(status) {
    case 'started':
      agentData.tasksStarted++;
      agentData.currentTask = task;
      state.tasks[task].startTime = timestamp;
      state.tasks[task].startTokens = prevAgentTokens;
      break;
      
    case 'completed':
      agentData.tasksCompleted++;
      if (state.tasks[task].startTime) {
        const duration = timestamp - state.tasks[task].startTime;
        agentData.taskTimes.push(duration);
      }
      // Calculate tokens used for this task
      const taskTokens = state.tasks[task].tokensSpent || tokensForThisUpdate;
      if (taskTokens > 0) {
        agentData.tokensPerTask.push(taskTokens);
        agentData.avgTokensPerTask = Math.round(
          agentData.tokensPerTask.reduce((a,b) => a+b, 0) / agentData.tokensPerTask.length
        );
        agentData.tokenEfficiency = Math.round(
          agentData.totalTokens / Math.max(1, agentData.tasksCompleted)
        );
      }
      agentData.currentTask = null;
      break;
      
    case 'failed':
      agentData.tasksFailed++;
      agentData.currentTask = null;
      break;
      
    case 'in_progress':
      // Just track the running total
      break;
  }
  
  agentData.path.push({ task, status, timestamp, tokens: tokensForThisUpdate });
  if (agentData.path.length > 50) agentData.path.shift();
  
  // Add to events
  state.events.push({ type: 'task', agent, task, status, timestamp, tokens: tokensForThisUpdate });
  if (state.events.length > 1000) state.events.shift();
}

// Process git commit
function processGitCommit(data) {
  const { agent, sha, message, files, timestamp, branch } = data;
  
  initAgent(agent);
  
  const agentData = state.agents[agent];
  const hour = new Date(timestamp * 1000).getHours();
  const date = new Date(timestamp * 1000).toISOString().split('T')[0];
  
  // Calculate metrics
  let linesAdded = 0;
  let linesRemoved = 0;
  let filesChanged = [];
  
  if (files && Array.isArray(files)) {
    files.forEach(file => {
      linesAdded += file.additions || 0;
      linesRemoved += file.deletions || 0;
      filesChanged.push(file.path);
      
      // Track hot files
      state.gitStats.hotFiles[file.path] = (state.gitStats.hotFiles[file.path] || 0) + 1;
      
      // Track file change frequency
      const ext = file.path.split('.').pop();
      state.gitStats.fileChangeFrequency[ext] = (state.gitStats.fileChangeFrequency[ext] || 0) + 1;
    });
  }
  
  // Update agent git stats
  agentData.commits++;
  agentData.linesAdded += linesAdded;
  agentData.linesRemoved += linesRemoved;
  agentData.lastCommit = { sha, timestamp, message };
  agentData.commitSizes.push(linesAdded + linesRemoved);
  filesChanged.forEach(f => agentData.filesChanged.add(f));
  
  // Update productivity metrics
  if (agentData.tasksCompleted > 0) {
    agentData.tasksPerCommit = (agentData.tasksCompleted / agentData.commits).toFixed(2);
    agentData.linesPerTask = (agentData.linesAdded / agentData.tasksCompleted).toFixed(0);
  }
  if (agentData.totalTokens > 0 && agentData.linesAdded > 0) {
    agentData.tokensPerLine = (agentData.totalTokens / agentData.linesAdded).toFixed(1);
  }
  
  // Update global git stats
  state.gitStats.totalCommits++;
  state.gitStats.totalLinesAdded += linesAdded;
  state.gitStats.totalLinesRemoved += linesRemoved;
  state.gitStats.hourlyActivity[hour]++;
  state.gitStats.dailyActivity[date] = (state.gitStats.dailyActivity[date] || 0) + 1;
  state.gitStats.commitSizes.push(linesAdded + linesRemoved);
  
  // Track agent commits
  state.gitStats.agentCommits[agent] = (state.gitStats.agentCommits[agent] || 0) + 1;
  
  // Add commit record
  const commit = {
    sha,
    agent,
    message: message || '',
    branch: branch || 'main',
    timestamp,
    linesAdded,
    linesRemoved,
    filesChanged: filesChanged.length,
    files: filesChanged
  };
  
  state.commits.push(commit);
  if (state.commits.length > 500) state.commits.shift();
  
  // Add to events
  state.events.push({ type: 'commit', ...commit });
  if (state.events.length > 1000) state.events.shift();
}

// Calculate agent metrics
function calculateAgentMetrics(agentId) {
  const agent = state.agents[agentId];
  if (!agent) return null;
  
  const taskTimes = agent.taskTimes || [];
  const medianTime = taskTimes.length > 0
    ? taskTimes.sort((a,b) => a-b)[Math.floor(taskTimes.length / 2)]
    : 0;
  
  const avgCommitSize = agent.commitSizes.length > 0
    ? agent.commitSizes.reduce((a,b) => a+b, 0) / agent.commitSizes.length
    : 0;
  
  // Calculate token cost (rough estimate)
  const tokenCost = (agent.totalTokens / 1000000) * 3; // $3 per million tokens
  
  return {
    ...agent,
    // Task metrics
    successRate: agent.tasksStarted > 0 
      ? ((agent.tasksCompleted / agent.tasksStarted) * 100).toFixed(1)
      : 0,
    avgTaskTime: taskTimes.length > 0
      ? taskTimes.reduce((a,b) => a+b, 0) / taskTimes.length
      : 0,
    medianTaskTime: medianTime,
    // Token metrics
    totalTokensK: (agent.totalTokens / 1000).toFixed(1) + 'k',
    tokenCost: '$' + tokenCost.toFixed(2),
    costPerTask: agent.tasksCompleted > 0 
      ? '$' + (tokenCost / agent.tasksCompleted).toFixed(2)
      : '$0',
    // Git metrics
    avgCommitSize: avgCommitSize.toFixed(0),
    codeChurn: agent.linesAdded + agent.linesRemoved,
    netLines: agent.linesAdded - agent.linesRemoved,
    filesChangedCount: agent.filesChanged.size,
    // Productivity score (composite)
    productivityScore: calculateProductivityScore(agent)
  };
}

// Calculate productivity score
function calculateProductivityScore(agent) {
  let score = 0;
  
  // Task completion (30 points)
  score += Math.min(agent.tasksCompleted * 1.5, 30);
  
  // Code contribution (25 points)
  score += Math.min(agent.commits * 2.5, 25);
  
  // Efficiency (20 points) - based on success rate
  const successRate = agent.tasksStarted > 0 
    ? (agent.tasksCompleted / agent.tasksStarted) * 100
    : 0;
  score += (successRate / 100) * 20;
  
  // Token efficiency (15 points) - fewer tokens per task is better
  if (agent.tokenEfficiency > 0) {
    const effScore = Math.max(0, 15 - (agent.tokenEfficiency / 10000) * 15);
    score += effScore;
  }
  
  // Activity (10 points) - based on recency
  const hoursSinceLastSeen = (Date.now() - agent.lastSeen) / 3600000;
  if (hoursSinceLastSeen < 1) score += 10;
  else if (hoursSinceLastSeen < 8) score += 5;
  
  return Math.round(score);
}

// Get leaderboard
function getLeaderboard() {
  return Object.keys(state.agents)
    .map(id => calculateAgentMetrics(id))
    .sort((a, b) => b.productivityScore - a.productivityScore);
}

// Get git insights
function getGitInsights() {
  const recentCommits = state.commits.slice(-20);
  const hotFiles = Object.entries(state.gitStats.hotFiles)
    .sort((a, b) => b[1] - a[1])
    .slice(0, 10);
  
  return {
    summary: {
      totalCommits: state.gitStats.totalCommits,
      totalLinesAdded: state.gitStats.totalLinesAdded,
      totalLinesRemoved: state.gitStats.totalLinesRemoved,
      netLines: state.gitStats.totalLinesAdded - state.gitStats.totalLinesRemoved,
      avgCommitSize: state.gitStats.commitSizes.length > 0
        ? (state.gitStats.commitSizes.reduce((a,b) => a+b, 0) / state.gitStats.commitSizes.length).toFixed(0)
        : 0
    },
    hotFiles: hotFiles.map(([path, count]) => ({ path, changes: count })),
    recentCommits,
    hourlyActivity: state.gitStats.hourlyActivity,
    topContributors: Object.entries(state.gitStats.agentCommits)
      .sort((a, b) => b[1] - a[1])
      .slice(0, 5)
      .map(([agent, commits]) => ({ agent, commits })),
    fileTypes: state.gitStats.fileChangeFrequency
  };
}

// Broadcast to WebSocket clients
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
  res.setHeader('Access-Control-Allow-Methods', 'GET, POST, OPTIONS');
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type');
  res.setHeader('Content-Type', 'application/json');
  
  if (req.method === 'POST' && req.url === '/event') {
    let body = '';
    
    req.on('data', chunk => {
      body += chunk.toString();
    });
    
    req.on('end', () => {
      try {
        const data = JSON.parse(body);
        
        switch(data.event) {
          case 'task_update':
            processTaskUpdate(data.data);
            broadcast('task-update', data.data);
            break;
            
          case 'git_commit':
            processGitCommit({ agent: data.agent, ...data.data });
            broadcast('git-commit', { agent: data.agent, ...data.data });
            break;
            
          default:
            // Store unknown events
            state.events.push(data);
            broadcast('event', data);
        }
        
        // Always broadcast updated stats
        broadcast('stats-update', {
          leaderboard: getLeaderboard(),
          gitInsights: getGitInsights()
        });
        
        res.writeHead(200);
        res.end(JSON.stringify({ ok: true }));
        
      } catch (error) {
        res.writeHead(400);
        res.end(JSON.stringify({ ok: false, error: error.message }));
      }
    });
    
  } else if (req.url === '/state') {
    res.writeHead(200);
    res.end(JSON.stringify(state));
    
  } else if (req.url === '/leaderboard') {
    res.writeHead(200);
    res.end(JSON.stringify(getLeaderboard()));
    
  } else if (req.url === '/git-insights') {
    res.writeHead(200);
    res.end(JSON.stringify(getGitInsights()));
    
  } else if (req.url === '/events') {
    res.writeHead(200);
    res.end(JSON.stringify(state.events.slice(-100)));
    
  } else if (req.url === '/ready-tasks') {
    res.writeHead(200);
    res.end(JSON.stringify(getReadyTasks()));
    
  } else if (req.method === 'POST' && req.url === '/claim-task') {
    let body = '';
    
    req.on('data', chunk => {
      body += chunk.toString();
    });
    
    req.on('end', () => {
      try {
        const { agent, task } = JSON.parse(body);
        
        // Atomic claim check
        const taskState = state.tasks[task];
        
        // Check if task is already claimed or in progress
        if (taskState && (taskState.status === 'started' || 
                          taskState.status === 'in_progress' || 
                          taskState.status === 'completed')) {
          res.writeHead(409); // Conflict
          res.end(JSON.stringify({ 
            ok: false, 
            error: 'Task already claimed',
            currentStatus: taskState.status,
            currentAgent: taskState.agent
          }));
          return;
        }
        
        // Check if task exists in DAG and dependencies are met
        const readyTasks = getReadyTasks();
        const isReady = readyTasks.ready.some(t => t.task === task);
        
        if (!isReady) {
          res.writeHead(400);
          res.end(JSON.stringify({ 
            ok: false, 
            error: 'Task not ready or does not exist',
            blocked: readyTasks.blocked.find(t => t.task === task)
          }));
          return;
        }
        
        // Claim the task atomically
        initAgent(agent);
        const timestamp = Math.floor(Date.now() / 1000);
        
        state.tasks[task] = {
          status: 'started',
          agent: agent,
          timestamp: timestamp,
          startTime: timestamp,
          claimedAt: Date.now()
        };
        
        state.agents[agent].tasksStarted++;
        state.agents[agent].currentTask = task;
        state.agents[agent].lastSeen = Date.now();
        
        // Broadcast the claim
        broadcast('task-claimed', {
          task,
          agent,
          timestamp
        });
        
        res.writeHead(200);
        res.end(JSON.stringify({ 
          ok: true,
          task,
          agent,
          message: 'Task successfully claimed'
        }));
        
        console.log(`${agent} claimed ${task}`);
        
      } catch (error) {
        res.writeHead(400);
        res.end(JSON.stringify({ ok: false, error: error.message }));
      }
    });
    
  } else {
    res.writeHead(404);
    res.end(JSON.stringify({ ok: false, error: 'Not found' }));
  }
});

// WebSocket server
const wss = new WebSocketServer({ server });

wss.on('connection', (ws) => {
  console.log('Client connected');
  wsClients.add(ws);
  
  // Send initial state
  ws.send(JSON.stringify({
    type: 'init',
    dag: state.dag,
    tasks: state.tasks,
    agents: Object.fromEntries(
      Object.entries(state.agents).map(([k, v]) => [k, calculateAgentMetrics(k)])
    ),
    leaderboard: getLeaderboard(),
    gitInsights: getGitInsights(),
    gitStats: state.gitStats,
    events: state.events.slice(-50)
  }));
  
  ws.on('close', () => wsClients.delete(ws));
});

// Load state and DAG on startup
loadState();
loadDAG();

// Save state every 30 seconds
setInterval(saveState, 30000);

// Graceful shutdown
process.on('SIGINT', () => {
  console.log('\nSaving state...');
  saveState();
  process.exit(0);
});

server.listen(PORT, () => {
  console.log(`
╔════════════════════════════════════════════════════════╗
║          DAG STATE SERVER - FULL ANALYTICS            ║
╠════════════════════════════════════════════════════════╣
║  POST http://localhost:${PORT}/claim-task              ║
║  POST http://localhost:${PORT}/event                   ║
║  GET  http://localhost:${PORT}/ready-tasks             ║
║  GET  http://localhost:${PORT}/leaderboard             ║
║  GET  http://localhost:${PORT}/git-insights            ║
║  GET  http://localhost:${PORT}/state                   ║
║  GET  http://localhost:${PORT}/events                  ║
║  WS   ws://localhost:${PORT}                           ║
╚════════════════════════════════════════════════════════╝

Event Types:

Claim Task (ATOMIC):
curl -X POST http://localhost:${PORT}/claim-task \\
  -H "Content-Type: application/json" \\
  -d '{"agent": "mr_clean", "task": "P1.T003"}'

Task Update (with total tokens):
curl -X POST http://localhost:${PORT}/event \\
  -H "Content-Type: application/json" \\
  -d '{
    "event": "task_update",
    "data": {
      "agent": "mr_clean",
      "task": "P1.T003",
      "status": "completed",
      "timestamp": ${Math.floor(Date.now()/1000)},
      "total_tokens": 1234567
    }
  }'

Git Commit:
curl -X POST http://localhost:${PORT}/event \\
  -H "Content-Type: application/json" \\
  -d '{
    "event": "git_commit",
    "agent": "mr_clean",
    "data": {
      "sha": "abc123",
      "message": "Fix authentication bug",
      "branch": "main",
      "timestamp": ${Math.floor(Date.now()/1000)},
      "files": [
        {
          "path": "src/auth.js",
          "additions": 45,
          "deletions": 12
        }
      ]
    }
  }'
  `);
});