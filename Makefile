BINARY := gostrap
GO     := go

.PHONY: build test lint fmt clean help

build: ## Build the binary
	$(GO) build -o $(BINARY) ./cmd/gostrap/

test: ## Run tests with race detector
	$(GO) test -race -count=1 ./...

lint: ## Run golangci-lint
	golangci-lint run

fmt: ## Format code
	$(GO) fmt ./...
	golangci-lint fmt

clean: ## Remove build artifacts
	rm -f $(BINARY)
	rm -f coverage.out

cover: ## Run tests with coverage report
	$(GO) test -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Open coverage.html in your browser"

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
