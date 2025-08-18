# DAG Viewer - Technical Specifications

## CLI Interface Specification

### Command Syntax
```bash
dag-viewer [options]
```

### Options
| Option | Alias | Type | Required | Default | Description |
|--------|-------|------|----------|---------|-------------|
| `--dag` | `-d` | string | Yes | - | Path to dag.json file |
| `--features` | `-f` | string | Yes | - | Path to features.json file |
| `--tasks` | `-t` | string | Yes | - | Path to tasks.json file |
| `--waves` | `-w` | string | Yes | - | Path to waves.json file |
| `--output` | `-o` | string | No | `./dag-visualization.html` | Output HTML file path |
| `--dir` | `-D` | string | No | - | Directory containing all JSON files (alternative to individual paths) |
| `--verbose` | `-v` | boolean | No | false | Enable verbose output |
| `--no-open` | - | boolean | No | false | Don't open the HTML file after generation |

### Alternative Usage
```bash
# Specify individual files
dag-viewer --dag ./dag.json --features ./features.json --tasks ./tasks.json --waves ./waves.json

# Specify directory containing all files
dag-viewer --dir ./examples --output ./output/visualization.html
```

## JSON Response Format

### Success Response
```json
{
  "success": true,
  "htmlPath": "/absolute/path/to/generated/file.html",
  "metadata": {
    "nodeCount": 23,
    "edgeCount": 19,
    "generatedAt": "2025-01-20T10:30:00Z",
    "toolVersion": "1.0.0"
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

### Error Codes
| Code | Description | Recoverable |
|------|-------------|-------------|
| `FILE_NOT_FOUND` | One or more input files not found | Yes - provide correct paths |
| `INVALID_JSON` | JSON parsing failed | Yes - fix JSON syntax |
| `SCHEMA_VALIDATION_ERROR` | JSON doesn't match T.A.S.K.S schema | Yes - fix data structure |
| `WRITE_PERMISSION_ERROR` | Cannot write output file | Yes - change output location |
| `CYCLIC_DEPENDENCY` | Graph contains cycles | No - fix in source data |
| `INTERNAL_ERROR` | Unexpected error | No - report bug |

## Data Schema Validation

### Required Fields Validation

#### tasks.json
```typescript
interface TasksSchema {
  meta: {
    min_confidence: number;
    notes?: string;
  };
  generated: {
    by: string;
    timestamp: string;
    contentHash: string;
  };
  tasks: Array<{
    id: string;
    feature_id: string;
    title: string;
    description: string;
    category: 'foundation' | 'implementation' | 'integration' | 'optimization';
    duration: {
      optimistic: number;
      mostLikely: number;
      pessimistic: number;
    };
    durationUnits: string;
  }>;
  dependencies: Array<{
    from: string;
    to: string;
    type: 'technical' | 'sequential' | 'infrastructure' | 'knowledge';
    confidence: number;
    isHard: boolean;
  }>;
}
```

#### dag.json
```typescript
interface DagSchema {
  ok: boolean;
  errors: string[];
  warnings: string[];
  metrics: {
    nodes: number;
    edges: number;
    edgeDensity: number;
    widthApprox: number;
    longestPath: number;
  };
  topo_order: string[];
  reduced_edges_sample: Array<[string, string]>;
  softDeps: Array<{
    from: string;
    to: string;
    type: string;
    confidence: number;
    isHard: boolean;
  }>;
}
```

## Graph Construction Specification

### Node Properties
```javascript
{
  id: "P1.T001",                    // Unique task ID
  label: "Create admin config",     // Task title (truncated)
  fullTitle: "Full task title",     // Complete title
  description: "Task description",  // Task description
  feature: "F001",                  // Associated feature ID
  category: "foundation",           // Task category
  duration: 3.0,                    // Most likely duration
  wave: 1,                          // Wave number
  metadata: {                       // Additional properties
    priority: "high",
    complexity: 7,
    skills: ["backend", "devops"]
  }
}
```

### Edge Properties
```javascript
{
  from: "P1.T001",           // Source node ID
  to: "P1.T003",            // Target node ID
  type: "technical",        // Dependency type
  isHard: true,            // Hard vs soft dependency
  confidence: 0.95,        // Confidence level
  label: "requires",       // Edge label (optional)
  style: "solid",          // Visual style
  weight: 1                // Layout weight
}
```

### Layout Configuration
```javascript
{
  rankdir: "TB",           // Top to Bottom layout
  nodesep: 50,            // Horizontal spacing between nodes
  ranksep: 80,            // Vertical spacing between ranks
  marginx: 20,            // Horizontal margin
  marginy: 20,            // Vertical margin
  acyclicer: "greedy",    // Algorithm for removing cycles
  ranker: "longest-path"  // Node ranking algorithm
}
```

## HTML Output Specification

### Structure
```html
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>T.A.S.K.S DAG Visualization</title>
  <style>/* Embedded CSS */</style>
  <script src="https://cdn.jsdelivr.net/npm/d3@7"></script>
  <script src="https://cdn.jsdelivr.net/npm/dagre@0.8.5"></script>
  <script src="https://cdn.jsdelivr.net/npm/dagre-d3@0.6.4"></script>
</head>
<body>
  <div id="controls"><!-- Interactive controls --></div>
  <svg id="graph"><!-- DAG visualization --></svg>
  <div id="details"><!-- Task details panel --></div>
  <script>
    // Embedded graph data
    const graphData = {/* Generated from input */};
    // Visualization code
  </script>
</body>
</html>
```

### Interactive Features
1. **Zoom and Pan**: Mouse wheel zoom, click and drag to pan
2. **Node Hover**: Show tooltip with basic info
3. **Node Click**: Display detailed information panel
4. **Edge Hover**: Highlight dependency path
5. **Search**: Find and highlight specific tasks
6. **Filter**: Show/hide by feature, category, or wave
7. **Legend**: Color and style explanations
8. **Export**: Save as SVG or PNG (optional)

### Visual Encoding
| Element | Property | Encoding |
|---------|----------|----------|
| Node Color | Feature | Categorical color scale |
| Node Size | Duration | Linear scale (min: 20px, max: 60px) |
| Node Shape | Category | Rectangle (foundation), Rounded (implementation), Diamond (integration) |
| Edge Style | Dependency Type | Solid (hard), Dashed (soft) |
| Edge Color | Confidence | Gradient from red (low) to green (high) |
| Edge Width | Importance | 1-3px based on confidence |

## Performance Requirements

### Benchmarks
| Metric | Target | Maximum |
|--------|--------|---------|
| File parsing | < 100ms | 500ms |
| Graph construction | < 200ms | 1000ms |
| Layout calculation | < 500ms | 2000ms |
| HTML generation | < 100ms | 500ms |
| Total execution | < 1s | 3s |

### Scalability Limits
| Resource | Soft Limit | Hard Limit |
|----------|------------|------------|
| Nodes | 100 | 1000 |
| Edges | 500 | 5000 |
| File size (each) | 1MB | 10MB |
| Output HTML size | 2MB | 20MB |

## Error Handling Specification

### Validation Levels
1. **File Level**: Check existence and readability
2. **JSON Level**: Validate syntax and structure
3. **Schema Level**: Ensure T.A.S.K.S compliance
4. **Graph Level**: Check for cycles and consistency
5. **Output Level**: Verify write permissions

### Error Recovery Strategies
| Error Type | Recovery Strategy |
|------------|------------------|
| Missing optional field | Use default value |
| Unknown task reference | Skip dependency, warn |
| Duplicate task ID | Use first occurrence, warn |
| Invalid duration | Use default (4 hours) |
| Cycle detected | Report cycle path, fail |

## Testing Requirements

### Unit Tests
- CLI argument parsing
- JSON file reading
- Schema validation
- Graph construction
- HTML generation

### Integration Tests
- End-to-end with example files
- Error scenarios
- Large graph handling
- Cross-platform compatibility

### Acceptance Criteria
- [ ] Generates valid HTML from example files
- [ ] Returns correct JSON response format
- [ ] Handles all error cases gracefully
- [ ] Visualization is interactive and readable
- [ ] Performance meets benchmarks
- [ ] Works on Node.js 16+