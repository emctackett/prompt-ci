package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"prompt-ci/internal/report"
	"prompt-ci/internal/runner"
	"prompt-ci/internal/suite"
)

var (
	suitePath   string
	fixturesDir string
	outDir      string
	mode        string
	failFast    bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "prompt-ci",
		Short: "A CLI tool for running LLM eval suites",
		Long:  "prompt-ci is a CLI tool that executes eval suites and produces machine-readable and human-readable artifacts.",
	}

	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate the suite file",
		Long:  "Validates that the suite file is internally consistent and all fixtures exist.",
		RunE:  runValidate,
	}
	validateCmd.Flags().StringVar(&suitePath, "suite", "", "Path to the suite file (required)")
	validateCmd.Flags().StringVar(&fixturesDir, "fixtures", "./fixtures", "Path to fixtures directory")
	validateCmd.MarkFlagRequired("suite")

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run the eval suite",
		Long:  "Runs the eval suite against fixtures and produces results.",
		RunE:  runRun,
	}
	runCmd.Flags().StringVar(&suitePath, "suite", "", "Path to the suite file (required)")
	runCmd.Flags().StringVar(&mode, "mode", "fixtures", "Run mode (only 'fixtures' supported in MVP)")
	runCmd.Flags().StringVar(&fixturesDir, "fixtures", "./fixtures", "Path to fixtures directory")
	runCmd.Flags().StringVar(&outDir, "out", "./out", "Output directory for artifacts")
	runCmd.Flags().BoolVar(&failFast, "fail-fast", false, "Stop on first failure")
	runCmd.MarkFlagRequired("suite")

	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(runCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runValidate(cmd *cobra.Command, args []string) error {
	// Parse suite
	s, err := suite.ParseFile(suitePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}

	// Validate suite
	if err := suite.ValidateSuite(s, fixturesDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}

	fmt.Printf("Suite '%s' is valid (%d cases)\n", s.Name, len(s.Cases))
	return nil
}

func runRun(cmd *cobra.Command, args []string) error {
	// Validate mode
	if mode != "fixtures" {
		fmt.Fprintf(os.Stderr, "Error: only 'fixtures' mode is supported in MVP\n")
		os.Exit(2)
	}

	// Parse suite
	s, err := suite.ParseFile(suitePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}

	// Validate suite
	if err := suite.ValidateSuite(s, fixturesDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}

	// Run suite
	results, hasError := runner.RunSuite(s, fixturesDir, failFast)

	// Create output directory
	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(2)
	}

	// Write reports
	if err := report.WriteResults(outDir, results); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing results.json: %v\n", err)
		os.Exit(2)
	}

	if err := report.WriteJUnit(outDir, s.Name, results); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing junit.xml: %v\n", err)
		os.Exit(2)
	}

	if err := report.WriteHTML(outDir, s.Name, results); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing report.html: %v\n", err)
		os.Exit(2)
	}

	if err := report.WriteTrace(outDir, s.Name, results); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing trace.json: %v\n", err)
		os.Exit(2)
	}

	// Calculate summary
	passed := 0
	failed := 0
	errors := 0
	for _, r := range results {
		switch r.Status {
		case suite.StatusPass:
			passed++
		case suite.StatusFail:
			failed++
		case suite.StatusError:
			errors++
		}
	}

	// Print summary (stable format for piping)
	fmt.Printf("prompt-ci: %d/%d cases passed\n", passed, len(results))

	// Determine exit code
	if hasError {
		os.Exit(2)
	}
	if failed > 0 {
		os.Exit(1)
	}
	return nil
}
