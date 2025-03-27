DEFAULT_GOAL := test

.PHONY: test
test:
	@echo "Running tests..."
	@go test ./...

.PHONY: lint
lint:
	@echo "Running linters..."
	@golangci-lint run ./...
