package canonjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
)

// canonicalize recursively sorts keys in maps to ensure a deterministic structure.
func canonicalize(data interface{}) interface{} {
	if m, ok := data.(map[string]interface{}); ok {
		sortedMap := make(map[string]interface{}, len(m))
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			sortedMap[k] = canonicalize(m[k])
		}
		return sortedMap
	} else if s, ok := data.([]interface{}); ok {
		for i, v := range s {
			s[i] = canonicalize(v)
		}
		return s
	}
	return data
}

// ToCanonicalJSON converts raw JSON bytes into a canonical form as per the v8 spec.
func ToCanonicalJSON(input []byte) ([]byte, error) {
	var data interface{}
	decoder := json.NewDecoder(bytes.NewReader(input))
	decoder.UseNumber()
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode json: %w", err)
	}

	canonicalData := canonicalize(data)

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(canonicalData); err != nil {
		return nil, fmt.Errorf("failed to encode canonical json: %w", err)
	}

	return buf.Bytes(), nil
}
