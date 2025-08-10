# Makefile for Go REST API Test Suite

# Variables
GO_VERSION := 1.21
BINARY_NAME := server
MAIN_PATH := ./cmd/server
TEST_TIMEOUT := 10m
COVERAGE_FILE := coverage.out

# Default target
.PHONY: help
help: ## Show this help message
	@echo 'Usage: make <target>'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
.PHONY: build
build: ## Build the server binary
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

.PHONY: run
run: ## Run the server
	go run $(MAIN_PATH)

# Test targets
.PHONY: test
test: ## Run unit tests
	go test -v -timeout=$(TEST_TIMEOUT) ./entities/... ./use_cases/... ./interfaces/...

.PHONY: test-unit
test-unit: ## Run unit tests with coverage
	go test -v -timeout=$(TEST_TIMEOUT) -coverprofile=$(COVERAGE_FILE) \
		./entities/... ./use_cases/... ./interfaces/...

.PHONY: test-integration
test-integration: ## Run integration tests
	INTEGRATION_TEST=true go test -v -timeout=$(TEST_TIMEOUT) ./tests/integration/...

.PHONY: test-e2e
test-e2e: ## Run E2E tests
	E2E_TEST=true go test -v -timeout=$(TEST_TIMEOUT) ./tests/e2e/...

.PHONY: test-all
test-all: test-unit test-integration test-e2e ## Run all tests

.PHONY: test-ginkgo
test-ginkgo: ## Run all tests using Ginkgo
	ginkgo -r --randomize-all --randomize-suites --fail-on-pending \
		--cover --coverprofile=$(COVERAGE_FILE) --race --trace

.PHONY: test-ginkgo-unit
test-ginkgo-unit: ## Run unit tests using Ginkgo
	ginkgo --randomize-all --randomize-suites --fail-on-pending \
		--cover --race --trace \
		entities use_cases interfaces

.PHONY: test-ginkgo-integration
test-ginkgo-integration: ## Run integration tests using Ginkgo
	INTEGRATION_TEST=true ginkgo --randomize-all --randomize-suites \
		--fail-on-pending --race --trace tests/integration

.PHONY: test-ginkgo-e2e
test-ginkgo-e2e: ## Run E2E tests using Ginkgo
	E2E_TEST=true ginkgo --randomize-all --randomize-suites \
		--fail-on-pending --trace tests/e2e

# Coverage targets
.PHONY: coverage
coverage: test-unit ## Generate and view test coverage
	go tool cover -html=$(COVERAGE_FILE) -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: coverage-func
coverage-func: test-unit ## Show coverage by function
	go tool cover -func=$(COVERAGE_FILE)

.PHONY: coverage-total
coverage-total: test-unit ## Show total coverage percentage
	@go tool cover -func=$(COVERAGE_FILE) | grep total | awk '{print "Total coverage: " $$3}'

# Quality targets
.PHONY: lint
lint: ## Run linter
	golangci-lint run

.PHONY: format
format: ## Format code
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: mod-tidy
mod-tidy: ## Tidy go modules
	go mod tidy

.PHONY: mod-download
mod-download: ## Download dependencies
	go mod download

# Development targets
.PHONY: deps
deps: ## Install development dependencies
	go install github.com/onsi/ginkgo/v2/ginkgo@latest
	go install github.com/matryer/moq@latest

.PHONY: generate
generate: ## Generate mocks and other generated files
	go generate ./...

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf bin/
	rm -f $(COVERAGE_FILE)
	rm -f coverage.html

# Docker targets (optional)
.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t $(BINARY_NAME) .

.PHONY: docker-run
docker-run: ## Run Docker container
	docker run -p 8080:8080 $(BINARY_NAME)

# CI targets
.PHONY: ci-test
ci-test: mod-tidy vet lint test-all ## Run all CI tests

.PHONY: ci-coverage
ci-coverage: test-unit coverage-total ## Generate coverage for CI

# Watch targets (requires fswatch or inotify-tools)
.PHONY: watch-test
watch-test: ## Watch files and run tests on changes
	@if command -v fswatch > /dev/null; then \
		fswatch -o . -e ".*" -i "\\.go$$" | xargs -n1 -I{} make test; \
	else \
		echo "fswatch not found. Install it with: brew install fswatch"; \
	fi

.PHONY: watch-test-unit
watch-test-unit: ## Watch files and run unit tests on changes
	@if command -v fswatch > /dev/null; then \
		fswatch -o . -e ".*" -i "\\.go$$" | xargs -n1 -I{} make test-unit; \
	else \
		echo "fswatch not found. Install it with: brew install fswatch"; \
	fi

# Benchmark targets
.PHONY: bench
bench: ## Run benchmarks
	go test -bench=. -benchmem ./...

.PHONY: bench-cpu
bench-cpu: ## Run CPU benchmarks with profiling
	go test -bench=. -benchmem -cpuprofile=cpu.prof ./...

.PHONY: bench-mem
bench-mem: ## Run memory benchmarks with profiling
	go test -bench=. -benchmem -memprofile=mem.prof ./...

# Performance targets
.PHONY: race
race: ## Run tests with race detection
	go test -race ./...

.PHONY: profile-cpu
profile-cpu: ## Profile CPU usage
	go tool pprof cpu.prof

.PHONY: profile-mem
profile-mem: ## Profile memory usage
	go tool pprof mem.prof

# Development server targets
.PHONY: dev
dev: ## Run development server with live reload
	air

.PHONY: install-air
install-air: ## Install air for live reload
	go install github.com/cosmtrek/air@latest