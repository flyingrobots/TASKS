#!/usr/bin/env node

import { program } from 'commander';
import path from 'path';
import { fileURLToPath } from 'url';
import { dirname } from 'path';
import { readFileSync } from 'fs';
import { orchestrate } from './index.js';
import chalk from 'chalk';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);
const packageJson = JSON.parse(readFileSync(path.join(__dirname, '..', 'package.json'), 'utf8'));

program
  .name('dag-viewer')
  .description('Generate interactive DAG visualization from T.A.S.K.S JSON files')
  .version(packageJson.version)
  .option('-d, --dag <path>', 'Path to dag.json file')
  .option('-f, --features <path>', 'Path to features.json file')
  .option('-t, --tasks <path>', 'Path to tasks.json file')
  .option('-w, --waves <path>', 'Path to waves.json file')
  .option('-D, --dir <path>', 'Directory containing all JSON files')
  .option('-o, --output <path>', 'Output HTML file path', './dag-visualization.html')
  .option('-v, --verbose', 'Enable verbose output', false)
  .option('--no-open', 'Don\'t open the HTML file after generation')
  .option('-r, --renderer <type>', 'Renderer to use: cytoscape (canvas), modern (d3), or classic', 'cytoscape')
  .option('--classic', 'Use classic HTML generator (deprecated, use --renderer classic instead)')
  .parse();

const options = program.opts();

async function main() {
  try {
    // Validate input options
    if (!options.dir && (!options.dag || !options.features || !options.tasks || !options.waves)) {
      console.error(chalk.red('Error: You must provide either --dir or all individual file paths (--dag, --features, --tasks, --waves)'));
      process.exit(1);
    }

    // Resolve file paths
    let filePaths;
    if (options.dir) {
      const dir = path.resolve(options.dir);
      filePaths = {
        dag: path.join(dir, 'dag.json'),
        features: path.join(dir, 'features.json'),
        tasks: path.join(dir, 'tasks.json'),
        waves: path.join(dir, 'waves.json')
      };
    } else {
      filePaths = {
        dag: path.resolve(options.dag),
        features: path.resolve(options.features),
        tasks: path.resolve(options.tasks),
        waves: path.resolve(options.waves)
      };
    }

    // Resolve output path
    const outputPath = path.resolve(options.output);

    if (options.verbose) {
      console.log(chalk.blue('Input files:'));
      console.log(`  DAG: ${filePaths.dag}`);
      console.log(`  Features: ${filePaths.features}`);
      console.log(`  Tasks: ${filePaths.tasks}`);
      console.log(`  Waves: ${filePaths.waves}`);
      console.log(`  Output: ${outputPath}`);
    }

    // Run the orchestrator
    const result = await orchestrate({
      filePaths,
      outputPath,
      verbose: options.verbose,
      open: options.open,
      renderer: options.renderer,
      classic: options.classic
    });

    // Output JSON result
    console.log(JSON.stringify(result, null, 2));

    // Exit with appropriate code
    process.exit(result.success ? 0 : 1);

  } catch (error) {
    // Handle unexpected errors
    const errorResult = {
      success: false,
      error: {
        code: 'INTERNAL_ERROR',
        message: error.message,
        details: {
          stack: options.verbose ? error.stack : undefined
        }
      }
    };
    
    console.log(JSON.stringify(errorResult, null, 2));
    process.exit(1);
  }
}

main();