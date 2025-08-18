import fs from 'fs-extra';
import path from 'path';

/**
 * Load and validate JSON files for T.A.S.K.S visualization
 */
export class DataLoader {
  constructor(verbose = false) {
    this.verbose = verbose;
  }

  /**
   * Load all required JSON files
   * @param {Object} filePaths - Paths to all JSON files
   * @returns {Object} Parsed and validated data
   */
  async loadAll(filePaths) {
    const errors = [];
    const data = {};

    // Load each file
    for (const [key, filePath] of Object.entries(filePaths)) {
      try {
        data[key] = await this.loadJsonFile(filePath);
        if (this.verbose) {
          console.log(`✓ Loaded ${key}.json`);
        }
      } catch (error) {
        errors.push({
          file: key,
          path: filePath,
          error: error.message
        });
      }
    }

    if (errors.length > 0) {
      throw new DataLoadError('Failed to load JSON files', errors);
    }

    // Validate schema
    this.validateSchema(data);

    return data;
  }

  /**
   * Load and parse a single JSON file
   * @param {string} filePath - Path to JSON file
   * @returns {Object} Parsed JSON data
   */
  async loadJsonFile(filePath) {
    // Check if file exists
    const exists = await fs.pathExists(filePath);
    if (!exists) {
      throw new Error(`File not found: ${filePath}`);
    }

    // Read file content
    let content;
    try {
      content = await fs.readFile(filePath, 'utf8');
    } catch (error) {
      throw new Error(`Cannot read file: ${error.message}`);
    }

    // Parse JSON
    let data;
    try {
      data = JSON.parse(content);
    } catch (error) {
      throw new Error(`Invalid JSON: ${error.message}`);
    }

    return data;
  }

  /**
   * Validate data against T.A.S.K.S schema
   * @param {Object} data - All loaded data
   */
  validateSchema(data) {
    const errors = [];

    // Validate tasks.json
    if (data.tasks) {
      if (!data.tasks.tasks || !Array.isArray(data.tasks.tasks)) {
        errors.push('tasks.json: Missing or invalid "tasks" array');
      }
      if (!data.tasks.dependencies || !Array.isArray(data.tasks.dependencies)) {
        errors.push('tasks.json: Missing or invalid "dependencies" array');
      }
      
      // Validate each task
      data.tasks.tasks?.forEach((task, index) => {
        if (!task.id) {
          errors.push(`tasks.json: Task at index ${index} missing "id"`);
        }
        if (!task.title) {
          errors.push(`tasks.json: Task ${task.id || index} missing "title"`);
        }
        if (!task.duration || typeof task.duration !== 'object') {
          errors.push(`tasks.json: Task ${task.id || index} missing or invalid "duration"`);
        }
      });

      // Validate dependencies
      data.tasks.dependencies?.forEach((dep, index) => {
        if (!dep.from) {
          errors.push(`tasks.json: Dependency at index ${index} missing "from"`);
        }
        if (!dep.to) {
          errors.push(`tasks.json: Dependency at index ${index} missing "to"`);
        }
        if (!dep.type) {
          errors.push(`tasks.json: Dependency ${dep.from}->${dep.to} missing "type"`);
        }
      });
    }

    // Validate dag.json
    if (data.dag) {
      if (typeof data.dag.ok !== 'boolean') {
        errors.push('dag.json: Missing or invalid "ok" field');
      }
      if (!data.dag.topo_order || !Array.isArray(data.dag.topo_order)) {
        errors.push('dag.json: Missing or invalid "topo_order" array');
      }
      if (!data.dag.metrics || typeof data.dag.metrics !== 'object') {
        errors.push('dag.json: Missing or invalid "metrics" object');
      }
    }

    // Validate features.json
    if (data.features) {
      if (!data.features.features || !Array.isArray(data.features.features)) {
        errors.push('features.json: Missing or invalid "features" array');
      }
      
      data.features.features?.forEach((feature, index) => {
        if (!feature.id) {
          errors.push(`features.json: Feature at index ${index} missing "id"`);
        }
        if (!feature.title) {
          errors.push(`features.json: Feature ${feature.id || index} missing "title"`);
        }
      });
    }

    // Validate waves.json
    if (data.waves) {
      if (!data.waves.waves || !Array.isArray(data.waves.waves)) {
        errors.push('waves.json: Missing or invalid "waves" array');
      }
      
      data.waves.waves?.forEach((wave, index) => {
        if (!wave.tasks || !Array.isArray(wave.tasks)) {
          errors.push(`waves.json: Wave ${wave.waveNumber || index} missing or invalid "tasks" array`);
        }
        if (typeof wave.waveNumber !== 'number') {
          errors.push(`waves.json: Wave at index ${index} missing "waveNumber"`);
        }
      });
    }

    if (errors.length > 0) {
      throw new SchemaValidationError('Schema validation failed', errors);
    }

    // Cross-validate references
    this.validateReferences(data);
  }

  /**
   * Validate cross-references between files
   * @param {Object} data - All loaded data
   */
  validateReferences(data) {
    const errors = [];
    const taskIds = new Set(data.tasks.tasks.map(t => t.id));
    const featureIds = new Set(data.features.features.map(f => f.id));

    // Check that all dependency references exist
    data.tasks.dependencies.forEach(dep => {
      if (!taskIds.has(dep.from)) {
        errors.push(`Invalid dependency: Task "${dep.from}" not found`);
      }
      if (!taskIds.has(dep.to)) {
        errors.push(`Invalid dependency: Task "${dep.to}" not found`);
      }
    });

    // Check that all task feature references exist
    data.tasks.tasks.forEach(task => {
      if (task.feature_id && !featureIds.has(task.feature_id)) {
        errors.push(`Task "${task.id}" references non-existent feature "${task.feature_id}"`);
      }
    });

    // Check that all wave task references exist
    data.waves.waves.forEach(wave => {
      wave.tasks.forEach(taskId => {
        if (!taskIds.has(taskId)) {
          errors.push(`Wave ${wave.waveNumber} references non-existent task "${taskId}"`);
        }
      });
    });

    if (errors.length > 0) {
      throw new ReferenceValidationError('Reference validation failed', errors);
    }
  }
}

// Custom error classes
export class DataLoadError extends Error {
  constructor(message, errors) {
    super(message);
    this.name = 'DataLoadError';
    this.errors = errors;
    this.code = 'FILE_NOT_FOUND';
  }
}

export class SchemaValidationError extends Error {
  constructor(message, errors) {
    super(message);
    this.name = 'SchemaValidationError';
    this.errors = errors;
    this.code = 'SCHEMA_VALIDATION_ERROR';
  }
}

export class ReferenceValidationError extends Error {
  constructor(message, errors) {
    super(message);
    this.name = 'ReferenceValidationError';
    this.errors = errors;
    this.code = 'SCHEMA_VALIDATION_ERROR';
  }
}