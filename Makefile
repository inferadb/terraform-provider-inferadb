# Copyright 2025 InferaDB
# SPDX-License-Identifier: Apache-2.0

default: install

# Build the provider
.PHONY: build
build:
	go build -v ./...

# Install the provider locally for testing
.PHONY: install
install: build
	go install -v ./...

# Run linting
.PHONY: lint
lint:
	golangci-lint run

# Generate code (currently just docs)
.PHONY: generate
generate:
	go generate ./...

# Format code
.PHONY: fmt
fmt:
	gofmt -s -w -e .

# Run unit tests
.PHONY: test
test:
	go test -v -cover -timeout=120s -parallel=4 ./...

# Run acceptance tests (requires TF_ACC=1 and INFERADB_* env vars)
# For local development, use: make test-env-up && source .env.test && make testacc
.PHONY: testacc
testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

# ============================================================================
# Test Environment Management
# ============================================================================

# Docker Compose project name for test environment
TEST_COMPOSE_PROJECT := terraform-provider-inferadb-test
TEST_COMPOSE_FILE := docker-compose.test.yml

# Start the test environment (FoundationDB + MailHog + Control)
.PHONY: test-env-up
test-env-up:
	@echo "Starting test environment..."
	docker compose -p $(TEST_COMPOSE_PROJECT) -f $(TEST_COMPOSE_FILE) up -d --build
	@echo "Test environment started. Waiting for services to be healthy..."
	@echo "Run 'make test-env-bootstrap' to register a test user and get credentials."

# Stop and remove the test environment
.PHONY: test-env-down
test-env-down:
	@echo "Stopping test environment..."
	docker compose -p $(TEST_COMPOSE_PROJECT) -f $(TEST_COMPOSE_FILE) down -v --remove-orphans
	@rm -f .env.test
	@echo "Test environment stopped and cleaned up."

# Restart the test environment
.PHONY: test-env-restart
test-env-restart: test-env-down test-env-up

# Show test environment logs
.PHONY: test-env-logs
test-env-logs:
	docker compose -p $(TEST_COMPOSE_PROJECT) -f $(TEST_COMPOSE_FILE) logs -f

# Show test environment status
.PHONY: test-env-status
test-env-status:
	docker compose -p $(TEST_COMPOSE_PROJECT) -f $(TEST_COMPOSE_FILE) ps

# Bootstrap test environment (register test user, export credentials)
.PHONY: test-env-bootstrap
test-env-bootstrap:
	@echo "Bootstrapping test environment..."
	./scripts/bootstrap-test-env.sh --endpoint http://localhost:9090
	@echo ""
	@echo "Bootstrap complete! To run acceptance tests:"
	@echo "  source .env.test && make testacc"

# Run acceptance tests with full environment setup
# This is the all-in-one target for CI or local development
.PHONY: testacc-local
testacc-local: test-env-up test-env-bootstrap
	@echo "Running acceptance tests..."
	@bash -c 'source .env.test && TF_ACC=1 go test -v -cover -timeout 120m ./...'
	@echo "Acceptance tests complete."

# Clean up everything (test env + build artifacts)
.PHONY: clean-all
clean-all: test-env-down clean
	@echo "All artifacts cleaned."

# Generate documentation
.PHONY: docs
docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

# Clean build artifacts
.PHONY: clean
clean:
	rm -f terraform-provider-inferadb

# Download dependencies
.PHONY: deps
deps:
	go mod download
	go mod tidy

# Verify dependencies
.PHONY: verify
verify:
	go mod verify

# Run all checks
.PHONY: check
check: fmt lint test

# Local development - install to ~/.terraform.d/plugins for manual testing
.PHONY: dev-install
dev-install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/inferadb/inferadb/0.0.1/$(shell go env GOOS)_$(shell go env GOARCH)
	cp terraform-provider-inferadb ~/.terraform.d/plugins/registry.terraform.io/inferadb/inferadb/0.0.1/$(shell go env GOOS)_$(shell go env GOARCH)/
