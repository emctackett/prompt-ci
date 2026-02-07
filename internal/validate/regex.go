package validate

import (
	"fmt"
	"regexp"
)

// ValidateRegex checks if content matches the expected regex pattern
func ValidateRegex(content string, expected interface{}) (bool, string) {
	pattern, ok := expected.(string)
	if !ok {
		return false, fmt.Sprintf("expected must be a string pattern, got %T", expected)
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, fmt.Sprintf("invalid regex pattern '%s': %v", pattern, err)
	}

	if re.MatchString(content) {
		return true, ""
	}

	return false, fmt.Sprintf("content does not match pattern '%s'", pattern)
}
