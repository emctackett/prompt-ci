package report

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"prompt-ci/internal/suite"
)

// TraceEntry represents a trace entry for a test case
type TraceEntry struct {
	CaseID        string   `json:"case_id"`
	Mode          string   `json:"mode"`
	FixturePath   string   `json:"fixture_path,omitempty"`
	CitationsFound []string `json:"citations_found,omitempty"`
	StartedAt     string   `json:"started_at"`
	CompletedAt   string   `json:"completed_at"`
}

// Trace represents the trace file structure
type Trace struct {
	SuiteName string       `json:"suite_name"`
	Mode      string       `json:"mode"`
	StartedAt string       `json:"started_at"`
	CompletedAt string     `json:"completed_at"`
	Entries   []TraceEntry `json:"entries"`
}

// WriteTrace writes the trace.json file (stub for MVP)
func WriteTrace(outDir, suiteName string, results []suite.Result) error {
	path := filepath.Join(outDir, "trace.json")

	now := time.Now().Format(time.RFC3339)
	trace := Trace{
		SuiteName:   suiteName,
		Mode:        "fixtures",
		StartedAt:   now,
		CompletedAt: now,
		Entries:     make([]TraceEntry, 0, len(results)),
	}

	for _, r := range results {
		entry := TraceEntry{
			CaseID:      r.ID,
			Mode:        "fixtures",
			StartedAt:   now,
			CompletedAt: now,
		}
		trace.Entries = append(trace.Entries, entry)
	}

	data, err := json.MarshalIndent(trace, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
