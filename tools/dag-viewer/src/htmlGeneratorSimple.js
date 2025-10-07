/**
 * Simple Live HTML Generator - just task status updates
 */

import HTMLGeneratorCytoscape from './htmlGeneratorCytoscape.js';

export default class HTMLGeneratorSimple extends HTMLGeneratorCytoscape {
  constructor(graphData, options = {}) {
    super(graphData, options);
    this.wsPort = options.wsPort || 8080;
  }

  generateVisualizationScript() {
    return `
    <script>
        (function() {
            // Status -> Color mapping
            const statusColors = {
                pending: '#94a3b8',      // gray
                started: '#22c55e',      // green
                in_progress: '#3b82f6',  // blue
                failed: '#ef4444',       // red
                blocked: '#f97316',      // orange
                completed: '#16a34a'     // darker green
            };

            // WebSocket connection
            let ws = null;

            function connect() {
                ws = new WebSocket('ws://localhost:${this.wsPort}');
                
                ws.onopen = () => {
                    console.log('Connected to state server');
                };
                
                ws.onmessage = (event) => {
                    const data = JSON.parse(event.data);
                    
                    if (data.type === 'init') {
                        // Initial state load
                        Object.entries(data.states).forEach(([taskId, info]) => {
                            updateNode(taskId, info.status);
                        });
                    } else {
                        // Single update
                        updateNode(data.taskId, data.status);
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

            function updateNode(taskId, status) {
                if (!cy) return;
                
                const node = cy.getElementById(taskId);
                if (node.length === 0) {
                    console.warn('Task not found in DAG:', taskId);
                    return;
                }
                
                const color = statusColors[status] || '#666';
                
                node.style({
                    'background-color': color,
                    'border-color': color
                });
                
                // Quick flash animation
                node.addClass('flash');
                setTimeout(() => node.removeClass('flash'), 500);
                
                console.log(\`Updated \${taskId} -> \${status}\`);
            }

            // Add flash animation
            const style = document.createElement('style');
            style.textContent = \`
                .flash {
                    animation: flash 0.5s ease-in-out;
                }
                @keyframes flash {
                    50% { opacity: 0.5; }
                }
            \`;
            document.head.appendChild(style);

            // Connect to WebSocket
            connect();

            // Initialize Cytoscape (get parent script without <script> tags)
            ${super.generateVisualizationScript().replace(/<script>|<\/script>/g, '')}
        })();
    </script>`;
  }
}