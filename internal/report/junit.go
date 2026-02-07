package report

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"

	"prompt-ci/internal/suite"
)

// JUnitTestSuites is the root element
type JUnitTestSuites struct {
	XMLName    xml.Name         `xml:"testsuites"`
	TestSuites []JUnitTestSuite `xml:"testsuite"`
}

// JUnitTestSuite represents a test suite
type JUnitTestSuite struct {
	XMLName   xml.Name        `xml:"testsuite"`
	Name      string          `xml:"name,attr"`
	Tests     int             `xml:"tests,attr"`
	Failures  int             `xml:"failures,attr"`
	Errors    int             `xml:"errors,attr"`
	Skipped   int             `xml:"skipped,attr"`
	Time      float64         `xml:"time,attr"`
	TestCases []JUnitTestCase `xml:"testcase"`
}

// JUnitTestCase represents a test case
type JUnitTestCase struct {
	XMLName   xml.Name      `xml:"testcase"`
	Name      string        `xml:"name,attr"`
	ClassName string        `xml:"classname,attr"`
	Time      float64       `xml:"time,attr"`
	Failure   *JUnitFailure `xml:"failure,omitempty"`
	Error     *JUnitError   `xml:"error,omitempty"`
	Skipped   *JUnitSkipped `xml:"skipped,omitempty"`
}

// JUnitFailure represents a test failure
type JUnitFailure struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Content string `xml:",chardata"`
}

// JUnitError represents a test error
type JUnitError struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Content string `xml:",chardata"`
}

// JUnitSkipped represents a skipped test
type JUnitSkipped struct {
	Message string `xml:"message,attr,omitempty"`
}

// WriteJUnit writes the junit.xml file
func WriteJUnit(outDir, suiteName string, results []suite.Result) error {
	path := filepath.Join(outDir, "junit.xml")

	var testCases []JUnitTestCase
	failures := 0
	errors := 0
	skipped := 0
	totalTime := 0.0

	for _, r := range results {
		tc := JUnitTestCase{
			Name:      r.ID,
			ClassName: "prompt-ci." + r.Validator,
			Time:      float64(r.DurationMS) / 1000.0,
		}
		totalTime += tc.Time

		switch r.Status {
		case suite.StatusFail:
			failures++
			tc.Failure = &JUnitFailure{
				Message: "Test case failed",
				Type:    "AssertionError",
				Content: strings.Join(r.FailureReasons, "\n"),
			}
		case suite.StatusError:
			errors++
			tc.Error = &JUnitError{
				Message: "Test case error",
				Type:    "RuntimeError",
				Content: strings.Join(r.FailureReasons, "\n"),
			}
		case suite.StatusSkip:
			skipped++
			tc.Skipped = &JUnitSkipped{
				Message: "Test case skipped",
			}
		}

		testCases = append(testCases, tc)
	}

	testSuites := JUnitTestSuites{
		TestSuites: []JUnitTestSuite{
			{
				Name:      suiteName,
				Tests:     len(results),
				Failures:  failures,
				Errors:    errors,
				Skipped:   skipped,
				Time:      totalTime,
				TestCases: testCases,
			},
		},
	}

	data, err := xml.MarshalIndent(testSuites, "", "  ")
	if err != nil {
		return err
	}

	// Add XML header
	output := []byte(xml.Header + string(data))
	return os.WriteFile(path, output, 0644)
}
