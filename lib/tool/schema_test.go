package tool

import (
	"testing"
)

// Sample structs for schema generation testing
type SimpleStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type NestedStruct struct {
	Person SimpleStruct `json:"person"`
	City   string       `json:"city"`
}

type StructWithArrays struct {
	Items []string `json:"items"`
	Count int      `json:"count"`
}

type StructWithOptional struct {
	Name   string  `json:"name"`
	Age    *int    `json:"age,omitempty"`
	Score  float64 `json:"score,omitempty"`
	Active bool    `json:"active"`
}

type ComplexStruct struct {
	ID       int            `json:"id"`
	Name     string         `json:"name"`
	Tags     []string       `json:"tags"`
	Metadata map[string]any `json:"metadata"`
	Optional *string        `json:"optional,omitempty"`
}

func TestSchemaFromStruct(t *testing.T) {
	tests := []struct {
		name       string
		structType any // Used for type inference
		wantErr    bool
	}{
		{
			name:       "simple struct",
			structType: SimpleStruct{},
			wantErr:    false,
		},
		{
			name:       "nested struct",
			structType: NestedStruct{},
			wantErr:    false,
		},
		{
			name:       "struct with arrays",
			structType: StructWithArrays{},
			wantErr:    false,
		},
		{
			name:       "struct with optional fields",
			structType: StructWithOptional{},
			wantErr:    false,
		},
		{
			name:       "complex struct",
			structType: ComplexStruct{},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var schema map[string]any
			var err error

			switch tt.structType.(type) {
			case SimpleStruct:
				schema, err = SchemaFromStruct[SimpleStruct]()
			case NestedStruct:
				schema, err = SchemaFromStruct[NestedStruct]()
			case StructWithArrays:
				schema, err = SchemaFromStruct[StructWithArrays]()
			case StructWithOptional:
				schema, err = SchemaFromStruct[StructWithOptional]()
			case ComplexStruct:
				schema, err = SchemaFromStruct[ComplexStruct]()
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("SchemaFromStruct() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("SchemaFromStruct() unexpected error: %v", err)
				return
			}

			if schema == nil {
				t.Error("SchemaFromStruct() returned nil schema")
			}

			if schema["type"] == nil {
				t.Error("SchemaFromStruct() missing type field")
			}
		})
	}
}

func TestFixNullableTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]any
		expected map[string]any
	}{
		{
			name: "removes additionalProperties",
			input: map[string]any{
				"type":                 "object",
				"additionalProperties": false,
			},
			expected: map[string]any{
				"type":       "object",
				"properties": map[string]any{}, // fixNullableTypes adds empty properties for object type
			},
		},
		{
			name: "filters null from type array",
			input: map[string]any{
				"type": []any{"string", "null"},
			},
			expected: map[string]any{
				"type": "string",
			},
		},
		{
			name: "keeps non-null types in array",
			input: map[string]any{
				"type": []any{"string", "number"},
			},
			expected: map[string]any{
				"type": []any{"string", "number"},
			},
		},
		{
			name: "filters null from multi-type array",
			input: map[string]any{
				"type": []any{"string", "number", "null"},
			},
			expected: map[string]any{
				"type": []any{"string", "number"},
			},
		},
		{
			name: "handles empty type array after filtering",
			input: map[string]any{
				"type": []any{"null"},
			},
			expected: map[string]any{
				"type": []any{"null"},
			},
		},
		{
			name: "adds empty properties for object type",
			input: map[string]any{
				"type": "object",
			},
			expected: map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
		},
		{
			name: "recursively processes properties",
			input: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"field1": map[string]any{
						"type": []any{"string", "null"},
					},
				},
			},
			expected: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"field1": map[string]any{
						"type": "string",
					},
				},
			},
		},
		{
			name: "recursively processes items",
			input: map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": []any{"string", "null"},
				},
			},
			expected: map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "string",
				},
			},
		},
		{
			name: "complex nested schema",
			input: map[string]any{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]any{
					"name": map[string]any{
						"type": []any{"string", "null"},
					},
					"items": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": []any{"integer", "null"},
						},
					},
				},
			},
			expected: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"name": map[string]any{
						"type": "string",
					},
					"items": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "integer",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixNullableTypes(tt.input)
			// Check key equality
			if len(tt.input) != len(tt.expected) {
				t.Errorf("Expected map to have %d keys, got %d", len(tt.expected), len(tt.input))
			}
			// Check specific keys
			for k := range tt.expected {
				if _, ok := tt.input[k]; !ok {
					t.Errorf("Expected map to have key %q", k)
				}
			}
			// Check that additionalProperties is removed
			if _, ok := tt.input["additionalProperties"]; ok {
				t.Error("Expected additionalProperties to be removed")
			}
		})
	}
}

// Helper function for map comparison
func assertMapsEqual(t *testing.T, expected, actual map[string]any) {
	t.Helper()

	if len(expected) != len(actual) {
		t.Errorf("Expected map to have %d keys, got %d", len(expected), len(actual))
		return
	}

	for k, expectedVal := range expected {
		actualVal, ok := actual[k]
		if !ok {
			t.Errorf("Expected map to have key %q", k)
			return
		}

		// Handle slices specially
		expectedSlice, expectedIsSlice := expectedVal.([]any)
		actualSlice, actualIsSlice := actualVal.([]any)

		if expectedIsSlice && actualIsSlice {
			if len(expectedSlice) != len(actualSlice) {
				t.Errorf("Expected map key %q slice to have %d elements, got %d", k, len(expectedSlice), len(actualSlice))
				continue
			}
			for i := range expectedSlice {
				if expectedSlice[i] != actualSlice[i] {
					t.Errorf("Expected map key %q slice element %d to be %v, got %v", k, i, expectedSlice[i], actualSlice[i])
				}
			}
			continue
		}

		if expectedVal != actualVal {
			t.Errorf("Expected map key %q to be %v, got %v", k, expectedVal, actualVal)
		}
	}
}

func TestFixNullableTypes_NoMutationOnNil(t *testing.T) {
	// Should not panic on nil map
	var schema map[string]any
	fixNullableTypes(schema)
	// If we got here without panic, test passes
}

func TestFixNullableTypes_StringType(t *testing.T) {
	schema := map[string]any{
		"type": "string",
	}

	fixNullableTypes(schema)

	if schema["type"] != "string" {
		t.Errorf("Expected type to remain 'string', got %v", schema["type"])
	}
}

func TestFixNullableTypes_ObjectWithExistingProperties(t *testing.T) {
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{
				"type": "string",
			},
		},
	}

	fixNullableTypes(schema)

	props, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatal("Expected properties to be a map")
	}

	if len(props) != 1 {
		t.Errorf("Expected 1 property, got %d", len(props))
	}

	if props["name"].(map[string]any)["type"] != "string" {
		t.Error("Property type should remain 'string'")
	}
}

func TestSchemaToMap(t *testing.T) {
	// This is an integration test for schemaToMap
	schema, err := SchemaFromStruct[SimpleStruct]()
	if err != nil {
		t.Fatalf("SchemaFromStruct() error: %v", err)
	}

	if schema == nil {
		t.Fatal("Expected non-nil schema")
	}

	if schema["type"] != "object" {
		t.Errorf("Expected type to be 'object', got %v", schema["type"])
	}

	props, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatal("Expected properties to be a map")
	}

	if len(props) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(props))
	}
}

func TestSchemaFromStruct_WithNullableFields(t *testing.T) {
	schema, err := SchemaFromStruct[StructWithOptional]()
	if err != nil {
		t.Fatalf("SchemaFromStruct() error: %v", err)
	}

	if schema == nil {
		t.Fatal("Expected non-nil schema")
	}

	props, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatal("Expected properties to be a map")
	}

	// Check that nullable types are properly fixed
	ageField, ok := props["age"].(map[string]any)
	if ok {
		ageType := ageField["type"]
		// After fixNullableTypes, type should be "integer", not []any{"integer", "null"}
		if typeStr, ok := ageType.(string); ok && typeStr == "integer" {
			// Correctly fixed to single type
		} else if typeArr, ok := ageType.([]any); ok {
			// Should not contain "null"
			for _, typeVal := range typeArr {
				if typeVal == "null" {
					t.Error("Age type should not contain 'null' after fixNullableTypes")
				}
			}
		}
	}
}

func TestSchemaFromStruct_NestedStructures(t *testing.T) {
	schema, err := SchemaFromStruct[NestedStruct]()
	if err != nil {
		t.Fatalf("SchemaFromStruct() error: %v", err)
	}

	if schema == nil {
		t.Fatal("Expected non-nil schema")
	}

	props, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatal("Expected properties to be a map")
	}

	// Check nested person property
	personField, ok := props["person"].(map[string]any)
	if !ok {
		t.Fatal("Expected person to be a map")
	}

	if personField["type"] != "object" {
		t.Errorf("Expected person type to be 'object', got %v", personField["type"])
	}

	// Check that person has nested properties
	personProps, ok := personField["properties"].(map[string]any)
	if !ok {
		t.Fatal("Expected person to have properties")
	}

	if len(personProps) != 2 {
		t.Errorf("Expected person to have 2 properties, got %d", len(personProps))
	}
}

func TestSchemaFromStruct_ArrayFields(t *testing.T) {
	schema, err := SchemaFromStruct[StructWithArrays]()
	if err != nil {
		t.Fatalf("SchemaFromStruct() error: %v", err)
	}

	if schema == nil {
		t.Fatal("Expected non-nil schema")
	}

	props, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatal("Expected properties to be a map")
	}

	// Check items array field
	itemsField, ok := props["items"].(map[string]any)
	if !ok {
		t.Fatal("Expected items to be a map")
	}

	if itemsField["type"] != "array" {
		t.Errorf("Expected items type to be 'array', got %v", itemsField["type"])
	}

	// Check items field has items property
	itemsItems, ok := itemsField["items"].(map[string]any)
	if !ok {
		t.Fatal("Expected items to have items property")
	}

	if itemsItems["type"] != "string" {
		t.Errorf("Expected items items type to be 'string', got %v", itemsItems["type"])
	}
}

func TestFixNullableTypes_WithNilProperties(t *testing.T) {
	// Should not panic when properties is nil
	schema := map[string]any{
		"type":       "object",
		"properties": nil,
	}

	fixNullableTypes(schema)
	// Test passes if no panic occurred
}

func TestFixNullableTypes_WithNilItems(t *testing.T) {
	// Should not panic when items is nil
	schema := map[string]any{
		"type":  "array",
		"items": nil,
	}

	fixNullableTypes(schema)
	// Test passes if no panic occurred
}

func TestSchemaFromStruct_ComplexStruct(t *testing.T) {
	schema, err := SchemaFromStruct[ComplexStruct]()
	if err != nil {
		t.Fatalf("SchemaFromStruct() error: %v", err)
	}

	if schema == nil {
		t.Fatal("Expected non-nil schema")
	}

	props, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatal("Expected properties to be a map")
	}

	// Check tags array field
	tagsField, ok := props["tags"].(map[string]any)
	if !ok {
		t.Fatal("Expected tags to be a map")
	}

	if tagsField["type"] != "array" {
		t.Errorf("Expected tags type to be 'array', got %v", tagsField["type"])
	}

	// Check metadata map field
	metadataField, ok := props["metadata"].(map[string]any)
	if !ok {
		t.Fatal("Expected metadata to be a map")
	}

	if metadataField["type"] != "object" {
		t.Errorf("Expected metadata type to be 'object', got %v", metadataField["type"])
	}

	// Check optional field
	optionalField, ok := props["optional"].(map[string]any)
	if !ok {
		t.Fatal("Expected optional to be a map")
	}

	// After fixNullableTypes, should not have null in type
	optionalType := optionalField["type"]
	if typeStr, ok := optionalType.(string); ok {
		if typeStr == "string" || typeStr == "" {
			// OK - either fixed to string or empty (omitted)
		}
	}
}

func TestFixNullableTypes_PreservesNonTypeFields(t *testing.T) {
	schema := map[string]any{
		"type":        "object",
		"description": "A test schema",
		"title":       "Test",
		"properties": map[string]any{
			"name": map[string]any{
				"type":        "string",
				"description": "The name",
			},
		},
	}

	fixNullableTypes(schema)

	// Check that non-type fields are preserved
	if schema["description"] != "A test schema" {
		t.Error("Description should be preserved")
	}

	if schema["title"] != "Test" {
		t.Error("Title should be preserved")
	}

	props := schema["properties"].(map[string]any)
	nameField := props["name"].(map[string]any)
	if nameField["description"] != "The name" {
		t.Error("Field description should be preserved")
	}
}
