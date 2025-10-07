# quality_analyzer Tool - Complete Implementation Specification

> **Zero-Questions Implementation Guide for Centralized Quality Checks and Plan Health Analysis**

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
def quality_analyzer(
    plan: Dict[str, Any],
    analysis_config: Dict[str, Any],
    quality_gates: Dict[str, Any] = None
) -> QualityAnalysisResult:
    """
    Performs comprehensive quality analysis on T.A.S.K.S. project plans.
    
    CRITICAL PURPOSE: Prevents poor-quality plans from entering execution phase.
    Detects MECE violations, sizing issues, dependency problems, and structural
    defects that would lead to project failure.
    
    Args:
        plan: Complete T.A.S.K.S. project plan with tasks, waves, and dependencies
        analysis_config: Configuration for quality checks and thresholds
        quality_gates: Optional quality gate definitions for pass/fail decisions
        
    Returns:
        QualityAnalysisResult with quality scores, issues, and recommendations
    """
```

### 1.2 Quality Philosophy

**PLAN HEALTH FIRST**: A high-quality plan is the foundation of successful execution. Poor planning leads to cascade failures, missed deadlines, and scope creep.

**MECE ENFORCEMENT**: Plans must be Mutually Exclusive and Collectively Exhaustive. No overlapping work, no gaps in coverage.

**PREDICTIVE FAILURE DETECTION**: Identify structural problems that will cause execution failures before they happen.

---

## 2. Input/Output Schemas

### 2.1 Input Schema

```json
{
  "plan": {
    "project_info": {
      "id": "string",
      "name": "string", 
      "description": "string",
      "scope": "string",
      "constraints": ["string"]
    },
    "tasks": [
      {
        "id": "string",
        "name": "string",
        "description": "string",
        "wave": "integer",
        "dependencies": ["string"],
        "estimated_effort": "integer",
        "complexity": "low|medium|high",
        "acceptance_criteria": ["object"],
        "interfaces": ["object"]
      }
    ],
    "waves": [
      {
        "id": "integer",
        "name": "string",
        "tasks": ["string"],
        "estimated_duration": "integer"
      }
    ]
  },
  "analysis_config": {
    "checks_enabled": {
      "mece_analysis": true,
      "dependency_analysis": true,
      "sizing_analysis": true,
      "complexity_analysis": true,
      "coverage_analysis": true,
      "structural_analysis": true
    },
    "thresholds": {
      "max_wave_size": 10,
      "max_dependency_depth": 5,
      "max_task_complexity_variance": 2,
      "min_acceptance_criteria_coverage": 0.8,
      "max_interface_coupling": 0.6
    },
    "quality_standards": "strict|moderate|permissive"
  },
  "quality_gates": {
    "overall_quality_threshold": 0.8,
    "critical_issue_tolerance": 0,
    "blocker_issue_tolerance": 1,
    "warning_issue_tolerance": 5
  }
}
```

### 2.2 Output Schema

```json
{
  "overall_quality_score": 0.85,
  "quality_grade": "B+",
  "passes_quality_gates": true,
  "analysis_timestamp": "2024-01-15T10:30:00Z",
  "plan_statistics": {
    "total_tasks": 24,
    "total_waves": 4,
    "total_dependencies": 18,
    "avg_task_complexity": 2.3,
    "estimated_total_effort": 120,
    "parallelization_potential": 0.7
  },
  "quality_dimensions": {
    "mece_compliance": {
      "score": 0.9,
      "grade": "A-",
      "issues": ["minor_overlap_in_wave_2"],
      "recommendations": ["clarify_task_boundaries"]
    },
    "dependency_health": {
      "score": 0.8,
      "grade": "B+",
      "issues": ["circular_dependency_risk"],
      "recommendations": ["break_circular_dependency_in_wave_3"]
    },
    "sizing_consistency": {
      "score": 0.85,
      "grade": "B+",
      "issues": ["oversized_task_in_wave_1"],
      "recommendations": ["decompose_large_tasks"]
    },
    "complexity_distribution": {
      "score": 0.9,
      "grade": "A-",
      "issues": [],
      "recommendations": ["maintain_current_distribution"]
    },
    "coverage_completeness": {
      "score": 0.75,
      "grade": "C+",
      "issues": ["missing_error_handling_tasks"],
      "recommendations": ["add_error_handling_coverage"]
    },
    "structural_integrity": {
      "score": 0.88,
      "grade": "B+",
      "issues": ["potential_bottleneck_in_wave_2"],
      "recommendations": ["parallelize_bottleneck_tasks"]
    }
  },
  "critical_issues": [
    {
      "id": "circular_dependency_wave_3",
      "severity": "critical",
      "category": "dependency_violation",
      "description": "Circular dependency detected between Task-7 and Task-9",
      "impact": "Execution deadlock, project failure",
      "recommendation": "Break circular dependency by introducing intermediate task",
      "affected_tasks": ["task-7", "task-9"],
      "auto_fixable": false
    }
  ],
  "quality_improvements": [
    {
      "category": "mece_compliance",
      "priority": "high",
      "description": "Clarify boundaries between overlapping tasks",
      "estimated_effort": "2 hours",
      "impact_on_quality": "+0.05"
    },
    {
      "category": "coverage_completeness", 
      "priority": "medium",
      "description": "Add comprehensive error handling tasks",
      "estimated_effort": "4 hours",
      "impact_on_quality": "+0.1"
    }
  ],
  "execution_risks": [
    {
      "risk_type": "scheduling_risk",
      "probability": 0.3,
      "impact": "high",
      "description": "Wave 2 bottleneck may delay subsequent waves",
      "mitigation": "Parallelize Task-5 and Task-6"
    },
    {
      "risk_type": "complexity_risk",
      "probability": 0.2,
      "impact": "medium", 
      "description": "High complexity variance may lead to estimation errors",
      "mitigation": "Break down high-complexity tasks"
    }
  ],
  "plan_health_metrics": {
    "dependency_density": 0.75,
    "wave_balance": 0.8,
    "complexity_variance": 1.2,
    "coverage_gaps": 3,
    "bottleneck_count": 1,
    "parallelization_efficiency": 0.7
  }
}
```

---

## 3. Algorithm Implementation

### 3.1 Core Quality Analysis Algorithm

```python
def analyze_plan_quality(plan, analysis_config, quality_gates):
    """
    Core quality analysis algorithm with multi-dimensional scoring.
    """
    
    # Initialize analysis results
    quality_dimensions = {}
    critical_issues = []
    execution_risks = []
    
    # Phase 1: MECE Compliance Analysis
    if analysis_config["checks_enabled"]["mece_analysis"]:
        mece_result = analyze_mece_compliance(plan, analysis_config)
        quality_dimensions["mece_compliance"] = mece_result
        critical_issues.extend(mece_result.get("critical_issues", []))
    
    # Phase 2: Dependency Health Analysis
    if analysis_config["checks_enabled"]["dependency_analysis"]:
        dependency_result = analyze_dependency_health(plan, analysis_config)
        quality_dimensions["dependency_health"] = dependency_result
        critical_issues.extend(dependency_result.get("critical_issues", []))
        execution_risks.extend(dependency_result.get("execution_risks", []))
    
    # Phase 3: Task Sizing Analysis
    if analysis_config["checks_enabled"]["sizing_analysis"]:
        sizing_result = analyze_sizing_consistency(plan, analysis_config)
        quality_dimensions["sizing_consistency"] = sizing_result
        execution_risks.extend(sizing_result.get("execution_risks", []))
    
    # Phase 4: Complexity Distribution Analysis
    if analysis_config["checks_enabled"]["complexity_analysis"]:
        complexity_result = analyze_complexity_distribution(plan, analysis_config)
        quality_dimensions["complexity_distribution"] = complexity_result
        execution_risks.extend(complexity_result.get("execution_risks", []))
    
    # Phase 5: Coverage Completeness Analysis
    if analysis_config["checks_enabled"]["coverage_analysis"]:
        coverage_result = analyze_coverage_completeness(plan, analysis_config)
        quality_dimensions["coverage_completeness"] = coverage_result
        critical_issues.extend(coverage_result.get("critical_issues", []))
    
    # Phase 6: Structural Integrity Analysis
    if analysis_config["checks_enabled"]["structural_analysis"]:
        structural_result = analyze_structural_integrity(plan, analysis_config)
        quality_dimensions["structural_integrity"] = structural_result
        execution_risks.extend(structural_result.get("execution_risks", []))
    
    # Phase 7: Overall Quality Score Calculation
    overall_score = calculate_overall_quality_score(quality_dimensions)
    
    # Phase 8: Quality Gate Evaluation
    quality_gate_result = evaluate_quality_gates(
        overall_score, critical_issues, quality_gates
    )
    
    # Phase 9: Improvement Recommendations
    improvements = generate_quality_improvements(
        quality_dimensions, critical_issues, analysis_config
    )
    
    # Phase 10: Plan Health Metrics
    health_metrics = calculate_plan_health_metrics(plan, quality_dimensions)
    
    return build_quality_analysis_result(
        overall_score,
        quality_dimensions,
        critical_issues,
        execution_risks,
        improvements,
        health_metrics,
        quality_gate_result
    )
```

### 3.2 MECE Compliance Analysis

```python
def analyze_mece_compliance(plan, config):
    """
    Analyzes Mutually Exclusive and Collectively Exhaustive compliance.
    """
    
    tasks = plan["tasks"]
    project_scope = plan["project_info"]["scope"]
    
    # Mutual Exclusivity Check
    exclusivity_issues = []
    for i, task_a in enumerate(tasks):
        for j, task_b in enumerate(tasks[i+1:], i+1):
            overlap_score = calculate_task_overlap(task_a, task_b)
            if overlap_score > config["thresholds"]["max_task_overlap"]:
                exclusivity_issues.append({
                    "type": "task_overlap",
                    "severity": "high" if overlap_score > 0.7 else "medium",
                    "task_a": task_a["id"],
                    "task_b": task_b["id"],
                    "overlap_score": overlap_score,
                    "overlap_areas": identify_overlap_areas(task_a, task_b)
                })
    
    # Collective Exhaustiveness Check
    coverage_analysis = analyze_scope_coverage(tasks, project_scope)
    coverage_gaps = coverage_analysis["gaps"]
    coverage_score = coverage_analysis["coverage_percentage"]
    
    # Wave-Level MECE Analysis
    wave_mece_issues = []
    for wave in plan["waves"]:
        wave_tasks = [t for t in tasks if t["wave"] == wave["id"]]
        wave_mece = analyze_wave_mece_compliance(wave_tasks, wave)
        wave_mece_issues.extend(wave_mece["issues"])
    
    # Calculate MECE Score
    exclusivity_score = max(0, 1 - (len(exclusivity_issues) * 0.2))
    exhaustiveness_score = coverage_score
    wave_mece_score = calculate_wave_mece_score(wave_mece_issues)
    
    overall_mece_score = (
        exclusivity_score * 0.4 + 
        exhaustiveness_score * 0.4 + 
        wave_mece_score * 0.2
    )
    
    return {
        "score": overall_mece_score,
        "grade": score_to_grade(overall_mece_score),
        "exclusivity_score": exclusivity_score,
        "exhaustiveness_score": exhaustiveness_score,
        "wave_mece_score": wave_mece_score,
        "issues": exclusivity_issues + wave_mece_issues,
        "coverage_gaps": coverage_gaps,
        "recommendations": generate_mece_recommendations(
            exclusivity_issues, coverage_gaps, wave_mece_issues
        )
    }
```

### 3.3 Dependency Health Analysis

```python
def analyze_dependency_health(plan, config):
    """
    Analyzes dependency structure for cycles, depth, and complexity.
    """
    
    tasks = plan["tasks"]
    dependency_graph = build_dependency_graph(tasks)
    
    # Circular Dependency Detection
    circular_dependencies = detect_circular_dependencies(dependency_graph)
    
    # Dependency Depth Analysis
    depth_analysis = analyze_dependency_depth(dependency_graph)
    max_depth = depth_analysis["max_depth"]
    critical_path = depth_analysis["critical_path"]
    
    # Dependency Density Analysis
    density_analysis = analyze_dependency_density(dependency_graph, tasks)
    
    # Bottleneck Detection
    bottlenecks = detect_dependency_bottlenecks(dependency_graph)
    
    # Fan-out/Fan-in Analysis
    fanout_analysis = analyze_dependency_fanout(dependency_graph)
    
    # Calculate Dependency Health Score
    circular_penalty = len(circular_dependencies) * 0.3
    depth_score = max(0, 1 - (max_depth - config["thresholds"]["max_dependency_depth"]) * 0.1)
    density_score = max(0, 1 - abs(density_analysis["density"] - 0.5) * 2)
    bottleneck_penalty = len(bottlenecks) * 0.1
    
    dependency_score = max(0, depth_score + density_score - circular_penalty - bottleneck_penalty)
    
    # Identify Critical Issues
    critical_issues = []
    if circular_dependencies:
        for cycle in circular_dependencies:
            critical_issues.append({
                "id": f"circular_dependency_{cycle['id']}",
                "severity": "critical",
                "category": "dependency_violation",
                "description": f"Circular dependency detected: {' -> '.join(cycle['tasks'])}",
                "impact": "Execution deadlock, project failure",
                "affected_tasks": cycle["tasks"],
                "auto_fixable": False
            })
    
    # Identify Execution Risks
    execution_risks = []
    if bottlenecks:
        for bottleneck in bottlenecks:
            execution_risks.append({
                "risk_type": "scheduling_risk",
                "probability": 0.6,
                "impact": "high",
                "description": f"Task {bottleneck['task_id']} is a critical bottleneck",
                "mitigation": f"Parallelize dependencies or break down {bottleneck['task_id']}"
            })
    
    return {
        "score": dependency_score,
        "grade": score_to_grade(dependency_score),
        "circular_dependencies": circular_dependencies,
        "max_dependency_depth": max_depth,
        "critical_path": critical_path,
        "dependency_density": density_analysis["density"],
        "bottlenecks": bottlenecks,
        "critical_issues": critical_issues,
        "execution_risks": execution_risks,
        "recommendations": generate_dependency_recommendations(
            circular_dependencies, bottlenecks, depth_analysis
        )
    }
```

### 3.4 Coverage Completeness Analysis

```python
def analyze_coverage_completeness(plan, config):
    """
    Analyzes functional and non-functional coverage completeness.
    """
    
    tasks = plan["tasks"]
    project_scope = plan["project_info"]["scope"]
    
    # Functional Coverage Analysis
    functional_coverage = analyze_functional_coverage(tasks, project_scope)
    
    # Non-Functional Coverage Analysis
    nf_coverage = analyze_non_functional_coverage(tasks)
    
    # Error Handling Coverage
    error_coverage = analyze_error_handling_coverage(tasks)
    
    # Testing Coverage
    testing_coverage = analyze_testing_coverage(tasks)
    
    # Security Coverage
    security_coverage = analyze_security_coverage(tasks)
    
    # Performance Coverage
    performance_coverage = analyze_performance_coverage(tasks)
    
    # Integration Coverage
    integration_coverage = analyze_integration_coverage(tasks)
    
    # Calculate Overall Coverage Score
    coverage_weights = {
        "functional": 0.3,
        "error_handling": 0.2,
        "testing": 0.2,
        "security": 0.1,
        "performance": 0.1,
        "integration": 0.1
    }
    
    overall_coverage = (
        functional_coverage["score"] * coverage_weights["functional"] +
        error_coverage["score"] * coverage_weights["error_handling"] +
        testing_coverage["score"] * coverage_weights["testing"] +
        security_coverage["score"] * coverage_weights["security"] +
        performance_coverage["score"] * coverage_weights["performance"] +
        integration_coverage["score"] * coverage_weights["integration"]
    )
    
    # Identify Critical Coverage Gaps
    critical_gaps = []
    all_coverage_areas = [
        functional_coverage, error_coverage, testing_coverage,
        security_coverage, performance_coverage, integration_coverage
    ]
    
    for coverage_area in all_coverage_areas:
        if coverage_area["score"] < config["thresholds"]["min_coverage_threshold"]:
            critical_gaps.extend(coverage_area.get("critical_gaps", []))
    
    return {
        "score": overall_coverage,
        "grade": score_to_grade(overall_coverage),
        "functional_coverage": functional_coverage,
        "error_handling_coverage": error_coverage,
        "testing_coverage": testing_coverage,
        "security_coverage": security_coverage,
        "performance_coverage": performance_coverage,
        "integration_coverage": integration_coverage,
        "critical_gaps": critical_gaps,
        "recommendations": generate_coverage_recommendations(critical_gaps)
    }
```

---

## 4. Error Handling

### 4.1 Quality Analysis Error Categories

```python
class QualityAnalysisError(Exception):
    """Base exception for quality analysis errors."""
    pass

class InvalidPlanStructureError(QualityAnalysisError):
    """Raised when plan structure is invalid or incomplete."""
    def __init__(self, missing_fields, invalid_fields):
        self.missing_fields = missing_fields
        self.invalid_fields = invalid_fields
        super().__init__(f"Invalid plan structure: missing {missing_fields}, invalid {invalid_fields}")

class CircularDependencyError(QualityAnalysisError):
    """Raised when unresolvable circular dependencies are detected."""
    def __init__(self, cycles):
        self.cycles = cycles
        super().__init__(f"Circular dependencies detected: {cycles}")

class InsufficientCoverageError(QualityAnalysisError):
    """Raised when coverage falls below critical thresholds."""
    def __init__(self, coverage_score, threshold, missing_areas):
        self.coverage_score = coverage_score
        self.threshold = threshold
        self.missing_areas = missing_areas
        super().__init__(f"Coverage {coverage_score} below threshold {threshold}, missing: {missing_areas}")

class QualityGateFailureError(QualityAnalysisError):
    """Raised when plan fails to meet quality gate requirements."""
    def __init__(self, failed_gates, overall_score):
        self.failed_gates = failed_gates
        self.overall_score = overall_score
        super().__init__(f"Quality gates failed: {failed_gates}, score: {overall_score}")
```

### 4.2 Error Recovery and Analysis Continuation

```python
def handle_analysis_error(error, plan, config):
    """
    Handles analysis errors and attempts partial analysis recovery.
    """
    
    if isinstance(error, InvalidPlanStructureError):
        # Attempt to continue with partial plan analysis
        partial_plan = sanitize_plan_structure(plan, error.missing_fields)
        return {
            "analysis_status": "partial",
            "error": str(error),
            "partial_results": analyze_plan_quality(partial_plan, config, None),
            "limitations": "Analysis limited due to missing plan fields"
        }
    
    elif isinstance(error, CircularDependencyError):
        # Continue analysis but flag critical issue
        return {
            "analysis_status": "completed_with_critical_issues",
            "error": str(error),
            "critical_dependency_issues": error.cycles,
            "recommendation": "Fix circular dependencies before proceeding with execution"
        }
    
    elif isinstance(error, InsufficientCoverageError):
        # Complete analysis but flag coverage gaps
        return {
            "analysis_status": "completed_with_warnings",
            "coverage_warning": str(error),
            "missing_coverage": error.missing_areas,
            "recommendation": "Address coverage gaps to improve plan quality"
        }
    
    else:
        # Unknown error - fail safe
        return {
            "analysis_status": "failed",
            "error": str(error),
            "recommendation": "Review plan structure and retry analysis"
        }
```

---

## 5. Test Cases

### 5.1 High-Quality Plan Test Cases

```python
HIGH_QUALITY_PLAN_TESTS = [
    {
        "name": "well_structured_plan",
        "plan": {
            "project_info": {
                "id": "proj-1",
                "name": "API Development",
                "scope": "Complete REST API with authentication and data management"
            },
            "tasks": [
                {
                    "id": "task-1",
                    "name": "Setup project structure",
                    "wave": 1,
                    "dependencies": [],
                    "complexity": "low",
                    "estimated_effort": 4
                },
                {
                    "id": "task-2", 
                    "name": "Implement authentication",
                    "wave": 2,
                    "dependencies": ["task-1"],
                    "complexity": "medium",
                    "estimated_effort": 8
                },
                {
                    "id": "task-3",
                    "name": "Create data models",
                    "wave": 2,
                    "dependencies": ["task-1"],
                    "complexity": "medium", 
                    "estimated_effort": 6
                }
            ],
            "waves": [
                {"id": 1, "tasks": ["task-1"], "estimated_duration": 4},
                {"id": 2, "tasks": ["task-2", "task-3"], "estimated_duration": 8}
            ]
        },
        "expected_scores": {
            "overall_quality": 0.85,
            "mece_compliance": 0.9,
            "dependency_health": 0.9,
            "coverage_completeness": 0.8
        },
        "expected_issues": [],
        "expected_grade": "B+"
    }
]
```

### 5.2 Poor-Quality Plan Test Cases

```python
POOR_QUALITY_PLAN_TESTS = [
    {
        "name": "circular_dependency_plan",
        "plan": {
            "tasks": [
                {
                    "id": "task-1",
                    "dependencies": ["task-2"],
                    "wave": 1
                },
                {
                    "id": "task-2", 
                    "dependencies": ["task-3"],
                    "wave": 1
                },
                {
                    "id": "task-3",
                    "dependencies": ["task-1"],
                    "wave": 1
                }
            ]
        },
        "expected_critical_issues": [
            {
                "category": "dependency_violation",
                "severity": "critical",
                "type": "circular_dependency"
            }
        ],
        "expected_quality_score": 0.2,
        "expected_grade": "F"
    },
    {
        "name": "overlapping_tasks_plan",
        "plan": {
            "tasks": [
                {
                    "id": "task-1",
                    "name": "Implement user authentication",
                    "description": "Create login and registration functionality"
                },
                {
                    "id": "task-2",
                    "name": "Build login system", 
                    "description": "Develop user login capabilities"
                }
            ]
        },
        "expected_issues": [
            {
                "type": "task_overlap",
                "severity": "high",
                "overlap_score": 0.8
            }
        ],
        "expected_mece_score": 0.4
    },
    {
        "name": "insufficient_coverage_plan",
        "plan": {
            "project_info": {
                "scope": "Complete e-commerce platform with payments, inventory, and user management"
            },
            "tasks": [
                {
                    "id": "task-1",
                    "name": "Setup basic structure",
                    "description": "Initialize project"
                }
            ]
        },
        "expected_coverage_score": 0.1,
        "expected_coverage_gaps": [
            "payment_processing",
            "inventory_management", 
            "user_management",
            "error_handling",
            "security",
            "testing"
        ]
    }
]
```

### 5.3 Edge Case Test Cases

```python
EDGE_CASE_TESTS = [
    {
        "name": "empty_plan",
        "plan": {"tasks": [], "waves": []},
        "expected_error": "InvalidPlanStructureError",
        "expected_analysis_status": "failed"
    },
    {
        "name": "single_task_plan",
        "plan": {
            "tasks": [{"id": "task-1", "name": "Do everything", "wave": 1}],
            "waves": [{"id": 1, "tasks": ["task-1"]}]
        },
        "expected_issues": ["insufficient_decomposition"],
        "expected_quality_score": 0.3
    },
    {
        "name": "massive_wave_plan",
        "plan": {
            "tasks": [{"id": f"task-{i}", "wave": 1} for i in range(50)],
            "waves": [{"id": 1, "tasks": [f"task-{i}" for i in range(50)]}]
        },
        "expected_issues": ["oversized_wave"],
        "expected_structural_score": 0.2
    }
]
```

---

## 6. Performance Requirements

### 6.1 Analysis Time Constraints

```python
PERFORMANCE_REQUIREMENTS = {
    "analysis_time_per_task": {
        "target": "< 10ms",
        "maximum": "< 50ms",
        "measurement": "wall_clock_time"
    },
    "total_analysis_time": {
        "small_plan": "< 100ms (up to 10 tasks)",
        "medium_plan": "< 500ms (up to 50 tasks)", 
        "large_plan": "< 2s (up to 200 tasks)",
        "scaling": "O(n²) for dependency analysis, O(n) for other checks"
    },
    "quality_score_calculation": {
        "target": "< 50ms",
        "maximum": "< 200ms",
        "complexity": "O(n) where n = number of quality dimensions"
    }
}
```

### 6.2 Memory Usage Constraints

```python
MEMORY_REQUIREMENTS = {
    "max_memory_per_analysis": "50MB",
    "dependency_graph_memory": "O(n²) where n = number of tasks",
    "analysis_result_memory": "O(n) where n = number of issues found",
    "cleanup_strategy": "immediate_cleanup_after_analysis"
}
```

### 6.3 Scalability Requirements

```python
SCALABILITY_REQUIREMENTS = {
    "max_tasks_per_plan": 500,
    "max_waves_per_plan": 20,
    "max_dependencies_per_task": 10,
    "concurrent_analysis_limit": 5,
    "batch_processing": "disabled_for_accuracy"
}
```

---

## 7. Library Dependencies

### 7.1 Core Dependencies

```python
CORE_DEPENDENCIES = {
    "python": ">=3.8",
    "networkx": "^2.6.0",       # Graph analysis for dependencies
    "numpy": "^1.21.0",         # Statistical calculations
    "scipy": "^1.7.0",          # Advanced statistical analysis
    "pydantic": "^1.8.0",       # Data validation
    "jsonschema": "^4.0.0",     # Schema validation
    "click": "^8.0.0"           # CLI interface
}
```

### 7.2 Optional Dependencies

```python
OPTIONAL_DEPENDENCIES = {
    "matplotlib": "^3.4.0",     # Quality score visualization
    "plotly": "^5.3.0",         # Interactive quality reports
    "pandas": "^1.3.0",         # Data analysis and reporting
    "jinja2": "^3.0.0",         # HTML report generation
    "pygraphviz": "^1.7.0"      # Dependency graph visualization
}
```

### 7.3 Development Dependencies

```python
DEVELOPMENT_DEPENDENCIES = {
    "pytest": "^6.2.0",        # Testing framework
    "pytest-cov": "^2.12.0",   # Coverage analysis
    "black": "^21.6.0",        # Code formatting
    "mypy": "^0.910",          # Type checking
    "pre-commit": "^2.15.0"    # Git hooks
}
```

---

## 8. Integration Examples

### 8.1 T.A.S.K.S Pipeline Integration

```python
# Example integration in T.A.S.K.S pipeline
def validate_plan_quality(plan_data, quality_config):
    """
    Validates plan quality before execution phase.
    """
    
    # Configure quality analysis
    analysis_config = {
        "checks_enabled": {
            "mece_analysis": True,
            "dependency_analysis": True,
            "sizing_analysis": True,
            "complexity_analysis": True,
            "coverage_analysis": True,
            "structural_analysis": True
        },
        "thresholds": quality_config.get("thresholds", DEFAULT_THRESHOLDS),
        "quality_standards": quality_config.get("standards", "strict")
    }
    
    # Define quality gates
    quality_gates = {
        "overall_quality_threshold": 0.8,
        "critical_issue_tolerance": 0,
        "blocker_issue_tolerance": 1,
        "warning_issue_tolerance": 5
    }
    
    # Perform quality analysis
    analysis_result = quality_analyzer(
        plan=plan_data,
        analysis_config=analysis_config,
        quality_gates=quality_gates
    )
    
    # Check quality gates
    if not analysis_result["passes_quality_gates"]:
        raise QualityGateFailureError(
            f"Plan quality score {analysis_result['overall_quality_score']} below threshold",
            failed_gates=analysis_result["failed_quality_gates"],
            critical_issues=analysis_result["critical_issues"]
        )
    
    # Enhance plan with quality metadata
    plan_data["quality_analysis"] = analysis_result
    plan_data["quality_validated"] = True
    plan_data["quality_score"] = analysis_result["overall_quality_score"]
    
    return plan_data
```

### 8.2 CLI Usage Examples

```bash
# Analyze plan quality with default settings
quality_analyzer analyze \
    --plan-file project_plan.json \
    --output-format json \
    --output-file quality_report.json

# Analyze with custom quality gates
quality_analyzer analyze \
    --plan-file project_plan.json \
    --quality-gates-file custom_gates.json \
    --standards strict \
    --generate-report

# Generate detailed HTML quality report
quality_analyzer report \
    --plan-file project_plan.json \
    --report-format html \
    --include-visualizations \
    --output-file quality_report.html

# Check if plan passes quality gates (exit code based)
quality_analyzer validate \
    --plan-file project_plan.json \
    --quality-threshold 0.8 \
    --fail-on-critical-issues
```

---

## 9. Quality Improvement Recommendations

### 9.1 Automated Fix Suggestions

```python
AUTOMATED_IMPROVEMENTS = {
    "task_decomposition": {
        "trigger": "oversized_tasks",
        "action": "suggest_task_breakdown",
        "implementation": "analyze_task_complexity_and_suggest_splits"
    },
    "dependency_optimization": {
        "trigger": "excessive_dependencies",
        "action": "suggest_dependency_reduction",
        "implementation": "identify_unnecessary_dependencies"
    },
    "wave_rebalancing": {
        "trigger": "unbalanced_waves",
        "action": "suggest_task_redistribution", 
        "implementation": "optimize_wave_distribution"
    },
    "coverage_enhancement": {
        "trigger": "coverage_gaps",
        "action": "suggest_missing_tasks",
        "implementation": "generate_coverage_task_templates"
    }
}
```

### 9.2 Progressive Quality Improvement

```python
def generate_quality_improvement_plan(analysis_result, target_score):
    """
    Generates step-by-step plan to improve quality score.
    """
    
    current_score = analysis_result["overall_quality_score"]
    score_gap = target_score - current_score
    
    # Prioritize improvements by impact/effort ratio
    improvements = []
    
    for dimension, dimension_result in analysis_result["quality_dimensions"].items():
        if dimension_result["score"] < target_score:
            improvement_potential = calculate_improvement_potential(dimension_result)
            implementation_effort = estimate_improvement_effort(dimension_result)
            
            improvements.append({
                "dimension": dimension,
                "current_score": dimension_result["score"],
                "potential_improvement": improvement_potential,
                "estimated_effort": implementation_effort,
                "priority": improvement_potential / implementation_effort,
                "specific_actions": dimension_result["recommendations"]
            })
    
    # Sort by priority (impact/effort ratio)
    improvements.sort(key=lambda x: x["priority"], reverse=True)
    
    return {
        "target_score": target_score,
        "current_score": current_score,
        "required_improvement": score_gap,
        "improvement_plan": improvements,
        "estimated_total_effort": sum(imp["estimated_effort"] for imp in improvements),
        "expected_final_score": estimate_final_score(improvements)
    }
```

---

This comprehensive specification provides a complete foundation for implementing the `quality_analyzer` tool, ensuring T.A.S.K.S. v2 plans meet high quality standards before execution, preventing downstream failures and ensuring project success.
