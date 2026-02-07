package suite

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ParseFile parses a suite file from the given path
func ParseFile(path string) (*Suite, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read suite file: %w", err)
	}

	return Parse(data)
}

// Parse parses a suite from YAML data
func Parse(data []byte) (*Suite, error) {
	var suite Suite
	if err := yaml.Unmarshal(data, &suite); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &suite, nil
}

// BuildDocIndex creates a map of doc_id -> chunk_ids for quick lookup
func BuildDocIndex(suite *Suite) map[string]map[string]bool {
	index := make(map[string]map[string]bool)
	for _, doc := range suite.Docs {
		index[doc.ID] = make(map[string]bool)
		for _, chunk := range doc.Chunks {
			index[doc.ID][chunk.ID] = true
		}
	}
	return index
}

// BuildSchemaIndex creates a map of schema names for quick lookup
func BuildSchemaIndex(suite *Suite) map[string]bool {
	index := make(map[string]bool)
	for name := range suite.Schemas {
		index[name] = true
	}
	return index
}

// BuildToolIndex creates a map of tool names for quick lookup
func BuildToolIndex(suite *Suite) map[string]bool {
	index := make(map[string]bool)
	for _, tool := range suite.Tools {
		index[tool.Name] = true
	}
	return index
}
