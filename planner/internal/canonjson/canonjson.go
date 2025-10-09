package canonjson

import (
    "bytes"
    "encoding/json"
    "fmt"
    "regexp"
    "sort"
    "strings"
)

// canonicalize recursively sorts keys in maps to ensure a deterministic structure.
func canonicalize(data interface{}) interface{} {
    // Normalize numbers to minimal decimal representation
    if num, ok := data.(json.Number); ok {
        s := string(num)
        if ns, err := canonicalizeNumber(s); err == nil {
            return json.Number(ns)
        }
        return num
    }
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

var numRe = regexp.MustCompile(`^(-)?(\d+)(?:\.(\d+))?(?:[eE]([+-]?\d+))?$`)

// canonicalizeNumber renders a JSON number string in a minimal decimal form:
// - lowercases exponent 'e'
// - trims leading zeros in exponent and integer part (leaves single zero if needed)
// - trims trailing zeros in fraction; removes dot if fraction becomes empty
// - drops exponent if it becomes e0
// - normalizes -0 and -0.0 to 0; preserves sign for non-zero fractional values
func canonicalizeNumber(in string) (string, error) {
    m := numRe.FindStringSubmatch(in)
    if m == nil { // not a standard number, return as-is
        return in, nil
    }
    sign := m[1]
    intPart := m[2]
    frac := m[3]
    exp := m[4]

    // Trim leading zeros on integer part
    intPart = strings.TrimLeft(intPart, "0")
    if intPart == "" {
        intPart = "0"
    }

    // Fraction: trim trailing zeros
    if frac != "" {
        frac = strings.TrimRight(frac, "0")
    }
    hasFrac := frac != ""

    // Exponent: normalize case and zeros
    if exp != "" {
        // remove leading plus
        exp = strings.TrimPrefix(exp, "+")
        // normalize -0 or +0 to 0
        neg := strings.HasPrefix(exp, "-")
        if neg {
            exp = strings.TrimPrefix(exp, "-")
        }
        exp = strings.TrimLeft(exp, "0")
        if exp == "" {
            exp = "0"
            neg = false
        }
        if neg {
            exp = "-" + exp
        }
        // drop exponent if zero
        if exp == "0" {
            exp = ""
        }
    }

    // Normalize negative zero
    allZero := (intPart == "0" && !hasFrac)
    if allZero {
        sign = ""
        exp = ""
    }

    // Recompose
    var b strings.Builder
    b.WriteString(sign)
    b.WriteString(intPart)
    if hasFrac {
        b.WriteString(".")
        b.WriteString(frac)
    }
    if exp != "" {
        b.WriteString("e")
        b.WriteString(exp)
    }
    return b.String(), nil
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
