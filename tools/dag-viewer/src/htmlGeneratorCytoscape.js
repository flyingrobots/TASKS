/**
 * HTML Generator using Cytoscape.js for Canvas-based rendering
 * Provides fast, smooth graph visualization with working filters
 */

export default class HTMLGeneratorCytoscape {
  constructor(graphData, options = {}) {
    this.graphData = graphData;
    this.verbose = options.verbose || false;
  }

  /**
   * Generate complete HTML with Cytoscape.js visualization
   * @returns {string} Complete HTML document
   */
  generate() {
    if (this.verbose) {
      console.log('Generating Cytoscape.js HTML visualization...');
    }

    const html = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>T.A.S.K.S DAG Visualization</title>
    <script src="https://cdn.tailwindcss.com"></script>
    ${this.generateStyles()}
</head>
<body class="bg-background text-foreground">
    <div class="h-screen flex flex-col">
        <!-- Header -->
        <header class="border-b bg-card px-6 py-4">
            <div class="flex items-center justify-between">
                <div>
                    <h1 class="text-2xl font-bold">T.A.S.K.S DAG Visualization</h1>
                    <div class="flex gap-4 mt-2 text-sm text-muted-foreground">
                        <span>Nodes: ${this.graphData.nodes.length}</span>
                        <span>Edges: ${this.graphData.edges.length}</span>
                        <span>Features: ${this.graphData.features.length}</span>
                        <span>Waves: ${this.graphData.config.maxWave || 0}</span>
                    </div>
                </div>
                <button id="theme-toggle" class="p-2 rounded-lg hover:bg-accent transition-colors">
                    <svg class="w-5 h-5 sun-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"></path>
                    </svg>
                    <svg class="w-5 h-5 moon-icon hidden" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"></path>
                    </svg>
                </button>
            </div>
        </header>

        <!-- Controls -->
        <div class="border-b bg-card/50 px-6 py-3">
            <div class="flex flex-wrap items-center gap-4">
                <input 
                    type="text" 
                    id="search" 
                    placeholder="Search tasks..." 
                    class="px-3 py-1.5 rounded-md border bg-background text-sm w-64 focus:outline-none focus:ring-2 focus:ring-ring"
                />
                
                <select id="feature-filter" class="px-3 py-1.5 rounded-md border bg-background text-sm focus:outline-none focus:ring-2 focus:ring-ring">
                    <option value="">All Features</option>
                    ${this.graphData.features.map(f => 
                        `<option value="${f.id}">${f.title}</option>`
                    ).join('')}
                </select>
                
                <select id="wave-filter" class="px-3 py-1.5 rounded-md border bg-background text-sm focus:outline-none focus:ring-2 focus:ring-ring">
                    <option value="">All Waves</option>
                    ${Array.from({length: this.graphData.config.maxWave || 0}, (_, i) => 
                        `<option value="${i + 1}">Wave ${i + 1}</option>`
                    ).join('')}
                </select>

                <div class="flex gap-2 ml-auto">
                    <button id="reset-zoom" class="px-3 py-1.5 rounded-md bg-accent hover:bg-accent/80 text-sm font-medium transition-colors">
                        Reset Zoom
                    </button>
                    <button id="fit-to-screen" class="px-3 py-1.5 rounded-md bg-accent hover:bg-accent/80 text-sm font-medium transition-colors">
                        Fit to Screen
                    </button>
                </div>
            </div>
            
            <!-- Active Filters -->
            <div id="active-filters" class="flex flex-wrap gap-2 mt-3 min-h-[28px]"></div>
        </div>

        <!-- Main Content Area -->
        <div class="flex-1 flex">
            <!-- Graph Container -->
            <div class="flex-1 relative">
                <div id="cy" class="absolute inset-0"></div>
            </div>

            <!-- Details Panel -->
            <div class="w-80 border-l bg-card p-6 overflow-y-auto">
                <h3 class="text-lg font-semibold mb-4">Task Details</h3>
                <div id="task-details" class="space-y-3 text-sm">
                    <p class="text-muted-foreground">Click on a task to see details</p>
                </div>
            </div>
        </div>

        <!-- Legend -->
        <div class="border-t bg-card/50 px-6 py-3">
            <div class="flex items-center gap-8 text-sm">
                <span class="font-medium">Legend:</span>
                <div class="flex items-center gap-2">
                    <div class="w-4 h-4 rounded" style="background: #22c55e"></div>
                    <span class="text-muted-foreground">High confidence</span>
                </div>
                <div class="flex items-center gap-2">
                    <div class="w-4 h-4 rounded" style="background: #f97316"></div>
                    <span class="text-muted-foreground">Medium confidence</span>
                </div>
                <div class="flex items-center gap-2">
                    <div class="w-4 h-4 rounded" style="background: #ef4444"></div>
                    <span class="text-muted-foreground">Low confidence</span>
                </div>
                <div class="flex items-center gap-2">
                    <div class="w-4 h-4 rounded border-2 border-dashed" style="border-color: #86efac"></div>
                    <span class="text-muted-foreground">Soft dependency</span>
                </div>
            </div>
        </div>
    </div>

    <script src="https://cdn.jsdelivr.net/npm/cytoscape@3.28.1/dist/cytoscape.min.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/dagre@0.8.5/dist/dagre.min.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/cytoscape-dagre@2.5.0/cytoscape-dagre.js"></script>
    
    <script>
        // Register dagre layout
        cytoscape.use(cytoscapeDagre);
        
        // Embed graph data
        const graphData = ${JSON.stringify(this.graphData, null, 2)};
    </script>
    
    ${this.generateVisualizationScript()}
</body>
</html>`;

    if (this.verbose) {
      console.log('✓ Generated Cytoscape.js HTML visualization');
    }

    return html;
  }

  /**
   * Generate CSS styles
   * @returns {string} Style tag content
   */
  generateStyles() {
    return `<style>
        /* Shadcn-inspired design tokens */
        :root {
            --background: 0 0% 100%;
            --foreground: 222.2 84% 4.9%;
            --card: 0 0% 100%;
            --card-foreground: 222.2 84% 4.9%;
            --primary: 222.2 47.4% 11.2%;
            --primary-foreground: 210 40% 98%;
            --secondary: 210 40% 96.1%;
            --secondary-foreground: 222.2 47.4% 11.2%;
            --muted: 210 40% 96.1%;
            --muted-foreground: 215.4 16.3% 46.9%;
            --accent: 210 40% 96.1%;
            --accent-foreground: 222.2 47.4% 11.2%;
            --border: 214.3 31.8% 91.4%;
            --ring: 222.2 84% 4.9%;
        }

        .dark {
            --background: 222.2 84% 4.9%;
            --foreground: 210 40% 98%;
            --card: 222.2 84% 4.9%;
            --card-foreground: 210 40% 98%;
            --primary: 210 40% 98%;
            --primary-foreground: 222.2 47.4% 11.2%;
            --secondary: 217.2 32.6% 17.5%;
            --secondary-foreground: 210 40% 98%;
            --muted: 217.2 32.6% 17.5%;
            --muted-foreground: 215 20.2% 65.1%;
            --accent: 217.2 32.6% 17.5%;
            --accent-foreground: 210 40% 98%;
            --border: 217.2 32.6% 17.5%;
            --ring: 212.7 26.8% 83.9%;
        }

        * {
            border-color: hsl(var(--border));
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            margin: 0;
            padding: 0;
        }

        /* Tailwind utilities */
        .bg-background { background-color: hsl(var(--background)); }
        .bg-card { background-color: hsl(var(--card)); }
        .bg-accent { background-color: hsl(var(--accent)); }
        .text-foreground { color: hsl(var(--foreground)); }
        .text-muted-foreground { color: hsl(var(--muted-foreground)); }
        .border { border-color: hsl(var(--border)); }
        .ring-ring { --tw-ring-color: hsl(var(--ring)); }

        /* Cytoscape container */
        #cy {
            width: 100%;
            height: 100%;
            display: block;
        }

        /* Filter badges */
        .filter-badge {
            display: inline-flex;
            align-items: center;
            gap: 0.5rem;
            padding: 0.25rem 0.75rem;
            background-color: hsl(var(--secondary));
            color: hsl(var(--secondary-foreground));
            border-radius: 9999px;
            font-size: 0.875rem;
            font-weight: 500;
            transition: all 0.2s;
        }

        .filter-badge button {
            display: flex;
            align-items: center;
            justify-content: center;
            width: 1rem;
            height: 1rem;
            border-radius: 50%;
            background: transparent;
            border: none;
            cursor: pointer;
            padding: 0;
            transition: background-color 0.2s;
        }

        .filter-badge button:hover {
            background-color: hsl(var(--muted));
        }

        /* Dark mode transitions */
        * {
            transition-property: background-color, border-color, color;
            transition-duration: 200ms;
        }
    </style>`;
  }

  /**
   * Generate feature colors for consistent visualization
   * @returns {Object} Feature ID to color mapping
   */
  generateFeatureColors() {
    const colors = ['#3b82f6', '#8b5cf6', '#ec4899', '#f59e0b', '#10b981', '#06b6d4'];
    const featureColors = {};
    
    this.graphData.features.forEach((feature, index) => {
      featureColors[feature.id] = colors[index % colors.length];
    });
    
    return featureColors;
  }

  /**
   * Generate Cytoscape.js visualization script
   * @returns {string} Script content
   */
  generateVisualizationScript() {
    return `<script>
        (function() {
            // Dark mode setup
            const html = document.documentElement;
            const themeToggle = document.getElementById('theme-toggle');
            const sunIcon = themeToggle.querySelector('.sun-icon');
            const moonIcon = themeToggle.querySelector('.moon-icon');
            
            const savedTheme = localStorage.getItem('theme');
            if (savedTheme === 'dark' || (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
                html.classList.add('dark');
                sunIcon.classList.add('hidden');
                moonIcon.classList.remove('hidden');
            }
            
            themeToggle.addEventListener('click', () => {
                html.classList.toggle('dark');
                const isDark = html.classList.contains('dark');
                localStorage.setItem('theme', isDark ? 'dark' : 'light');
                
                sunIcon.classList.toggle('hidden', isDark);
                moonIcon.classList.toggle('hidden', !isDark);
                
                // Update Cytoscape styles for dark mode
                updateGraphStyles();
            });

            // Feature colors
            const featureColors = ${JSON.stringify(this.generateFeatureColors())};
            
            // Active filters
            const activeFilters = {
                search: null,
                feature: null,
                wave: null
            };

            // Convert graph data to Cytoscape format
            const elements = [];
            
            // Add nodes
            graphData.nodes.forEach(node => {
                elements.push({
                    data: {
                        id: node.id,
                        label: node.label,
                        fullTitle: node.fullTitle,
                        description: node.description,
                        feature: node.feature,
                        featureName: node.featureName,
                        category: node.category,
                        wave: node.wave,
                        priority: node.priority,
                        duration: node.duration,
                        durationUnits: node.durationUnits,
                        skills: node.skills
                    },
                    classes: \`category-\${node.category}\`
                });
            });
            
            // Add edges
            graphData.edges.forEach(edge => {
                // Map numeric confidence to category
                let confidenceLevel = 'low';
                if (edge.confidence >= 0.8) {
                    confidenceLevel = 'high';
                } else if (edge.confidence >= 0.6) {
                    confidenceLevel = 'medium';
                }
                
                elements.push({
                    data: {
                        id: \`edge-\${edge.from}-\${edge.to}\`,
                        source: edge.from,
                        target: edge.to,
                        type: edge.type,
                        confidence: edge.confidence,
                        confidenceLevel: confidenceLevel,
                        reason: edge.reason,
                        isHard: edge.isHard !== false,
                        isSoft: edge.isHard === false || edge.type === 'soft_dependency'
                    },
                    classes: (edge.isHard === false || edge.type === 'soft_dependency') ? 'soft-edge' : 'hard-edge'
                });
            });

            // Initialize Cytoscape
            const cy = cytoscape({
                container: document.getElementById('cy'),
                elements: elements,
                style: getGraphStyles(),
                layout: {
                    name: 'dagre',
                    rankDir: graphData.config.rankdir || 'TB',
                    nodeSep: 50,
                    rankSep: 100,
                    padding: 30
                },
                minZoom: 0.1,
                maxZoom: 3,
                wheelSensitivity: 0.2
            });

            // Get styles based on theme
            function getGraphStyles() {
                const isDark = html.classList.contains('dark');
                const bgColor = isDark ? '#020617' : '#ffffff';
                const nodeColor = isDark ? '#1e293b' : '#f8fafc';
                const textColor = isDark ? '#f1f5f9' : '#0f172a';
                const borderColor = isDark ? '#334155' : '#cbd5e1';
                const selectedColor = '#3b82f6';
                
                return [
                    {
                        selector: 'node',
                        style: {
                            'background-color': data => featureColors[data.feature] || nodeColor,
                            'label': 'data(label)',
                            'text-valign': 'center',
                            'text-halign': 'center',
                            'color': textColor,
                            'font-size': '12px',
                            'font-weight': '500',
                            'border-width': 2,
                            'border-color': borderColor,
                            'width': 40,
                            'height': 40,
                            'text-wrap': 'wrap',
                            'text-max-width': '100px',
                            'overlay-padding': '6px'
                        }
                    },
                    {
                        selector: 'node.highlighted',
                        style: {
                            'border-color': selectedColor,
                            'border-width': 3,
                            'background-color': selectedColor,
                            'z-index': 999
                        }
                    },
                    {
                        selector: 'node.selected',
                        style: {
                            'border-color': '#fbbf24',
                            'border-width': 5,
                            'width': 50,
                            'height': 50,
                            'background-blacken': -0.2,
                            'overlay-color': '#fbbf24',
                            'overlay-padding': '12px',
                            'overlay-opacity': 0.3,
                            'z-index': 1000
                        }
                    },
                    {
                        selector: 'node.dimmed',
                        style: {
                            'opacity': 0.2
                        }
                    },
                    {
                        selector: 'edge',
                        style: {
                            'width': 2,
                            'line-color': function(ele) {
                                const level = ele.data('confidenceLevel');
                                if (level === 'high') return '#22c55e';
                                if (level === 'medium') return '#f97316';
                                return '#ef4444';
                            },
                            'target-arrow-color': function(ele) {
                                const level = ele.data('confidenceLevel');
                                if (level === 'high') return '#22c55e';
                                if (level === 'medium') return '#f97316';
                                return '#ef4444';
                            },
                            'target-arrow-shape': 'triangle',
                            'curve-style': 'bezier',
                            'arrow-scale': 1.2
                        }
                    },
                    {
                        selector: 'edge.soft-edge',
                        style: {
                            'line-style': 'dashed',
                            'line-dash-pattern': [6, 3]
                        }
                    },
                    {
                        selector: 'edge.dimmed',
                        style: {
                            'opacity': 0.2
                        }
                    },
                    {
                        selector: 'edge.selected',
                        style: {
                            'width': 4,
                            'z-index': 999,
                            'overlay-color': '#fbbf24',
                            'overlay-padding': '6px',
                            'overlay-opacity': 0.3
                        }
                    },
                    {
                        selector: ':selected',
                        style: {
                            'border-color': selectedColor,
                            'border-width': 3
                        }
                    }
                ];
            }

            // Update styles when theme changes
            function updateGraphStyles() {
                cy.style(getGraphStyles());
            }

            // Apply filters
            function applyFilters() {
                // Reset all nodes and edges
                cy.elements().removeClass('dimmed highlighted');
                
                // Build filter conditions
                let hasActiveFilter = false;
                let filteredNodes = cy.nodes();
                
                if (activeFilters.search) {
                    hasActiveFilter = true;
                    const query = activeFilters.search.toLowerCase();
                    filteredNodes = filteredNodes.filter(node => {
                        const data = node.data();
                        return data.fullTitle.toLowerCase().includes(query) ||
                               data.id.toLowerCase().includes(query) ||
                               (data.description && data.description.toLowerCase().includes(query));
                    });
                }
                
                if (activeFilters.feature) {
                    hasActiveFilter = true;
                    filteredNodes = filteredNodes.filter(node => 
                        node.data('feature') === activeFilters.feature
                    );
                }
                
                if (activeFilters.wave !== null) {
                    hasActiveFilter = true;
                    filteredNodes = filteredNodes.filter(node => 
                        node.data('wave') === activeFilters.wave
                    );
                }
                
                // Apply visual changes
                if (hasActiveFilter) {
                    // Dim all elements first
                    cy.elements().addClass('dimmed');
                    
                    // Highlight filtered nodes and their edges
                    filteredNodes.removeClass('dimmed').addClass('highlighted');
                    
                    // Highlight edges connected to filtered nodes
                    filteredNodes.connectedEdges().removeClass('dimmed');
                }
                
                updateFilterBadges();
            }

            // Filter badge management
            function updateFilterBadges() {
                const container = document.getElementById('active-filters');
                container.innerHTML = '';
                
                if (activeFilters.search) {
                    addFilterBadge('Search', activeFilters.search, () => {
                        activeFilters.search = null;
                        document.getElementById('search').value = '';
                        applyFilters();
                    });
                }
                
                if (activeFilters.feature) {
                    const feature = graphData.features.find(f => f.id === activeFilters.feature);
                    addFilterBadge('Feature', feature ? feature.title : activeFilters.feature, () => {
                        activeFilters.feature = null;
                        document.getElementById('feature-filter').value = '';
                        applyFilters();
                    });
                }
                
                if (activeFilters.wave !== null) {
                    addFilterBadge('Wave', activeFilters.wave, () => {
                        activeFilters.wave = null;
                        document.getElementById('wave-filter').value = '';
                        applyFilters();
                    });
                }
            }

            function addFilterBadge(type, value, onRemove) {
                const container = document.getElementById('active-filters');
                const badge = document.createElement('div');
                badge.className = 'filter-badge';
                badge.innerHTML = \`
                    <span>\${type}: \${value}</span>
                    <button type="button">
                        <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                        </svg>
                    </button>
                \`;
                badge.querySelector('button').addEventListener('click', onRemove);
                container.appendChild(badge);
            }

            // Event handlers
            document.getElementById('search').addEventListener('input', (e) => {
                const query = e.target.value.trim();
                activeFilters.search = query || null;
                applyFilters();
            });

            document.getElementById('feature-filter').addEventListener('change', (e) => {
                activeFilters.feature = e.target.value || null;
                applyFilters();
            });

            document.getElementById('wave-filter').addEventListener('change', (e) => {
                const value = e.target.value;
                activeFilters.wave = value ? parseInt(value) : null;
                applyFilters();
            });

            document.getElementById('reset-zoom').addEventListener('click', () => {
                cy.fit();
            });

            document.getElementById('fit-to-screen').addEventListener('click', () => {
                cy.fit();
            });

            // Helper function to format duration
            function formatDuration(duration, units) {
                if (!duration) return 'N/A';
                
                let hours = duration;
                if (units === 'days') {
                    hours = duration * 8; // Assuming 8-hour work day
                } else if (units === 'weeks') {
                    hours = duration * 40; // Assuming 40-hour work week
                }
                
                if (hours < 1) {
                    const minutes = Math.round(hours * 60);
                    return \`\${minutes}m\`;
                } else if (hours === Math.floor(hours)) {
                    return \`\${hours}hr\`;
                } else {
                    const wholeHours = Math.floor(hours);
                    const minutes = Math.round((hours - wholeHours) * 60);
                    if (minutes === 0) {
                        return \`\${wholeHours}hr\`;
                    } else if (minutes === 60) {
                        return \`\${wholeHours + 1}hr\`;
                    } else {
                        return \`\${wholeHours}hr \${minutes}m\`;
                    }
                }
            }
            
            // Node click handler
            cy.on('tap', 'node', function(evt) {
                const node = evt.target;
                const data = node.data();
                
                // Clear previous selection
                cy.nodes().removeClass('selected');
                
                // Add selected class to clicked node
                node.addClass('selected');
                
                // Update details panel
                const detailsHtml = \`
                    <div class="space-y-3">
                        <div>
                            <div class="text-xs font-medium text-muted-foreground mb-1">ID</div>
                            <div class="font-mono text-sm">\${data.id}</div>
                        </div>
                        <div>
                            <div class="text-xs font-medium text-muted-foreground mb-1">Title</div>
                            <div>\${data.fullTitle}</div>
                        </div>
                        <div>
                            <div class="text-xs font-medium text-muted-foreground mb-1">Description</div>
                            <div>\${data.description || 'N/A'}</div>
                        </div>
                        <div>
                            <div class="text-xs font-medium text-muted-foreground mb-1">Feature</div>
                            <div>\${data.featureName || 'N/A'}</div>
                        </div>
                        <div class="grid grid-cols-2 gap-3">
                            <div>
                                <div class="text-xs font-medium text-muted-foreground mb-1">Category</div>
                                <div class="capitalize">\${data.category}</div>
                            </div>
                            <div>
                                <div class="text-xs font-medium text-muted-foreground mb-1">Wave</div>
                                <div>\${data.wave}</div>
                            </div>
                        </div>
                        <div class="grid grid-cols-2 gap-3">
                            <div>
                                <div class="text-xs font-medium text-muted-foreground mb-1">Priority</div>
                                <div class="capitalize">\${data.priority}</div>
                            </div>
                            <div>
                                <div class="text-xs font-medium text-muted-foreground mb-1">Duration</div>
                                <div>\${formatDuration(data.duration, data.durationUnits)}</div>
                            </div>
                        </div>
                        \${data.skills && data.skills.length > 0 ? \`
                        <div>
                            <div class="text-xs font-medium text-muted-foreground mb-1">Required Skills</div>
                            <div class="flex flex-wrap gap-1">
                                \${data.skills.map(skill => 
                                    \`<span class="px-2 py-0.5 bg-accent text-xs rounded">\${skill}</span>\`
                                ).join('')}
                            </div>
                        </div>
                        \` : ''}
                    </div>
                \`;
                
                document.getElementById('task-details').innerHTML = detailsHtml;
            });
            
            // Edge click handler
            cy.on('tap', 'edge', function(evt) {
                const edge = evt.target;
                const data = edge.data();
                
                // Clear all selections
                cy.nodes().removeClass('selected');
                cy.edges().removeClass('selected');
                
                // Add selected class to clicked edge
                edge.addClass('selected');
                
                // Format confidence percentage
                const confidencePct = Math.round((data.confidence || 0) * 100);
                
                // Update details panel with edge information
                const detailsHtml = \`
                    <div class="space-y-3">
                        <div>
                            <div class="text-xs font-medium text-muted-foreground mb-1">Dependency Type</div>
                            <div class="flex items-center gap-2">
                                <span class="font-medium">\${data.isHard ? 'Hard' : 'Soft'} Dependency</span>
                                \${data.isSoft ? '<span class="text-xs px-2 py-0.5 bg-accent rounded">Soft</span>' : ''}
                            </div>
                        </div>
                        <div>
                            <div class="text-xs font-medium text-muted-foreground mb-1">From → To</div>
                            <div class="font-mono text-sm">\${data.source} → \${data.target}</div>
                        </div>
                        <div>
                            <div class="text-xs font-medium text-muted-foreground mb-1">Confidence</div>
                            <div class="flex items-center gap-2">
                                <div class="flex items-center gap-1">
                                    <div class="w-3 h-3 rounded-full" style="background: \${
                                        data.confidenceLevel === 'high' ? '#22c55e' : 
                                        data.confidenceLevel === 'medium' ? '#f97316' : '#ef4444'
                                    }"></div>
                                    <span class="capitalize">\${data.confidenceLevel}</span>
                                </div>
                                <span class="text-sm text-muted-foreground">(\${confidencePct}%)</span>
                            </div>
                        </div>
                        <div>
                            <div class="text-xs font-medium text-muted-foreground mb-1">Type</div>
                            <div class="capitalize">\${data.type.replace(/_/g, ' ')}</div>
                        </div>
                        \${data.reason ? \`
                        <div>
                            <div class="text-xs font-medium text-muted-foreground mb-1">Reason</div>
                            <div>\${data.reason}</div>
                        </div>
                        \` : ''}
                    </div>
                \`;
                
                document.getElementById('task-details').innerHTML = detailsHtml;
                
                // Stop propagation to prevent background tap
                evt.stopPropagation();
            });
            
            // Click on background to clear selection
            cy.on('tap', function(evt) {
                if (evt.target === cy) {
                    cy.nodes().removeClass('selected');
                    cy.edges().removeClass('selected');
                    document.getElementById('task-details').innerHTML = '<p class="text-muted-foreground">Click on a task or edge to see details</p>';
                }
            });

            // Make Cytoscape instance available globally for debugging
            window.cy = cy;
        })();
    </script>`;
  }
}