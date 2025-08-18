/**
 * Modern HTML Generator with Shadcn-inspired UI for T.A.S.K.S DAG
 */
export class HTMLGeneratorModern {
  constructor(graphData, verbose = false) {
    this.graphData = graphData;
    this.verbose = verbose;
  }

  /**
   * Generate complete HTML document with modern UI
   * @returns {string} HTML content
   */
  generate() {
    const html = `<!DOCTYPE html>
<html lang="en" class="h-full">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>T.A.S.K.S DAG Visualization</title>
    <script src="https://cdn.tailwindcss.com"></script>
    ${this.generateStyles()}
</head>
<body class="h-full bg-background text-foreground">
    <div class="flex h-screen flex-col">
        <!-- Header -->
        <header class="border-b bg-card px-6 py-4">
            <div class="flex items-center justify-between">
                <div>
                    <h1 class="text-2xl font-semibold tracking-tight">T.A.S.K.S DAG Visualization</h1>
                    <div class="mt-1 flex gap-6 text-sm text-muted-foreground">
                        <span>Nodes: ${this.graphData.nodes.length}</span>
                        <span>Edges: ${this.graphData.edges.length}</span>
                        <span>Features: ${this.graphData.features.length}</span>
                        <span>Waves: ${this.graphData.waves.length}</span>
                    </div>
                </div>
                <button id="theme-toggle" class="rounded-lg border bg-background p-2 hover:bg-accent">
                    <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path class="sun-icon" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"></path>
                        <path class="moon-icon hidden" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"></path>
                    </svg>
                </button>
            </div>
        </header>

        <!-- Filters -->
        <div class="border-b bg-card/50 px-6 py-4">
            <div class="flex flex-wrap items-center gap-4">
                <div class="relative">
                    <input 
                        type="text" 
                        id="search" 
                        placeholder="Search tasks..."
                        class="h-10 w-64 rounded-md border bg-background px-3 text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
                    >
                </div>
                
                <select 
                    id="feature-filter"
                    class="h-10 rounded-md border bg-background px-3 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
                >
                    <option value="">All Features</option>
                    ${this.graphData.features.map(f => 
                      `<option value="${f.id}">${f.title}</option>`
                    ).join('')}
                </select>
                
                <select 
                    id="wave-filter"
                    class="h-10 rounded-md border bg-background px-3 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
                >
                    <option value="">All Waves</option>
                    ${this.graphData.waves.map(w => 
                      `<option value="${w.waveNumber}">Wave ${w.waveNumber}</option>`
                    ).join('')}
                </select>

                <div class="ml-auto flex gap-2">
                    <button id="reset-zoom" class="h-10 rounded-md border bg-background px-4 text-sm font-medium hover:bg-accent">
                        Reset Zoom
                    </button>
                    <button id="fit-to-screen" class="h-10 rounded-md border bg-background px-4 text-sm font-medium hover:bg-accent">
                        Fit to Screen
                    </button>
                </div>
            </div>
            
            <!-- Active Filters Badges -->
            <div id="active-filters" class="mt-3 flex flex-wrap gap-2"></div>
        </div>

        <!-- Main Content -->
        <div class="flex flex-1 overflow-hidden">
            <!-- Graph -->
            <div class="flex-1 bg-background">
                <svg id="graph" class="h-full w-full"></svg>
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
                    <svg width="30" height="2" class="overflow-visible">
                        <line x1="0" y1="1" x2="30" y2="1" stroke="#22c55e" stroke-width="2"/>
                    </svg>
                    <span class="text-muted-foreground">High confidence</span>
                </div>
                <div class="flex items-center gap-2">
                    <svg width="30" height="2" class="overflow-visible">
                        <line x1="0" y1="1" x2="30" y2="1" stroke="#f97316" stroke-width="2"/>
                    </svg>
                    <span class="text-muted-foreground">Medium confidence</span>
                </div>
                <div class="flex items-center gap-2">
                    <svg width="30" height="2" class="overflow-visible">
                        <line x1="0" y1="1" x2="30" y2="1" stroke="#ef4444" stroke-width="2"/>
                    </svg>
                    <span class="text-muted-foreground">Low confidence</span>
                </div>
                <div class="flex items-center gap-2">
                    <svg width="30" height="2" class="overflow-visible">
                        <line x1="0" y1="1" x2="30" y2="1" stroke="#86efac" stroke-width="2" stroke-dasharray="5,5"/>
                    </svg>
                    <span class="text-muted-foreground">Soft dependency</span>
                </div>
            </div>
        </div>
    </div>

    <script src="https://cdn.jsdelivr.net/npm/d3@7.8.5/dist/d3.min.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/dagre@0.8.5/dist/dagre.min.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/dagre-d3@0.6.4/dist/dagre-d3.min.js"></script>
    
    <script>
        // Embed graph data
        const graphData = ${JSON.stringify(this.graphData, null, 2)};
    </script>
    
    ${this.generateVisualizationScript()}
</body>
</html>`;

    if (this.verbose) {
      console.log('✓ Generated modern HTML visualization');
    }

    return html;
  }

  /**
   * Generate CSS styles with Shadcn design tokens
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
            --popover: 0 0% 100%;
            --popover-foreground: 222.2 84% 4.9%;
            --primary: 222.2 47.4% 11.2%;
            --primary-foreground: 210 40% 98%;
            --secondary: 210 40% 96.1%;
            --secondary-foreground: 222.2 47.4% 11.2%;
            --muted: 210 40% 96.1%;
            --muted-foreground: 215.4 16.3% 46.9%;
            --accent: 210 40% 96.1%;
            --accent-foreground: 222.2 47.4% 11.2%;
            --destructive: 0 84.2% 60.2%;
            --destructive-foreground: 210 40% 98%;
            --border: 214.3 31.8% 91.4%;
            --input: 214.3 31.8% 91.4%;
            --ring: 222.2 84% 4.9%;
            --radius: 0.5rem;
        }

        .dark {
            --background: 222.2 84% 4.9%;
            --foreground: 210 40% 98%;
            --card: 222.2 84% 4.9%;
            --card-foreground: 210 40% 98%;
            --popover: 222.2 84% 4.9%;
            --popover-foreground: 210 40% 98%;
            --primary: 210 40% 98%;
            --primary-foreground: 222.2 47.4% 11.2%;
            --secondary: 217.2 32.6% 17.5%;
            --secondary-foreground: 210 40% 98%;
            --muted: 217.2 32.6% 17.5%;
            --muted-foreground: 215 20.2% 65.1%;
            --accent: 217.2 32.6% 17.5%;
            --accent-foreground: 210 40% 98%;
            --destructive: 0 62.8% 30.6%;
            --destructive-foreground: 210 40% 98%;
            --border: 217.2 32.6% 17.5%;
            --input: 217.2 32.6% 17.5%;
            --ring: 212.7 26.8% 83.9%;
        }

        * {
            border-color: hsl(var(--border));
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
        }

        /* Tailwind utilities for our color scheme */
        .bg-background { background-color: hsl(var(--background)); }
        .bg-card { background-color: hsl(var(--card)); }
        .bg-accent { background-color: hsl(var(--accent)); }
        .text-foreground { color: hsl(var(--foreground)); }
        .text-muted-foreground { color: hsl(var(--muted-foreground)); }
        .border { border-color: hsl(var(--border)); }
        .ring-ring { --tw-ring-color: hsl(var(--ring)); }

        /* Graph styles */
        .node {
            cursor: pointer;
            transition: all 0.2s;
        }

        .node rect {
            stroke: hsl(var(--border));
            stroke-width: 1.5px;
            rx: 6;
            ry: 6;
        }

        .node text {
            font-size: 12px;
            font-weight: 500;
            fill: hsl(var(--foreground));
        }

        .node.highlighted rect {
            stroke: hsl(var(--primary));
            stroke-width: 2.5px;
            filter: drop-shadow(0 4px 6px rgba(0, 0, 0, 0.1));
        }

        .node.dimmed {
            opacity: 0.2;
        }

        .edgePath path {
            stroke-width: 1.5px;
            fill: none;
            transition: all 0.2s;
        }

        .edgePath.highlighted path {
            stroke-width: 2.5px;
        }

        .edgePath.dimmed {
            opacity: 0.2;
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

        .filter-badge:hover {
            background-color: hsl(var(--accent));
        }

        .filter-badge button {
            display: flex;
            align-items: center;
            justify-content: center;
            width: 1rem;
            height: 1rem;
            border-radius: 50%;
            background-color: hsl(var(--muted));
            transition: background-color 0.2s;
        }

        .filter-badge button:hover {
            background-color: hsl(var(--destructive));
            color: hsl(var(--destructive-foreground));
        }

        /* Tooltip */
        .tooltip {
            position: absolute;
            padding: 0.5rem 0.75rem;
            background-color: hsl(var(--popover));
            color: hsl(var(--popover-foreground));
            border: 1px solid hsl(var(--border));
            border-radius: 0.375rem;
            font-size: 0.875rem;
            box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
            pointer-events: none;
            z-index: 50;
        }

        /* Transitions */
        * {
            transition-property: background-color, border-color, color, fill, stroke;
            transition-duration: 200ms;
        }

        /* Dark mode transitions */
        html.transitioning,
        html.transitioning *,
        html.transitioning *:before,
        html.transitioning *:after {
            transition: all 200ms !important;
            transition-delay: 0 !important;
        }
    </style>`;
  }

  /**
   * Generate modern visualization JavaScript with fixed filtering
   * @returns {string} Script tag content
   */
  generateVisualizationScript() {
    return `<script>
        (function() {
            // Dark mode setup
            const html = document.documentElement;
            const themeToggle = document.getElementById('theme-toggle');
            const sunIcon = themeToggle.querySelector('.sun-icon');
            const moonIcon = themeToggle.querySelector('.moon-icon');
            
            // Check for saved theme preference
            const savedTheme = localStorage.getItem('theme');
            if (savedTheme === 'dark' || (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
                html.classList.add('dark');
                sunIcon.classList.add('hidden');
                moonIcon.classList.remove('hidden');
            }
            
            themeToggle.addEventListener('click', () => {
                html.classList.add('transitioning');
                html.classList.toggle('dark');
                const isDark = html.classList.contains('dark');
                localStorage.setItem('theme', isDark ? 'dark' : 'light');
                
                sunIcon.classList.toggle('hidden', isDark);
                moonIcon.classList.toggle('hidden', !isDark);
                
                setTimeout(() => {
                    html.classList.remove('transitioning');
                }, 200);
            });

            // Feature color mapping
            const featureColors = ${JSON.stringify(this.generateFeatureColors())};
            
            // Active filters state
            const activeFilters = {
                search: null,
                feature: null,
                wave: null
            };
            
            // Initialize dagre graph
            const g = new dagreD3.graphlib.Graph()
                .setGraph(graphData.config)
                .setDefaultEdgeLabel(() => ({}));

            // Add nodes
            graphData.nodes.forEach(node => {
                g.setNode(node.id, {
                    label: node.label,
                    class: \`category-\${node.category}\`,
                    style: \`fill: \${featureColors[node.feature] || featureColors.default}\`,
                    shape: getNodeShape(node.category),
                    width: Math.max(150, node.label.length * 7),
                    height: 40,
                    data: node
                });
            });

            // Add edges
            graphData.edges.forEach(edge => {
                g.setEdge(edge.from, edge.to, {
                    style: \`stroke: \${edge.color}; stroke-dasharray: \${edge.style === 'dashed' ? '5,5' : '0'}\`,
                    arrowheadStyle: \`fill: \${edge.color}; stroke: \${edge.color}\`,
                    curve: d3.curveBasis,
                    data: edge
                });
            });

            // Function to get node shape based on category
            function getNodeShape(category) {
                switch(category) {
                    case 'foundation': return 'rect';
                    case 'integration': return 'diamond';
                    default: return 'rect';
                }
            }

            // Create renderer
            const render = new dagreD3.render();

            // Set up SVG
            const svg = d3.select("#graph");
            const svgGroup = svg.append("g");

            // Render the graph
            render(svgGroup, g);

            // Get graph dimensions
            const graphWidth = g.graph().width;
            const graphHeight = g.graph().height;
            const width = svg.node().getBoundingClientRect().width;
            const height = svg.node().getBoundingClientRect().height;
            const initialScale = Math.min(width / graphWidth, height / graphHeight) * 0.9;
            
            // Set up zoom behavior
            const zoom = d3.zoom()
                .scaleExtent([0.1, 4])
                .on("zoom", (event) => {
                    svgGroup.attr("transform", event.transform);
                });

            svg.call(zoom);

            // Initial transform
            svg.call(
                zoom.transform,
                d3.zoomIdentity
                    .translate(width / 2 - graphWidth * initialScale / 2, 20)
                    .scale(initialScale)
            );

            // Add interactivity
            setupInteractivity();

            // Set up controls
            setupControls();

            function setupInteractivity() {
                // Node click handler
                svgGroup.selectAll(".node").on("click", function(event, nodeId) {
                    const node = g.node(nodeId);
                    if (node && node.data) {
                        showTaskDetails(node.data);
                        highlightNode(nodeId);
                    }
                });

                // Node hover
                svgGroup.selectAll(".node")
                    .on("mouseenter", function(event, nodeId) {
                        const node = g.node(nodeId);
                        if (node && node.data) {
                            showTooltip(event, node.data);
                        }
                    })
                    .on("mouseleave", hideTooltip);
            }

            function setupControls() {
                // Search functionality
                document.getElementById("search").addEventListener("input", (e) => {
                    const query = e.target.value.trim();
                    if (query) {
                        activeFilters.search = query;
                    } else {
                        activeFilters.search = null;
                    }
                    applyAllFilters();
                });

                // Feature filter
                document.getElementById("feature-filter").addEventListener("change", (e) => {
                    const featureId = e.target.value;
                    if (featureId) {
                        activeFilters.feature = featureId;
                    } else {
                        activeFilters.feature = null;
                    }
                    applyAllFilters();
                });

                // Wave filter
                document.getElementById("wave-filter").addEventListener("change", (e) => {
                    const waveValue = e.target.value;
                    if (waveValue) {
                        activeFilters.wave = parseInt(waveValue);
                    } else {
                        activeFilters.wave = null;
                    }
                    applyAllFilters();
                });

                // Reset zoom
                document.getElementById("reset-zoom").addEventListener("click", () => {
                    svg.transition().duration(750).call(
                        zoom.transform,
                        d3.zoomIdentity
                            .translate(width / 2 - graphWidth * initialScale / 2, 20)
                            .scale(initialScale)
                    );
                });

                // Fit to screen
                document.getElementById("fit-to-screen").addEventListener("click", () => {
                    const bounds = svgGroup.node().getBBox();
                    const fullWidth = svg.node().getBoundingClientRect().width;
                    const fullHeight = svg.node().getBoundingClientRect().height;
                    const width = bounds.width;
                    const height = bounds.height;
                    const midX = bounds.x + width / 2;
                    const midY = bounds.y + height / 2;
                    const scale = 0.9 / Math.max(width / fullWidth, height / fullHeight);
                    
                    svg.transition().duration(750).call(
                        zoom.transform,
                        d3.zoomIdentity
                            .translate(fullWidth / 2 - scale * midX, fullHeight / 2 - scale * midY)
                            .scale(scale)
                    );
                });
            }

            function applyAllFilters() {
                // Apply combined filters
                filterNodes(node => {
                    if (!node || !node.data) return false;
                    
                    // Check search filter
                    if (activeFilters.search) {
                        const query = activeFilters.search.toLowerCase();
                        const matchesSearch = 
                            node.data.fullTitle.toLowerCase().includes(query) ||
                            node.data.id.toLowerCase().includes(query);
                        if (!matchesSearch) return false;
                    }
                    
                    // Check feature filter
                    if (activeFilters.feature) {
                        if (node.data.feature !== activeFilters.feature) return false;
                    }
                    
                    // Check wave filter  
                    if (activeFilters.wave !== null && activeFilters.wave !== undefined) {
                        if (node.data.wave !== activeFilters.wave) return false;
                    }
                    
                    return true;
                });
                
                updateFilterBadges();
            }

            function filterNodes(predicate) {
                const matchingNodes = new Set();
                
                // First pass: identify matching nodes
                svgGroup.selectAll(".node").each(function(nodeId) {
                    const node = g.node(nodeId);
                    const matches = predicate(node);
                    if (matches) {
                        matchingNodes.add(nodeId);
                    }
                    d3.select(this).classed("dimmed", !matches);
                    d3.select(this).classed("highlighted", matches);
                });

                // Update edges: dim those that don't connect matching nodes
                svgGroup.selectAll(".edgePath").each(function(d) {
                    const edgeData = d.v && d.w ? d : null;
                    if (edgeData) {
                        const isRelevant = matchingNodes.has(edgeData.v) || matchingNodes.has(edgeData.w);
                        d3.select(this).classed("dimmed", !isRelevant);
                        d3.select(this).classed("highlighted", matchingNodes.has(edgeData.v) && matchingNodes.has(edgeData.w));
                    } else {
                        d3.select(this).classed("dimmed", true);
                    }
                });
                
                // If no filters active, clear all highlighting
                if (!activeFilters.search && !activeFilters.feature && activeFilters.wave === null) {
                    clearFilter();
                }
            }

            function clearFilter() {
                svgGroup.selectAll(".node").classed("dimmed", false).classed("highlighted", false);
                svgGroup.selectAll(".edgePath").classed("dimmed", false).classed("highlighted", false);
            }

            function updateFilterBadges() {
                const container = document.getElementById("active-filters");
                container.innerHTML = "";
                
                // Add search badge
                if (activeFilters.search) {
                    addFilterBadge("Search", activeFilters.search, () => {
                        activeFilters.search = null;
                        document.getElementById("search").value = "";
                        applyAllFilters();
                    });
                }
                
                // Add feature badge
                if (activeFilters.feature) {
                    const feature = graphData.features.find(f => f.id === activeFilters.feature);
                    addFilterBadge("Feature", feature ? feature.title : activeFilters.feature, () => {
                        activeFilters.feature = null;
                        document.getElementById("feature-filter").value = "";
                        applyAllFilters();
                    });
                }
                
                // Add wave badge
                if (activeFilters.wave !== null && activeFilters.wave !== undefined) {
                    addFilterBadge("Wave", activeFilters.wave, () => {
                        activeFilters.wave = null;
                        document.getElementById("wave-filter").value = "";
                        applyAllFilters();
                    });
                }
            }

            function addFilterBadge(type, value, onRemove) {
                const container = document.getElementById("active-filters");
                const badge = document.createElement("div");
                badge.className = "filter-badge";
                badge.innerHTML = \`
                    <span>\${type}: \${value}</span>
                    <button type="button">
                        <svg class="h-3 w-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                        </svg>
                    </button>
                \`;
                badge.querySelector("button").addEventListener("click", onRemove);
                container.appendChild(badge);
            }

            function highlightNode(nodeId) {
                svgGroup.selectAll(".node").classed("highlighted", false);
                svgGroup.select(\`.node[id="node-\${nodeId}"]\`).classed("highlighted", true);
            }

            function showTaskDetails(task) {
                const detailsHtml = \`
                    <div class="space-y-3">
                        <div>
                            <div class="text-xs font-medium text-muted-foreground mb-1">ID</div>
                            <div class="font-mono text-sm">\${task.id}</div>
                        </div>
                        <div>
                            <div class="text-xs font-medium text-muted-foreground mb-1">Title</div>
                            <div>\${task.fullTitle}</div>
                        </div>
                        <div>
                            <div class="text-xs font-medium text-muted-foreground mb-1">Description</div>
                            <div>\${task.description || 'N/A'}</div>
                        </div>
                        <div>
                            <div class="text-xs font-medium text-muted-foreground mb-1">Feature</div>
                            <div>\${task.featureName || 'N/A'}</div>
                        </div>
                        <div class="grid grid-cols-2 gap-3">
                            <div>
                                <div class="text-xs font-medium text-muted-foreground mb-1">Category</div>
                                <div class="capitalize">\${task.category}</div>
                            </div>
                            <div>
                                <div class="text-xs font-medium text-muted-foreground mb-1">Wave</div>
                                <div>\${task.wave || 'N/A'}</div>
                            </div>
                        </div>
                        <div>
                            <div class="text-xs font-medium text-muted-foreground mb-1">Duration</div>
                            <div>\${task.duration.toFixed(1)} \${task.durationUnits}</div>
                        </div>
                        <div>
                            <div class="text-xs font-medium text-muted-foreground mb-1">Skills</div>
                            <div>\${task.skills ? task.skills.join(', ') : 'N/A'}</div>
                        </div>
                    </div>
                \`;
                
                document.getElementById("task-details").innerHTML = detailsHtml;
            }

            let tooltip = null;

            function showTooltip(event, task) {
                hideTooltip();
                
                tooltip = d3.select("body")
                    .append("div")
                    .attr("class", "tooltip")
                    .html(\`
                        <div class="font-semibold">\${task.id}</div>
                        <div class="text-xs mt-1">\${task.fullTitle}</div>
                        <div class="text-xs text-muted-foreground mt-1">Duration: \${task.duration.toFixed(1)} \${task.durationUnits}</div>
                    \`)
                    .style("left", (event.pageX + 10) + "px")
                    .style("top", (event.pageY - 28) + "px");
            }

            function hideTooltip() {
                if (tooltip) {
                    tooltip.remove();
                    tooltip = null;
                }
            }
        })();
    </script>`;
  }

  /**
   * Generate feature color mapping with modern palette
   * @returns {Object} Feature to color map
   */
  generateFeatureColors() {
    const colors = [
      '#3b82f6', // blue
      '#10b981', // emerald
      '#f59e0b', // amber
      '#8b5cf6', // violet
      '#ec4899', // pink
      '#14b8a6', // teal
      '#f97316', // orange
      '#6366f1', // indigo
      '#84cc16', // lime
      '#06b6d4', // cyan
    ];

    const featureColors = {};
    this.graphData.features.forEach((feature, index) => {
      featureColors[feature.id] = colors[index % colors.length];
    });
    featureColors.default = '#6b7280'; // gray

    return featureColors;
  }
}