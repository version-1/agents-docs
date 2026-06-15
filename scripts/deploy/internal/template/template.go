package template

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	varPattern           = regexp.MustCompile(`\{\{([a-zA-Z0-9_]+(?:\.[a-zA-Z0-9_]+)*)\}\}`)
	trailingCommaPattern = regexp.MustCompile(`,\s*([\]}])`)
)

// Vars holds the parsed local config values as a nested map.
type Vars map[string]any

// LoadVars reads a JSON file and returns the parsed variables.
// Trailing commas in arrays and objects are tolerated.
func LoadVars(path string) (Vars, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read local config %q: %w", path, err)
	}

	cleaned := trailingCommaPattern.ReplaceAll(b, []byte("$1"))

	var v Vars
	if err := json.Unmarshal(cleaned, &v); err != nil {
		return nil, fmt.Errorf("parse local config %q: %w", path, err)
	}
	return v, nil
}

// Expand replaces all {{key}} and {{a.b.c}} patterns in content with values
// from vars. Returns an error if any variable cannot be resolved.
func Expand(content string, vars Vars) (string, error) {
	var expandErr error

	result := varPattern.ReplaceAllStringFunc(content, func(match string) string {
		if expandErr != nil {
			return match
		}

		key := varPattern.FindStringSubmatch(match)[1]
		val, err := resolve(vars, key)
		if err != nil {
			expandErr = err
			return match
		}
		return val
	})

	if expandErr != nil {
		return "", expandErr
	}
	return result, nil
}

// resolve traverses the nested map using a dot-separated key path.
func resolve(vars Vars, key string) (string, error) {
	parts := strings.Split(key, ".")
	var current any = map[string]any(vars)

	for i, part := range parts {
		m, ok := current.(map[string]any)
		if !ok {
			traversed := strings.Join(parts[:i], ".")
			return "", fmt.Errorf("template variable {{%s}}: %q is not an object", key, traversed)
		}

		val, exists := m[part]
		if !exists {
			return "", fmt.Errorf("template variable {{%s}}: key %q not found", key, part)
		}
		current = val
	}

	return toTOMLValue(current, key)
}

// toTOMLValue converts a JSON-decoded value to its TOML string representation.
// Supported types: string, float64 (number), bool, []any (flat array of primitives).
// Objects (map[string]any) are not supported as leaf values.
func toTOMLValue(v any, key string) (string, error) {
	switch val := v.(type) {
	case string:
		return val, nil
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val)), nil
		}
		return fmt.Sprintf("%g", val), nil
	case bool:
		if val {
			return "true", nil
		}
		return "false", nil
	case []any:
		elems := make([]string, 0, len(val))
		for i, elem := range val {
			s, err := toTOMLPrimitive(elem, key, i)
			if err != nil {
				return "", err
			}
			elems = append(elems, s)
		}
		return "[" + strings.Join(elems, ", ") + "]", nil
	case map[string]any:
		return "", fmt.Errorf("template variable {{%s}}: objects cannot be used as values", key)
	default:
		return "", fmt.Errorf("template variable {{%s}}: unsupported type %T", key, v)
	}
}

// toTOMLPrimitive converts a single array element to its TOML representation.
func toTOMLPrimitive(v any, key string, index int) (string, error) {
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("%q", val), nil
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val)), nil
		}
		return fmt.Sprintf("%g", val), nil
	case bool:
		if val {
			return "true", nil
		}
		return "false", nil
	default:
		return "", fmt.Errorf("template variable {{%s}}[%d]: unsupported array element type %T", key, index, v)
	}
}

// ExpandFile reads a file, expands template variables, and returns the result.
func ExpandFile(path string, vars Vars) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read template file %q: %w", path, err)
	}

	expanded, err := Expand(string(b), vars)
	if err != nil {
		return nil, fmt.Errorf("expand %q: %w", path, err)
	}
	return []byte(expanded), nil
}
