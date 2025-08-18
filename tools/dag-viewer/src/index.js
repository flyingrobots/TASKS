import fs from 'fs-extra';
import path from 'path';
import { exec } from 'child_process';
import { promisify } from 'util';
import { DataLoader } from './dataLoader.js';
import { GraphBuilder } from './graphBuilder.js';
import { HTMLGenerator } from './htmlGenerator.js';
import { HTMLGeneratorModern } from './htmlGeneratorModern.js';
import HTMLGeneratorCytoscape from './htmlGeneratorCytoscape.js';

const execAsync = promisify(exec);

/**
 * Main orchestrator for DAG viewer tool
 * @param {Object} options - Configuration options
 * @returns {Object} JSON result indicating success or failure
 */
export async function orchestrate(options) {
  const startTime = Date.now();
  
  try {
    // Step 1: Load and validate data
    if (options.verbose) {
      console.log('Loading JSON files...');
    }
    
    const dataLoader = new DataLoader(options.verbose);
    const data = await dataLoader.loadAll(options.filePaths);
    
    // Step 2: Build graph structure
    if (options.verbose) {
      console.log('Building graph structure...');
    }
    
    const graphBuilder = new GraphBuilder(data, options.verbose);
    const graphData = graphBuilder.build();
    
    // Step 3: Generate HTML
    if (options.verbose) {
      console.log('Generating HTML visualization...');
    }
    
    // Choose generator based on options
    let GeneratorClass;
    if (options.renderer === 'cytoscape') {
      GeneratorClass = HTMLGeneratorCytoscape;
    } else if (options.classic) {
      GeneratorClass = HTMLGenerator;
    } else {
      GeneratorClass = HTMLGeneratorModern;
    }
    
    const htmlGenerator = new GeneratorClass(graphData, { verbose: options.verbose });
    const htmlContent = htmlGenerator.generate();
    
    // Step 4: Write output file
    if (options.verbose) {
      console.log(`Writing output to ${options.outputPath}...`);
    }
    
    await writeOutputFile(options.outputPath, htmlContent);
    
    // Step 5: Optionally open the file
    if (options.open && !process.env.CI) {
      await openFile(options.outputPath);
    }
    
    const endTime = Date.now();
    const duration = endTime - startTime;
    
    // Return success response
    return {
      success: true,
      htmlPath: options.outputPath,
      metadata: {
        nodeCount: graphData.nodes.length,
        edgeCount: graphData.edges.length,
        featureCount: graphData.features.length,
        waveCount: graphData.waves.length,
        generatedAt: new Date().toISOString(),
        toolVersion: '1.0.0',
        duration: `${duration}ms`
      }
    };
    
  } catch (error) {
    // Handle known error types
    if (error.name === 'DataLoadError') {
      return {
        success: false,
        error: {
          code: 'FILE_NOT_FOUND',
          message: error.message,
          details: {
            errors: error.errors
          }
        }
      };
    }
    
    if (error.name === 'SchemaValidationError' || error.name === 'ReferenceValidationError') {
      return {
        success: false,
        error: {
          code: 'SCHEMA_VALIDATION_ERROR',
          message: error.message,
          details: {
            errors: error.errors
          }
        }
      };
    }
    
    if (error.code === 'EACCES' || error.code === 'EPERM') {
      return {
        success: false,
        error: {
          code: 'WRITE_PERMISSION_ERROR',
          message: `Cannot write to output file: ${options.outputPath}`,
          details: {
            outputPath: options.outputPath,
            error: error.message
          }
        }
      };
    }
    
    // Handle unexpected errors
    return {
      success: false,
      error: {
        code: 'INTERNAL_ERROR',
        message: error.message,
        details: {
          stack: options.verbose ? error.stack : undefined
        }
      }
    };
  }
}

/**
 * Write HTML content to file
 * @param {string} outputPath - Path to output file
 * @param {string} content - HTML content
 */
async function writeOutputFile(outputPath, content) {
  // Ensure directory exists
  const dir = path.dirname(outputPath);
  await fs.ensureDir(dir);
  
  // Write file
  await fs.writeFile(outputPath, content, 'utf8');
}

/**
 * Open HTML file in default browser
 * @param {string} filePath - Path to HTML file
 */
async function openFile(filePath) {
  const platform = process.platform;
  let command;
  
  if (platform === 'darwin') {
    command = `open "${filePath}"`;
  } else if (platform === 'win32') {
    command = `start "" "${filePath}"`;
  } else {
    command = `xdg-open "${filePath}"`;
  }
  
  try {
    await execAsync(command);
  } catch (error) {
    // Ignore errors from opening file
    console.warn('Could not open file automatically:', error.message);
  }
}