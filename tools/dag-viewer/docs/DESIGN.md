# DAG Viewer - Design Overview

## Executive Summary

The DAG Viewer is a command-line tool designed for LLM invocation that transforms T.A.S.K.S specification JSON files into an interactive, visual representation of the task dependency graph. It generates a standalone HTML file with an embedded DAG visualization using the dagre layout algorithm and D3.js for rendering.

## System Architecture

### High-Level Architecture

```
┌─────────────┐    ┌──────────────┐    ┌────────────────┐    ┌──────────────┐    ┌──────────┐
│  CLI Input  │───►│ Data Loader  │───►│ Graph Builder  │───►│ HTML Generator│───►│  Output  │
│  (paths)    │    │  (JSON)      │    │  (dagre)       │    │  (template)   │    │  (HTML)  │
└─────────────┘    └──────────────┘    └────────────────┘    └──────────────┘    └──────────┘
                           │                     │                      │
                           ▼                     ▼                      ▼
                   ┌──────────────┐    ┌────────────────┐    ┌──────────────┐
                   │  Validator   │    │  Transformer   │    │  Renderer    │
                   └──────────────┘    └────────────────┘    └──────────────┘
```

### Component Responsibilities

#### 1. CLI Interface (`cli.js`)
- Parse command-line arguments
- Validate input file paths
- Handle output path specification
- Return JSON response for LLM consumption

#### 2. Data Layer (`dataLoader.js`)
- Read JSON files from filesystem
- Parse and validate JSON structure
- Validate against T.A.S.K.S schema
- Aggregate data from multiple files

#### 3. Transform Layer (`graphBuilder.js`)
- Convert T.A.S.K.S data model to dagre graph format
- Map tasks to nodes with metadata
- Map dependencies to edges with types
- Calculate layout coordinates

#### 4. Visualization Layer (`htmlGenerator.js`)
- Generate HTML template
- Embed graph data
- Include D3.js and dagre-d3 libraries
- Add interactive features and styling

#### 5. Orchestrator (`index.js`)
- Coordinate pipeline execution
- Handle errors gracefully
- Generate JSON response
- Write output files

## Data Flow

### Input Processing
1. CLI receives file paths or directory path
2. Validator checks file existence and readability
3. DataLoader reads and parses JSON files
4. Schema validator ensures T.A.S.K.S compliance

### Graph Construction
1. Extract nodes from tasks.json
2. Enrich nodes with feature information
3. Extract edges from dependencies
4. Add soft dependencies from dag.json
5. Apply wave information for grouping

### Visualization Generation
1. Initialize dagre graph
2. Add nodes with labels and metadata
3. Add edges with styles and weights
4. Calculate automatic layout
5. Generate SVG representation
6. Embed in HTML template

### Output Production
1. Write HTML file to specified location
2. Generate JSON response with success/failure
3. Return response to stdout for LLM

## Key Design Patterns

### Pipeline Pattern
The tool follows a linear pipeline architecture where data flows through distinct transformation stages. Each stage has a single responsibility and produces output for the next stage.

### Builder Pattern
The GraphBuilder uses the builder pattern to construct the graph incrementally, adding nodes, edges, and metadata in a fluent interface.

### Template Pattern
The HTMLGenerator uses a template pattern with placeholders for dynamic content injection.

### Error Handler Pattern
Centralized error handling with structured JSON responses ensures consistent error reporting for LLM consumption.

## Technology Stack Rationale

### Node.js
- Native JavaScript environment for dagre library
- Excellent JSON handling
- Fast file I/O operations
- Wide ecosystem support

### dagre
- Purpose-built for DAG layout
- Automatic positioning algorithm
- Handles complex graph structures
- Well-maintained and documented

### D3.js + dagre-d3
- Industry standard for data visualization
- SVG-based rendering for scalability
- Rich interaction capabilities
- Extensive customization options

### Commander.js
- Declarative CLI interface
- Automatic help generation
- Type coercion and validation
- Industry standard for Node CLI tools

## Scalability Considerations

### Performance
- Lazy loading of large graphs
- Progressive rendering for better UX
- Efficient graph algorithms (O(V + E) complexity)

### Maintainability
- Modular architecture
- Clear separation of concerns
- Comprehensive error handling
- Extensive documentation

### Extensibility
- Plugin architecture for custom visualizations
- Configurable styling and themes
- Support for additional metadata
- Export to multiple formats (future)

## Security Considerations

- Input validation to prevent injection attacks
- No execution of user-provided code
- Sandboxed HTML output
- No external network requests in generated HTML
- File system access limited to specified paths