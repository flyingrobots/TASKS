/**
 * Build graph structure from T.A.S.K.S data for dagre visualization
 */
export class GraphBuilder {
  constructor(data, verbose = false) {
    this.data = data;
    this.verbose = verbose;
    this.nodes = [];
    this.edges = [];
    this.featureMap = new Map();
    this.waveMap = new Map();
  }

  /**
   * Build the complete graph structure
   * @returns {Object} Graph data for visualization
   */
  build() {
    // Build feature lookup map
    this.buildFeatureMap();
    
    // Build wave lookup map
    this.buildWaveMap();
    
    // Process nodes (tasks)
    this.buildNodes();
    
    // Process edges (dependencies)
    this.buildEdges();
    
    // Add soft dependencies
    this.addSoftDependencies();

    if (this.verbose) {
      console.log(`✓ Built graph with ${this.nodes.length} nodes and ${this.edges.length} edges`);
    }

    // Calculate max wave number
    const maxWave = this.data.waves.waves.reduce((max, wave) => 
      Math.max(max, wave.waveNumber), 0);
    
    return {
      nodes: this.nodes,
      edges: this.edges,
      features: Array.from(this.featureMap.values()),
      waves: this.data.waves.waves,
      metrics: this.data.dag.metrics,
      config: {
        ...this.getLayoutConfig(),
        maxWave: maxWave
      }
    };
  }

  /**
   * Build feature lookup map
   */
  buildFeatureMap() {
    this.data.features.features.forEach(feature => {
      this.featureMap.set(feature.id, feature);
    });
  }

  /**
   * Build wave lookup map
   */
  buildWaveMap() {
    this.data.waves.waves.forEach(wave => {
      wave.tasks.forEach(taskId => {
        this.waveMap.set(taskId, wave.waveNumber);
      });
    });
  }

  /**
   * Build nodes from tasks
   */
  buildNodes() {
    this.data.tasks.tasks.forEach(task => {
      const feature = this.featureMap.get(task.feature_id);
      const waveNumber = this.waveMap.get(task.id);
      
      // Calculate estimated duration (PERT)
      const duration = task.duration ? 
        (task.duration.optimistic + 4 * task.duration.mostLikely + task.duration.pessimistic) / 6 :
        4; // Default 4 hours

      const node = {
        id: task.id,
        label: this.truncateLabel(task.title, 30),
        fullTitle: task.title,
        description: task.description,
        feature: feature?.id,
        featureName: feature?.title,
        category: task.category,
        duration: duration,
        durationUnits: task.durationUnits || 'hours',
        wave: waveNumber,
        priority: feature?.priority,
        skills: task.skillsRequired,
        metadata: {
          complexity: task.complexity,
          interfaces_produced: task.interfaces_produced,
          interfaces_consumed: task.interfaces_consumed,
          acceptance_checks: task.acceptance_checks
        }
      };

      this.nodes.push(node);
    });
  }

  /**
   * Build edges from dependencies
   */
  buildEdges() {
    // Add hard dependencies
    this.data.tasks.dependencies
      .filter(dep => dep.isHard !== false && dep.confidence >= (this.data.tasks.meta?.min_confidence || 0.7))
      .forEach(dep => {
        const edge = {
          id: `${dep.from}-${dep.to}`,
          from: dep.from,
          to: dep.to,
          type: dep.type,
          isHard: true,
          confidence: dep.confidence,
          reason: dep.reason,
          style: 'solid',
          weight: dep.confidence,
          color: this.getEdgeColor(dep.confidence)
        };

        this.edges.push(edge);
      });
  }

  /**
   * Add soft dependencies from dag.json
   */
  addSoftDependencies() {
    if (this.data.dag.softDeps) {
      this.data.dag.softDeps.forEach(dep => {
        const edge = {
          id: `${dep.from}-${dep.to}-soft`,
          from: dep.from,
          to: dep.to,
          type: dep.type,
          isHard: false,
          confidence: dep.confidence,
          reason: dep.reason,
          style: 'dashed',
          weight: dep.confidence * 0.5, // Lower weight for soft deps
          color: this.getEdgeColor(dep.confidence, true)
        };

        this.edges.push(edge);
      });
    }
  }

  /**
   * Get color for edge based on confidence
   * @param {number} confidence - Confidence level (0-1)
   * @param {boolean} soft - Whether this is a soft dependency
   * @returns {string} Color code
   */
  getEdgeColor(confidence, soft = false) {
    if (soft) {
      // Lighter colors for soft dependencies
      if (confidence >= 0.8) return '#90EE90'; // Light green
      if (confidence >= 0.6) return '#FFD700'; // Gold
      return '#FFB6C1'; // Light red
    } else {
      // Stronger colors for hard dependencies
      if (confidence >= 0.8) return '#228B22'; // Forest green
      if (confidence >= 0.6) return '#FFA500'; // Orange
      return '#DC143C'; // Crimson
    }
  }

  /**
   * Truncate label to specified length
   * @param {string} text - Text to truncate
   * @param {number} maxLength - Maximum length
   * @returns {string} Truncated text
   */
  truncateLabel(text, maxLength) {
    if (text.length <= maxLength) return text;
    return text.substring(0, maxLength - 3) + '...';
  }

  /**
   * Get layout configuration for dagre
   * @returns {Object} Layout config
   */
  getLayoutConfig() {
    return {
      rankdir: 'TB',        // Top to bottom
      nodesep: 50,         // Horizontal spacing
      ranksep: 80,         // Vertical spacing
      marginx: 20,
      marginy: 20,
      acyclicer: 'greedy',
      ranker: 'longest-path'
    };
  }

  /**
   * Get feature color palette
   * @returns {Object} Feature ID to color mapping
   */
  getFeatureColors() {
    const colors = [
      '#FF6B6B', '#4ECDC4', '#45B7D1', '#96CEB4', '#FECA57',
      '#DDA0DD', '#98D8C8', '#FFD93D', '#6BCB77', '#FF6B9D',
      '#C44569', '#2E86AB', '#A23B72', '#F18F01', '#574B90'
    ];

    const featureColors = {};
    let colorIndex = 0;

    this.featureMap.forEach((feature, id) => {
      featureColors[id] = colors[colorIndex % colors.length];
      colorIndex++;
    });

    // Default color for tasks without features
    featureColors['default'] = '#B0B0B0';

    return featureColors;
  }
}