.PHONY: build
build:
	@echo "Building the project..."

.PHONY: test
unit-test:
	@echo "Running tests..."
	@go test ./...

.PHONY: run-example
run-example:
	@echo "Running example..."
	@./scripts/run_examples.sh
