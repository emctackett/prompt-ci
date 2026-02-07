.PHONY: build test validate run clean demo-fail-grounding demo-fail-schema

# Build the CLI
build:
	go build -o prompt-ci ./cmd/prompt-ci

# Run validation
validate: build
	./prompt-ci validate --suite eval-suite.yaml

# Run the eval suite
run: build
	./prompt-ci run --suite eval-suite.yaml --mode fixtures --fixtures ./fixtures --out ./out

# Run tests
test:
	go test ./...

# Clean build artifacts and output
clean:
	rm -f prompt-ci
	rm -rf out/

# Demo: Grounding failure
# This removes citations from a fixture to demonstrate how grounding validation fails
demo-fail-grounding: build
	@echo "=== Demo: Grounding Failure ==="
	@echo "Original fixture:"
	@cat fixtures/grounding/grounding_budget_flag.out.txt
	@echo ""
	@echo "Creating broken fixture (removing citation)..."
	@cp fixtures/grounding/grounding_budget_flag.out.txt fixtures/grounding/grounding_budget_flag.out.txt.bak
	@echo "The default value for the --budget flag is 100000 millicents which equals one dollar." > fixtures/grounding/grounding_budget_flag.out.txt
	@echo ""
	@echo "Running eval suite (expect failure)..."
	-./prompt-ci run --suite eval-suite.yaml --mode fixtures --fixtures ./fixtures --out ./out
	@echo ""
	@echo "Failure reasons from results.json:"
	@grep -A 10 '"grounding_budget_flag"' out/results.json | head -15
	@echo ""
	@echo "Restoring original fixture..."
	@mv fixtures/grounding/grounding_budget_flag.out.txt.bak fixtures/grounding/grounding_budget_flag.out.txt
	@echo ""
	@echo "Running eval suite (expect pass)..."
	./prompt-ci run --suite eval-suite.yaml --mode fixtures --fixtures ./fixtures --out ./out
	@echo "=== Demo Complete ==="

# Demo: Schema failure
# This adds an extra property to a JSON fixture to demonstrate additionalProperties: false
demo-fail-schema: build
	@echo "=== Demo: Schema Failure ==="
	@echo "Original fixture:"
	@cat fixtures/schema/schema_valid_assertion.out.json
	@echo ""
	@echo "Creating broken fixture (adding extra property)..."
	@cp fixtures/schema/schema_valid_assertion.out.json fixtures/schema/schema_valid_assertion.out.json.bak
	@echo '{"type": "contains", "expected": "hello world", "extra_field": "not allowed"}' > fixtures/schema/schema_valid_assertion.out.json
	@echo ""
	@echo "Running eval suite (expect failure)..."
	-./prompt-ci run --suite eval-suite.yaml --mode fixtures --fixtures ./fixtures --out ./out
	@echo ""
	@echo "Failure reasons from results.json:"
	@grep -A 10 '"schema_valid_assertion"' out/results.json | head -15
	@echo ""
	@echo "Restoring original fixture..."
	@mv fixtures/schema/schema_valid_assertion.out.json.bak fixtures/schema/schema_valid_assertion.out.json
	@echo ""
	@echo "Running eval suite (expect pass)..."
	./prompt-ci run --suite eval-suite.yaml --mode fixtures --fixtures ./fixtures --out ./out
	@echo "=== Demo Complete ==="

# Full demo: run both failure demos
demo: demo-fail-grounding demo-fail-schema
