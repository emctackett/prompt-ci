# prompt-ci

A CLI tool for running LLM eval suites with deterministic, grounded validation. Executes test cases against fixtures and produces machine-readable (JUnit, JSON) and human-readable (HTML) reports.

## Features

- **Deterministic evaluation**: Run eval suites against fixtures without making LLM API calls
- **Multiple validators**: Contains, regex, JSON schema, and grounding (citation) validation
- **Grounding enforcement**: Ensures responses cite source documentation correctly
- **Strict schema validation**: Enforces `additionalProperties: false` for JSON schemas
- **CI-friendly output**: JUnit XML, JSON results, and HTML reports

## Installation

```bash
# Clone and build
git clone <repository-url>
cd prompt-ci
go build -o prompt-ci ./cmd/prompt-ci
```

Or use Make:

```bash
make build
```

## Quick Start

```bash
# Validate your eval suite
./prompt-ci validate --suite eval-suite.yaml

# Run the eval suite against fixtures
./prompt-ci run --suite eval-suite.yaml --mode fixtures --fixtures ./fixtures --out ./out

# View results
open out/report.html
```

## CLI Commands

### `prompt-ci validate`

Validates that a suite file is internally consistent and all fixtures exist.

```bash
./prompt-ci validate --suite <path> [--fixtures <dir>]
```

| Flag | Description | Default |
|------|-------------|---------|
| `--suite` | Path to suite YAML file (required) | - |
| `--fixtures` | Path to fixtures directory | `./fixtures` |

**Exit codes:**
- `0` - Suite is valid
- `2` - Invalid suite or missing fixtures

### `prompt-ci run`

Runs the eval suite against fixtures and generates reports.

```bash
./prompt-ci run --suite <path> [--mode fixtures] [--fixtures <dir>] [--out <dir>] [--fail-fast]
```

| Flag | Description | Default |
|------|-------------|---------|
| `--suite` | Path to suite YAML file (required) | - |
| `--mode` | Run mode (only `fixtures` supported) | `fixtures` |
| `--fixtures` | Path to fixtures directory | `./fixtures` |
| `--out` | Output directory for artifacts | `./out` |
| `--fail-fast` | Stop on first failure | `false` |

**Exit codes:**
- `0` - All cases passed
- `1` - One or more cases failed
- `2` - Runtime or configuration error

## Eval Suite Format

```yaml
suite_name: my-eval-suite
capability: mixed

docs:
  - id: glossary
    title: Glossary
    chunks:
      - id: c1
        text: |
          Documentation content here...

cases:
  - id: test_case_id
    prompt: "The prompt that would be sent to an LLM"
    assertions:
      - type: contains
        expected: "expected substring"
      - type: regex
        expected: "pattern.*to match"
      - type: json_schema
        expected:
          type: object
          properties:
            field: { type: string }
          required: ["field"]
          additionalProperties: false
```

## Validators

### `contains`
Checks if the response contains an expected substring.

```yaml
- type: contains
  expected: "expected text"
```

### `regex`
Checks if the response matches a regular expression (RE2 syntax).

```yaml
- type: regex
  expected: "pattern\\s+to\\s+match"
```

### `json_schema`
Validates that the response contains valid JSON matching a schema. The schema **must** include `additionalProperties: false` for object types.

```yaml
- type: json_schema
  expected:
    type: object
    properties:
      name: { type: string }
    required: ["name"]
    additionalProperties: false
```

### `grounding`
Automatically validates citation requirements for grounding cases (case IDs starting with `grounding_`):

- Citations must use format `[doc:<doc_id>#<chunk_id>]`
- All citations must reference valid doc/chunk pairs defined in the suite
- Sentences with 6+ tokens must include at least one citation

## Fixtures

Fixtures are organized by case type:

```
fixtures/
├── grounding/          # Text files with citations
│   └── <case_id>.out.txt
├── schema/             # JSON files matching schemas
│   └── <case_id>.out.json
└── tool/               # Tool call fixtures
    └── <case_id>.out.{json,txt}
```

Case types are determined by ID prefix:
- `grounding_*` → `fixtures/grounding/*.out.txt`
- `schema_*` → `fixtures/schema/*.out.json`
- `tool_*` → `fixtures/tool/*.out.{json,txt}`

## Output Artifacts

After running, the `--out` directory contains:

| File | Description |
|------|-------------|
| `results.json` | Per-case results with status, validator type, duration, and failure reasons |
| `junit.xml` | JUnit XML format for CI integration |
| `report.html` | Interactive HTML report with expandable failure details |
| `trace.json` | Execution trace (stub for future live mode) |

### results.json format

```json
[
  {
    "id": "case_id",
    "status": "PASS",
    "validator": "grounding",
    "duration_ms": 1,
    "failure_reasons": []
  }
]
```

Status values: `PASS`, `FAIL`, `ERROR`, `SKIP`

## Demo Failure Modes

Demonstrate how validation catches issues:

```bash
# Show grounding failure (missing citation)
make demo-fail-grounding

# Show schema failure (extra property violates additionalProperties: false)
make demo-fail-schema

# Run both demos
make demo
```

## Project Structure

```
prompt-ci/
├── cmd/prompt-ci/
│   └── main.go              # CLI entry point
├── internal/
│   ├── suite/               # Suite parsing and validation
│   │   ├── types.go         # Data structures
│   │   ├── parser.go        # YAML parsing
│   │   └── validate.go      # Suite validation
│   ├── validate/            # Validators
│   │   ├── validator.go     # Validator interface
│   │   ├── contains.go      # Contains validator
│   │   ├── regex.go         # Regex validator
│   │   ├── json_schema.go   # JSON Schema validator
│   │   └── grounding.go     # Citation/grounding validator
│   ├── runner/              # Test execution
│   │   ├── runner.go        # Suite runner
│   │   └── fixture.go       # Fixture loading
│   └── report/              # Report generation
│       ├── results.go       # results.json
│       ├── junit.go         # junit.xml
│       ├── html.go          # report.html
│       └── trace.go         # trace.json
├── fixtures/                # Test fixtures
├── eval-suite.yaml          # Example eval suite
├── Makefile                 # Build and demo targets
└── SPEC.md                  # Full specification
```

## Make Targets

| Target | Description |
|--------|-------------|
| `make build` | Build the CLI binary |
| `make validate` | Validate the eval suite |
| `make run` | Run the eval suite |
| `make test` | Run Go tests |
| `make clean` | Remove build artifacts |
| `make demo` | Run failure mode demos |

## Requirements

- Go 1.21 or later

## Dependencies

- `github.com/spf13/cobra` - CLI framework
- `gopkg.in/yaml.v3` - YAML parsing
- `github.com/santhosh-tekuri/jsonschema/v5` - JSON Schema validation
