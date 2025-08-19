# task_validator Tool - Complete Implementation Specification

> **Zero-Questions Implementation Guide for Graph Validation and DAG Construction**

## Table of Contents

1. [API Specification](#1-api-specification)
2. [Input/Output Schemas](#2-inputoutput-schemas)
3. [Algorithm Implementation](#3-algorithm-implementation)
4. [Error Handling](#4-error-handling)
5. [Test Cases](#5-test-cases)
6. [Performance Requirements](#6-performance-requirements)
7. [Library Dependencies](#7-library-dependencies)

---

## 1. API Specification

### 1.1 Function Signature

```python
def validate_and_build_dag(
    tasks_json_data: dict,
    min_confidence: float = 0.7,
    validation_config: dict = None
) -> dict:
    """
    Validates task dependencies and builds a DAG with comprehensive metrics.
    
    Args:
        tasks_json_data: Complete tasks.json object (see schema below)
        min_confidence: Minimum confidence threshold for hard dependencies [0.0-1.0]
        validation_config: Optional configuration overrides
        
    Returns:
        Complete dag.json object with validation results
        
    Raises:
        ValidationError: Input schema validation failure
        CycleDetectedError: Dependency cycles found (includes cycle path)
        ConfigurationError: Invalid configuration parameters
        ProcessingError: Unexpected processing failure
    """
```

### 1.2 Command Line Interface

```bash
python -m task_validator \
    --input tasks.json \
    --output dag.json \
    --min-confidence 0.7 \
    --config validation_config.json \
    --verbose
```

### 1.3 Validation Configuration Schema

```json
{
  "validation_config": {
    "min_confidence": 0.7,
    "max_nodes": 10000,
    "max_edges": 50000,
    "enable_transitive_reduction": true,
    "enable_quality_checks": true,
    "quality_thresholds": {
      "max_edge_density": 0.1,
      "max_isolated_tasks": 10,
      "min_verb_first_percentage": 0.8,
      "max_mece_overlap_threshold": 0.8
    },
    "performance_limits": {
      "max_processing_time_seconds": 300,
      "max_memory_mb": 1024
    }
  }
}
```

---

## 2. Input/Output Schemas

### 2.1 Input Schema: tasks.json

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["meta", "tasks"],
  "properties": {
    "meta": {
      "type": "object",
      "required": ["min_confidence"],
      "properties": {
        "min_confidence": {
          "type": "number",
          "minimum": 0.0,
          "maximum": 1.0
        },
        "notes": {"type": "string"},
        "autonormalization": {
          "type": "object",
          "properties": {
            "split": {"type": "array", "items": {"type": "string"}},
            "merged": {"type": "array", "items": {"type": "string"}}
          }
        }
      }
    },
    "tasks": {
      "type": "array",
      "minItems": 1,
      "maxItems": 10000,
      "items": {
        "type": "object",
        "required": ["id", "title", "duration", "dependencies"],
        "properties": {
          "id": {
            "type": "string",
            "pattern": "^P\\d+\\.T\\d+$",
            "description": "Task ID format: P1.T001"
          },
          "feature_id": {
            "type": "string",
            "pattern": "^F\\d+$",
            "description": "Feature ID format: F001"
          },
          "title": {
            "type": "string",
            "minLength": 5,
            "maxLength": 200
          },
          "duration": {
            "type": "object",
            "required": ["optimistic", "mostLikely", "pessimistic"],
            "properties": {
              "optimistic": {"type": "number", "minimum": 0.1},
              "mostLikely": {"type": "number", "minimum": 0.1},
              "pessimistic": {"type": "number", "minimum": 0.1}
            }
          },
          "dependencies": {
            "type": "array",
            "items": {
              "type": "object",
              "required": ["from", "to", "type", "confidence", "isHard"],
              "properties": {
                "from": {"type": "string", "pattern": "^P\\d+\\.T\\d+$"},
                "to": {"type": "string", "pattern": "^P\\d+\\.T\\d+$"},
                "type": {
                  "type": "string",
                  "enum": ["technical", "sequential", "infrastructure", "knowledge"]
                },
                "confidence": {"type": "number", "minimum": 0.0, "maximum": 1.0},
                "isHard": {"type": "boolean"},
                "evidence": {
                  "type": "object",
                  "properties": {
                    "quote": {"type": "string"},
                    "reasoning": {"type": "string"}
                  }
                }
              }
            }
          },
          "interfaces_produced": {
            "type": "array",
            "items": {"type": "string", "pattern": "^[A-Za-z][A-Za-z0-9]*:v\\d+$"}
          },
          "interfaces_consumed": {
            "type": "array",
            "items": {"type": "string", "pattern": "^[A-Za-z][A-Za-z0-9]*:v\\d+$"}
          }
        }
      }
    }
  }
}
```

### 2.2 Output Schema: dag.json

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["ok", "generated", "metrics"],
  "properties": {
    "ok": {"type": "boolean"},
    "errors": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["code", "message"],
        "properties": {
          "code": {
            "type": "string",
            "enum": ["CYCLE_DETECTED", "VALIDATION_ERROR", "PROCESSING_ERROR"]
          },
          "message": {"type": "string"},
          "details": {"type": "object"},
          "cycle_path": {
            "type": "array",
            "items": {"type": "string"}
          },
          "suggested_fixes": {
            "type": "array",
            "items": {"type": "string"}
          }
        }
      }
    },
    "warnings": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "code": {"type": "string"},
          "message": {"type": "string"},
          "affected_tasks": {"type": "array", "items": {"type": "string"}}
        }
      }
    },
    "generated": {
      "type": "object",
      "required": ["by", "timestamp", "contentHash"],
      "properties": {
        "by": {"type": "string"},
        "timestamp": {"type": "string", "format": "date-time"},
        "contentHash": {"type": "string", "pattern": "^[a-f0-9]{64}$"}
      }
    },
    "metrics": {
      "type": "object",
      "required": ["nodes", "edges", "edgeDensity", "longestPath"],
      "properties": {
        "minConfidenceApplied": {"type": "number"},
        "keptByType": {
          "type": "object",
          "properties": {
            "technical": {"type": "integer"},
            "sequential": {"type": "integer"},
            "infrastructure": {"type": "integer"},
            "knowledge": {"type": "integer"}
          }
        },
        "droppedByType": {
          "type": "object",
          "properties": {
            "technical": {"type": "integer"},
            "sequential": {"type": "integer"},
            "infrastructure": {"type": "integer"},
            "knowledge": {"type": "integer"}
          }
        },
        "nodes": {"type": "integer", "minimum": 0},
        "edges": {"type": "integer", "minimum": 0},
        "edgeDensity": {"type": "number", "minimum": 0.0, "maximum": 1.0},
        "widthApprox": {"type": "integer", "minimum": 1},
        "widthMethod": {"type": "string", "enum": ["kahn_layer_max", "antichain_max"]},
        "longestPath": {"type": "integer", "minimum": 0},
        "isolatedTasks": {"type": "integer", "minimum": 0},
        "lowConfidenceDeps": {"type": "integer", "minimum": 0},
        "softDeps": {"type": "integer", "minimum": 0},
        "verbFirstPercentage": {"type": "number", "minimum": 0.0, "maximum": 1.0},
        "meceOverlapSuspects": {"type": "integer", "minimum": 0}
      }
    },
    "topo_order": {
      "type": "array",
      "items": {"type": "string"}
    },
    "reduced_edges": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["from", "to", "type"],
        "properties": {
          "from": {"type": "string"},
          "to": {"type": "string"},
          "type": {"type": "string"},
          "confidence": {"type": "number"}
        }
      }
    },
    "lowConfidenceDeps": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "from": {"type": "string"},
          "to": {"type": "string"},
          "confidence": {"type": "number"},
          "reason": {"type": "string"}
        }
      }
    },
    "softDeps": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "from": {"type": "string"},
          "to": {"type": "string"},
          "type": {"type": "string"},
          "reasoning": {"type": "string"}
        }
      }
    }
  }
}
```

---

## 3. Algorithm Implementation

### 3.1 Core Algorithm Flow

```python
def validate_and_build_dag(tasks_json_data, min_confidence=0.7, validation_config=None):
    """Complete implementation with all error handling."""
    
    # Step 1: Input Validation
    config = _merge_config(validation_config)
    _validate_input_schema(tasks_json_data)
    _validate_config(config)
    
    # Step 2: Extract and Filter Dependencies
    all_dependencies = _extract_dependencies(tasks_json_data)
    hard_dependencies = _filter_hard_dependencies(all_dependencies, min_confidence)
    
    # Step 3: Build Graph
    graph = _build_adjacency_list(hard_dependencies)
    
    # Step 4: Cycle Detection (CRITICAL)
    cycles = _detect_cycles_dfs(graph)
    if cycles:
        return _generate_cycle_error_response(cycles, hard_dependencies)
    
    # Step 5: Topological Sort
    topo_order = _topological_sort_kahn(graph)
    
    # Step 6: Transitive Reduction
    if config['enable_transitive_reduction']:
        reduced_edges = _transitive_reduction_gries_herman(graph)
    else:
        reduced_edges = hard_dependencies
    
    # Step 7: Compute Metrics
    metrics = _compute_comprehensive_metrics(
        graph, reduced_edges, all_dependencies, tasks_json_data, config
    )
    
    # Step 8: Generate Content Hash
    content_hash = _generate_content_hash(reduced_edges, metrics)
    
    # Step 9: Build Response
    return _build_success_response(
        topo_order, reduced_edges, metrics, content_hash, 
        all_dependencies, min_confidence
    )
```

### 3.2 Cycle Detection Algorithm (DFS + Stack)

```python
def _detect_cycles_dfs(graph):
    """
    Detect cycles using DFS with recursion stack.
    Returns list of cycles, each cycle is a list of node IDs.
    """
    WHITE, GRAY, BLACK = 0, 1, 2
    color = {node: WHITE for node in graph}
    parent = {node: None for node in graph}
    cycles = []
    
    def dfs_visit(node, path):
        """DFS visit with cycle detection."""
        if color[node] == GRAY:
            # Back edge found - cycle detected
            cycle_start = path.index(node)
            cycle = path[cycle_start:] + [node]
            cycles.append(cycle)
            return
        
        if color[node] == BLACK:
            return
        
        color[node] = GRAY
        path.append(node)
        
        for neighbor in graph.get(node, []):
            dfs_visit(neighbor, path)
        
        path.pop()
        color[node] = BLACK
    
    # Check all nodes to handle disconnected components
    for node in graph:
        if color[node] == WHITE:
            dfs_visit(node, [])
    
    return cycles

def _generate_cycle_error_response(cycles, dependencies):
    """Generate detailed cycle error with fix suggestions."""
    cycle_errors = []
    
    for cycle in cycles:
        # Generate human-readable cycle path
        cycle_path = " → ".join(cycle)
        
        # Suggest fixes based on cycle analysis
        fixes = _suggest_cycle_fixes(cycle, dependencies)
        
        cycle_errors.append({
            "code": "CYCLE_DETECTED",
            "message": f"Dependency cycle detected: {cycle_path}",
            "cycle_path": cycle,
            "suggested_fixes": fixes
        })
    
    return {
        "ok": False,
        "errors": cycle_errors,
        "warnings": [],
        "generated": _generate_metadata(),
        "metrics": {},
        "cycle_break_suggestions": _generate_cycle_break_suggestions(cycles)
    }

def _suggest_cycle_fixes(cycle, dependencies):
    """Suggest specific fixes for detected cycles."""
    fixes = []
    
    # Analyze cycle for common patterns
    if len(cycle) == 3:  # Simple 2-node cycle
        fixes.append(f"Remove dependency between {cycle[0]} and {cycle[1]}")
        fixes.append(f"Split {cycle[0]} into interface and implementation phases")
    
    # Look for tasks that could be split
    for task in cycle[:-1]:  # Exclude duplicate at end
        if _is_splittable_task(task, dependencies):
            fixes.append(f"Split {task} into sequential subtasks")
    
    # Suggest interface insertion
    if len(cycle) > 3:
        middle_idx = len(cycle) // 2
        fixes.append(f"Insert interface task between {cycle[middle_idx-1]} and {cycle[middle_idx]}")
    
    return fixes

def _is_splittable_task(task_id, dependencies):
    """Determine if a task can be reasonably split."""
    # Tasks with many dependencies are good candidates for splitting
    in_degree = sum(1 for dep in dependencies if dep['to'] == task_id)
    out_degree = sum(1 for dep in dependencies if dep['from'] == task_id)
    return (in_degree + out_degree) > 2
```

### 3.3 Topological Sort (Kahn's Algorithm)

```python
def _topological_sort_kahn(graph):
    """
    Kahn's algorithm for topological sorting.
    Returns list of nodes in topological order.
    """
    # Calculate in-degrees
    in_degree = {node: 0 for node in graph}
    for node in graph:
        for neighbor in graph[node]:
            in_degree[neighbor] = in_degree.get(neighbor, 0) + 1
    
    # Initialize queue with nodes having in-degree 0
    from collections import deque
    queue = deque([node for node, degree in in_degree.items() if degree == 0])
    topo_order = []
    
    while queue:
        node = queue.popleft()
        topo_order.append(node)
        
        # Remove edges from current node
        for neighbor in graph.get(node, []):
            in_degree[neighbor] -= 1
            if in_degree[neighbor] == 0:
                queue.append(neighbor)
    
    # Verify all nodes were processed (should not happen if no cycles)
    if len(topo_order) != len(in_degree):
        raise ProcessingError("Topological sort failed - possible undetected cycle")
    
    return topo_order
```

### 3.4 Transitive Reduction (Gries-Herman Algorithm)

```python
def _transitive_reduction_gries_herman(graph):
    """
    Gries-Herman algorithm for transitive reduction.
    Returns minimal set of edges preserving reachability.
    """
    nodes = list(graph.keys())
    n = len(nodes)
    node_to_idx = {node: i for i, node in enumerate(nodes)}
    
    # Build adjacency matrix
    adj_matrix = [[False] * n for _ in range(n)]
    for node in graph:
        node_idx = node_to_idx[node]
        for neighbor in graph[node]:
            neighbor_idx = node_to_idx[neighbor]
            adj_matrix[node_idx][neighbor_idx] = True
    
    # Compute transitive closure
    reachable = [row[:] for row in adj_matrix]  # Copy matrix
    
    for k in range(n):
        for i in range(n):
            for j in range(n):
                reachable[i][j] = reachable[i][j] or (reachable[i][k] and reachable[k][j])
    
    # Remove transitive edges
    reduced_matrix = [row[:] for row in adj_matrix]  # Copy original
    
    for i in range(n):
        for j in range(n):
            if adj_matrix[i][j]:  # If edge exists
                # Check if there's an intermediate path
                for k in range(n):
                    if k != i and k != j and reachable[i][k] and reachable[k][j]:
                        reduced_matrix[i][j] = False  # Remove transitive edge
                        break
    
    # Convert back to edge list
    reduced_edges = []
    for i in range(n):
        for j in range(n):
            if reduced_matrix[i][j]:
                reduced_edges.append({
                    "from": nodes[i],
                    "to": nodes[j],
                    "type": _get_edge_type(nodes[i], nodes[j], graph),
                    "confidence": _get_edge_confidence(nodes[i], nodes[j], graph)
                })
    
    return reduced_edges

def _get_edge_type(from_node, to_node, original_dependencies):
    """Get original edge type from dependency list."""
    for dep in original_dependencies:
        if dep['from'] == from_node and dep['to'] == to_node:
            return dep['type']
    return "unknown"

def _get_edge_confidence(from_node, to_node, original_dependencies):
    """Get original edge confidence from dependency list."""
    for dep in original_dependencies:
        if dep['from'] == from_node and dep['to'] == to_node:
            return dep['confidence']
    return 1.0
```

### 3.5 Comprehensive Metrics Computation

```python
def _compute_comprehensive_metrics(graph, reduced_edges, all_dependencies, tasks_data, config):
    """Compute all DAG quality metrics."""
    
    nodes = set()
    for edge in reduced_edges:
        nodes.add(edge['from'])
        nodes.add(edge['to'])
    
    # Add isolated nodes
    for task in tasks_data['tasks']:
        nodes.add(task['id'])
    
    metrics = {
        "minConfidenceApplied": config.get('min_confidence', 0.7),
        "nodes": len(nodes),
        "edges": len(reduced_edges),
        "edgeDensity": _calculate_edge_density(len(nodes), len(reduced_edges)),
        "widthApprox": _calculate_width_approximation(graph),
        "widthMethod": "kahn_layer_max",
        "longestPath": _calculate_longest_path(graph, tasks_data),
        "isolatedTasks": _count_isolated_tasks(graph, nodes),
        "lowConfidenceDeps": _count_low_confidence_deps(all_dependencies, config['min_confidence']),
        "softDeps": _count_soft_deps(all_dependencies),
        "verbFirstPercentage": _calculate_verb_first_percentage(tasks_data),
        "meceOverlapSuspects": _count_mece_overlap_suspects(tasks_data),
        "keptByType": _count_kept_dependencies_by_type(reduced_edges),
        "droppedByType": _count_dropped_dependencies_by_type(all_dependencies, reduced_edges)
    }
    
    return metrics

def _calculate_edge_density(num_nodes, num_edges):
    """Calculate edge density: edges / (nodes * (nodes - 1))."""
    if num_nodes <= 1:
        return 0.0
    max_edges = num_nodes * (num_nodes - 1)
    return num_edges / max_edges

def _calculate_width_approximation(graph):
    """Calculate maximum width using Kahn's layering."""
    # Use Kahn's algorithm to generate layers
    in_degree = {node: 0 for node in graph}
    for node in graph:
        for neighbor in graph[node]:
            in_degree[neighbor] = in_degree.get(neighbor, 0) + 1
    
    from collections import deque
    current_layer = deque([node for node, degree in in_degree.items() if degree == 0])
    max_width = 0
    
    while current_layer:
        layer_size = len(current_layer)
        max_width = max(max_width, layer_size)
        
        next_layer = deque()
        for _ in range(layer_size):
            node = current_layer.popleft()
            for neighbor in graph.get(node, []):
                in_degree[neighbor] -= 1
                if in_degree[neighbor] == 0:
                    next_layer.append(neighbor)
        
        current_layer = next_layer
    
    return max_width

def _calculate_longest_path(graph, tasks_data):
    """Calculate longest path using task durations."""
    # Get task durations
    task_durations = {}
    for task in tasks_data['tasks']:
        duration = task['duration']
        # Use PERT expected value: (a + 4m + b) / 6
        expected = (duration['optimistic'] + 4 * duration['mostLikely'] + duration['pessimistic']) / 6
        task_durations[task['id']] = expected
    
    # Use dynamic programming to find longest path
    memo = {}
    
    def longest_path_from(node):
        if node in memo:
            return memo[node]
        
        if node not in graph or not graph[node]:
            # Leaf node
            memo[node] = task_durations.get(node, 0)
            return memo[node]
        
        max_path = 0
        for neighbor in graph[node]:
            path_length = longest_path_from(neighbor)
            max_path = max(max_path, path_length)
        
        memo[node] = task_durations.get(node, 0) + max_path
        return memo[node]
    
    # Find maximum among all possible starting nodes
    max_path_length = 0
    for node in graph:
        path_length = longest_path_from(node)
        max_path_length = max(max_path_length, path_length)
    
    return max_path_length

def _calculate_verb_first_percentage(tasks_data):
    """Calculate percentage of tasks with verb-first titles."""
    verb_patterns = [
        r'^(add|create|build|implement|develop|design|configure|setup|install)',
        r'^(update|modify|change|refactor|improve|enhance|optimize)',
        r'^(remove|delete|cleanup|fix|resolve|debug|troubleshoot)',
        r'^(test|validate|verify|check|ensure|confirm)',
        r'^(deploy|release|publish|migrate|integrate|connect)'
    ]
    
    import re
    verb_first_count = 0
    total_tasks = len(tasks_data['tasks'])
    
    for task in tasks_data['tasks']:
        title_lower = task['title'].lower()
        for pattern in verb_patterns:
            if re.match(pattern, title_lower):
                verb_first_count += 1
                break
    
    return verb_first_count / total_tasks if total_tasks > 0 else 0.0

def _count_mece_overlap_suspects(tasks_data):
    """Count potential MECE (Mutually Exclusive, Collectively Exhaustive) overlaps."""
    # Group tasks by feature
    features = {}
    for task in tasks_data['tasks']:
        feature_id = task.get('feature_id', 'unknown')
        if feature_id not in features:
            features[feature_id] = []
        features[feature_id].append(task)
    
    overlap_suspects = 0
    overlap_threshold = 0.8  # 80% token overlap threshold
    
    for feature_id, tasks in features.items():
        if len(tasks) < 2:
            continue
        
        # Compare each pair of tasks in the feature
        for i in range(len(tasks)):
            for j in range(i + 1, len(tasks)):
                similarity = _calculate_title_similarity(tasks[i]['title'], tasks[j]['title'])
                if similarity > overlap_threshold:
                    overlap_suspects += 1
    
    return overlap_suspects

def _calculate_title_similarity(title1, title2):
    """Calculate token overlap similarity between two task titles."""
    import re
    
    # Tokenize and normalize
    tokens1 = set(re.findall(r'\b\w+\b', title1.lower()))
    tokens2 = set(re.findall(r'\b\w+\b', title2.lower()))
    
    # Remove common stop words
    stop_words = {'a', 'an', 'the', 'and', 'or', 'but', 'in', 'on', 'at', 'to', 'for', 'of', 'with', 'by'}
    tokens1 -= stop_words
    tokens2 -= stop_words
    
    if not tokens1 or not tokens2:
        return 0.0
    
    # Calculate Jaccard similarity
    intersection = len(tokens1 & tokens2)
    union = len(tokens1 | tokens2)
    
    return intersection / union if union > 0 else 0.0
```

---

## 4. Error Handling

### 4.1 Custom Exception Classes

```python
class TaskValidatorError(Exception):
    """Base exception for task validator errors."""
    pass

class ValidationError(TaskValidatorError):
    """Raised when input validation fails."""
    def __init__(self, message, details=None):
        super().__init__(message)
        self.details = details or {}

class CycleDetectedError(TaskValidatorError):
    """Raised when dependency cycles are detected."""
    def __init__(self, cycles, suggestions=None):
        self.cycles = cycles
        self.suggestions = suggestions or []
        super().__init__(f"Dependency cycles detected: {len(cycles)} cycle(s)")

class ConfigurationError(TaskValidatorError):
    """Raised when configuration is invalid."""
    pass

class ProcessingError(TaskValidatorError):
    """Raised when unexpected processing errors occur."""
    pass
```

### 4.2 Input Validation Functions

```python
def _validate_input_schema(tasks_json_data):
    """Validate input against JSON schema."""
    import jsonschema
    
    try:
        jsonschema.validate(tasks_json_data, TASKS_JSON_SCHEMA)
    except jsonschema.ValidationError as e:
        raise ValidationError(f"Input schema validation failed: {e.message}", {
            "path": list(e.absolute_path),
            "invalid_value": e.instance,
            "constraint": e.validator
        })

def _validate_config(config):
    """Validate configuration parameters."""
    if not isinstance(config['min_confidence'], (int, float)):
        raise ConfigurationError("min_confidence must be a number")
    
    if not 0.0 <= config['min_confidence'] <= 1.0:
        raise ConfigurationError("min_confidence must be between 0.0 and 1.0")
    
    if config['max_nodes'] <= 0:
        raise ConfigurationError("max_nodes must be positive")
    
    if config['max_edges'] <= 0:
        raise ConfigurationError("max_edges must be positive")

def _validate_task_references(tasks_json_data):
    """Validate that all dependency references point to existing tasks."""
    task_ids = {task['id'] for task in tasks_json_data['tasks']}
    errors = []
    
    for task in tasks_json_data['tasks']:
        for dep in task.get('dependencies', []):
            if dep['from'] not in task_ids:
                errors.append(f"Dependency 'from' task {dep['from']} does not exist")
            if dep['to'] not in task_ids:
                errors.append(f"Dependency 'to' task {dep['to']} does not exist")
    
    if errors:
        raise ValidationError("Invalid task references in dependencies", {"errors": errors})
```

### 4.3 Error Response Generation

```python
def _generate_error_response(error_type, message, details=None, cycles=None):
    """Generate standardized error response."""
    response = {
        "ok": False,
        "errors": [{
            "code": error_type,
            "message": message,
            "details": details or {}
        }],
        "warnings": [],
        "generated": _generate_metadata(),
        "metrics": {}
    }
    
    if cycles:
        response["errors"][0]["cycle_path"] = cycles[0] if cycles else []
        response["cycle_break_suggestions"] = _generate_cycle_break_suggestions(cycles)
    
    return response

def _generate_metadata():
    """Generate standard metadata block."""
    from datetime import datetime
    return {
        "by": "task_validator v2.0",
        "timestamp": datetime.utcnow().isoformat() + "Z",
        "contentHash": ""  # Will be filled later
    }
```

---

## 5. Test Cases

### 5.1 Valid DAG Test Case

```python
def test_valid_dag():
    """Test with a valid DAG structure."""
    input_data = {
        "meta": {"min_confidence": 0.7},
        "tasks": [
            {
                "id": "P1.T001",
                "title": "Create database schema",
                "duration": {"optimistic": 2, "mostLikely": 4, "pessimistic": 6},
                "dependencies": []
            },
            {
                "id": "P1.T002", 
                "title": "Implement API endpoints",
                "duration": {"optimistic": 4, "mostLikely": 8, "pessimistic": 12},
                "dependencies": [
                    {
                        "from": "P1.T001",
                        "to": "P1.T002",
                        "type": "technical",
                        "confidence": 0.9,
                        "isHard": true
                    }
                ]
            },
            {
                "id": "P1.T003",
                "title": "Create user interface",
                "duration": {"optimistic": 3, "mostLikely": 6, "pessimistic": 9},
                "dependencies": [
                    {
                        "from": "P1.T002",
                        "to": "P1.T003", 
                        "type": "technical",
                        "confidence": 0.8,
                        "isHard": true
                    }
                ]
            }
        ]
    }
    
    result = validate_and_build_dag(input_data)
    
    assert result["ok"] == True
    assert result["metrics"]["nodes"] == 3
    assert result["metrics"]["edges"] == 2
    assert len(result["topo_order"]) == 3
    assert result["topo_order"][0] == "P1.T001"  # Should be first
```

### 5.2 Cycle Detection Test Case

```python
def test_cycle_detection():
    """Test cycle detection with a simple cycle."""
    input_data = {
        "meta": {"min_confidence": 0.7},
        "tasks": [
            {
                "id": "P1.T001",
                "title": "Task A",
                "duration": {"optimistic": 1, "mostLikely": 2, "pessimistic": 3},
                "dependencies": [
                    {
                        "from": "P1.T002",
                        "to": "P1.T001",
                        "type": "technical", 
                        "confidence": 0.9,
                        "isHard": true
                    }
                ]
            },
            {
                "id": "P1.T002",
                "title": "Task B", 
                "duration": {"optimistic": 1, "mostLikely": 2, "pessimistic": 3},
                "dependencies": [
                    {
                        "from": "P1.T001",
                        "to": "P1.T002",
                        "type": "sequential",
                        "confidence": 0.8,
                        "isHard": true
                    }
                ]
            }
        ]
    }
    
    result = validate_and_build_dag(input_data)
    
    assert result["ok"] == False
    assert len(result["errors"]) == 1
    assert result["errors"][0]["code"] == "CYCLE_DETECTED"
    assert "P1.T001" in result["errors"][0]["cycle_path"]
    assert "P1.T002" in result["errors"][0]["cycle_path"]
    assert len(result["errors"][0]["suggested_fixes"]) > 0
```

### 5.3 Low Confidence Dependencies Test

```python
def test_low_confidence_filtering():
    """Test filtering of low confidence dependencies."""
    input_data = {
        "meta": {"min_confidence": 0.7},
        "tasks": [
            {
                "id": "P1.T001",
                "title": "Task A",
                "duration": {"optimistic": 1, "mostLikely": 2, "pessimistic": 3},
                "dependencies": []
            },
            {
                "id": "P1.T002",
                "title": "Task B",
                "duration": {"optimistic": 1, "mostLikely": 2, "pessimistic": 3},
                "dependencies": [
                    {
                        "from": "P1.T001",
                        "to": "P1.T002",
                        "type": "technical",
                        "confidence": 0.5,  # Below threshold
                        "isHard": true
                    }
                ]
            }
        ]
    }
    
    result = validate_and_build_dag(input_data, min_confidence=0.7)
    
    assert result["ok"] == True
    assert result["metrics"]["edges"] == 0  # Should be filtered out
    assert result["metrics"]["lowConfidenceDeps"] == 1
    assert len(result["lowConfidenceDeps"]) == 1
```

### 5.4 Performance Test Case

```python
def test_large_dag_performance():
    """Test performance with large DAG (1000 nodes)."""
    import time
    
    # Generate large test case
    tasks = []
    for i in range(1000):
        task = {
            "id": f"P1.T{i:03d}",
            "title": f"Task {i}",
            "duration": {"optimistic": 1, "mostLikely": 2, "pessimistic": 3},
            "dependencies": []
        }
        
        # Add dependencies to previous tasks (creates a chain)
        if i > 0:
            task["dependencies"].append({
                "from": f"P1.T{i-1:03d}",
                "to": f"P1.T{i:03d}",
                "type": "sequential",
                "confidence": 0.9,
                "isHard": True
            })
        
        tasks.append(task)
    
    input_data = {
        "meta": {"min_confidence": 0.7},
        "tasks": tasks
    }
    
    start_time = time.time()
    result = validate_and_build_dag(input_data)
    end_time = time.time()
    
    assert result["ok"] == True
    assert result["metrics"]["nodes"] == 1000
    assert (end_time - start_time) < 5.0  # Should complete in under 5 seconds
```

---

## 6. Performance Requirements

### 6.1 Complexity Requirements

| Operation | Required Complexity | Maximum Runtime (1000 tasks) |
|-----------|-------------------|-------------------------------|
| Input Validation | O(n) | 50ms |
| Cycle Detection | O(V + E) | 100ms |
| Topological Sort | O(V + E) | 50ms |
| Transitive Reduction | O(V³) worst case | 500ms |
| Metrics Computation | O(V + E) | 100ms |
| Total Pipeline | O(V³) worst case | 800ms |

### 6.2 Memory Requirements

| Component | Memory Complexity | Maximum Memory (1000 tasks) |
|-----------|------------------|------------------------------|
| Adjacency List | O(V + E) | 100KB |
| Transitive Reduction Matrix | O(V²) | 1MB |
| Topological Sort Queue | O(V) | 10KB |
| Total Memory Usage | O(V²) | 2MB |

### 6.3 Scalability Limits

```python
DEFAULT_LIMITS = {
    "max_nodes": 10000,
    "max_edges": 50000,
    "max_processing_time_seconds": 300,
    "max_memory_mb": 1024
}

def _check_performance_limits(tasks_json_data, config):
    """Check if input exceeds performance limits."""
    num_tasks = len(tasks_json_data['tasks'])
    num_deps = sum(len(task.get('dependencies', [])) for task in tasks_json_data['tasks'])
    
    if num_tasks > config['max_nodes']:
        raise ValidationError(f"Too many tasks: {num_tasks} > {config['max_nodes']}")
    
    if num_deps > config['max_edges']:
        raise ValidationError(f"Too many dependencies: {num_deps} > {config['max_edges']}")
```

---

## 7. Library Dependencies

### 7.1 Required Python Packages

```txt
# Core dependencies
jsonschema>=4.0.0          # JSON schema validation
networkx>=2.8.0            # Graph algorithms (alternative implementation)
numpy>=1.21.0              # Numerical computations
hashlib                    # Content hashing (built-in)
collections                # Data structures (built-in)
datetime                   # Timestamp generation (built-in)
re                         # Regular expressions (built-in)

# Development dependencies
pytest>=7.0.0              # Testing framework
pytest-cov>=4.0.0          # Coverage reporting
black>=22.0.0              # Code formatting
mypy>=0.991                # Type checking
```

### 7.2 Alternative Graph Library Implementation

```python
# Option 1: Pure Python (recommended for T.A.S.K.S.)
def _build_adjacency_list_pure(dependencies):
    """Pure Python adjacency list implementation."""
    graph = {}
    for dep in dependencies:
        from_node = dep['from']
        to_node = dep['to']
        
        if from_node not in graph:
            graph[from_node] = []
        graph[from_node].append(to_node)
        
        # Ensure all nodes exist in graph
        if to_node not in graph:
            graph[to_node] = []
    
    return graph

# Option 2: NetworkX (for complex analysis)
def _build_networkx_graph(dependencies):
    """NetworkX implementation for advanced graph operations."""
    import networkx as nx
    
    G = nx.DiGraph()
    for dep in dependencies:
        G.add_edge(dep['from'], dep['to'], 
                  type=dep['type'], 
                  confidence=dep['confidence'])
    
    return G

def _detect_cycles_networkx(G):
    """Use NetworkX for cycle detection."""
    import networkx as nx
    try:
        cycles = list(nx.simple_cycles(G))
        return cycles
    except nx.NetworkXNoCycle:
        return []
```

### 7.3 Content Hashing Implementation

```python
import hashlib
import json

def _generate_content_hash(reduced_edges, metrics):
    """Generate SHA256 content hash for deterministic output."""
    
    # Create normalized content for hashing
    hash_content = {
        "edges": sorted(reduced_edges, key=lambda x: (x['from'], x['to'])),
        "metrics": {k: v for k, v in sorted(metrics.items()) if k != 'timestamp'}
    }
    
    # Serialize to JSON with consistent formatting
    json_str = json.dumps(hash_content, sort_keys=True, separators=(',', ':'))
    
    # Generate SHA256 hash
    hash_obj = hashlib.sha256(json_str.encode('utf-8'))
    return hash_obj.hexdigest()
```

### 7.4 Configuration Management

```python
def _merge_config(user_config):
    """Merge user configuration with defaults."""
    default_config = {
        "min_confidence": 0.7,
        "max_nodes": 10000,
        "max_edges": 50000,
        "enable_transitive_reduction": True,
        "enable_quality_checks": True,
        "quality_thresholds": {
            "max_edge_density": 0.1,
            "max_isolated_tasks": 10,
            "min_verb_first_percentage": 0.8,
            "max_mece_overlap_threshold": 0.8
        },
        "performance_limits": {
            "max_processing_time_seconds": 300,
            "max_memory_mb": 1024
        }
    }
    
    if user_config:
        # Deep merge configuration
        import copy
        config = copy.deepcopy(default_config)
        _deep_update(config, user_config)
        return config
    
    return default_config

def _deep_update(base_dict, update_dict):
    """Deep update dictionary with nested values."""
    for key, value in update_dict.items():
        if isinstance(value, dict) and key in base_dict and isinstance(base_dict[key], dict):
            _deep_update(base_dict[key], value)
        else:
            base_dict[key] = value
```

---

## Summary

This specification provides complete implementation details for the `task_validator` tool with:

- **Complete API specification** with type hints and error handling
- **Full JSON schemas** for input validation and output structure
- **Detailed algorithm implementations** for all core operations
- **Comprehensive test cases** covering normal and edge cases
- **Performance requirements** and scalability limits
- **Library dependencies** and alternative implementations

The implementation is designed to be **zero-questions buildable** by any competent Python developer, with clear specifications for all edge cases, error conditions, and performance requirements.

**Key Implementation Notes**:
1. Use DFS with recursion stack for cycle detection (most reliable)
2. Implement Gries-Herman algorithm for transitive reduction (optimal for sparse graphs)
3. Include comprehensive quality metrics for enterprise usage
4. Provide detailed error messages with fix suggestions for human intervention
5. Ensure deterministic output through content hashing

This tool specification serves as the foundation for the T.A.S.K.S. v2 validation pipeline.
