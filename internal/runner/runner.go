package runner

import (
	"fmt"
	"time"

	"prompt-ci/internal/suite"
	"prompt-ci/internal/validate"
)

// RunSuite runs all cases in the suite and returns results
// Returns the results and a boolean indicating if any errors occurred
func RunSuite(s *suite.Suite, fixturesDir string, failFast bool) ([]suite.Result, bool) {
	var results []suite.Result
	hasError := false

	for _, c := range s.Cases {
		result := runCase(c, s, fixturesDir)
		results = append(results, result)

		if result.Status == suite.StatusError {
			hasError = true
		}

		if failFast && (result.Status == suite.StatusFail || result.Status == suite.StatusError) {
			break
		}
	}

	return results, hasError
}

// runCase runs a single test case
func runCase(c suite.Case, s *suite.Suite, fixturesDir string) suite.Result {
	start := time.Now()

	// Load fixture
	content, err := LoadFixture(fixturesDir, c.ID)
	if err != nil {
		return suite.Result{
			ID:             c.ID,
			Status:         suite.StatusError,
			Validator:      getValidatorType(c),
			DurationMS:     time.Since(start).Milliseconds(),
			FailureReasons: []string{fmt.Sprintf("fixture error: %v", err)},
		}
	}

	// Run all assertions
	var failures []string
	for _, assertion := range c.Assertions {
		passed, reason := validate.ValidateAssertion(content, assertion, s)
		if !passed {
			failures = append(failures, fmt.Sprintf("[%s] %s", assertion.Type, reason))
		}
	}

	// For grounding cases, also validate citations
	caseType := suite.GetCaseType(c.ID)
	if caseType == suite.CaseTypeGrounding {
		passed, groundingFailures := validate.ValidateGrounding(content, s)
		if !passed {
			for _, f := range groundingFailures {
				failures = append(failures, fmt.Sprintf("[grounding] %s", f))
			}
		}
	}

	// Determine status
	status := suite.StatusPass
	if len(failures) > 0 {
		status = suite.StatusFail
	}

	return suite.Result{
		ID:             c.ID,
		Status:         status,
		Validator:      getValidatorType(c),
		DurationMS:     time.Since(start).Milliseconds(),
		FailureReasons: failures,
	}
}

// getValidatorType determines the primary validator type for a case
func getValidatorType(c suite.Case) string {
	caseType := suite.GetCaseType(c.ID)
	switch caseType {
	case suite.CaseTypeGrounding:
		return "grounding"
	case suite.CaseTypeSchema:
		return "json_schema"
	case suite.CaseTypeTool:
		return "tool"
	default:
		// Determine from assertions
		for _, a := range c.Assertions {
			if a.Type == "json_schema" {
				return "json_schema"
			}
		}
		return "contains"
	}
}
