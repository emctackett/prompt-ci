package suite

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ValidationError represents a suite validation error
type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("suite validation failed with %d errors:\n  - %s",
		len(e.Errors), strings.Join(e.Errors, "\n  - "))
}

// ValidateSuite validates the suite structure and returns any errors
func ValidateSuite(suite *Suite, fixturesDir string) error {
	var errors []string

	// Check suite name
	if suite.Name == "" {
		errors = append(errors, "suite_name is required")
	}

	// Check cases exist
	if len(suite.Cases) == 0 {
		errors = append(errors, "cases array must contain at least one test case")
	}

	// Build indices for lookups
	docIndex := BuildDocIndex(suite)
	schemaIndex := BuildSchemaIndex(suite)
	toolIndex := BuildToolIndex(suite)

	// Track case IDs for uniqueness
	caseIDs := make(map[string]bool)
	caseIDPattern := regexp.MustCompile(`^[a-z0-9_]{1,32}$`)

	for i, c := range suite.Cases {
		// Validate case ID format
		if !caseIDPattern.MatchString(c.ID) {
			errors = append(errors, fmt.Sprintf("case[%d]: id '%s' does not match pattern ^[a-z0-9_]{1,32}$", i, c.ID))
		}

		// Check for duplicate IDs
		if caseIDs[c.ID] {
			errors = append(errors, fmt.Sprintf("case[%d]: duplicate case id '%s'", i, c.ID))
		}
		caseIDs[c.ID] = true

		// Check assertions exist
		if len(c.Assertions) == 0 {
			errors = append(errors, fmt.Sprintf("case[%d] '%s': must have at least one assertion", i, c.ID))
		}

		// Validate each assertion
		for j, a := range c.Assertions {
			if err := validateAssertion(a, i, j, c.ID, schemaIndex, toolIndex); err != nil {
				errors = append(errors, err.Error())
			}
		}

		// Check fixture exists
		fixturePath := GetFixturePath(fixturesDir, c.ID)
		if _, err := os.Stat(fixturePath); os.IsNotExist(err) {
			errors = append(errors, fmt.Sprintf("case[%d] '%s': fixture file not found at %s", i, c.ID, fixturePath))
		}
	}

	// Validate grounding config if present
	if suite.Grounding.CitationFormat != "" {
		for _, docID := range suite.Grounding.ValidDocIDs {
			if _, exists := docIndex[docID]; !exists {
				errors = append(errors, fmt.Sprintf("grounding: valid_doc_ids references non-existent doc '%s'", docID))
			}
		}
	}

	// Validate json_schema assertions have additionalProperties: false
	for i, c := range suite.Cases {
		for j, a := range c.Assertions {
			if a.Type == "json_schema" {
				if err := validateSchemaHasAdditionalPropertiesFalse(a.Expected, i, j, c.ID); err != nil {
					errors = append(errors, err.Error())
				}
			}
		}
	}

	if len(errors) > 0 {
		return &ValidationError{Errors: errors}
	}

	return nil
}

func validateAssertion(a Assertion, caseIdx, assertIdx int, caseID string, schemaIndex, toolIndex map[string]bool) error {
	validTypes := map[string]bool{
		"exact_match":         true,
		"contains":            true,
		"regex":               true,
		"json_schema":         true,
		"semantic_similarity": true,
		"llm_judge":           true,
	}

	if !validTypes[a.Type] {
		return fmt.Errorf("case[%d] '%s' assertion[%d]: unknown type '%s'", caseIdx, caseID, assertIdx, a.Type)
	}

	if a.Expected == nil {
		return fmt.Errorf("case[%d] '%s' assertion[%d]: expected is required", caseIdx, caseID, assertIdx)
	}

	return nil
}

func validateSchemaHasAdditionalPropertiesFalse(expected interface{}, caseIdx, assertIdx int, caseID string) error {
	schema, ok := expected.(map[string]interface{})
	if !ok {
		return nil // Not a map, can't validate structure
	}

	// Check if this is an object type
	if schemaType, exists := schema["type"]; exists {
		if schemaType == "object" {
			// Check additionalProperties
			if addProps, exists := schema["additionalProperties"]; exists {
				if addProps != false {
					return fmt.Errorf("case[%d] '%s' assertion[%d]: json_schema must have additionalProperties: false", caseIdx, caseID, assertIdx)
				}
			} else {
				return fmt.Errorf("case[%d] '%s' assertion[%d]: json_schema must have additionalProperties: false", caseIdx, caseID, assertIdx)
			}
		}
	}

	return nil
}

// GetFixturePath returns the expected fixture path for a case ID
func GetFixturePath(fixturesDir, caseID string) string {
	caseType := GetCaseType(caseID)
	var subdir string
	var ext string

	switch caseType {
	case CaseTypeGrounding:
		subdir = "grounding"
		ext = ".out.txt"
	case CaseTypeSchema:
		subdir = "schema"
		ext = ".out.json"
	case CaseTypeTool:
		subdir = "tool"
		// Tool cases can be .json or .txt depending on the case
		if strings.HasSuffix(caseID, "_secret") || strings.HasSuffix(caseID, "_behavior") {
			ext = ".out.txt"
		} else {
			ext = ".out.json"
		}
	default:
		subdir = "grounding"
		ext = ".out.txt"
	}

	return filepath.Join(fixturesDir, subdir, caseID+ext)
}
