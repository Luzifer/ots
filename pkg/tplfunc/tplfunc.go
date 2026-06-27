// Package tplfunc provides common Go template helper functions.
package tplfunc

import (
	"encoding/json"
	"fmt"
	"maps"
	"strings"
	"text/template"
)

// FuncMap returns template helper functions for Sprig-compatible templates.
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"dict":               dict,
		"list":               list,
		"mustToJson":         mustToJSON,
		"mustMergeOverwrite": mustMergeOverwrite,
		"replace":            replace,
	}
}

func dict(v ...any) map[string]any {
	dict := make(map[string]any)

	for len(v) > 0 {
		key := strval(v[0])
		if len(v) == 1 {
			dict[key] = ""
			break
		}
		dict[key] = v[1]
		v = v[2:]
	}
	return dict
}

func list(args ...any) []any { return args }

func mustMergeOverwrite(srcs ...map[string]any) (out map[string]any) {
	out = make(map[string]any)

	for _, src := range srcs {
		maps.Insert(out, maps.All(src))
	}

	return out
}

func mustToJSON(v any) (string, error) {
	output, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("marshalling data: %w", err)
	}
	return string(output), nil
}

func replace(oldSubstr, newSubstr, src string) string {
	return strings.ReplaceAll(src, oldSubstr, newSubstr)
}

func strval(v any) string {
	switch v := v.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case error:
		return v.Error()
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}
