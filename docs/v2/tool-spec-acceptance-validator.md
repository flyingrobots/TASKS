# acceptance_validator Tool - Complete Implementation Specification

> **Zero-Questions Implementation Guide for Acceptance Criteria Validation and Machine-Checkable Test Generation**

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
def acceptance_validator(
    criteria: List[Dict[str, Any]],
    context: Dict[str, Any],
    validation_level: str = "strict"
) -> AcceptanceValidationResult:
    """
    Validates acceptance criteria are machine-checkable and executable.
    
    CRITICAL PURPOSE: Prevents unrunnable acceptance criteria that would make 
    task completion verification impossible. Ensures every criterion can be 
    automatically validated through deterministic checks.
    
    Args:
        criteria: List of acceptance criteria objects
        context: Task and project context for validation
        validation_level: "strict", "moderate", or "permissive"
        
    Returns:
        AcceptanceValidationResult with validation status and executable test commands
    """
```

### 1.2 Validation Philosophy

**EXECUTION GUARANTEE**: Every acceptance criterion must be reducible to executable commands that return clear pass/fail results. No subjective criteria, no manual checks, no "user satisfaction" metrics.

**AUTOMATION FIRST**: If a human would need to manually verify something, it's not a valid acceptance criterion for T.A.S.K.S.

---

## 2. Input/Output Schemas

### 2.1 Input Schema

```json
{
  "criteria": [
    {
      "id": "string",
      "description": "string",
      "type": "file_exists|command_succeeds|file_contains|api_responds|test_passes|performance_meets|schema_validates",
      "target": "string",
      "parameters": {},
      "priority": "must_have|should_have|nice_to_have",
      "timeout_seconds": 30
    }
  ],
  "context": {
    "task_id": "string",
    "project_root": "string",
    "environment": "development|staging|production",
    "dependencies": ["string"],
    "working_directory": "string",
    "available_tools": ["string"]
  },
  "validation_level": "strict|moderate|permissive"
}
```

### 2.2 Output Schema

```json
{
  "is_valid": true,
  "validation_level": "strict",
  "summary": {
    "total_criteria": 8,
    "valid_criteria": 7,
    "invalid_criteria": 1,
    "executable_criteria": 7,
    "machine_checkable": 7
  },
  "criteria_analysis": [
    {
      "criterion_id": "api_endpoint_works",
      "is_valid": true,
      "is_executable": true,
      "is_machine_checkable": true,
      "executable_command": "curl -f http://localhost:3000/api/health",
      "expected_result": "exit_code_0",
      "validation_issues": [],
      "recommendations": []
    }
  ],
  "executable_test_suite": {
    "setup_commands": ["string"],
    "test_commands": ["string"], 
    "cleanup_commands": ["string"],
    "total_estimated_runtime": 45
  },
  "validation_issues": [
    {
      "criterion_id": "user_satisfaction",
      "issue_type": "not_machine_checkable",
      "severity": "error",
      "message": "Criterion 'user satisfaction' cannot be automatically validated",
      "suggestion": "Replace with specific UI interaction tests using automated testing tools"
    }
  ],
  "coverage_analysis": {
    "functional_coverage": 0.85,
    "error_case_coverage": 0.6,
    "performance_coverage": 0.4,
    "security_coverage": 0.3,
    "missing_coverage_areas": ["input_validation", "error_handling"]
  }
}
```

---

## 3. Algorithm Implementation

### 3.1 Core Validation Algorithm

```python
def validate_acceptance_criteria(criteria, context, validation_level):
    """
    Core validation algorithm following strict machine-checkability principles.
    """
    
    results = []
    executable_commands = []
    
    for criterion in criteria:
        # Phase 1: Structural Validation
        structural_result = validate_criterion_structure(criterion)
        
        # Phase 2: Machine-Checkability Analysis  
        checkability_result = analyze_machine_checkability(criterion, context)
        
        # Phase 3: Executable Command Generation
        executable_result = generate_executable_command(criterion, context)
        
        # Phase 4: Dependency and Environment Validation
        dependency_result = validate_dependencies(criterion, context)
        
        # Phase 5: Performance and Timeout Analysis
        performance_result = analyze_performance_requirements(criterion)
        
        # Combine results
        criterion_result = combine_validation_results(
            structural_result,
            checkability_result, 
            executable_result,
            dependency_result,
            performance_result
        )
        
        results.append(criterion_result)
        
        if criterion_result["is_executable"]:
            executable_commands.append(criterion_result["executable_command"])
    
    # Phase 6: Coverage Analysis
    coverage_analysis = analyze_coverage_completeness(results, context)
    
    # Phase 7: Test Suite Generation
    test_suite = generate_executable_test_suite(executable_commands, context)
    
    return build_final_result(results, coverage_analysis, test_suite, validation_level)
```

### 3.2 Machine-Checkability Analysis

```python
def analyze_machine_checkability(criterion, context):
    """
    Determines if criterion can be automatically validated.
    """
    
    MACHINE_CHECKABLE_PATTERNS = {
        "file_operations": [
            r"file.*exists",
            r"directory.*created", 
            r".*\.py.*contains",
            r"config.*has.*value"
        ],
        "command_execution": [
            r"command.*succeeds",
            r"script.*runs.*without.*error",
            r"test.*passes",
            r"build.*completes"
        ],
        "api_interactions": [
            r"endpoint.*returns.*\d+",
            r"api.*responds.*with",
            r"service.*available",
            r"health.*check.*passes"
        ],
        "performance_metrics": [
            r"response.*time.*under.*\d+",
            r"memory.*usage.*below.*\d+",
            r"cpu.*utilization.*under.*\d+"
        ]
    }
    
    NON_CHECKABLE_PATTERNS = [
        r"user.*satisfaction",
        r"looks.*good",
        r"feels.*responsive", 
        r"code.*quality.*improved",
        r"team.*agrees",
        r"stakeholder.*approval"
    ]
    
    description = criterion["description"].lower()
    
    # Check for non-machine-checkable patterns
    for pattern in NON_CHECKABLE_PATTERNS:
        if re.search(pattern, description):
            return {
                "is_machine_checkable": False,
                "issue_type": "subjective_criterion",
                "problematic_phrase": re.search(pattern, description).group(),
                "suggestion": generate_objective_alternative(pattern, criterion)
            }
    
    # Check for machine-checkable patterns
    for category, patterns in MACHINE_CHECKABLE_PATTERNS.items():
        for pattern in patterns:
            if re.search(pattern, description):
                return {
                    "is_machine_checkable": True,
                    "checkability_category": category,
                    "matched_pattern": pattern,
                    "confidence": calculate_checkability_confidence(criterion, context)
                }
    
    # Ambiguous case - attempt to infer checkability
    return analyze_ambiguous_criterion(criterion, context)
```

### 3.3 Executable Command Generation

```python
def generate_executable_command(criterion, context):
    """
    Converts validated criterion into executable shell command.
    """
    
    criterion_type = criterion.get("type", "inferred")
    target = criterion.get("target", "")
    parameters = criterion.get("parameters", {})
    
    command_generators = {
        "file_exists": lambda: f'test -f "{target}"',
        "directory_exists": lambda: f'test -d "{target}"',
        "file_contains": lambda: f'grep -q "{parameters.get("pattern", "")}" "{target}"',
        "command_succeeds": lambda: f'{target}',
        "api_responds": lambda: f'curl -f -s "{target}" > /dev/null',
        "test_passes": lambda: f'{target} --json | jq -e ".success == true"',
        "performance_meets": lambda: generate_performance_check(criterion, context),
        "schema_validates": lambda: f'jsonschema -i "{target}" "{parameters.get("schema_file", "")}"'
    }
    
    if criterion_type in command_generators:
        base_command = command_generators[criterion_type]()
        
        # Add timeout wrapper
        timeout = criterion.get("timeout_seconds", 30)
        command_with_timeout = f'timeout {timeout} {base_command}'
        
        # Add error handling
        robust_command = f'{command_with_timeout} && echo "PASS" || echo "FAIL"'
        
        return {
            "executable_command": robust_command,
            "base_command": base_command,
            "timeout_seconds": timeout,
            "expected_output": "PASS",
            "is_executable": True
        }
    
    # Attempt to infer command from description
    return infer_command_from_description(criterion, context)
```

### 3.4 Coverage Analysis Implementation

```python
def analyze_coverage_completeness(validation_results, context):
    """
    Analyzes acceptance criteria coverage across multiple dimensions.
    """
    
    coverage_dimensions = {
        "functional": analyze_functional_coverage,
        "error_case": analyze_error_coverage,
        "performance": analyze_performance_coverage,
        "security": analyze_security_coverage,
        "integration": analyze_integration_coverage,
        "edge_case": analyze_edge_case_coverage
    }
    
    coverage_scores = {}
    missing_areas = []
    
    for dimension, analyzer in coverage_dimensions.items():
        score, gaps = analyzer(validation_results, context)
        coverage_scores[dimension] = score
        missing_areas.extend(gaps)
    
    # Identify critical coverage gaps
    critical_gaps = identify_critical_gaps(missing_areas, context)
    
    # Generate coverage improvement recommendations
    recommendations = generate_coverage_recommendations(coverage_scores, critical_gaps)
    
    return {
        "coverage_scores": coverage_scores,
        "overall_coverage": calculate_overall_coverage(coverage_scores),
        "missing_areas": missing_areas,
        "critical_gaps": critical_gaps,
        "recommendations": recommendations
    }
```

---

## 4. Error Handling

### 4.1 Validation Error Categories

```python
class AcceptanceValidationError(Exception):
    """Base exception for acceptance validation errors."""
    pass

class NonExecutableCriterionError(AcceptanceValidationError):
    """Raised when criterion cannot be converted to executable command."""
    def __init__(self, criterion_id, reason):
        self.criterion_id = criterion_id
        self.reason = reason
        super().__init__(f"Criterion {criterion_id} is not executable: {reason}")

class SubjectiveCriterionError(AcceptanceValidationError):
    """Raised when criterion requires subjective human judgment."""
    def __init__(self, criterion_id, problematic_phrase):
        self.criterion_id = criterion_id
        self.problematic_phrase = problematic_phrase
        super().__init__(f"Criterion {criterion_id} contains subjective element: {problematic_phrase}")

class MissingDependencyError(AcceptanceValidationError):
    """Raised when criterion requires unavailable tools or dependencies."""
    def __init__(self, criterion_id, missing_dependencies):
        self.criterion_id = criterion_id
        self.missing_dependencies = missing_dependencies
        super().__init__(f"Criterion {criterion_id} requires missing dependencies: {missing_dependencies}")

class InvalidTargetError(AcceptanceValidationError):
    """Raised when criterion target is invalid or inaccessible."""
    def __init__(self, criterion_id, target, reason):
        self.criterion_id = criterion_id
        self.target = target
        self.reason = reason
        super().__init__(f"Criterion {criterion_id} has invalid target {target}: {reason}")
```

### 4.2 Error Recovery Strategies

```python
def handle_validation_error(error, criterion, context, validation_level):
    """
    Implements error recovery strategies based on validation level.
    """
    
    if validation_level == "strict":
        # Strict mode: Any error invalidates the criterion
        return {
            "is_valid": False,
            "error": str(error),
            "recovery_attempted": False,
            "suggestion": generate_fix_suggestion(error, criterion)
        }
    
    elif validation_level == "moderate":
        # Moderate mode: Attempt recovery for certain error types
        if isinstance(error, MissingDependencyError):
            recovery_result = attempt_dependency_recovery(error, context)
            if recovery_result["success"]:
                return {
                    "is_valid": True,
                    "error": str(error),
                    "recovery_attempted": True,
                    "recovery_method": recovery_result["method"],
                    "modified_criterion": recovery_result["modified_criterion"]
                }
        
        return {
            "is_valid": False,
            "error": str(error),
            "recovery_attempted": True,
            "recovery_failed": True,
            "suggestion": generate_fix_suggestion(error, criterion)
        }
    
    elif validation_level == "permissive":
        # Permissive mode: Attempt to salvage partial functionality
        partial_result = extract_partial_functionality(criterion, error, context)
        return {
            "is_valid": partial_result["has_partial_functionality"],
            "error": str(error),
            "recovery_attempted": True,
            "partial_functionality": partial_result,
            "warning": "Criterion has limited automated validation capability"
        }
```

---

## 5. Test Cases

### 5.1 Valid Acceptance Criteria Test Cases

```python
VALID_CRITERIA_TESTS = [
    {
        "name": "file_existence_check",
        "criterion": {
            "id": "config_file_exists",
            "description": "Configuration file config/app.json exists",
            "type": "file_exists",
            "target": "config/app.json",
            "priority": "must_have"
        },
        "expected_command": 'test -f "config/app.json" && echo "PASS" || echo "FAIL"',
        "expected_result": {"is_valid": True, "is_executable": True}
    },
    {
        "name": "api_endpoint_test",
        "criterion": {
            "id": "health_endpoint_responds",
            "description": "Health check endpoint returns 200 status",
            "type": "api_responds",
            "target": "http://localhost:3000/health",
            "parameters": {"expected_status": 200},
            "timeout_seconds": 10
        },
        "expected_command": 'timeout 10 curl -f -s "http://localhost:3000/health" > /dev/null && echo "PASS" || echo "FAIL"',
        "expected_result": {"is_valid": True, "is_executable": True}
    },
    {
        "name": "test_suite_execution",
        "criterion": {
            "id": "unit_tests_pass",
            "description": "All unit tests pass successfully",
            "type": "command_succeeds",
            "target": "npm test",
            "priority": "must_have",
            "timeout_seconds": 120
        },
        "expected_command": 'timeout 120 npm test && echo "PASS" || echo "FAIL"',
        "expected_result": {"is_valid": True, "is_executable": True}
    },
    {
        "name": "performance_criterion",
        "criterion": {
            "id": "response_time_under_500ms",
            "description": "API response time is under 500ms",
            "type": "performance_meets",
            "target": "http://localhost:3000/api/data",
            "parameters": {"max_response_time_ms": 500, "measurement_count": 10},
            "priority": "should_have"
        },
        "expected_result": {"is_valid": True, "is_executable": True}
    }
]
```

### 5.2 Invalid Acceptance Criteria Test Cases

```python
INVALID_CRITERIA_TESTS = [
    {
        "name": "subjective_user_satisfaction",
        "criterion": {
            "id": "user_happy",
            "description": "Users are satisfied with the new interface",
            "type": "subjective",
            "priority": "should_have"
        },
        "expected_error": "SubjectiveCriterionError",
        "expected_suggestion": "Replace with specific UI interaction tests using automated testing tools"
    },
    {
        "name": "vague_quality_improvement",
        "criterion": {
            "id": "code_quality_better",
            "description": "Code quality is improved",
            "type": "quality",
            "priority": "nice_to_have"
        },
        "expected_error": "SubjectiveCriterionError",
        "expected_suggestion": "Replace with specific metrics like test coverage, complexity scores, or linting results"
    },
    {
        "name": "missing_executable_target",
        "criterion": {
            "id": "something_works",
            "description": "The system works correctly",
            "type": "functional",
            "priority": "must_have"
        },
        "expected_error": "NonExecutableCriterionError",
        "expected_suggestion": "Specify concrete behavioral criteria with testable outcomes"
    },
    {
        "name": "invalid_file_path",
        "criterion": {
            "id": "nonexistent_file",
            "description": "File /invalid/path/file.txt contains configuration",
            "type": "file_contains",
            "target": "/invalid/path/file.txt",
            "parameters": {"pattern": "config"}
        },
        "expected_error": "InvalidTargetError",
        "expected_suggestion": "Verify file path exists or will be created by previous tasks"
    }
]
```

### 5.3 Coverage Analysis Test Cases

```python
COVERAGE_ANALYSIS_TESTS = [
    {
        "name": "comprehensive_coverage",
        "criteria": [
            {"type": "file_exists", "category": "functional"},
            {"type": "test_passes", "category": "functional"},
            {"type": "api_responds", "category": "functional"},
            {"type": "command_succeeds", "category": "error_case", "target": "npm test -- --error-scenarios"},
            {"type": "performance_meets", "category": "performance"},
            {"type": "schema_validates", "category": "integration"}
        ],
        "expected_coverage": {
            "functional": 0.9,
            "error_case": 0.3,
            "performance": 0.5,
            "security": 0.0,
            "integration": 0.4
        },
        "expected_gaps": ["security_testing", "edge_case_handling"]
    },
    {
        "name": "minimal_coverage",
        "criteria": [
            {"type": "file_exists", "category": "functional"}
        ],
        "expected_coverage": {
            "functional": 0.2,
            "error_case": 0.0,
            "performance": 0.0,
            "security": 0.0,
            "integration": 0.0
        },
        "expected_gaps": ["error_handling", "performance_validation", "security_checks", "integration_testing"]
    }
]
```

---

## 6. Performance Requirements

### 6.1 Execution Time Constraints

```python
PERFORMANCE_REQUIREMENTS = {
    "validation_time_per_criterion": {
        "target": "< 100ms",
        "maximum": "< 500ms",
        "measurement": "wall_clock_time"
    },
    "total_validation_time": {
        "target": "< 1s for 10 criteria", 
        "maximum": "< 5s for 50 criteria",
        "scaling": "linear"
    },
    "command_generation_time": {
        "target": "< 50ms per command",
        "maximum": "< 200ms per command",
        "complexity": "O(1) per criterion"
    },
    "coverage_analysis_time": {
        "target": "< 200ms",
        "maximum": "< 1s", 
        "complexity": "O(n) where n = number of criteria"
    }
}
```

### 6.2 Memory Usage Constraints

```python
MEMORY_REQUIREMENTS = {
    "max_memory_per_validation": "10MB",
    "max_total_memory": "50MB",
    "memory_scaling": "O(n) where n = number of criteria",
    "cleanup_strategy": "immediate_cleanup_after_validation"
}
```

### 6.3 Scalability Requirements

```python
SCALABILITY_REQUIREMENTS = {
    "max_criteria_per_task": 100,
    "max_concurrent_validations": 10,
    "batch_processing_size": 20,
    "timeout_handling": "fail_fast_with_partial_results"
}
```

---

## 7. Library Dependencies

### 7.1 Core Dependencies

```python
CORE_DEPENDENCIES = {
    "python": ">=3.8",
    "jsonschema": "^4.0.0",    # Schema validation
    "pydantic": "^1.8.0",      # Data validation
    "click": "^8.0.0",         # CLI interface
    "requests": "^2.25.0",     # HTTP requests for API testing
    "psutil": "^5.8.0",        # System resource monitoring
}
```

### 7.2 Optional Dependencies

```python
OPTIONAL_DEPENDENCIES = {
    "pytest": "^6.0.0",        # Test framework integration
    "coverage": "^6.0.0",      # Coverage analysis
    "jq": "system_dependency", # JSON processing for complex validations
    "curl": "system_dependency", # HTTP testing
    "timeout": "system_dependency" # Command timeout handling
}
```

### 7.3 Development Dependencies

```python
DEVELOPMENT_DEPENDENCIES = {
    "black": "^21.0.0",        # Code formatting
    "mypy": "^0.812",          # Type checking
    "pytest-cov": "^2.12.0",  # Test coverage
    "pre-commit": "^2.15.0"    # Git hooks
}
```

---

## 8. Integration Examples

### 8.1 T.A.S.K.S Pipeline Integration

```python
# Example integration in T.A.S.K.S pipeline
def validate_task_acceptance_criteria(task_spec):
    """
    Validates acceptance criteria during task processing.
    """
    
    acceptance_criteria = task_spec.get("acceptance_criteria", [])
    context = {
        "task_id": task_spec["id"],
        "project_root": task_spec.get("project_root", "."),
        "environment": "development",
        "working_directory": task_spec.get("working_directory", "."),
        "available_tools": ["npm", "curl", "test", "python", "node"]
    }
    
    # Validate acceptance criteria
    validation_result = acceptance_validator(
        criteria=acceptance_criteria,
        context=context,
        validation_level="strict"
    )
    
    if not validation_result["is_valid"]:
        raise TaskValidationError(
            f"Task {task_spec['id']} has invalid acceptance criteria",
            validation_errors=validation_result["validation_issues"]
        )
    
    # Enhance task with executable test suite
    task_spec["executable_tests"] = validation_result["executable_test_suite"]
    task_spec["coverage_analysis"] = validation_result["coverage_analysis"]
    
    return task_spec
```

### 8.2 CLI Usage Example

```bash
# Validate acceptance criteria for a task
acceptance_validator validate \
    --criteria-file task_criteria.json \
    --context-file task_context.json \
    --validation-level strict \
    --output-format json

# Generate executable test suite
acceptance_validator generate-tests \
    --criteria-file task_criteria.json \
    --output-file acceptance_tests.sh \
    --include-setup \
    --include-cleanup

# Analyze coverage gaps
acceptance_validator analyze-coverage \
    --criteria-file task_criteria.json \
    --report-format markdown \
    --output-file coverage_report.md
```

---

## 9. Future Enhancements

### 9.1 Planned Extensions

1. **AI-Powered Criterion Enhancement**: Use LLM to suggest improvements for vague criteria
2. **Integration Test Generation**: Automatically generate integration tests from acceptance criteria
3. **Performance Regression Detection**: Historical performance tracking for performance criteria
4. **Security Criterion Templates**: Pre-built templates for common security requirements
5. **Visual Test Reporting**: HTML reports with execution graphs and coverage visualizations

### 9.2 Advanced Features

1. **Criterion Dependency Analysis**: Detect dependencies between acceptance criteria
2. **Parallel Test Execution**: Execute independent criteria in parallel for faster validation
3. **Smart Timeout Adjustment**: Dynamically adjust timeouts based on system performance
4. **Rollback Capability**: Undo changes if acceptance criteria fail during execution
5. **Criterion Versioning**: Track changes to acceptance criteria over time

---

This specification provides a complete, implementable foundation for the `acceptance_validator` tool, ensuring all acceptance criteria in T.A.S.K.S. v2 are machine-checkable and executable, eliminating the risk of unverifiable task completion.
