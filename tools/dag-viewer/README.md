# DAG Viewer for T.A.S.K.S

A command-line tool that generates interactive DAG (Directed Acyclic Graph) visualizations from T.A.S.K.S specification JSON files. Designed for LLM invocation with structured JSON output.

## Overview

The DAG Viewer transforms T.A.S.K.S project planning data (tasks, features, dependencies, and waves) into an interactive HTML visualization using the dagre layout algorithm and D3.js for rendering. The output is a standalone HTML file that can be viewed in any modern web browser without requiring a server.

## Installation

```bash
cd tools/dag-viewer
npm install
```

## Usage

### Basic Usage

Using individual file paths:
```bash
dag-viewer --dag ./dag.json --features ./features.json --tasks ./tasks.json --waves ./waves.json
```

Using a directory containing all JSON files:
```bash
dag-viewer --dir ./examples
```

Specifying output location:
```bash
dag-viewer --dir ./examples --output ./visualizations/my-dag.html
```

### Command-Line Options

| Option | Alias | Description | Default |
|--------|-------|-------------|---------|
| `--dag` | `-d` | Path to dag.json file | Required* |
| `--features` | `-f` | Path to features.json file | Required* |
| `--tasks` | `-t` | Path to tasks.json file | Required* |
| `--waves` | `-w` | Path to waves.json file | Required* |
| `--dir` | `-D` | Directory containing all JSON files | Alternative to individual paths |
| `--output` | `-o` | Output HTML file path | `./dag-visualization.html` |
| `--verbose` | `-v` | Enable verbose output | `false` |
| `--no-open` | | Don't open HTML file after generation | `false` |

*Either provide all individual file paths OR use `--dir` to specify a directory

### Examples

Quick visualization of example files:
```bash
npm run example
```

This runs:
```bash
node src/cli.js --dir ../../examples --output ./output/example.html
```

LLM-friendly invocation with JSON output:
```bash
dag-viewer --dir /path/to/tasks/files 2>/dev/null
```

## Output Format

### Success Response
```json
{
  "success": true,
  "htmlPath": "/absolute/path/to/generated/visualization.html",
  "metadata": {
    "nodeCount": 23,
    "edgeCount": 19,
    "featureCount": 3,
    "waveCount": 14,
    "generatedAt": "2025-01-20T10:30:00Z",
    "toolVersion": "1.0.0",
    "duration": "245ms"
  }
}
```

### Error Response
```json
{
  "success": false,
  "error": {
    "code": "FILE_NOT_FOUND",
    "message": "Cannot find file: dag.json",
    "details": {
      "missingFile": "./dag.json",
      "searchPath": "/current/working/directory"
    }
  }
}
```

## Visualization Features

### Interactive Elements
- **Zoom and Pan**: Use mouse wheel to zoom, click and drag to pan
- **Node Selection**: Click on any task node to view detailed information
- **Hover Tooltips**: Hover over nodes for quick information
- **Search**: Find tasks by ID or title
- **Filters**: Filter by feature or wave
- **Reset View**: Return to initial zoom and position

### Visual Encoding
- **Node Colors**: Different colors for each feature
- **Node Shapes**: Rectangle (foundation), Rounded (implementation), Diamond (integration)
- **Edge Styles**: Solid lines for hard dependencies, dashed for soft dependencies
- **Edge Colors**: Green (high confidence), Orange (medium), Red (low)

### Information Panel
The right panel displays detailed information for selected tasks:
- Task ID and full title
- Description
- Associated feature
- Category and duration
- Wave assignment
- Required skills
- Acceptance checks

## Error Codes

| Code | Description | Solution |
|------|-------------|----------|
| `FILE_NOT_FOUND` | Input file not found | Check file paths |
| `INVALID_JSON` | JSON parsing failed | Fix JSON syntax |
| `SCHEMA_VALIDATION_ERROR` | Invalid T.A.S.K.S schema | Ensure data matches specification |
| `WRITE_PERMISSION_ERROR` | Cannot write output | Check directory permissions |
| `INTERNAL_ERROR` | Unexpected error | Check verbose output for details |

## Development

### Project Structure
```
dag-viewer/
├── src/
│   ├── cli.js           # CLI argument parser
│   ├── index.js         # Main orchestrator
│   ├── dataLoader.js    # JSON file loading and validation
│   ├── graphBuilder.js  # DAG construction logic
│   └── htmlGenerator.js # HTML visualization generator
├── docs/
│   ├── DESIGN.md               # Architecture overview
│   ├── TECHNICAL_SPEC.md       # Technical specifications
│   └── DECISION_RATIONALE.md   # Design decisions
├── output/              # Generated HTML files
├── package.json
└── README.md
```

### Running Tests
```bash
npm test
```

### Linting
```bash
npm run lint
```

## Requirements

- Node.js 16.0.0 or higher
- Modern web browser for viewing output

## Design Philosophy

This tool follows these principles:
1. **LLM-Friendly**: Structured JSON I/O for easy integration
2. **Zero Dependencies**: Generated HTML files are self-contained
3. **Stateless**: Each invocation is independent
4. **Minimal Configuration**: Works out of the box with sensible defaults
5. **Error Resilience**: Graceful error handling with actionable messages

## License

MIT

## Contributing

See the documentation in the `docs/` directory for detailed information about the architecture and design decisions.