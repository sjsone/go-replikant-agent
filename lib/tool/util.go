package tool

import "encoding/json"

// ParseArgsToMap parses JSON arguments string into a map.
func ParseArgsToMap(argsStr string) map[string]any {
	if argsStr == "" {
		return make(map[string]any)
	}
	var result map[string]any
	if err := json.Unmarshal([]byte(argsStr), &result); err != nil {
		return make(map[string]any)
	}
	return result
}
