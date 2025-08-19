# evidence_validator Tool - Complete Implementation Specification

> **Zero-Questions Implementation Guide for Evidence Validation and Quote Verification**

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
def validate_evidence(
    source_document: str,
    evidence_claims: List[Dict],
    validation_config: Dict = None
) -> Dict:
    """
    Validates that LLM-extracted evidence claims actually exist in source document.
    
    Args:
        source_document: Complete source document text (markdown, plain text, etc.)
        evidence_claims: List of evidence claims to validate (see schema below)
        validation_config: Optional configuration overrides
        
    Returns:
        Complete evidence_validation.json object with validation results
        
    Raises:
        ValidationError: Input schema validation failure
        DocumentParsingError: Source document cannot be parsed
        ConfigurationError: Invalid configuration parameters
        ProcessingError: Unexpected processing failure
    """
```

### 1.2 Command Line Interface

```bash
python -m evidence_validator \
    --source source_document.md \
    --evidence evidence_claims.json \
    --output evidence_validation.json \
    --config validation_config.json \
    --verbose
```

### 1.3 Validation Configuration Schema

```json
{
  "validation_config": {
    "fuzzy_matching": {
      "enabled": true,
      "similarity_threshold": 0.85,
      "max_edit_distance": 3,
      "algorithm": "levenshtein"
    },
    "exact_matching": {
      "enabled": true,
      "case_sensitive": false,
      "normalize_whitespace": true
    },
    "semantic_matching": {
      "enabled": true,
      "chunk_size": 500,
      "overlap_size": 100,
      "similarity_threshold": 0.75
    },
    "context_validation": {
      "enabled": true,
      "context_window": 200,
      "section_awareness": true
    },
    "performance_limits": {
      "max_document_size_mb": 50,
      "max_claims": 1000,
      "max_processing_time_seconds": 120
    }
  }
}
```

---

## 2. Input/Output Schemas

### 2.1 Input Schema: evidence_claims.json

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["claims"],
  "properties": {
    "source_metadata": {
      "type": "object",
      "properties": {
        "filename": {"type": "string"},
        "format": {"type": "string", "enum": ["markdown", "plain_text", "html"]},
        "encoding": {"type": "string", "default": "utf-8"}
      }
    },
    "claims": {
      "type": "array",
      "minItems": 1,
      "maxItems": 1000,
      "items": {
        "type": "object",
        "required": ["id", "task_id", "quote", "evidence_type"],
        "properties": {
          "id": {
            "type": "string",
            "pattern": "^EV\\d+$",
            "description": "Evidence claim ID format: EV001"
          },
          "task_id": {
            "type": "string",
            "pattern": "^P\\d+\\.T\\d+$",
            "description": "Task ID that references this evidence"
          },
          "quote": {
            "type": "string",
            "minLength": 10,
            "maxLength": 2000,
            "description": "The claimed text from source document"
          },
          "evidence_type": {
            "type": "string",
            "enum": ["direct_quote", "paraphrase", "section_reference", "concept_reference"],
            "description": "Type of evidence claim"
          },
          "context_hint": {
            "type": "string",
            "description": "Optional section/paragraph hint for faster matching"
          },
          "expected_section": {
            "type": "string", 
            "description": "Expected document section (header, paragraph number, etc.)"
          },
          "confidence_threshold": {
            "type": "number",
            "minimum": 0.0,
            "maximum": 1.0,
            "default": 0.8,
            "description": "Minimum confidence required for validation"
          }
        }
      }
    }
  }
}
```

### 2.2 Output Schema: evidence_validation.json

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["ok", "generated", "validation_summary"],
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
            "enum": ["VALIDATION_ERROR", "DOCUMENT_PARSING_ERROR", "PROCESSING_ERROR"]
          },
          "message": {"type": "string"},
          "details": {"type": "object"},
          "affected_claims": {
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
          "claim_id": {"type": "string"}
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
    "document_metadata": {
      "type": "object",
      "properties": {
        "size_bytes": {"type": "integer"},
        "line_count": {"type": "integer"},
        "section_count": {"type": "integer"},
        "estimated_reading_time_minutes": {"type": "number"}
      }
    },
    "validation_summary": {
      "type": "object",
      "required": ["total_claims", "validated_claims", "failed_claims"],
      "properties": {
        "total_claims": {"type": "integer", "minimum": 0},
        "validated_claims": {"type": "integer", "minimum": 0},
        "failed_claims": {"type": "integer", "minimum": 0},
        "ambiguous_claims": {"type": "integer", "minimum": 0},
        "average_confidence": {"type": "number", "minimum": 0.0, "maximum": 1.0},
        "validation_rate": {"type": "number", "minimum": 0.0, "maximum": 1.0}
      }
    },
    "validated_claims": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["claim_id", "validation_status", "confidence_score"],
        "properties": {
          "claim_id": {"type": "string"},
          "validation_status": {
            "type": "string",
            "enum": ["VALIDATED", "FAILED", "AMBIGUOUS", "LOW_CONFIDENCE"]
          },
          "confidence_score": {"type": "number", "minimum": 0.0, "maximum": 1.0},
          "match_details": {
            "type": "object",
            "properties": {
              "match_type": {
                "type": "string",
                "enum": ["exact", "fuzzy", "semantic", "contextual"]
              },
              "start_position": {"type": "integer"},
              "end_position": {"type": "integer"},
              "line_number": {"type": "integer"},
              "section": {"type": "string"},
              "matched_text": {"type": "string"},
              "similarity_score": {"type": "number"},
              "edit_distance": {"type": "integer"}
            }
          },
          "alternative_matches": {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "matched_text": {"type": "string"},
                "confidence_score": {"type": "number"},
                "position": {"type": "integer"}
              }
            }
          }
        }
      }
    },
    "failed_claims": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["claim_id", "failure_reason"],
        "properties": {
          "claim_id": {"type": "string"},
          "failure_reason": {
            "type": "string",
            "enum": ["NOT_FOUND", "LOW_CONFIDENCE", "MULTIPLE_MATCHES", "CONTEXT_MISMATCH"]
          },
          "attempted_matches": {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "candidate_text": {"type": "string"},
                "confidence_score": {"type": "number"},
                "rejection_reason": {"type": "string"}
              }
            }
          },
          "suggestions": {
            "type": "array",
            "items": {"type": "string"}
          }
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
def validate_evidence(source_document, evidence_claims, validation_config=None):
    """Complete implementation with all validation strategies."""
    
    # Step 1: Input Validation
    config = _merge_config(validation_config)
    _validate_input_schema(evidence_claims)
    _validate_document_size(source_document, config)
    
    # Step 2: Document Preprocessing
    parsed_document = _parse_document(source_document, config)
    document_chunks = _chunk_document(parsed_document, config)
    
    # Step 3: Build Search Indexes
    exact_index = _build_exact_search_index(parsed_document)
    fuzzy_index = _build_fuzzy_search_index(parsed_document) 
    semantic_index = _build_semantic_index(document_chunks, config)
    
    # Step 4: Validate Each Claim
    validation_results = []
    for claim in evidence_claims['claims']:
        result = _validate_single_claim(
            claim, parsed_document, exact_index, fuzzy_index, 
            semantic_index, config
        )
        validation_results.append(result)
    
    # Step 5: Generate Summary Statistics
    summary = _generate_validation_summary(validation_results)
    
    # Step 6: Generate Content Hash
    content_hash = _generate_content_hash(validation_results, summary)
    
    # Step 7: Build Response
    return _build_validation_response(
        validation_results, summary, content_hash, parsed_document
    )
```

### 3.2 Document Parsing and Preprocessing

```python
def _parse_document(source_document, config):
    """Parse document into structured format with sections and metadata."""
    
    document_info = {
        "raw_text": source_document,
        "lines": source_document.split('\n'),
        "sections": [],
        "metadata": {}
    }
    
    # Detect document format
    if source_document.strip().startswith('#'):
        document_info.update(_parse_markdown(source_document))
    elif '<html' in source_document.lower():
        document_info.update(_parse_html(source_document))
    else:
        document_info.update(_parse_plain_text(source_document))
    
    # Add positional metadata
    _add_positional_metadata(document_info)
    
    return document_info

def _parse_markdown(text):
    """Parse Markdown document into sections."""
    import re
    
    sections = []
    current_section = None
    line_number = 0
    
    for line in text.split('\n'):
        line_number += 1
        
        # Detect headers
        header_match = re.match(r'^(#{1,6})\s+(.+)$', line)
        if header_match:
            # Save previous section
            if current_section:
                sections.append(current_section)
            
            # Start new section
            level = len(header_match.group(1))
            title = header_match.group(2)
            current_section = {
                "type": "section",
                "level": level,
                "title": title,
                "start_line": line_number,
                "content": [],
                "subsections": []
            }
        else:
            # Add to current section
            if current_section:
                current_section["content"].append(line)
            else:
                # Content before first header
                if not sections:
                    sections.append({
                        "type": "preamble",
                        "title": "Document Preamble",
                        "start_line": 1,
                        "content": [],
                        "subsections": []
                    })
                sections[-1]["content"].append(line)
    
    # Add final section
    if current_section:
        sections.append(current_section)
    
    return {"sections": sections, "format": "markdown"}

def _add_positional_metadata(document_info):
    """Add character positions and search-friendly indexes."""
    
    lines = document_info["lines"]
    char_positions = []
    current_pos = 0
    
    for line in lines:
        char_positions.append({
            "start": current_pos,
            "end": current_pos + len(line),
            "text": line
        })
        current_pos += len(line) + 1  # +1 for newline
    
    document_info["char_positions"] = char_positions
    document_info["total_chars"] = current_pos
```

### 3.3 Exact Text Matching

```python
def _build_exact_search_index(parsed_document):
    """Build index for exact text matching."""
    
    text = parsed_document["raw_text"]
    
    # Normalize text if configured
    normalized_text = text.lower()
    normalized_text = re.sub(r'\s+', ' ', normalized_text)  # Normalize whitespace
    
    return {
        "original_text": text,
        "normalized_text": normalized_text,
        "char_positions": parsed_document["char_positions"]
    }

def _exact_match_search(claim_text, exact_index, config):
    """Perform exact text matching with optional normalization."""
    
    search_text = claim_text
    target_text = exact_index["original_text"]
    
    if not config["exact_matching"]["case_sensitive"]:
        search_text = search_text.lower()
        target_text = target_text.lower()
    
    if config["exact_matching"]["normalize_whitespace"]:
        search_text = re.sub(r'\s+', ' ', search_text.strip())
        target_text = re.sub(r'\s+', ' ', target_text)
    
    matches = []
    start = 0
    
    while True:
        pos = target_text.find(search_text, start)
        if pos == -1:
            break
        
        # Find line number
        line_num = _find_line_number(pos, exact_index["char_positions"])
        
        match = {
            "match_type": "exact",
            "start_position": pos,
            "end_position": pos + len(search_text),
            "line_number": line_num,
            "matched_text": exact_index["original_text"][pos:pos + len(search_text)],
            "confidence_score": 1.0,
            "similarity_score": 1.0,
            "edit_distance": 0
        }
        matches.append(match)
        start = pos + 1
    
    return matches

def _find_line_number(char_position, char_positions):
    """Find line number for given character position."""
    for i, line_info in enumerate(char_positions):
        if line_info["start"] <= char_position <= line_info["end"]:
            return i + 1
    return 1
```

### 3.4 Fuzzy Text Matching

```python
def _build_fuzzy_search_index(parsed_document):
    """Build index optimized for fuzzy matching."""
    
    # Create overlapping windows for better fuzzy matching
    text = parsed_document["raw_text"]
    windows = []
    window_size = 200  # characters
    step_size = 100    # characters
    
    for i in range(0, len(text) - window_size + 1, step_size):
        window = {
            "text": text[i:i + window_size],
            "start_pos": i,
            "end_pos": i + window_size
        }
        windows.append(window)
    
    # Add remaining text as final window
    if len(text) > window_size:
        windows.append({
            "text": text[-(window_size//2):],
            "start_pos": len(text) - (window_size//2),
            "end_pos": len(text)
        })
    
    return {
        "windows": windows,
        "full_text": text
    }

def _fuzzy_match_search(claim_text, fuzzy_index, config):
    """Perform fuzzy text matching using edit distance."""
    from difflib import SequenceMatcher
    import Levenshtein  # or implement edit distance
    
    claim_normalized = claim_text.lower().strip()
    similarity_threshold = config["fuzzy_matching"]["similarity_threshold"]
    max_edit_distance = config["fuzzy_matching"]["max_edit_distance"]
    
    matches = []
    
    for window in fuzzy_index["windows"]:
        window_text = window["text"].lower()
        
        # Use sliding window within each text window
        claim_len = len(claim_normalized)
        
        for i in range(len(window_text) - claim_len + 1):
            candidate = window_text[i:i + claim_len]
            
            # Calculate similarity
            if config["fuzzy_matching"]["algorithm"] == "levenshtein":
                edit_dist = Levenshtein.distance(claim_normalized, candidate)
                similarity = 1.0 - (edit_dist / max(len(claim_normalized), len(candidate)))
            else:
                similarity = SequenceMatcher(None, claim_normalized, candidate).ratio()
            
            if similarity >= similarity_threshold and edit_dist <= max_edit_distance:
                absolute_pos = window["start_pos"] + i
                line_num = _find_line_number(absolute_pos, fuzzy_index["char_positions"])
                
                match = {
                    "match_type": "fuzzy",
                    "start_position": absolute_pos,
                    "end_position": absolute_pos + claim_len,
                    "line_number": line_num,
                    "matched_text": fuzzy_index["full_text"][absolute_pos:absolute_pos + claim_len],
                    "confidence_score": similarity,
                    "similarity_score": similarity,
                    "edit_distance": edit_dist
                }
                matches.append(match)
    
    # Sort by confidence and remove duplicates
    matches.sort(key=lambda x: x["confidence_score"], reverse=True)
    return _deduplicate_matches(matches)

def _deduplicate_matches(matches):
    """Remove overlapping matches, keeping highest confidence."""
    if not matches:
        return matches
    
    unique_matches = [matches[0]]
    
    for match in matches[1:]:
        is_duplicate = False
        for existing in unique_matches:
            # Check for significant overlap
            overlap = min(match["end_position"], existing["end_position"]) - \
                     max(match["start_position"], existing["start_position"])
            total_span = max(match["end_position"], existing["end_position"]) - \
                        min(match["start_position"], existing["start_position"])
            
            overlap_ratio = overlap / total_span if total_span > 0 else 0
            
            if overlap_ratio > 0.5:  # 50% overlap threshold
                is_duplicate = True
                break
        
        if not is_duplicate:
            unique_matches.append(match)
    
    return unique_matches
```

### 3.5 Semantic Matching

```python
def _build_semantic_index(document_chunks, config):
    """Build semantic search index using embeddings."""
    # Note: This would typically use sentence-transformers or similar
    # For this spec, we'll outline the approach
    
    if not config["semantic_matching"]["enabled"]:
        return None
    
    chunk_embeddings = []
    
    # In a real implementation, use a model like:
    # from sentence_transformers import SentenceTransformer
    # model = SentenceTransformer('all-MiniLM-L6-v2')
    
    for chunk in document_chunks:
        # embedding = model.encode(chunk["text"])
        # For spec purposes, simulate with TF-IDF
        embedding = _simulate_embedding(chunk["text"])
        
        chunk_embeddings.append({
            "chunk": chunk,
            "embedding": embedding,
            "text": chunk["text"]
        })
    
    return {
        "embeddings": chunk_embeddings,
        "chunk_size": config["semantic_matching"]["chunk_size"]
    }

def _semantic_match_search(claim_text, semantic_index, config):
    """Perform semantic similarity matching."""
    if not semantic_index:
        return []
    
    # claim_embedding = model.encode(claim_text)
    claim_embedding = _simulate_embedding(claim_text)
    similarity_threshold = config["semantic_matching"]["similarity_threshold"]
    
    matches = []
    
    for chunk_data in semantic_index["embeddings"]:
        # similarity = cosine_similarity(claim_embedding, chunk_data["embedding"])
        similarity = _simulate_cosine_similarity(claim_embedding, chunk_data["embedding"])
        
        if similarity >= similarity_threshold:
            chunk = chunk_data["chunk"]
            
            match = {
                "match_type": "semantic",
                "start_position": chunk["start_pos"],
                "end_position": chunk["end_pos"],
                "line_number": chunk.get("line_number", 1),
                "matched_text": chunk["text"][:200] + "..." if len(chunk["text"]) > 200 else chunk["text"],
                "confidence_score": similarity,
                "similarity_score": similarity,
                "edit_distance": -1  # Not applicable for semantic matching
            }
            matches.append(match)
    
    matches.sort(key=lambda x: x["confidence_score"], reverse=True)
    return matches[:5]  # Return top 5 semantic matches

def _simulate_embedding(text):
    """Simulate text embedding for specification purposes."""
    # In real implementation, use proper embeddings
    # This is just for the spec
    words = text.lower().split()
    return hash(tuple(sorted(words))) % 1000000  # Simplified simulation

def _simulate_cosine_similarity(embed1, embed2):
    """Simulate cosine similarity for specification purposes."""
    # Simplified similarity based on hash proximity
    return 1.0 / (1.0 + abs(embed1 - embed2) / 1000000)
```

### 3.6 Claim Validation Logic

```python
def _validate_single_claim(claim, parsed_document, exact_index, fuzzy_index, semantic_index, config):
    """Validate a single evidence claim using all matching strategies."""
    
    claim_id = claim["id"]
    quote = claim["quote"]
    evidence_type = claim["evidence_type"]
    confidence_threshold = claim.get("confidence_threshold", 0.8)
    
    all_matches = []
    
    # Try exact matching first
    if config["exact_matching"]["enabled"]:
        exact_matches = _exact_match_search(quote, exact_index, config)
        all_matches.extend(exact_matches)
    
    # Try fuzzy matching if no exact matches or evidence type allows it
    if (config["fuzzy_matching"]["enabled"] and 
        (not all_matches or evidence_type in ["paraphrase", "concept_reference"])):
        fuzzy_matches = _fuzzy_match_search(quote, fuzzy_index, config)
        all_matches.extend(fuzzy_matches)
    
    # Try semantic matching for paraphrases and concept references
    if (config["semantic_matching"]["enabled"] and semantic_index and
        evidence_type in ["paraphrase", "concept_reference", "section_reference"]):
        semantic_matches = _semantic_match_search(quote, semantic_index, config)
        all_matches.extend(semantic_matches)
    
    # Context validation
    if config["context_validation"]["enabled"] and claim.get("context_hint"):
        all_matches = _filter_by_context(all_matches, claim["context_hint"], parsed_document)
    
    # Select best match and determine validation status
    return _determine_validation_result(claim_id, all_matches, confidence_threshold)

def _filter_by_context(matches, context_hint, parsed_document):
    """Filter matches based on context hints."""
    if not context_hint:
        return matches
    
    context_words = set(context_hint.lower().split())
    filtered_matches = []
    
    for match in matches:
        # Get surrounding context
        start_pos = max(0, match["start_position"] - 200)
        end_pos = min(len(parsed_document["raw_text"]), match["end_position"] + 200)
        context = parsed_document["raw_text"][start_pos:end_pos].lower()
        
        # Check if context contains hint words
        context_words_found = sum(1 for word in context_words if word in context)
        context_score = context_words_found / len(context_words) if context_words else 0
        
        if context_score > 0.3:  # 30% of context words found
            match["context_score"] = context_score
            match["confidence_score"] *= (1 + context_score * 0.2)  # Boost confidence
            filtered_matches.append(match)
    
    return filtered_matches

def _determine_validation_result(claim_id, matches, confidence_threshold):
    """Determine final validation result for a claim."""
    
    if not matches:
        return {
            "claim_id": claim_id,
            "validation_status": "FAILED",
            "failure_reason": "NOT_FOUND",
            "confidence_score": 0.0,
            "attempted_matches": []
        }
    
    # Sort matches by confidence
    matches.sort(key=lambda x: x["confidence_score"], reverse=True)
    best_match = matches[0]
    
    # Check for ambiguous matches (multiple high-confidence matches)
    high_confidence_matches = [m for m in matches if m["confidence_score"] >= confidence_threshold]
    
    if len(high_confidence_matches) > 1:
        # Check if they're in different locations (not duplicates)
        unique_locations = set((m["start_position"], m["end_position"]) for m in high_confidence_matches)
        if len(unique_locations) > 1:
            return {
                "claim_id": claim_id,
                "validation_status": "AMBIGUOUS",
                "confidence_score": best_match["confidence_score"],
                "match_details": best_match,
                "alternative_matches": [
                    {
                        "matched_text": m["matched_text"],
                        "confidence_score": m["confidence_score"],
                        "position": m["start_position"]
                    } for m in high_confidence_matches[1:4]  # Top 3 alternatives
                ]
            }
    
    # Determine final status
    if best_match["confidence_score"] >= confidence_threshold:
        status = "VALIDATED"
    else:
        status = "LOW_CONFIDENCE"
    
    return {
        "claim_id": claim_id,
        "validation_status": status,
        "confidence_score": best_match["confidence_score"],
        "match_details": best_match,
        "alternative_matches": [
            {
                "matched_text": m["matched_text"],
                "confidence_score": m["confidence_score"],
                "position": m["start_position"]
            } for m in matches[1:4] if m["confidence_score"] >= 0.5
        ]
    }
```

---

## 4. Error Handling

### 4.1 Custom Exception Classes

```python
class EvidenceValidatorError(Exception):
    """Base exception for evidence validator errors."""
    pass

class ValidationError(EvidenceValidatorError):
    """Raised when input validation fails."""
    def __init__(self, message, details=None):
        super().__init__(message)
        self.details = details or {}

class DocumentParsingError(EvidenceValidatorError):
    """Raised when source document cannot be parsed."""
    def __init__(self, message, document_format=None):
        super().__init__(message)
        self.document_format = document_format

class ConfigurationError(EvidenceValidatorError):
    """Raised when configuration is invalid."""
    pass

class ProcessingError(EvidenceValidatorError):
    """Raised when unexpected processing errors occur."""
    pass
```

### 4.2 Input Validation Functions

```python
def _validate_input_schema(evidence_claims):
    """Validate input against JSON schema."""
    import jsonschema
    
    try:
        jsonschema.validate(evidence_claims, EVIDENCE_CLAIMS_SCHEMA)
    except jsonschema.ValidationError as e:
        raise ValidationError(f"Input schema validation failed: {e.message}", {
            "path": list(e.absolute_path),
            "invalid_value": e.instance,
            "constraint": e.validator
        })

def _validate_document_size(source_document, config):
    """Validate document size limits."""
    max_size_mb = config["performance_limits"]["max_document_size_mb"]
    doc_size_mb = len(source_document.encode('utf-8')) / (1024 * 1024)
    
    if doc_size_mb > max_size_mb:
        raise ValidationError(f"Document too large: {doc_size_mb:.2f}MB > {max_size_mb}MB")

def _validate_claims_count(evidence_claims, config):
    """Validate number of claims doesn't exceed limits."""
    max_claims = config["performance_limits"]["max_claims"]
    claim_count = len(evidence_claims.get("claims", []))
    
    if claim_count > max_claims:
        raise ValidationError(f"Too many claims: {claim_count} > {max_claims}")
```

---

## 5. Test Cases

### 5.1 Valid Evidence Test Case

```python
def test_exact_match_validation():
    """Test exact quote validation."""
    source_doc = """
    # Project Requirements
    
    The system must implement user authentication with OAuth 2.0 protocol.
    All API endpoints should return JSON responses with appropriate status codes.
    The database schema must support user roles and permissions.
    """
    
    evidence_claims = {
        "claims": [
            {
                "id": "EV001",
                "task_id": "P1.T001",
                "quote": "implement user authentication with OAuth 2.0 protocol",
                "evidence_type": "direct_quote",
                "confidence_threshold": 0.8
            }
        ]
    }
    
    result = validate_evidence(source_doc, evidence_claims)
    
    assert result["ok"] == True
    assert len(result["validated_claims"]) == 1
    assert result["validated_claims"][0]["validation_status"] == "VALIDATED"
    assert result["validated_claims"][0]["confidence_score"] == 1.0
    assert result["validated_claims"][0]["match_details"]["match_type"] == "exact"
```

### 5.2 Fuzzy Matching Test Case

```python
def test_fuzzy_match_validation():
    """Test fuzzy matching for paraphrased quotes."""
    source_doc = """
    The application needs to support user authentication using OAuth 2.0.
    """
    
    evidence_claims = {
        "claims": [
            {
                "id": "EV002", 
                "task_id": "P1.T002",
                "quote": "application must implement OAuth 2.0 authentication",
                "evidence_type": "paraphrase",
                "confidence_threshold": 0.7
            }
        ]
    }
    
    result = validate_evidence(source_doc, evidence_claims)
    
    assert result["ok"] == True
    validated = result["validated_claims"][0]
    assert validated["validation_status"] in ["VALIDATED", "LOW_CONFIDENCE"]
    assert validated["match_details"]["match_type"] == "fuzzy"
    assert validated["confidence_score"] > 0.7
```

### 5.3 Failed Validation Test Case

```python
def test_evidence_not_found():
    """Test handling of evidence not found in document."""
    source_doc = """
    The system will use PostgreSQL database.
    Users can create accounts and login.
    """
    
    evidence_claims = {
        "claims": [
            {
                "id": "EV003",
                "task_id": "P1.T003", 
                "quote": "implement microservices architecture",
                "evidence_type": "direct_quote",
                "confidence_threshold": 0.8
            }
        ]
    }
    
    result = validate_evidence(source_doc, evidence_claims)
    
    assert result["ok"] == True  # Process completes successfully
    assert len(result["failed_claims"]) == 1
    assert result["failed_claims"][0]["failure_reason"] == "NOT_FOUND"
    assert result["validation_summary"]["validation_rate"] == 0.0
```

### 5.4 Ambiguous Matches Test Case

```python
def test_ambiguous_matches():
    """Test handling of multiple potential matches."""
    source_doc = """
    # Section 1
    Users need to authenticate with the system.
    
    # Section 2  
    The authentication system must be secure.
    
    # Section 3
    Authentication should use industry standards.
    """
    
    evidence_claims = {
        "claims": [
            {
                "id": "EV004",
                "task_id": "P1.T004",
                "quote": "authentication",
                "evidence_type": "concept_reference",
                "confidence_threshold": 0.8
            }
        ]
    }
    
    result = validate_evidence(source_doc, evidence_claims)
    
    assert result["ok"] == True
    validated = result["validated_claims"][0]
    assert validated["validation_status"] == "AMBIGUOUS"
    assert len(validated["alternative_matches"]) > 1
```

---

## 6. Performance Requirements

### 6.1 Complexity Requirements

| Operation | Required Complexity | Maximum Runtime (1000 claims, 1MB doc) |
|-----------|-------------------|---------------------------|
| Document Parsing | O(n) | 100ms |
| Exact Search | O(nm) worst case | 500ms |
| Fuzzy Search | O(n²m) worst case | 2000ms |
| Semantic Search | O(km) where k=chunks | 1000ms |
| Total Pipeline | O(n²m) worst case | 3000ms |

### 6.2 Memory Requirements

| Component | Memory Complexity | Maximum Memory (1MB doc) |
|-----------|------------------|--------------------------|
| Document Storage | O(n) | 2MB |
| Search Indexes | O(n) | 5MB |
| Embeddings | O(kd) where d=dimensions | 10MB |
| Total Memory Usage | O(n + kd) | 20MB |

### 6.3 Scalability Limits

```python
DEFAULT_LIMITS = {
    "max_document_size_mb": 50,
    "max_claims": 1000,
    "max_processing_time_seconds": 120,
    "max_memory_mb": 100
}
```

---

## 7. Library Dependencies

### 7.1 Required Python Packages

```txt
# Core dependencies
jsonschema>=4.0.0          # JSON schema validation
python-Levenshtein>=0.20.0 # Edit distance calculations
difflib                    # Built-in sequence matching
re                         # Regular expressions (built-in)
hashlib                    # Content hashing (built-in)

# Document processing
markdown>=3.4.0            # Markdown parsing
beautifulsoup4>=4.11.0     # HTML parsing (if needed)
nltk>=3.7                  # Text processing utilities

# Semantic search (optional)
sentence-transformers>=2.2.0  # Embeddings
numpy>=1.21.0              # Numerical operations
scikit-learn>=1.1.0        # Cosine similarity

# Development dependencies
pytest>=7.0.0              # Testing framework
pytest-cov>=4.0.0          # Coverage reporting
```

### 7.2 Alternative Implementation Options

```python
# Option 1: Basic implementation (no semantic search)
BASIC_DEPENDENCIES = [
    "jsonschema",
    "python-Levenshtein", 
    "markdown"
]

# Option 2: Full implementation with semantic search
FULL_DEPENDENCIES = [
    "jsonschema",
    "python-Levenshtein",
    "sentence-transformers",
    "numpy",
    "scikit-learn",
    "markdown",
    "beautifulsoup4"
]
```

---

## Summary

The `evidence_validator` tool provides comprehensive validation of LLM-extracted evidence claims against source documents using multiple matching strategies:

- **Exact matching** for direct quotes with normalization options
- **Fuzzy matching** for paraphrases using edit distance algorithms  
- **Semantic matching** for concept references using embeddings
- **Context validation** using surrounding text analysis

**Key Features:**
1. Multi-strategy validation with confidence scoring
2. Ambiguous match detection and alternative suggestions
3. Performance optimization for large documents
4. Comprehensive error handling and detailed feedback
5. Deterministic output through content hashing

This tool is **critical** for preventing hallucination in the T.A.S.K.S. pipeline by ensuring all evidence claims are grounded in actual source content.
