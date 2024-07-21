.PHONY: build
build:
	@echo "Building the project..."

.PHONY: test
unit-test:
	@echo "Running tests..."
	@go test ./...

.PHONY: run-examples
run-examples:
	@echo "Running example..."
	@./scripts/run_examples.sh
