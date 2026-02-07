package report

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"prompt-ci/internal/suite"
)

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>prompt-ci Report - {{.SuiteName}}</title>
    <style>
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; padding: 20px; max-width: 1200px; margin: 0 auto; background: #f5f5f5; }
        h1 { margin-bottom: 20px; color: #333; }
        .summary { display: flex; gap: 20px; margin-bottom: 30px; }
        .summary-card { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); flex: 1; text-align: center; }
        .summary-card h2 { font-size: 2em; margin-bottom: 5px; }
        .summary-card.pass h2 { color: #22c55e; }
        .summary-card.fail h2 { color: #ef4444; }
        .summary-card.error h2 { color: #f59e0b; }
        .summary-card.total h2 { color: #3b82f6; }
        table { width: 100%; border-collapse: collapse; background: white; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        th, td { padding: 12px 15px; text-align: left; border-bottom: 1px solid #eee; }
        th { background: #f8f9fa; font-weight: 600; color: #333; }
        tr:hover { background: #f8f9fa; }
        .status { padding: 4px 12px; border-radius: 20px; font-size: 0.85em; font-weight: 500; }
        .status-pass { background: #dcfce7; color: #166534; }
        .status-fail { background: #fee2e2; color: #991b1b; }
        .status-error { background: #fef3c7; color: #92400e; }
        .status-skip { background: #e5e7eb; color: #374151; }
        .details { cursor: pointer; }
        .details:hover { text-decoration: underline; }
        .failure-reasons { display: none; background: #fef2f2; padding: 15px; margin: 10px 0; border-radius: 4px; font-family: monospace; font-size: 0.9em; white-space: pre-wrap; }
        .failure-reasons.show { display: block; }
    </style>
</head>
<body>
    <h1>prompt-ci Report: {{.SuiteName}}</h1>

    <div class="summary">
        <div class="summary-card pass">
            <h2>{{.Passed}}</h2>
            <p>Passed</p>
        </div>
        <div class="summary-card fail">
            <h2>{{.Failed}}</h2>
            <p>Failed</p>
        </div>
        <div class="summary-card error">
            <h2>{{.Errors}}</h2>
            <p>Errors</p>
        </div>
        <div class="summary-card total">
            <h2>{{.Total}}</h2>
            <p>Total</p>
        </div>
    </div>

    <table>
        <thead>
            <tr>
                <th>Case ID</th>
                <th>Validator</th>
                <th>Status</th>
                <th>Duration</th>
                <th>Details</th>
            </tr>
        </thead>
        <tbody>
            {{range .Results}}
            <tr>
                <td>{{.ID}}</td>
                <td>{{.Validator}}</td>
                <td><span class="status status-{{.StatusLower}}">{{.Status}}</span></td>
                <td>{{.DurationMS}}ms</td>
                <td>
                    {{if .FailureReasons}}
                    <span class="details" onclick="this.nextElementSibling.classList.toggle('show')">Show failures</span>
                    <div class="failure-reasons">{{.FailureReasonsText}}</div>
                    {{else}}
                    -
                    {{end}}
                </td>
            </tr>
            {{end}}
        </tbody>
    </table>
</body>
</html>`

type htmlData struct {
	SuiteName string
	Passed    int
	Failed    int
	Errors    int
	Total     int
	Results   []htmlResult
}

type htmlResult struct {
	ID                 string
	Status             string
	StatusLower        string
	Validator          string
	DurationMS         int64
	FailureReasons     []string
	FailureReasonsText string
}

// WriteHTML writes the report.html file
func WriteHTML(outDir, suiteName string, results []suite.Result) error {
	path := filepath.Join(outDir, "report.html")

	data := htmlData{
		SuiteName: suiteName,
		Total:     len(results),
	}

	for _, r := range results {
		switch r.Status {
		case suite.StatusPass:
			data.Passed++
		case suite.StatusFail:
			data.Failed++
		case suite.StatusError:
			data.Errors++
		}

		hr := htmlResult{
			ID:                 r.ID,
			Status:             string(r.Status),
			StatusLower:        strings.ToLower(string(r.Status)),
			Validator:          r.Validator,
			DurationMS:         r.DurationMS,
			FailureReasons:     r.FailureReasons,
			FailureReasonsText: strings.Join(r.FailureReasons, "\n"),
		}
		data.Results = append(data.Results, hr)
	}

	tmpl, err := template.New("report").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}
