# secret_redactor Tool - Complete Implementation Specification

> **Zero-Questions Implementation Guide for Security-Critical Secret Detection and Redaction**

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
def redact_secrets(
    content_data: Dict,
    redaction_config: Dict = None,
    custom_patterns: List[Dict] = None
) -> Dict:
    """
    Detects and redacts sensitive information from project artifacts.
    
    Args:
        content_data: Complete project content (tasks.json, markdown, etc.)
        redaction_config: Optional redaction configuration overrides
        custom_patterns: Optional custom secret detection patterns
        
    Returns:
        Complete redaction_report.json object with redacted content and findings
        
    Raises:
        ValidationError: Input schema validation failure
        RedactionError: Critical redaction processing failure
        ConfigurationError: Invalid configuration parameters
        ProcessingError: Unexpected processing failure
    """
```

### 1.2 Command Line Interface

```bash
python -m secret_redactor \
    --input project_artifacts.json \
    --output redacted_artifacts.json \
    --report redaction_report.json \
    --config redaction_config.json \
    --patterns custom_patterns.json \
    --strict-mode \
    --verbose
```

### 1.3 Redaction Configuration Schema

```json
{
  "redaction_config": {
    "secret_detection": {
      "enable_entropy_analysis": true,
      "entropy_threshold": 4.5,
      "min_secret_length": 8,
      "max_secret_length": 128,
      "enable_pattern_matching": true,
      "enable_keyword_detection": true
    },
    "redaction_strategies": {
      "full_redaction": "[REDACTED]",
      "partial_redaction": true,
      "preserve_prefix_chars": 2,
      "preserve_suffix_chars": 2,
      "mask_character": "*",
      "structured_redaction": true
    },
    "content_types": {
      "scan_json": true,
      "scan_markdown": true,
      "scan_plain_text": true,
      "scan_code_blocks": true,
      "scan_urls": true,
      "scan_file_paths": true
    },
    "sensitivity_levels": {
      "api_keys": "high",
      "passwords": "critical",
      "tokens": "high", 
      "private_keys": "critical",
      "database_urls": "high",
      "email_addresses": "medium",
      "ip_addresses": "low",
      "internal_hostnames": "medium"
    },
    "whitelist": {
      "enable_whitelist": true,
      "whitelisted_patterns": [],
      "whitelisted_domains": ["example.com", "localhost"],
      "whitelisted_prefixes": ["demo-", "test-", "mock-"]
    },
    "performance_limits": {
      "max_content_size_mb": 10,
      "max_processing_time_seconds": 30,
      "max_patterns": 500
    }
  }
}
```

---

## 2. Input/Output Schemas

### 2.1 Input Schema: content_data.json

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["artifacts"],
  "properties": {
    "artifacts": {
      "type": "array",
      "minItems": 1,
      "maxItems": 1000,
      "items": {
        "type": "object",
        "required": ["id", "type", "content"],
        "properties": {
          "id": {
            "type": "string",
            "description": "Unique identifier for the artifact"
          },
          "type": {
            "type": "string",
            "enum": ["tasks_json", "markdown", "plain_text", "code", "configuration", "documentation"],
            "description": "Type of content for appropriate scanning"
          },
          "content": {
            "type": "string",
            "description": "The actual content to scan for secrets"
          },
          "filename": {
            "type": "string",
            "description": "Original filename for context"
          },
          "metadata": {
            "type": "object",
            "properties": {
              "created_by": {"type": "string"},
              "created_at": {"type": "string", "format": "date-time"},
              "sensitivity_level": {"type": "string", "enum": ["public", "internal", "confidential", "restricted"]}
            }
          }
        }
      }
    },
    "scanning_context": {
      "type": "object",
      "properties": {
        "project_name": {"type": "string"},
        "organization": {"type": "string"},
        "environment": {"type": "string", "enum": ["development", "staging", "production"]},
        "compliance_requirements": {"type": "array", "items": {"type": "string"}}
      }
    }
  }
}
```

### 2.2 Input Schema: custom_patterns.json (optional)

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "patterns": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["name", "pattern", "severity"],
        "properties": {
          "name": {
            "type": "string",
            "description": "Human-readable name for the pattern"
          },
          "pattern": {
            "type": "string",
            "description": "Regular expression pattern"
          },
          "severity": {
            "type": "string",
            "enum": ["low", "medium", "high", "critical"],
            "description": "Severity level of detected secrets"
          },
          "description": {
            "type": "string",
            "description": "Description of what this pattern detects"
          },
          "confidence_modifier": {
            "type": "number",
            "minimum": 0.0,
            "maximum": 2.0,
            "default": 1.0,
            "description": "Multiplier for confidence scoring"
          },
          "context_keywords": {
            "type": "array",
            "items": {"type": "string"},
            "description": "Keywords that increase detection confidence"
          },
          "false_positive_patterns": {
            "type": "array",
            "items": {"type": "string"},
            "description": "Patterns that indicate false positives"
          }
        }
      }
    }
  }
}
```

### 2.3 Output Schema: redaction_report.json

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["ok", "generated", "redaction_summary", "redacted_artifacts"],
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
            "enum": ["VALIDATION_ERROR", "REDACTION_ERROR", "PATTERN_ERROR", "PROCESSING_ERROR"]
          },
          "message": {"type": "string"},
          "details": {"type": "object"},
          "affected_artifacts": {"type": "array", "items": {"type": "string"}}
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
          "artifact_id": {"type": "string"},
          "severity": {"type": "string", "enum": ["low", "medium", "high"]}
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
    "redaction_summary": {
      "type": "object",
      "required": ["total_artifacts", "artifacts_with_secrets", "total_secrets_found"],
      "properties": {
        "total_artifacts": {"type": "integer", "minimum": 0},
        "artifacts_with_secrets": {"type": "integer", "minimum": 0},
        "total_secrets_found": {"type": "integer", "minimum": 0},
        "secrets_redacted": {"type": "integer", "minimum": 0},
        "high_confidence_detections": {"type": "integer", "minimum": 0},
        "potential_false_positives": {"type": "integer", "minimum": 0},
        "redaction_rate": {"type": "number", "minimum": 0.0, "maximum": 1.0},
        "by_severity": {
          "type": "object",
          "properties": {
            "critical": {"type": "integer", "minimum": 0},
            "high": {"type": "integer", "minimum": 0},
            "medium": {"type": "integer", "minimum": 0},
            "low": {"type": "integer", "minimum": 0}
          }
        },
        "by_detection_method": {
          "type": "object",
          "properties": {
            "pattern_matching": {"type": "integer", "minimum": 0},
            "entropy_analysis": {"type": "integer", "minimum": 0},
            "keyword_detection": {"type": "integer", "minimum": 0},
            "custom_patterns": {"type": "integer", "minimum": 0}
          }
        }
      }
    },
    "redacted_artifacts": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["id", "original_content", "redacted_content"],
        "properties": {
          "id": {"type": "string"},
          "original_content": {"type": "string"},
          "redacted_content": {"type": "string"},
          "redaction_applied": {"type": "boolean"},
          "secrets_found": {"type": "integer", "minimum": 0},
          "redaction_locations": {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "start_position": {"type": "integer"},
                "end_position": {"type": "integer"},
                "original_text": {"type": "string"},
                "redacted_text": {"type": "string"},
                "redaction_reason": {"type": "string"}
              }
            }
          }
        }
      }
    },
    "detected_secrets": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["id", "artifact_id", "detection_method", "severity", "confidence_score"],
        "properties": {
          "id": {"type": "string"},
          "artifact_id": {"type": "string"},
          "secret_type": {
            "type": "string",
            "enum": [
              "api_key", "password", "token", "private_key", "certificate",
              "database_url", "connection_string", "email", "phone_number",
              "ssn", "credit_card", "ip_address", "hostname", "aws_key",
              "gcp_key", "azure_key", "github_token", "slack_token", "generic_secret"
            ]
          },
          "detection_method": {
            "type": "string",
            "enum": ["pattern_matching", "entropy_analysis", "keyword_detection", "custom_pattern"]
          },
          "severity": {
            "type": "string",
            "enum": ["low", "medium", "high", "critical"]
          },
          "confidence_score": {"type": "number", "minimum": 0.0, "maximum": 1.0},
          "location": {
            "type": "object",
            "properties": {
              "start_position": {"type": "integer"},
              "end_position": {"type": "integer"},
              "line_number": {"type": "integer"},
              "column_number": {"type": "integer"},
              "context": {"type": "string", "description": "Surrounding text for context"}
            }
          },
          "pattern_matched": {"type": "string"},
          "entropy_score": {"type": "number", "minimum": 0.0},
          "redaction_applied": {"type": "boolean"},
          "redaction_strategy": {"type": "string"},
          "false_positive_risk": {
            "type": "string",
            "enum": ["low", "medium", "high"]
          },
          "recommendations": {
            "type": "array",
            "items": {"type": "string"}
          }
        }
      }
    },
    "false_positive_analysis": {
      "type": "object",
      "properties": {
        "total_potential_false_positives": {"type": "integer"},
        "false_positive_patterns": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "pattern": {"type": "string"},
              "count": {"type": "integer"},
              "examples": {"type": "array", "items": {"type": "string"}}
            }
          }
        },
        "recommendations": {
          "type": "array",
          "items": {"type": "string"}
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
def redact_secrets(content_data, redaction_config=None, custom_patterns=None):
    """Complete implementation with multi-strategy secret detection and redaction."""
    
    # Step 1: Input Validation
    config = _merge_config(redaction_config)
    _validate_input_schema(content_data)
    _validate_custom_patterns(custom_patterns)
    
    # Step 2: Load Detection Patterns
    detection_patterns = _load_detection_patterns(config, custom_patterns)
    
    # Step 3: Initialize Detection Engines
    pattern_engine = _initialize_pattern_engine(detection_patterns, config)
    entropy_engine = _initialize_entropy_engine(config)
    keyword_engine = _initialize_keyword_engine(config)
    
    # Step 4: Process Each Artifact
    redacted_artifacts = []
    all_detected_secrets = []
    
    for artifact in content_data["artifacts"]:
        # Scan for secrets using multiple detection methods
        detected_secrets = _scan_artifact_for_secrets(
            artifact, pattern_engine, entropy_engine, keyword_engine, config
        )
        
        # Apply redaction strategies
        redacted_artifact = _apply_redaction_to_artifact(
            artifact, detected_secrets, config
        )
        
        redacted_artifacts.append(redacted_artifact)
        all_detected_secrets.extend(detected_secrets)
    
    # Step 5: False Positive Analysis
    false_positive_analysis = _analyze_false_positives(all_detected_secrets, config)
    
    # Step 6: Generate Summary Statistics
    summary = _generate_redaction_summary(redacted_artifacts, all_detected_secrets)
    
    # Step 7: Generate Content Hash
    content_hash = _generate_content_hash(redacted_artifacts, summary)
    
    # Step 8: Build Response
    return _build_redaction_response(
        redacted_artifacts, all_detected_secrets, false_positive_analysis,
        summary, content_hash
    )
```

### 3.2 Detection Pattern Loading

```python
def _load_detection_patterns(config, custom_patterns):
    """Load built-in and custom detection patterns."""
    
    patterns = []
    
    # Load built-in patterns
    builtin_patterns = _get_builtin_patterns()
    patterns.extend(builtin_patterns)
    
    # Load custom patterns if provided
    if custom_patterns:
        validated_custom = _validate_custom_patterns(custom_patterns)
        patterns.extend(validated_custom)
    
    # Filter patterns by configuration
    enabled_patterns = []
    for pattern in patterns:
        if _is_pattern_enabled(pattern, config):
            enabled_patterns.append(pattern)
    
    return enabled_patterns

def _get_builtin_patterns():
    """Get comprehensive built-in secret detection patterns."""
    
    return [
        # AWS Access Keys
        {
            "name": "AWS Access Key ID",
            "pattern": r"(?i)(?:aws.{0,20})?(?:access.{0,20})?key.{0,20}[:=\\s]\\s*([A-Z0-9]{20})",
            "severity": "critical",
            "secret_type": "aws_key",
            "confidence_modifier": 1.0,
            "context_keywords": ["aws", "amazon", "access", "key", "secret"],
            "false_positive_patterns": [r"EXAMPLE", r"SAMPLE", r"TEST", r"DEMO"]
        },
        {
            "name": "AWS Secret Access Key", 
            "pattern": r"(?i)(?:aws.{0,20})?(?:secret.{0,20})?(?:access.{0,20})?key.{0,20}[:=\\s]\\s*([A-Za-z0-9/+=]{40})",
            "severity": "critical",
            "secret_type": "aws_key",
            "confidence_modifier": 1.0,
            "context_keywords": ["aws", "secret", "access", "key"],
            "false_positive_patterns": [r"EXAMPLE", r"SAMPLE"]
        },
        
        # API Keys - Generic
        {
            "name": "Generic API Key",
            "pattern": r"(?i)(?:api.{0,20})?key.{0,20}[:=\\s]\\s*([A-Za-z0-9]{32,})",
            "severity": "high",
            "secret_type": "api_key",
            "confidence_modifier": 0.8,
            "context_keywords": ["api", "key", "token", "auth"],
            "false_positive_patterns": [r"example", r"sample", r"placeholder"]
        },
        
        # GitHub Tokens
        {
            "name": "GitHub Personal Access Token",
            "pattern": r"\\bghp_[A-Za-z0-9]{36}\\b",
            "severity": "high",
            "secret_type": "github_token",
            "confidence_modifier": 1.0,
            "context_keywords": ["github", "token", "pat"],
            "false_positive_patterns": []
        },
        {
            "name": "GitHub OAuth Token",
            "pattern": r"\\bgho_[A-Za-z0-9]{36}\\b",
            "severity": "high",
            "secret_type": "github_token",
            "confidence_modifier": 1.0,
            "context_keywords": ["github", "oauth"],
            "false_positive_patterns": []
        },
        
        # Slack Tokens
        {
            "name": "Slack Bot Token",
            "pattern": r"\\bxoxb-[0-9]{11,}-[0-9]{11,}-[A-Za-z0-9]{24}\\b",
            "severity": "high",
            "secret_type": "slack_token",
            "confidence_modifier": 1.0,
            "context_keywords": ["slack", "bot", "token"],
            "false_positive_patterns": []
        },
        {
            "name": "Slack User Token",
            "pattern": r"\\bxoxp-[0-9]{11,}-[0-9]{11,}-[0-9]{11,}-[A-Za-z0-9]{32}\\b",
            "severity": "critical",
            "secret_type": "slack_token",
            "confidence_modifier": 1.0,
            "context_keywords": ["slack", "user", "token"],
            "false_positive_patterns": []
        },
        
        # Database URLs
        {
            "name": "Database Connection String",
            "pattern": r"(?i)(?:mongodb|mysql|postgresql|postgres|redis)://[^\\s\"']+",
            "severity": "high",
            "secret_type": "database_url",
            "confidence_modifier": 0.9,
            "context_keywords": ["database", "db", "connection", "url"],
            "false_positive_patterns": [r"localhost", r"example\\.com", r"127\\.0\\.0\\.1"]
        },
        
        # Private Keys
        {
            "name": "RSA Private Key",
            "pattern": r"-----BEGIN RSA PRIVATE KEY-----[\\s\\S]*?-----END RSA PRIVATE KEY-----",
            "severity": "critical",
            "secret_type": "private_key",
            "confidence_modifier": 1.0,
            "context_keywords": ["private", "key", "rsa"],
            "false_positive_patterns": [r"EXAMPLE", r"SAMPLE"]
        },
        {
            "name": "Generic Private Key",
            "pattern": r"-----BEGIN PRIVATE KEY-----[\\s\\S]*?-----END PRIVATE KEY-----",
            "severity": "critical",
            "secret_type": "private_key",
            "confidence_modifier": 1.0,
            "context_keywords": ["private", "key"],
            "false_positive_patterns": [r"EXAMPLE", r"SAMPLE"]
        },
        
        # Email Addresses
        {
            "name": "Email Address",
            "pattern": r"\\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}\\b",
            "severity": "medium",
            "secret_type": "email",
            "confidence_modifier": 0.7,
            "context_keywords": ["email", "mail", "contact"],
            "false_positive_patterns": [r"example\\.com", r"test\\.com", r"sample\\.com"]
        },
        
        # IP Addresses
        {
            "name": "IP Address",
            "pattern": r"\\b(?:[0-9]{1,3}\\.){3}[0-9]{1,3}\\b",
            "severity": "low",
            "secret_type": "ip_address",
            "confidence_modifier": 0.6,
            "context_keywords": ["ip", "address", "server", "host"],
            "false_positive_patterns": [r"127\\.0\\.0\\.1", r"0\\.0\\.0\\.0", r"255\\.255\\.255\\.255"]
        },
        
        # Phone Numbers (US format)
        {
            "name": "US Phone Number",
            "pattern": r"\\b(?:\\+1[-\\s]?)?\\(?[2-9][0-8][0-9]\\)?[-\\s]?[2-9][0-9]{2}[-\\s]?[0-9]{4}\\b",
            "severity": "medium",
            "secret_type": "phone_number",
            "confidence_modifier": 0.7,
            "context_keywords": ["phone", "telephone", "mobile", "cell"],
            "false_positive_patterns": [r"555-?0000", r"123-?4567"]
        },
        
        # Credit Card Numbers (Luhn algorithm would be applied separately)
        {
            "name": "Credit Card Number",
            "pattern": r"\\b(?:4[0-9]{12}(?:[0-9]{3})?|5[1-5][0-9]{14}|3[47][0-9]{13}|3[0-9]{13}|6(?:011|5[0-9]{2})[0-9]{12})\\b",
            "severity": "critical",
            "secret_type": "credit_card",
            "confidence_modifier": 0.9,
            "context_keywords": ["card", "credit", "payment", "visa", "mastercard"],
            "false_positive_patterns": [r"4111111111111111", r"5555555555554444"]
        },
        
        # Generic Passwords
        {
            "name": "Password Field",
            "pattern": r"(?i)(?:password|passwd|pwd)\\s*[:=]\\s*[\"']?([^\\s\"']{8,})[\"']?",
            "severity": "high",
            "secret_type": "password",
            "confidence_modifier": 0.8,
            "context_keywords": ["password", "passwd", "pwd", "auth"],
            "false_positive_patterns": [r"password", r"\\*+", r"example", r"placeholder"]
        }
    ]

def _is_pattern_enabled(pattern, config):
    """Check if a pattern should be enabled based on configuration."""
    
    # Check sensitivity level configuration
    secret_type = pattern.get("secret_type", "generic_secret")
    severity = pattern.get("severity", "medium")
    
    if secret_type in config["sensitivity_levels"]:
        required_level = config["sensitivity_levels"][secret_type]
        if severity != required_level and severity not in ["high", "critical"]:
            return False
    
    return True
```

### 3.3 Pattern-Based Detection Engine

```python
def _initialize_pattern_engine(detection_patterns, config):
    """Initialize the pattern-based detection engine."""
    import re
    
    compiled_patterns = []
    
    for pattern_def in detection_patterns:
        try:
            compiled_pattern = {
                "name": pattern_def["name"],
                "regex": re.compile(pattern_def["pattern"], re.MULTILINE | re.DOTALL),
                "severity": pattern_def["severity"],
                "secret_type": pattern_def.get("secret_type", "generic_secret"),
                "confidence_modifier": pattern_def.get("confidence_modifier", 1.0),
                "context_keywords": pattern_def.get("context_keywords", []),
                "false_positive_patterns": [
                    re.compile(fp_pattern, re.IGNORECASE) 
                    for fp_pattern in pattern_def.get("false_positive_patterns", [])
                ]
            }
            compiled_patterns.append(compiled_pattern)
        except re.error as e:
            # Log pattern compilation error but continue
            print(f"Warning: Failed to compile pattern '{pattern_def['name']}': {e}")
    
    return {
        "patterns": compiled_patterns,
        "config": config
    }

def _detect_secrets_by_patterns(content, pattern_engine, artifact_id):
    """Detect secrets using compiled regex patterns."""
    
    detected_secrets = []
    patterns = pattern_engine["patterns"]
    
    for pattern_info in patterns:
        regex = pattern_info["regex"]
        matches = regex.finditer(content)
        
        for match in matches:
            # Extract the secret value
            secret_value = match.group(1) if match.groups() else match.group(0)
            
            # Check for false positives
            if _is_false_positive(secret_value, pattern_info["false_positive_patterns"]):
                continue
            
            # Calculate confidence score
            confidence = _calculate_pattern_confidence(
                secret_value, match, content, pattern_info
            )
            
            # Get location information
            location = _get_location_info(match, content)
            
            secret = {
                "id": _generate_secret_id(),
                "artifact_id": artifact_id,
                "secret_type": pattern_info["secret_type"],
                "detection_method": "pattern_matching",
                "severity": pattern_info["severity"],
                "confidence_score": confidence,
                "location": location,
                "pattern_matched": pattern_info["name"],
                "secret_value": secret_value,  # Will be redacted before output
                "false_positive_risk": _assess_false_positive_risk(secret_value, pattern_info)
            }
            
            detected_secrets.append(secret)
    
    return detected_secrets

def _is_false_positive(secret_value, false_positive_patterns):
    """Check if detected value matches false positive patterns."""
    
    for fp_pattern in false_positive_patterns:
        if fp_pattern.search(secret_value):
            return True
    
    return False

def _calculate_pattern_confidence(secret_value, match, content, pattern_info):
    """Calculate confidence score for pattern-based detection."""
    
    base_confidence = 0.7
    confidence_modifier = pattern_info["confidence_modifier"]
    
    # Adjust confidence based on context keywords
    context_start = max(0, match.start() - 50)
    context_end = min(len(content), match.end() + 50)
    context = content[context_start:context_end].lower()
    
    keyword_bonus = 0
    for keyword in pattern_info["context_keywords"]:
        if keyword.lower() in context:
            keyword_bonus += 0.1
    
    # Adjust confidence based on secret length and complexity
    length_bonus = min(0.2, len(secret_value) / 100.0)
    
    # Entropy bonus for high-entropy strings
    entropy_bonus = min(0.15, (_calculate_entropy(secret_value) - 3.0) / 10.0)
    
    final_confidence = base_confidence + keyword_bonus + length_bonus + entropy_bonus
    final_confidence *= confidence_modifier
    
    return min(1.0, max(0.0, final_confidence))

def _get_location_info(match, content):
    """Get location information for a detected secret."""
    
    start_pos = match.start()
    end_pos = match.end()
    
    # Calculate line and column numbers
    lines_before = content[:start_pos].count('\\n')
    line_number = lines_before + 1
    
    last_newline = content.rfind('\\n', 0, start_pos)
    column_number = start_pos - last_newline
    
    # Get surrounding context
    context_start = max(0, start_pos - 100)
    context_end = min(len(content), end_pos + 100)
    context = content[context_start:context_end]
    
    return {
        "start_position": start_pos,
        "end_position": end_pos,
        "line_number": line_number,
        "column_number": column_number,
        "context": context
    }
```

### 3.4 Entropy-Based Detection Engine

```python
def _initialize_entropy_engine(config):
    """Initialize the entropy-based detection engine."""
    
    return {
        "entropy_threshold": config["secret_detection"]["entropy_threshold"],
        "min_length": config["secret_detection"]["min_secret_length"],
        "max_length": config["secret_detection"]["max_secret_length"],
        "enabled": config["secret_detection"]["enable_entropy_analysis"]
    }

def _detect_secrets_by_entropy(content, entropy_engine, artifact_id):
    """Detect secrets using entropy analysis for high-randomness strings."""
    
    if not entropy_engine["enabled"]:
        return []
    
    detected_secrets = []
    words = _extract_words_from_content(content)
    
    for word_info in words:
        word = word_info["word"]
        position = word_info["position"]
        
        # Check length constraints
        if len(word) < entropy_engine["min_length"] or len(word) > entropy_engine["max_length"]:
            continue
        
        # Calculate entropy
        entropy_score = _calculate_entropy(word)
        
        if entropy_score >= entropy_engine["entropy_threshold"]:
            # Check if it looks like a secret (not just random text)
            if _looks_like_secret(word):
                confidence = _calculate_entropy_confidence(word, entropy_score, content, position)
                
                location = {
                    "start_position": position,
                    "end_position": position + len(word),
                    "line_number": _get_line_number(content, position),
                    "column_number": _get_column_number(content, position),
                    "context": _get_context(content, position, len(word))
                }
                
                secret = {
                    "id": _generate_secret_id(),
                    "artifact_id": artifact_id,
                    "secret_type": "generic_secret",
                    "detection_method": "entropy_analysis",
                    "severity": "medium",
                    "confidence_score": confidence,
                    "location": location,
                    "entropy_score": entropy_score,
                    "secret_value": word,
                    "false_positive_risk": _assess_entropy_false_positive_risk(word)
                }
                
                detected_secrets.append(secret)
    
    return detected_secrets

def _calculate_entropy(string):
    """Calculate Shannon entropy of a string."""
    import math
    from collections import Counter
    
    if not string:
        return 0
    
    # Count character frequencies
    char_counts = Counter(string)
    string_length = len(string)
    
    # Calculate entropy
    entropy = 0
    for count in char_counts.values():
        probability = count / string_length
        if probability > 0:
            entropy -= probability * math.log2(probability)
    
    return entropy

def _extract_words_from_content(content):
    """Extract potential secret words from content."""
    import re
    
    # Pattern to match potential secrets (alphanumeric strings with some symbols)
    pattern = r'\\b[A-Za-z0-9+/=_-]{8,}\\b'
    
    words = []
    for match in re.finditer(pattern, content):
        words.append({
            "word": match.group(0),
            "position": match.start()
        })
    
    return words

def _looks_like_secret(word):
    """Heuristic to determine if a high-entropy word looks like a secret."""
    
    # Check for common secret characteristics
    has_mixed_case = any(c.isupper() for c in word) and any(c.islower() for c in word)
    has_numbers = any(c.isdigit() for c in word)
    has_special_chars = any(c in "+/=_-" for c in word)
    
    # Secrets typically have mixed case, numbers, or special characters
    secret_indicators = sum([has_mixed_case, has_numbers, has_special_chars])
    
    # Exclude common English words (simple check)
    common_words = {"password", "username", "example", "sample", "test", "demo"}
    if word.lower() in common_words:
        return False
    
    # Check for base64-like patterns
    if re.match(r'^[A-Za-z0-9+/]*={0,2}$', word) and len(word) % 4 == 0:
        return True
    
    return secret_indicators >= 1

def _calculate_entropy_confidence(word, entropy_score, content, position):
    """Calculate confidence for entropy-based detection."""
    
    # Base confidence from entropy score
    base_confidence = min(0.8, (entropy_score - 4.0) / 4.0)
    
    # Context analysis
    context = _get_context(content, position, len(word))
    context_lower = context.lower()
    
    # Look for secret-related keywords nearby
    secret_keywords = ["key", "token", "secret", "password", "auth", "api", "access"]
    keyword_bonus = 0.1 if any(keyword in context_lower for keyword in secret_keywords) else 0
    
    # Length bonus
    length_bonus = min(0.1, (len(word) - 16) / 100.0)
    
    final_confidence = base_confidence + keyword_bonus + length_bonus
    return min(1.0, max(0.0, final_confidence))

def _get_line_number(content, position):
    """Get line number for a position in content."""
    return content[:position].count('\\n') + 1

def _get_column_number(content, position):
    """Get column number for a position in content."""
    last_newline = content.rfind('\\n', 0, position)
    return position - last_newline

def _get_context(content, position, length):
    """Get surrounding context for a position."""
    start = max(0, position - 50)
    end = min(len(content), position + length + 50)
    return content[start:end]
```

### 3.5 Keyword-Based Detection Engine

```python
def _initialize_keyword_engine(config):
    """Initialize keyword-based detection for context-aware detection."""
    
    keyword_patterns = [
        # API-related keywords
        {
            "keywords": ["api_key", "apikey", "api-key"],
            "value_pattern": r"[:=\\s]\\s*[\"']?([A-Za-z0-9]{16,})[\"']?",
            "severity": "high",
            "secret_type": "api_key"
        },
        # Password keywords
        {
            "keywords": ["password", "passwd", "pwd", "pass"],
            "value_pattern": r"[:=\\s]\\s*[\"']?([^\\s\"']{6,})[\"']?",
            "severity": "high",
            "secret_type": "password"
        },
        # Token keywords
        {
            "keywords": ["token", "access_token", "auth_token", "bearer"],
            "value_pattern": r"[:=\\s]\\s*[\"']?([A-Za-z0-9._-]{20,})[\"']?",
            "severity": "high",
            "secret_type": "token"
        },
        # Database keywords
        {
            "keywords": ["db_password", "database_password", "db_pass"],
            "value_pattern": r"[:=\\s]\\s*[\"']?([^\\s\"']{6,})[\"']?",
            "severity": "critical",
            "secret_type": "password"
        }
    ]
    
    return {
        "patterns": keyword_patterns,
        "enabled": config["secret_detection"]["enable_keyword_detection"]
    }

def _detect_secrets_by_keywords(content, keyword_engine, artifact_id):
    """Detect secrets using keyword-based context analysis."""
    
    if not keyword_engine["enabled"]:
        return []
    
    detected_secrets = []
    
    for pattern_info in keyword_engine["patterns"]:
        keywords = pattern_info["keywords"]
        value_pattern = pattern_info["value_pattern"]
        
        for keyword in keywords:
            # Create pattern that looks for keyword followed by value
            full_pattern = f"(?i){re.escape(keyword)}{value_pattern}"
            
            try:
                regex = re.compile(full_pattern)
                matches = regex.finditer(content)
                
                for match in matches:
                    secret_value = match.group(1) if match.groups() else match.group(0)
                    
                    # Skip obviously fake values
                    if _is_obvious_placeholder(secret_value):
                        continue
                    
                    confidence = _calculate_keyword_confidence(secret_value, keyword, content, match)
                    location = _get_location_info(match, content)
                    
                    secret = {
                        "id": _generate_secret_id(),
                        "artifact_id": artifact_id,
                        "secret_type": pattern_info["secret_type"],
                        "detection_method": "keyword_detection",
                        "severity": pattern_info["severity"],
                        "confidence_score": confidence,
                        "location": location,
                        "pattern_matched": f"Keyword: {keyword}",
                        "secret_value": secret_value,
                        "false_positive_risk": _assess_keyword_false_positive_risk(secret_value, keyword)
                    }
                    
                    detected_secrets.append(secret)
                    
            except re.error:
                # Skip invalid patterns
                continue
    
    return detected_secrets

def _is_obvious_placeholder(value):
    """Check if value is obviously a placeholder."""
    
    placeholder_patterns = [
        r"^\\*+$",  # All asterisks
        r"^x+$",    # All x's
        r"^[0]+$",  # All zeros
        r"(?i)^(password|secret|key|token|example|sample|test|demo|placeholder)$",
        r"(?i)^(your|my|the)_",
        r"(?i)_(here|here|placeholder|example)$"
    ]
    
    for pattern in placeholder_patterns:
        if re.match(pattern, value):
            return True
    
    return False

def _calculate_keyword_confidence(secret_value, keyword, content, match):
    """Calculate confidence for keyword-based detection."""
    
    base_confidence = 0.6
    
    # Higher confidence for longer, more complex values
    length_bonus = min(0.2, len(secret_value) / 50.0)
    
    # Entropy bonus
    entropy_bonus = min(0.15, (_calculate_entropy(secret_value) - 2.0) / 10.0)
    
    # Context bonus for security-related terms
    context = content[max(0, match.start() - 100):match.end() + 100]
    security_terms = ["auth", "secure", "credential", "login", "access"]
    context_bonus = 0.1 if any(term in context.lower() for term in security_terms) else 0
    
    final_confidence = base_confidence + length_bonus + entropy_bonus + context_bonus
    return min(1.0, max(0.0, final_confidence))
```

### 3.6 Redaction Strategies

```python
def _apply_redaction_to_artifact(artifact, detected_secrets, config):
    """Apply redaction strategies to an artifact."""
    
    content = artifact["content"]
    redaction_locations = []
    
    # Sort secrets by position (descending) to avoid position shifts during redaction
    secrets_by_position = sorted(detected_secrets, key=lambda s: s["location"]["start_position"], reverse=True)
    
    redacted_content = content
    
    for secret in secrets_by_position:
        location = secret["location"]
        start_pos = location["start_position"]
        end_pos = location["end_position"]
        
        # Determine redaction strategy based on severity and configuration
        redaction_strategy = _select_redaction_strategy(secret, config)
        redacted_text = _apply_redaction_strategy(
            secret["secret_value"], redaction_strategy, config
        )
        
        # Apply redaction
        redacted_content = (
            redacted_content[:start_pos] + 
            redacted_text + 
            redacted_content[end_pos:]
        )
        
        # Record redaction location
        redaction_locations.append({
            "start_position": start_pos,
            "end_position": end_pos,
            "original_text": secret["secret_value"],
            "redacted_text": redacted_text,
            "redaction_reason": f"{secret['secret_type']} detected via {secret['detection_method']}"
        })
        
        # Mark secret as redacted
        secret["redaction_applied"] = True
        secret["redaction_strategy"] = redaction_strategy
        
        # Remove secret value from output for security
        del secret["secret_value"]
    
    return {
        "id": artifact["id"],
        "original_content": content,
        "redacted_content": redacted_content,
        "redaction_applied": len(redaction_locations) > 0,
        "secrets_found": len(detected_secrets),
        "redaction_locations": redaction_locations
    }

def _select_redaction_strategy(secret, config):
    """Select appropriate redaction strategy based on secret properties."""
    
    severity = secret["severity"]
    secret_type = secret["secret_type"]
    
    # Critical secrets get full redaction
    if severity == "critical":
        return "full_redaction"
    
    # High severity secrets get partial redaction if configured
    if severity == "high" and config["redaction_strategies"]["partial_redaction"]:
        return "partial_redaction"
    
    # Structured redaction for certain types
    if (secret_type in ["email", "ip_address", "phone_number"] and 
        config["redaction_strategies"]["structured_redaction"]):
        return "structured_redaction"
    
    return "full_redaction"

def _apply_redaction_strategy(secret_value, strategy, config):
    """Apply the selected redaction strategy to a secret value."""
    
    if strategy == "full_redaction":
        return config["redaction_strategies"]["full_redaction"]
    
    elif strategy == "partial_redaction":
        prefix_chars = config["redaction_strategies"]["preserve_prefix_chars"]
        suffix_chars = config["redaction_strategies"]["preserve_suffix_chars"]
        mask_char = config["redaction_strategies"]["mask_character"]
        
        if len(secret_value) <= prefix_chars + suffix_chars + 2:
            return config["redaction_strategies"]["full_redaction"]
        
        prefix = secret_value[:prefix_chars]
        suffix = secret_value[-suffix_chars:] if suffix_chars > 0 else ""
        mask_length = len(secret_value) - prefix_chars - suffix_chars
        mask = mask_char * mask_length
        
        return f"{prefix}{mask}{suffix}"
    
    elif strategy == "structured_redaction":
        return _apply_structured_redaction(secret_value, config)
    
    else:
        return config["redaction_strategies"]["full_redaction"]

def _apply_structured_redaction(secret_value, config):
    """Apply structured redaction for specific data types."""
    
    # Email: preserve domain
    if "@" in secret_value:
        parts = secret_value.split("@")
        if len(parts) == 2:
            username_mask = "*" * min(8, len(parts[0]))
            return f"{username_mask}@{parts[1]}"
    
    # IP Address: redact last octet
    if re.match(r"^\\d+\\.\\d+\\.\\d+\\.\\d+$", secret_value):
        parts = secret_value.split(".")
        return f"{parts[0]}.{parts[1]}.{parts[2]}.***"
    
    # Phone: redact middle digits
    phone_digits = re.sub(r"\\D", "", secret_value)
    if len(phone_digits) >= 10:
        return f"{phone_digits[:3]}-***-{phone_digits[-4:]}"
    
    # Default to partial redaction
    return _apply_redaction_strategy(secret_value, "partial_redaction", config)
```

---

## 4. Error Handling

### 4.1 Custom Exception Classes

```python
class SecretRedactorError(Exception):
    """Base exception for secret redactor errors."""
    pass

class ValidationError(SecretRedactorError):
    """Raised when input validation fails."""
    def __init__(self, message, details=None):
        super().__init__(message)
        self.details = details or {}

class RedactionError(SecretRedactorError):
    """Raised when critical redaction processing fails."""
    def __init__(self, message, artifact_id=None):
        super().__init__(message)
        self.artifact_id = artifact_id

class PatternError(SecretRedactorError):
    """Raised when pattern compilation or processing fails."""
    def __init__(self, pattern_name, error_details):
        self.pattern_name = pattern_name
        self.error_details = error_details
        super().__init__(f"Pattern error in '{pattern_name}': {error_details}")

class ConfigurationError(SecretRedactorError):
    """Raised when configuration is invalid."""
    pass

class ProcessingError(SecretRedactorError):
    """Raised when unexpected processing errors occur."""
    pass
```

### 4.2 Input Validation Functions

```python
def _validate_input_schema(content_data):
    """Validate input content schema."""
    import jsonschema
    
    try:
        jsonschema.validate(content_data, CONTENT_DATA_SCHEMA)
    except jsonschema.ValidationError as e:
        raise ValidationError(f"Content schema validation failed: {e.message}", {
            "path": list(e.absolute_path),
            "invalid_value": e.instance,
            "constraint": e.validator
        })

def _validate_custom_patterns(custom_patterns):
    """Validate custom pattern definitions."""
    if not custom_patterns:
        return  # Optional parameter
    
    import jsonschema
    import re
    
    try:
        jsonschema.validate(custom_patterns, CUSTOM_PATTERNS_SCHEMA)
    except jsonschema.ValidationError as e:
        raise ValidationError(f"Custom patterns schema validation failed: {e.message}")
    
    # Validate regex patterns can be compiled
    for pattern in custom_patterns.get("patterns", []):
        try:
            re.compile(pattern["pattern"])
        except re.error as e:
            raise PatternError(pattern["name"], f"Invalid regex: {e}")

def _validate_content_size(content, config):
    """Validate content size doesn't exceed limits."""
    max_size_mb = config["performance_limits"]["max_content_size_mb"]
    content_size_mb = len(content.encode('utf-8')) / (1024 * 1024)
    
    if content_size_mb > max_size_mb:
        raise ValidationError(f"Content too large: {content_size_mb:.2f}MB > {max_size_mb}MB")
```

---

## 5. Test Cases

### 5.1 API Key Detection Test Case

```python
def test_api_key_detection():
    """Test detection of various API key formats."""
    content_data = {
        "artifacts": [
            {
                "id": "test-001",
                "type": "configuration",
                "content": '''
                API_KEY=abc123def456ghi789jkl012mno345pqr678
                aws_access_key_id = AKIAIOSFODNN7EXAMPLE
                github_token: ghp_1234567890abcdef1234567890abcdef12345678
                '''
            }
        ]
    }
    
    result = redact_secrets(content_data)
    
    assert result["ok"] == True
    assert result["redaction_summary"]["total_secrets_found"] >= 3
    
    # Check that different secret types were detected
    secret_types = {secret["secret_type"] for secret in result["detected_secrets"]}
    assert "api_key" in secret_types
    assert "aws_key" in secret_types
    assert "github_token" in secret_types
    
    # Verify redaction was applied
    redacted_artifact = result["redacted_artifacts"][0]
    assert redacted_artifact["redaction_applied"] == True
    assert "AKIAIOSFODNN7EXAMPLE" not in redacted_artifact["redacted_content"]
```

### 5.2 Entropy-Based Detection Test Case

```python
def test_entropy_detection():
    """Test entropy-based detection of high-randomness strings."""
    content_data = {
        "artifacts": [
            {
                "id": "test-002",
                "type": "plain_text",
                "content": '''
                Here is some normal text with a random secret: kX9mL2pQ8vN5wR7tY1eE3rT6uI0oP4sD9fG8hJ2kL5nM
                This should be detected as high entropy.
                But this normal sentence should not be flagged.
                '''
            }
        ]
    }
    
    config = {
        "secret_detection": {
            "enable_entropy_analysis": True,
            "entropy_threshold": 4.0
        }
    }
    
    result = redact_secrets(content_data, redaction_config=config)
    
    assert result["ok"] == True
    assert result["redaction_summary"]["total_secrets_found"] >= 1
    
    # Check that entropy detection was used
    detection_methods = {secret["detection_method"] for secret in result["detected_secrets"]}
    assert "entropy_analysis" in detection_methods
```

### 5.3 False Positive Handling Test Case

```python
def test_false_positive_handling():
    """Test handling of common false positives."""
    content_data = {
        "artifacts": [
            {
                "id": "test-003",
                "type": "documentation",
                "content": '''
                # Configuration Examples
                
                API_KEY=your_api_key_here
                SECRET_TOKEN=EXAMPLE_TOKEN_NOT_REAL
                password: placeholder_password
                
                # Real configuration (should be detected)
                REAL_API_KEY=sk_live_51H1234567890abcdef1234567890abcdef12345678
                '''
            }
        ]
    }
    
    result = redact_secrets(content_data)
    
    assert result["ok"] == True
    
    # Check false positive analysis
    fp_analysis = result["false_positive_analysis"]
    assert fp_analysis["total_potential_false_positives"] >= 0
    
    # Real secrets should still be detected
    assert result["redaction_summary"]["total_secrets_found"] >= 1
    
    # Example/placeholder values should not be in high-confidence detections
    high_confidence_secrets = [
        secret for secret in result["detected_secrets"] 
        if secret["confidence_score"] > 0.8
    ]
    
    # Should have at least one high-confidence detection (the real API key)
    assert len(high_confidence_secrets) >= 1
```

### 5.4 Redaction Strategy Test Case

```python
def test_redaction_strategies():
    """Test different redaction strategies based on secret type and severity."""
    content_data = {
        "artifacts": [
            {
                "id": "test-004",
                "type": "configuration",
                "content": '''
                # Critical secrets (should be fully redacted)
                private_key: -----BEGIN PRIVATE KEY-----\\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIB...\\n-----END PRIVATE KEY-----
                
                # High severity (partial redaction)
                api_token: abc123def456ghi789
                
                # Medium severity (structured redaction)
                email: user@company.com
                ip_address: 192.168.1.100
                '''
            }
        ]
    }
    
    config = {
        "redaction_strategies": {
            "partial_redaction": True,
            "structured_redaction": True,
            "preserve_prefix_chars": 3,
            "preserve_suffix_chars": 3
        }
    }
    
    result = redact_secrets(content_data, redaction_config=config)
    
    assert result["ok"] == True
    
    redacted_content = result["redacted_artifacts"][0]["redacted_content"]
    
    # Private key should be fully redacted
    assert "-----BEGIN PRIVATE KEY-----" not in redacted_content
    assert "[REDACTED]" in redacted_content
    
    # Email should use structured redaction (preserve domain)
    assert "@company.com" in redacted_content
    assert "user@company.com" not in redacted_content
```

### 5.5 Custom Pattern Test Case

```python
def test_custom_patterns():
    """Test custom pattern detection."""
    content_data = {
        "artifacts": [
            {
                "id": "test-005",
                "type": "code",
                "content": '''
                # Custom internal secret format
                INTERNAL_AUTH_CODE=COMP_auth_1234567890abcdef
                SERVICE_ID=SVC_98765432109876543210
                '''
            }
        ]
    }
    
    custom_patterns = {
        "patterns": [
            {
                "name": "Internal Auth Code",
                "pattern": r"COMP_auth_[A-Za-z0-9]{16}",
                "severity": "high",
                "description": "Company internal authentication code"
            },
            {
                "name": "Service ID", 
                "pattern": r"SVC_[0-9]{20}",
                "severity": "medium",
                "description": "Internal service identifier"
            }
        ]
    }
    
    result = redact_secrets(content_data, custom_patterns=custom_patterns)
    
    assert result["ok"] == True
    assert result["redaction_summary"]["total_secrets_found"] >= 2
    
    # Check custom pattern detection
    detection_methods = {secret["detection_method"] for secret in result["detected_secrets"]}
    assert "custom_pattern" in detection_methods or "pattern_matching" in detection_methods
    
    # Verify custom secrets were redacted
    redacted_content = result["redacted_artifacts"][0]["redacted_content"]
    assert "COMP_auth_1234567890abcdef" not in redacted_content
    assert "SVC_98765432109876543210" not in redacted_content
```

---

## 6. Performance Requirements

### 6.1 Complexity Requirements

| Operation | Required Complexity | Maximum Runtime (10MB content) |
|-----------|-------------------|------------------------|
| Pattern Matching | O(nm) where n=content, m=patterns | 1000ms |
| Entropy Analysis | O(n) | 500ms |
| Keyword Detection | O(n) | 300ms |
| Redaction Application | O(k) where k=secrets | 100ms |
| Total Pipeline | O(nm) worst case | 2000ms |

### 6.2 Memory Requirements

| Component | Memory Complexity | Maximum Memory (10MB content) |
|-----------|------------------|-------------------------------|
| Content Storage | O(n) | 20MB |
| Pattern Engine | O(m) where m=patterns | 5MB |
| Detection Results | O(k) where k=secrets | 10MB |
| Total Memory Usage | O(n + m + k) | 50MB |

### 6.3 Scalability Limits

```python
DEFAULT_LIMITS = {
    "max_content_size_mb": 10,
    "max_artifacts": 1000,
    "max_patterns": 500,
    "max_processing_time_seconds": 30,
    "max_memory_mb": 100
}
```

---

## 7. Library Dependencies

### 7.1 Required Python Packages

```txt
# Core dependencies
jsonschema>=4.0.0          # JSON schema validation
re                         # Regular expressions (built-in)
hashlib                    # Content hashing (built-in)
math                       # Mathematical functions (built-in)
collections                # Data structures (built-in)

# Optional enhancements
python-Levenshtein>=0.20.0 # String similarity for false positive detection
nltk>=3.7                  # Natural language processing for context analysis

# Development dependencies
pytest>=7.0.0              # Testing framework
pytest-cov>=4.0.0          # Coverage reporting
```

### 7.2 Alternative Implementation Options

```python
# Option 1: Basic implementation (regex-only)
BASIC_DEPENDENCIES = [
    "jsonschema"
]

# Option 2: Enhanced implementation (with NLP)
ENHANCED_DEPENDENCIES = [
    "jsonschema",
    "python-Levenshtein",
    "nltk"
]
```

---

## Summary

The `secret_redactor` tool provides comprehensive, multi-strategy detection and redaction of sensitive information:

- **Pattern Matching** - Extensive library of secret patterns (API keys, tokens, credentials)
- **Entropy Analysis** - Detection of high-randomness strings that may be secrets
- **Keyword Detection** - Context-aware detection using secret-related keywords
- **Smart Redaction** - Multiple redaction strategies based on content type and severity
- **False Positive Mitigation** - Intelligent filtering of common placeholder values

**Key Features:**
1. Multi-strategy detection with confidence scoring and false positive analysis
2. Configurable redaction strategies (full, partial, structured)
3. Comprehensive built-in patterns plus custom pattern support
4. Security-first design with performance optimization
5. Detailed reporting and audit trails
6. Deterministic output through content hashing

This tool is **security-critical** for ensuring no sensitive information leaks into shareable T.A.S.K.S. artifacts, providing a robust last line of defense against data exposure.
