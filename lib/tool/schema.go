package tool

import (
	"encoding/json"

	"github.com/google/jsonschema-go/jsonschema"
)

// SchemaFromStruct generates a JSON schema from a Go struct
// and converts it to map[string]any for compatibility.
func SchemaFromStruct[T any]() (map[string]any, error) {
	schema, err := jsonschema.For[T](nil)
	if err != nil {
		return nil, err
	}
	return schemaToMap(schema)
}

// schemaToMap converts a jsonschema.Schema to map[string]any
// and post-processes it to fix type unions that include null.
func schemaToMap(s *jsonschema.Schema) (map[string]any, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	// Fix type unions that include null (e.g., ["null", "array"] -> "array")
	fixNullableTypes(result)
	return result, nil
}

// fixNullableTypes recursively fixes type unions that include null
// and removes additionalProperties constraints.
func fixNullableTypes(schema map[string]any) {
	// Remove additionalProperties - not all APIs support strict mode
	delete(schema, "additionalProperties")

	// Check if type is an array (type union)
	if typeVal, ok := schema["type"]; ok {
		if typeArray, ok := typeVal.([]any); ok {
			// Filter out "null" from the type array
			filtered := make([]any, 0, len(typeArray))
			for _, t := range typeArray {
				if str, ok := t.(string); ok && str != "null" {
					filtered = append(filtered, str)
				}
			}
			// If only one type remains, set it directly
			if len(filtered) == 1 {
				schema["type"] = filtered[0]
			} else if len(filtered) > 1 {
				schema["type"] = filtered
			}
			// If filtered is empty, keep the original (edge case)
		} else if typeStr, ok := typeVal.(string); ok && typeStr == "object" {
			// Ensure object types have a properties field
			if _, hasProps := schema["properties"]; !hasProps {
				schema["properties"] = make(map[string]any)
			}
		}
	}

	// Recursively process properties
	if props, ok := schema["properties"].(map[string]any); ok {
		for _, prop := range props {
			if propMap, ok := prop.(map[string]any); ok {
				fixNullableTypes(propMap)
			}
		}
	}

	// Recursively process items (for array types)
	if items, ok := schema["items"].(map[string]any); ok {
		fixNullableTypes(items)
	}
}
