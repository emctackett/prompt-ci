package runner

import (
	"fmt"
	"os"

	"prompt-ci/internal/suite"
)

// LoadFixture loads the fixture content for a test case
func LoadFixture(fixturesDir string, caseID string) (string, error) {
	path := suite.GetFixturePath(fixturesDir, caseID)

	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read fixture %s: %w", path, err)
	}

	return string(data), nil
}
