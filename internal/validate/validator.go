package validate

import (
	"prompt-ci/internal/suite"
)

// Validator is the interface for all validators
type Validator interface {
	Validate(content string, expected interface{}) (bool, string)
}

// ValidateAssertion validates a single assertion against content
func ValidateAssertion(content string, assertion suite.Assertion, s *suite.Suite) (bool, string) {
	switch assertion.Type {
	case "contains":
		return ValidateContains(content, assertion.Expected)
	case "regex":
		return ValidateRegex(content, assertion.Expected)
	case "json_schema":
		return ValidateJSONSchema(content, assertion.Expected)
	case "exact_match":
		return ValidateExactMatch(content, assertion.Expected)
	default:
		return false, "unknown assertion type: " + assertion.Type
	}
}

// ValidateGrounding validates grounding requirements for a response
func ValidateGrounding(content string, s *suite.Suite) (bool, []string) {
	return ValidateGroundingCitations(content, s)
}
