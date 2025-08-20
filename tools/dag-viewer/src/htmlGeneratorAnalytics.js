/**
 * Analytics HTML Generator - Agent tracking, leaderboards, and paths
 */

import HTMLGeneratorCytoscape from './htmlGeneratorCytoscape.js';

export default class HTMLGeneratorAnalytics extends HTMLGeneratorCytoscape {
  constructor(graphData, options = {}) {
    super(graphData, options);
    this.wsPort = options.wsPort || 8080;
  }

  generateVisualizationScript() {
    return `
    <script>
        (function() {
            // State management
            let currentTab = 'task-details';
            let selectedAgent = null;
            let agentStats = {};
            let leaderboard = [];
            let taskStates = {};
            let agentPaths = {};
            
            // Status colors
            const statusColors = {
                pending: '#94a3b8',
                started: '#22c55e',
                in_progress: '#3b82f6',
                failed: '#ef4444',
                blocked: '#f97316',
                completed: '#16a34a'
            };
            
            // Agent colors (for paths)
            const agentColors = [
                '#8b5cf6', '#ec4899', '#06b6d4', '#10b981', '#f59e0b',
                '#6366f1', '#84cc16', '#f43f5e', '#0ea5e9', '#a855f7'
            ];
            const agentColorMap = {};
            let colorIndex = 0;
            
            function getAgentColor(agentId) {
                if (!agentColorMap[agentId]) {
                    agentColorMap[agentId] = agentColors[colorIndex % agentColors.length];
                    colorIndex++;
                }
                return agentColorMap[agentId];
            }

            // Update side panel with tabs
            function initializeSidePanel() {
                const detailsPanel = document.querySelector('.w-80.border-l');
                if (!detailsPanel) return;
                
                detailsPanel.innerHTML = \`
                    <div class="flex flex-col h-full">
                        <!-- Tabs -->
                        <div class="border-b flex">
                            <button onclick="switchTab('task-details')" 
                                    class="tab-btn flex-1 px-4 py-3 text-sm font-medium tab-active" 
                                    data-tab="task-details">
                                Task Details
                            </button>
                            <button onclick="switchTab('agent-stats')" 
                                    class="tab-btn flex-1 px-4 py-3 text-sm font-medium" 
                                    data-tab="agent-stats">
                                Agent Stats
                            </button>
                            <button onclick="switchTab('leaderboard')" 
                                    class="tab-btn flex-1 px-4 py-3 text-sm font-medium" 
                                    data-tab="leaderboard">
                                Leaderboard
                            </button>
                        </div>
                        
                        <!-- Content -->
                        <div class="flex-1 overflow-y-auto p-6">
                            <div id="task-details" class="tab-content">
                                <h3 class="text-lg font-semibold mb-4">Task Details</h3>
                                <div class="space-y-3 text-sm">
                                    <p class="text-muted-foreground">Click on a task to see details</p>
                                </div>
                            </div>
                            
                            <div id="agent-stats" class="tab-content hidden">
                                <h3 class="text-lg font-semibold mb-4">Agent Statistics</h3>
                                <select id="agent-selector" class="w-full px-3 py-2 mb-4 rounded-md border bg-background text-sm">
                                    <option value="">Select an agent...</option>
                                </select>
                                <div id="agent-details" class="space-y-3 text-sm"></div>
                            </div>
                            
                            <div id="leaderboard" class="tab-content hidden">
                                <h3 class="text-lg font-semibold mb-4">Agent Leaderboard</h3>
                                <div id="leaderboard-content" class="space-y-2"></div>
                            </div>
                        </div>
                    </div>
                \`;
                
                // Add tab styles
                const style = document.createElement('style');
                style.textContent += \`
                    .tab-btn {
                        position: relative;
                        transition: all 0.2s;
                        border-bottom: 2px solid transparent;
                    }
                    .tab-btn:hover {
                        background-color: hsl(var(--accent));
                    }
                    .tab-btn.tab-active {
                        border-bottom-color: hsl(var(--primary));
                        background-color: hsl(var(--accent));
                    }
                    .agent-badge {
                        display: inline-flex;
                        align-items: center;
                        gap: 0.25rem;
                        padding: 0.125rem 0.5rem;
                        border-radius: 9999px;
                        font-size: 0.75rem;
                        font-weight: 500;
                    }
                    .stat-card {
                        padding: 0.75rem;
                        border-radius: 0.5rem;
                        background-color: hsl(var(--accent));
                        margin-bottom: 0.5rem;
                    }
                    .leaderboard-item {
                        display: flex;
                        align-items: center;
                        padding: 0.75rem;
                        border-radius: 0.5rem;
                        background-color: hsl(var(--card));
                        border: 1px solid hsl(var(--border));
                        transition: all 0.2s;
                    }
                    .leaderboard-item:hover {
                        transform: translateX(4px);
                        box-shadow: 0 2px 8px rgba(0,0,0,0.1);
                    }
                    .rank-badge {
                        width: 2rem;
                        height: 2rem;
                        display: flex;
                        align-items: center;
                        justify-content: center;
                        border-radius: 50%;
                        font-weight: bold;
                        margin-right: 0.75rem;
                    }
                    .rank-1 { background: linear-gradient(135deg, #fbbf24, #f59e0b); color: white; }
                    .rank-2 { background: linear-gradient(135deg, #cbd5e1, #94a3b8); color: white; }
                    .rank-3 { background: linear-gradient(135deg, #c9916a, #a0522d); color: white; }
                \`;
                document.head.appendChild(style);
            }

            // Switch tabs
            window.switchTab = function(tabName) {
                currentTab = tabName;
                
                // Update tab buttons
                document.querySelectorAll('.tab-btn').forEach(btn => {
                    if (btn.dataset.tab === tabName) {
                        btn.classList.add('tab-active');
                    } else {
                        btn.classList.remove('tab-active');
                    }
                });
                
                // Update content
                document.querySelectorAll('.tab-content').forEach(content => {
                    if (content.id === tabName) {
                        content.classList.remove('hidden');
                    } else {
                        content.classList.add('hidden');
                    }
                });
                
                // Refresh content
                if (tabName === 'leaderboard') {
                    updateLeaderboard();
                } else if (tabName === 'agent-stats') {
                    updateAgentSelector();
                }
            };

            // Update leaderboard display
            function updateLeaderboard() {
                const content = document.getElementById('leaderboard-content');
                if (!content) return;
                
                content.innerHTML = leaderboard.slice(0, 10).map((agent, index) => \`
                    <div class="leaderboard-item">
                        <div class="rank-badge rank-\${index + 1}">\${index + 1}</div>
                        <div class="flex-1">
                            <div class="flex items-center justify-between mb-1">
                                <span class="font-semibold" style="color: \${getAgentColor(agent.name)}">
                                    \${agent.name}
                                </span>
                                <span class="text-xs text-muted-foreground">
                                    \${agent.status || 'idle'}
                                </span>
                            </div>
                            <div class="grid grid-cols-2 gap-2 text-xs text-muted-foreground">
                                <div>Tasks: \${agent.tasksCompleted}</div>
                                <div>Success: \${agent.successRate}%</div>
                                <div>Med. Time: \${formatTime(agent.medianTime)}</div>
                                <div>Efficiency: \${agent.efficiency}</div>
                            </div>
                        </div>
                    </div>
                \`).join('');
            }

            // Update agent selector
            function updateAgentSelector() {
                const selector = document.getElementById('agent-selector');
                if (!selector) return;
                
                const currentValue = selector.value;
                selector.innerHTML = '<option value="">Select an agent...</option>' +
                    Object.keys(agentStats).map(agentId => 
                        \`<option value="\${agentId}" \${agentId === currentValue ? 'selected' : ''}>
                            \${agentId} (\${agentStats[agentId].tasksCompleted} tasks)
                        </option>\`
                    ).join('');
                
                selector.onchange = (e) => {
                    if (e.target.value) {
                        showAgentDetails(e.target.value);
                    }
                };
                
                if (currentValue) {
                    showAgentDetails(currentValue);
                }
            }

            // Show agent details
            function showAgentDetails(agentId) {
                const agent = agentStats[agentId];
                if (!agent) return;
                
                const details = document.getElementById('agent-details');
                if (!details) return;
                
                // Highlight agent's current task
                if (cy) {
                    cy.nodes().removeClass('agent-highlight');
                    if (agent.currentTask) {
                        cy.getElementById(agent.currentTask).addClass('agent-highlight');
                    }
                    
                    // Show agent path
                    highlightAgentPath(agentId);
                }
                
                details.innerHTML = \`
                    <div class="stat-card">
                        <div class="flex items-center justify-between mb-2">
                            <span class="font-semibold text-lg" style="color: \${getAgentColor(agentId)}">
                                \${agentId}
                            </span>
                            <span class="agent-badge" style="background: \${statusColors[agent.status] || '#666'}20; color: \${statusColors[agent.status] || '#666'}">
                                \${agent.status || 'idle'}
                            </span>
                        </div>
                        \${agent.currentTask ? \`
                            <div class="text-xs text-muted-foreground mb-2">
                                Current: \${agent.currentTask}
                            </div>
                        \` : ''}
                    </div>
                    
                    <div class="grid grid-cols-2 gap-2">
                        <div class="stat-card">
                            <div class="text-2xl font-bold">\${agent.tasksCompleted}</div>
                            <div class="text-xs text-muted-foreground">Tasks Completed</div>
                        </div>
                        <div class="stat-card">
                            <div class="text-2xl font-bold">\${agent.successRate}%</div>
                            <div class="text-xs text-muted-foreground">Success Rate</div>
                        </div>
                        <div class="stat-card">
                            <div class="text-2xl font-bold">\${formatTime(agent.medianTime)}</div>
                            <div class="text-xs text-muted-foreground">Median Time</div>
                        </div>
                        <div class="stat-card">
                            <div class="text-2xl font-bold">\${agent.efficiency}</div>
                            <div class="text-xs text-muted-foreground">Efficiency</div>
                        </div>
                    </div>
                    
                    <div class="mt-4">
                        <h4 class="font-semibold mb-2">Recent Path</h4>
                        <div class="space-y-1 max-h-48 overflow-y-auto">
                            \${(agent.path || []).slice(-10).reverse().map(p => \`
                                <div class="flex items-center gap-2 text-xs">
                                    <div class="w-2 h-2 rounded-full" style="background: \${statusColors[p.status]}"></div>
                                    <span class="font-mono">\${p.task}</span>
                                    <span class="text-muted-foreground">\${p.status}</span>
                                </div>
                            \`).join('')}
                        </div>
                    </div>
                \`;
            }

            // Highlight agent path on graph
            function highlightAgentPath(agentId) {
                if (!cy) return;
                
                const agent = agentStats[agentId];
                if (!agent || !agent.path) return;
                
                // Clear previous highlights
                cy.edges().removeClass('agent-path');
                
                // Create edge highlights for the path
                const color = getAgentColor(agentId);
                const tasks = agent.path.map(p => p.task);
                
                for (let i = 0; i < tasks.length - 1; i++) {
                    const source = cy.getElementById(tasks[i]);
                    const target = cy.getElementById(tasks[i + 1]);
                    
                    if (source.length && target.length) {
                        // Find or create edge between them
                        const edges = source.edgesWith(target);
                        edges.addClass('agent-path');
                        edges.style({
                            'line-color': color,
                            'target-arrow-color': color,
                            'width': 4,
                            'z-index': 1000
                        });
                    }
                }
            }

            // Format time
            function formatTime(ms) {
                if (!ms) return '0s';
                if (ms < 1000) return ms + 'ms';
                if (ms < 60000) return (ms / 1000).toFixed(1) + 's';
                if (ms < 3600000) return (ms / 60000).toFixed(1) + 'm';
                return (ms / 3600000).toFixed(1) + 'h';
            }

            // WebSocket connection
            let ws = null;
            
            function connect() {
                ws = new WebSocket('ws://localhost:${this.wsPort}');
                
                ws.onopen = () => {
                    console.log('Connected to analytics server');
                };
                
                ws.onmessage = (event) => {
                    const data = JSON.parse(event.data);
                    
                    switch(data.type) {
                        case 'init':
                            taskStates = data.taskStates || {};
                            agentStats = data.agentStats || {};
                            leaderboard = data.leaderboard || [];
                            
                            // Update all nodes
                            Object.entries(taskStates).forEach(([taskId, info]) => {
                                updateNode(taskId, info.status, info.agent);
                            });
                            
                            // Update displays
                            updateLeaderboard();
                            updateAgentSelector();
                            break;
                            
                        case 'task-update':
                            updateNode(data.task, data.status, data.agent);
                            taskStates[data.task] = { status: data.status, agent: data.agent };
                            break;
                            
                        case 'agent-update':
                            agentStats[data.agent] = data.stats;
                            if (currentTab === 'agent-stats' && document.getElementById('agent-selector')?.value === data.agent) {
                                showAgentDetails(data.agent);
                            }
                            break;
                            
                        case 'leaderboard-update':
                            leaderboard = data.leaderboard || [];
                            if (currentTab === 'leaderboard') {
                                updateLeaderboard();
                            }
                            break;
                    }
                };
                
                ws.onerror = (err) => {
                    console.error('WebSocket error:', err);
                };
                
                ws.onclose = () => {
                    console.log('Disconnected, reconnecting in 3s...');
                    setTimeout(connect, 3000);
                };
            }

            function updateNode(taskId, status, agent) {
                if (!cy) return;
                
                const node = cy.getElementById(taskId);
                if (node.length === 0) return;
                
                const color = statusColors[status] || '#666';
                
                node.style({
                    'background-color': color,
                    'border-color': color
                });
                
                // Add agent indicator
                if (agent && status === 'started' || status === 'in_progress') {
                    node.style({
                        'border-width': 4,
                        'border-color': getAgentColor(agent)
                    });
                }
                
                // Flash animation
                node.addClass('flash');
                setTimeout(() => node.removeClass('flash'), 500);
            }

            // Animation styles
            const style = document.createElement('style');
            style.textContent = \`
                .flash {
                    animation: flash 0.5s ease-in-out;
                }
                @keyframes flash {
                    50% { opacity: 0.5; }
                }
                .agent-highlight {
                    box-shadow: 0 0 20px rgba(139, 92, 246, 0.6);
                }
            \`;
            document.head.appendChild(style);

            // Initialize
            initializeSidePanel();
            connect();

            // Initialize Cytoscape (get parent script without <script> tags)
            ${super.generateVisualizationScript().replace(/<script>|<\/script>/g, '')}

            // Override node click to update task details
            if (typeof cy !== 'undefined') {
                cy.on('tap', 'node', function(evt) {
                    const node = evt.target;
                    const data = node.data();
                    const state = taskStates[data.id];
                    
                    if (currentTab === 'task-details') {
                        const detailsEl = document.querySelector('#task-details .space-y-3');
                        if (detailsEl) {
                            detailsEl.innerHTML = \`
                                <div><span class="font-medium">ID:</span> \${data.id}</div>
                                <div><span class="font-medium">Title:</span> \${data.fullTitle}</div>
                                \${state ? \`
                                    <div class="flex items-center gap-2">
                                        <span class="font-medium">Status:</span>
                                        <span class="agent-badge" style="background: \${statusColors[state.status]}20; color: \${statusColors[state.status]}">
                                            \${state.status}
                                        </span>
                                    </div>
                                    \${state.agent ? \`
                                        <div class="flex items-center gap-2">
                                            <span class="font-medium">Agent:</span>
                                            <span style="color: \${getAgentColor(state.agent)}">
                                                \${state.agent}
                                            </span>
                                        </div>
                                    \` : ''}
                                \` : ''}
                                \${data.feature ? \`<div><span class="font-medium">Feature:</span> \${data.feature}</div>\` : ''}
                                \${data.wave ? \`<div><span class="font-medium">Wave:</span> \${data.wave}</div>\` : ''}
                                \${data.description ? \`<div><span class="font-medium">Description:</span> \${data.description}</div>\` : ''}
                            \`;
                        }
                    }
                });
            }
        })();
    </script>`;
  }
}