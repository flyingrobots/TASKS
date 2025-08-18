/**
 * Generate HTML visualization for T.A.S.K.S DAG
 */
export class HTMLGenerator {
  constructor(graphData, verbose = false) {
    this.graphData = graphData;
    this.verbose = verbose;
  }

  /**
   * Generate complete HTML document
   * @returns {string} HTML content
   */
  generate() {
    const html = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>T.A.S.K.S DAG Visualization</title>
    ${this.generateStyles()}
</head>
<body>
    <div id="container">
        <header id="header">
            <h1>T.A.S.K.S DAG Visualization</h1>
            <div id="metrics">
                <span>Nodes: ${this.graphData.nodes.length}</span>
                <span>Edges: ${this.graphData.edges.length}</span>
                <span>Features: ${this.graphData.features.length}</span>
                <span>Waves: ${this.graphData.waves.length}</span>
            </div>
        </header>
        
        <div id="controls">
            <div class="control-group">
                <label>Search: <input type="text" id="search" placeholder="Search tasks..."></label>
            </div>
            <div class="control-group">
                <label>Feature: 
                    <select id="feature-filter">
                        <option value="">All Features</option>
                        ${this.graphData.features.map(f => 
                          `<option value="${f.id}">${f.title}</option>`
                        ).join('')}
                    </select>
                </label>
            </div>
            <div class="control-group">
                <label>Wave: 
                    <select id="wave-filter">
                        <option value="">All Waves</option>
                        ${this.graphData.waves.map(w => 
                          `<option value="${w.waveNumber}">Wave ${w.waveNumber}</option>`
                        ).join('')}
                    </select>
                </label>
            </div>
            <div class="control-group">
                <button id="reset-zoom">Reset Zoom</button>
                <button id="fit-to-screen">Fit to Screen</button>
                <button id="clear-filters">Clear Filters</button>
            </div>
        </div>
        
        <div id="main">
            <svg id="graph"></svg>
            <div id="details-panel">
                <h3>Task Details</h3>
                <div id="task-details">
                    <p>Click on a task to see details</p>
                </div>
            </div>
        </div>
        
        <div id="legend">
            <h4>Legend</h4>
            <div class="legend-item">
                <svg width="30" height="20">
                    <line x1="0" y1="10" x2="30" y2="10" stroke="#228B22" stroke-width="2"/>
                </svg>
                <span>Hard dependency (high confidence)</span>
            </div>
            <div class="legend-item">
                <svg width="30" height="20">
                    <line x1="0" y1="10" x2="30" y2="10" stroke="#FFA500" stroke-width="2"/>
                </svg>
                <span>Hard dependency (medium confidence)</span>
            </div>
            <div class="legend-item">
                <svg width="30" height="20">
                    <line x1="0" y1="10" x2="30" y2="10" stroke="#DC143C" stroke-width="2"/>
                </svg>
                <span>Hard dependency (low confidence)</span>
            </div>
            <div class="legend-item">
                <svg width="30" height="20">
                    <line x1="0" y1="10" x2="30" y2="10" stroke="#90EE90" stroke-width="2" stroke-dasharray="5,5"/>
                </svg>
                <span>Soft dependency</span>
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
      console.log('✓ Generated HTML visualization');
    }

    return html;
  }

  /**
   * Generate CSS styles
   * @returns {string} Style tag content
   */
  generateStyles() {
    return `<style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            background: #f5f5f5;
            color: #333;
        }

        #container {
            display: flex;
            flex-direction: column;
            height: 100vh;
        }

        #header {
            background: #2c3e50;
            color: white;
            padding: 1rem;
            display: flex;
            justify-content: space-between;
            align-items: center;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }

        #header h1 {
            font-size: 1.5rem;
            font-weight: 500;
        }

        #metrics {
            display: flex;
            gap: 2rem;
            font-size: 0.9rem;
        }

        #metrics span {
            opacity: 0.9;
        }

        #controls {
            background: white;
            padding: 1rem;
            display: flex;
            gap: 1rem;
            align-items: center;
            border-bottom: 1px solid #ddd;
            flex-wrap: wrap;
        }

        .control-group {
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }

        .control-group label {
            font-size: 0.9rem;
            font-weight: 500;
        }

        input[type="text"], select {
            padding: 0.25rem 0.5rem;
            border: 1px solid #ddd;
            border-radius: 4px;
            font-size: 0.9rem;
        }

        button {
            padding: 0.25rem 0.75rem;
            background: #3498db;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 0.9rem;
        }

        button:hover {
            background: #2980b9;
        }

        #main {
            flex: 1;
            display: flex;
            overflow: hidden;
        }

        #graph {
            flex: 1;
            background: white;
            border: 1px solid #ddd;
        }

        #details-panel {
            width: 300px;
            background: white;
            border-left: 1px solid #ddd;
            padding: 1rem;
            overflow-y: auto;
        }

        #details-panel h3 {
            margin-bottom: 1rem;
            color: #2c3e50;
        }

        #task-details {
            font-size: 0.9rem;
            line-height: 1.6;
        }

        #task-details .detail-item {
            margin-bottom: 0.75rem;
            padding-bottom: 0.75rem;
            border-bottom: 1px solid #eee;
        }

        #task-details .detail-label {
            font-weight: 600;
            color: #555;
            display: block;
            margin-bottom: 0.25rem;
        }

        #task-details .detail-value {
            color: #666;
        }

        #legend {
            background: white;
            padding: 1rem;
            border-top: 1px solid #ddd;
            display: flex;
            gap: 2rem;
            align-items: center;
            flex-wrap: wrap;
        }

        #legend h4 {
            color: #2c3e50;
            font-size: 0.9rem;
        }

        .legend-item {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            font-size: 0.85rem;
        }

        .node {
            cursor: pointer;
        }

        .node rect {
            stroke: #333;
            stroke-width: 1.5px;
            rx: 5;
            ry: 5;
        }

        .node text {
            font-size: 12px;
            pointer-events: none;
        }

        .node.highlighted rect {
            stroke: #ff6b6b;
            stroke-width: 3px;
        }

        .node.dimmed {
            opacity: 0.3;
        }

        .edge {
            fill: none;
        }

        .edge.highlighted {
            stroke-width: 3px;
        }

        .edge.dimmed {
            opacity: 0.3;
        }

        .edgePath path {
            stroke-width: 1.5px;
            fill: none;
        }

        .edgeLabel {
            background-color: white;
            padding: 2px 4px;
            font-size: 10px;
        }

        .tooltip {
            position: absolute;
            text-align: left;
            padding: 8px;
            font-size: 12px;
            background: rgba(0, 0, 0, 0.8);
            color: white;
            border-radius: 4px;
            pointer-events: none;
            z-index: 100;
        }

        @media (max-width: 768px) {
            #main {
                flex-direction: column;
            }
            
            #details-panel {
                width: 100%;
                border-left: none;
                border-top: 1px solid #ddd;
                max-height: 200px;
            }
        }
    </style>`;
  }

  /**
   * Generate visualization JavaScript
   * @returns {string} Script tag content
   */
  generateVisualizationScript() {
    return `<script>
        (function() {
            // Feature color mapping
            const featureColors = ${JSON.stringify(this.generateFeatureColors())};
            
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

            // Center and fit the graph
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
                    const query = e.target.value.toLowerCase();
                    if (query) {
                        filterNodes(node => {
                            if (!node || !node.data) return false;
                            return node.data.fullTitle.toLowerCase().includes(query) ||
                                   node.data.id.toLowerCase().includes(query);
                        });
                    } else {
                        clearFilter();
                    }
                });

                // Feature filter
                document.getElementById("feature-filter").addEventListener("change", (e) => {
                    const featureId = e.target.value;
                    if (featureId) {
                        filterNodes(node => {
                            if (!node || !node.data) return false;
                            return node.data.feature === featureId;
                        });
                    } else {
                        clearFilter();
                    }
                });

                // Wave filter
                document.getElementById("wave-filter").addEventListener("change", (e) => {
                    const waveNumber = parseInt(e.target.value);
                    if (!isNaN(waveNumber)) {
                        filterNodes(node => {
                            if (!node || !node.data) return false;
                            return node.data.wave === waveNumber;
                        });
                    } else {
                        clearFilter();
                    }
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

                // Clear all filters
                document.getElementById("clear-filters").addEventListener("click", () => {
                    clearFilter();
                    document.getElementById("search").value = "";
                    document.getElementById("feature-filter").value = "";
                    document.getElementById("wave-filter").value = "";
                });
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
            }

            function clearFilter() {
                svgGroup.selectAll(".node").classed("dimmed", false).classed("highlighted", false);
                svgGroup.selectAll(".edgePath").classed("dimmed", false);
            }

            function highlightNode(nodeId) {
                svgGroup.selectAll(".node").classed("highlighted", false);
                svgGroup.selectAll(\`.node[id="node-\${nodeId}"]\`).classed("highlighted", true);
            }

            function showTaskDetails(task) {
                const detailsHtml = \`
                    <div class="detail-item">
                        <span class="detail-label">ID:</span>
                        <span class="detail-value">\${task.id}</span>
                    </div>
                    <div class="detail-item">
                        <span class="detail-label">Title:</span>
                        <span class="detail-value">\${task.fullTitle}</span>
                    </div>
                    <div class="detail-item">
                        <span class="detail-label">Description:</span>
                        <span class="detail-value">\${task.description || 'N/A'}</span>
                    </div>
                    <div class="detail-item">
                        <span class="detail-label">Feature:</span>
                        <span class="detail-value">\${task.featureName || 'N/A'}</span>
                    </div>
                    <div class="detail-item">
                        <span class="detail-label">Category:</span>
                        <span class="detail-value">\${task.category}</span>
                    </div>
                    <div class="detail-item">
                        <span class="detail-label">Duration:</span>
                        <span class="detail-value">\${task.duration.toFixed(1)} \${task.durationUnits}</span>
                    </div>
                    <div class="detail-item">
                        <span class="detail-label">Wave:</span>
                        <span class="detail-value">\${task.wave || 'N/A'}</span>
                    </div>
                    <div class="detail-item">
                        <span class="detail-label">Skills:</span>
                        <span class="detail-value">\${task.skills ? task.skills.join(', ') : 'N/A'}</span>
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
                        <strong>\${task.id}</strong><br>
                        \${task.fullTitle}<br>
                        Duration: \${task.duration.toFixed(1)} \${task.durationUnits}
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
   * Generate feature color mapping
   * @returns {Object} Feature to color map
   */
  generateFeatureColors() {
    const colors = [
      '#FF6B6B', '#4ECDC4', '#45B7D1', '#96CEB4', '#FECA57',
      '#DDA0DD', '#98D8C8', '#FFD93D', '#6BCB77', '#FF6B9D',
      '#C44569', '#2E86AB', '#A23B72', '#F18F01', '#574B90'
    ];

    const featureColors = {};
    this.graphData.features.forEach((feature, index) => {
      featureColors[feature.id] = colors[index % colors.length];
    });
    featureColors.default = '#B0B0B0';

    return featureColors;
  }
}