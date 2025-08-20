/**
 * HTML Generator with Live WebSocket Updates for DAG State
 * Extends Cytoscape.js visualization with real-time state updates
 */

import HTMLGeneratorCytoscape from './htmlGeneratorCytoscape.js';

export default class HTMLGeneratorLive extends HTMLGeneratorCytoscape {
  constructor(graphData, options = {}) {
    super(graphData, options);
    this.wsPort = options.wsPort || 8080;
    this.wsHost = options.wsHost || 'localhost';
  }

  /**
   * Generate visualization script with WebSocket connection
   * @returns {string} JavaScript code for visualization
   */
  generateVisualizationScript() {
    return `
    <script>
        (function() {
            // State color mapping
            const stateColors = {
                pending: '#94a3b8',      // slate-400
                started: '#60a5fa',      // blue-400
                in_progress: '#3b82f6',  // blue-500
                failed: '#ef4444',       // red-500
                blocked: '#f97316',      // orange-500
                completed: '#22c55e'     // green-500
            };

            const stateBorderColors = {
                pending: '#64748b',      // slate-500
                started: '#3b82f6',      // blue-500
                in_progress: '#2563eb',  // blue-600
                failed: '#dc2626',       // red-600
                blocked: '#ea580c',      // orange-600
                completed: '#16a34a'     // green-600
            };

            // WebSocket connection
            let ws = null;
            let reconnectInterval = null;
            let isConnected = false;

            // Create status indicator
            const statusIndicator = document.createElement('div');
            statusIndicator.className = 'fixed bottom-4 right-4 px-3 py-2 rounded-lg text-sm font-medium flex items-center gap-2 z-50';
            statusIndicator.innerHTML = \`
                <div class="w-2 h-2 rounded-full" id="ws-status-dot"></div>
                <span id="ws-status-text">Connecting...</span>
            \`;
            document.body.appendChild(statusIndicator);

            function updateConnectionStatus(connected, message = null) {
                isConnected = connected;
                const dot = document.getElementById('ws-status-dot');
                const text = document.getElementById('ws-status-text');
                
                if (connected) {
                    statusIndicator.className = 'fixed bottom-4 right-4 px-3 py-2 rounded-lg text-sm font-medium flex items-center gap-2 z-50 bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200';
                    dot.className = 'w-2 h-2 rounded-full bg-green-500 animate-pulse';
                    text.textContent = 'Live';
                } else {
                    statusIndicator.className = 'fixed bottom-4 right-4 px-3 py-2 rounded-lg text-sm font-medium flex items-center gap-2 z-50 bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200';
                    dot.className = 'w-2 h-2 rounded-full bg-red-500';
                    text.textContent = message || 'Disconnected';
                }
            }

            function connectWebSocket() {
                if (ws && ws.readyState === WebSocket.OPEN) {
                    return;
                }

                try {
                    ws = new WebSocket('ws://${this.wsHost}:${this.wsPort}');
                    
                    ws.onopen = () => {
                        console.log('WebSocket connected');
                        updateConnectionStatus(true);
                        
                        // Clear reconnect interval if it exists
                        if (reconnectInterval) {
                            clearInterval(reconnectInterval);
                            reconnectInterval = null;
                        }
                    };
                    
                    ws.onmessage = (event) => {
                        try {
                            const message = JSON.parse(event.data);
                            
                            if (message.type === 'state-update') {
                                updateNodeStates(message);
                            } else if (message.type === 'pong') {
                                // Handle pong response
                            }
                        } catch (error) {
                            console.error('Error processing WebSocket message:', error);
                        }
                    };
                    
                    ws.onerror = (error) => {
                        console.error('WebSocket error:', error);
                        updateConnectionStatus(false, 'Connection error');
                    };
                    
                    ws.onclose = () => {
                        console.log('WebSocket disconnected');
                        updateConnectionStatus(false, 'Reconnecting...');
                        
                        // Set up reconnection
                        if (!reconnectInterval) {
                            reconnectInterval = setInterval(() => {
                                console.log('Attempting to reconnect...');
                                connectWebSocket();
                            }, 3000);
                        }
                    };
                    
                } catch (error) {
                    console.error('Failed to connect WebSocket:', error);
                    updateConnectionStatus(false, 'Connection failed');
                }
            }

            function updateNodeStates(message) {
                if (!cy) return;
                
                const states = message.fullState || message.taskState || {};
                
                Object.entries(states).forEach(([taskId, stateData]) => {
                    const node = cy.getElementById(taskId);
                    if (node.length === 0) return;
                    
                    const state = typeof stateData === 'string' ? stateData : stateData.state;
                    const color = stateColors[state];
                    const borderColor = stateBorderColors[state];
                    
                    if (color) {
                        // Update node colors
                        node.style({
                            'background-color': color,
                            'border-color': borderColor
                        });
                        
                        // Add visual feedback for state change
                        node.addClass('state-changed');
                        setTimeout(() => node.removeClass('state-changed'), 1000);
                        
                        // Update node data
                        node.data('taskState', state);
                        
                        // Update details panel if this node is selected
                        if (node.hasClass('selected')) {
                            updateDetailsPanel(node);
                        }
                    }
                });
                
                // Update stats
                updateStateStats(states);
            }

            function updateStateStats(states) {
                const stateCounts = {
                    pending: 0,
                    started: 0,
                    in_progress: 0,
                    failed: 0,
                    blocked: 0,
                    completed: 0
                };
                
                Object.values(states).forEach(stateData => {
                    const state = typeof stateData === 'string' ? stateData : stateData.state;
                    if (stateCounts.hasOwnProperty(state)) {
                        stateCounts[state]++;
                    }
                });
                
                // Add state stats to header if not exists
                let statsElement = document.getElementById('state-stats');
                if (!statsElement) {
                    const headerDiv = document.querySelector('header .flex');
                    if (headerDiv) {
                        statsElement = document.createElement('div');
                        statsElement.id = 'state-stats';
                        statsElement.className = 'flex gap-3 mt-2 text-xs';
                        headerDiv.appendChild(statsElement);
                    }
                }
                
                if (statsElement) {
                    statsElement.innerHTML = Object.entries(stateCounts)
                        .filter(([_, count]) => count > 0)
                        .map(([state, count]) => \`
                            <span class="flex items-center gap-1">
                                <div class="w-3 h-3 rounded" style="background: \${stateColors[state]}"></div>
                                <span class="text-muted-foreground">\${state}: \${count}</span>
                            </span>
                        \`).join('');
                }
            }

            function updateDetailsPanel(node) {
                const detailsEl = document.getElementById('task-details');
                if (!detailsEl) return;
                
                const data = node.data();
                const taskState = data.taskState || 'unknown';
                
                // Find existing state element or create new one
                let stateElement = detailsEl.querySelector('.task-state');
                if (!stateElement) {
                    // Add state to existing details
                    const existingContent = detailsEl.innerHTML;
                    if (!existingContent.includes('Click on a task')) {
                        detailsEl.innerHTML = \`
                            <div class="task-state mb-3">
                                <span class="font-medium">State:</span>
                                <span class="inline-flex items-center gap-1 ml-2 px-2 py-1 rounded text-xs font-medium" 
                                      style="background: \${stateColors[taskState]}20; color: \${stateColors[taskState]}">
                                    <div class="w-2 h-2 rounded-full" style="background: \${stateColors[taskState]}"></div>
                                    \${taskState}
                                </span>
                            </div>
                            \${existingContent}
                        \`;
                    }
                } else {
                    // Update existing state element
                    stateElement.innerHTML = \`
                        <span class="font-medium">State:</span>
                        <span class="inline-flex items-center gap-1 ml-2 px-2 py-1 rounded text-xs font-medium" 
                              style="background: \${stateColors[taskState]}20; color: \${stateColors[taskState]}">
                            <div class="w-2 h-2 rounded-full" style="background: \${stateColors[taskState]}"></div>
                            \${taskState}
                        </span>
                    \`;
                }
            }

            // Add state change animation style
            const style = document.createElement('style');
            style.textContent = \`
                .state-changed {
                    animation: stateChange 1s ease-in-out;
                }
                
                @keyframes stateChange {
                    0% { opacity: 1; }
                    50% { opacity: 0.5; transform: scale(1.1); }
                    100% { opacity: 1; transform: scale(1); }
                }
            \`;
            document.head.appendChild(style);

            // Initialize WebSocket connection
            connectWebSocket();

            // Send periodic ping to keep connection alive
            setInterval(() => {
                if (ws && ws.readyState === WebSocket.OPEN) {
                    ws.send(JSON.stringify({ type: 'ping' }));
                }
            }, 30000);

            // Initialize Cytoscape (get parent script without <script> tags)
            ${super.generateVisualizationScript().replace(/<script>|<\/script>/g, '')}

            // Modify the details panel click handler to include state
            const originalNodeClick = cy.on;
            cy.on('tap', 'node', function(evt) {
                const node = evt.target;
                updateDetailsPanel(node);
            });

        })();
    </script>`;
  }
}