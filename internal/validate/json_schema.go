package validate

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

// ValidateJSONSchema validates content against a JSON schema
func ValidateJSONSchema(content string, expected interface{}) (bool, string) {
	// First, try to extract JSON from the content
	jsonContent := extractJSON(content)
	if jsonContent == "" {
		return false, "no valid JSON found in content"
	}

	// Parse the content as JSON
	var contentData interface{}
	if err := json.Unmarshal([]byte(jsonContent), &contentData); err != nil {
		return false, fmt.Sprintf("content is not valid JSON: %v", err)
	}

	// Convert expected schema to JSON
	schemaBytes, err := json.Marshal(expected)
	if err != nil {
		return false, fmt.Sprintf("failed to marshal schema: %v", err)
	}

	// Compile the schema
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", strings.NewReader(string(schemaBytes))); err != nil {
		return false, fmt.Sprintf("failed to add schema resource: %v", err)
	}

	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return false, fmt.Sprintf("failed to compile schema: %v", err)
	}

	// Validate content against schema
	if err := schema.Validate(contentData); err != nil {
		return false, fmt.Sprintf("JSON schema validation failed: %v", err)
	}

	return true, ""
}

// extractJSON attempts to extract a JSON object or array from content
// This handles cases where JSON might be embedded in other text
func extractJSON(content string) string {
	content = strings.TrimSpace(content)

	// Try to find JSON object
	start := strings.Index(content, "{")
	if start != -1 {
		// Find matching closing brace
		depth := 0
		for i := start; i < len(content); i++ {
			switch content[i] {
			case '{':
				depth++
			case '}':
				depth--
				if depth == 0 {
					return content[start : i+1]
				}
			}
		}
	}

	// Try to find JSON array
	start = strings.Index(content, "[")
	if start != -1 {
		depth := 0
		for i := start; i < len(content); i++ {
			switch content[i] {
			case '[':
				depth++
			case ']':
				depth--
				if depth == 0 {
					return content[start : i+1]
				}
			}
		}
	}

	// Return content as-is if no JSON structure found
	return content
}
