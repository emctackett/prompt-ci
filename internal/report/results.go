package report

import (
	"encoding/json"
	"os"
	"path/filepath"

	"prompt-ci/internal/suite"
)

// WriteResults writes the results.json file
func WriteResults(outDir string, results []suite.Result) error {
	path := filepath.Join(outDir, "results.json")

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
