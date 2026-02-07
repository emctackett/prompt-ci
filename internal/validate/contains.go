package validate

import (
	"fmt"
	"strings"
)

// ValidateContains checks if content contains the expected string
func ValidateContains(content string, expected interface{}) (bool, string) {
	expectedStr, ok := expected.(string)
	if !ok {
		return false, fmt.Sprintf("expected must be a string, got %T", expected)
	}

	if strings.Contains(content, expectedStr) {
		return true, ""
	}

	return false, fmt.Sprintf("content does not contain expected string '%s'", expectedStr)
}

// ValidateExactMatch checks if content exactly matches expected
func ValidateExactMatch(content string, expected interface{}) (bool, string) {
	expectedStr, ok := expected.(string)
	if !ok {
		return false, fmt.Sprintf("expected must be a string, got %T", expected)
	}

	if strings.TrimSpace(content) == strings.TrimSpace(expectedStr) {
		return true, ""
	}

	return false, fmt.Sprintf("content does not exactly match expected")
}
