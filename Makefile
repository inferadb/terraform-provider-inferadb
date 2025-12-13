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

# Run acceptance tests (requires TF_ACC=1)
.PHONY: testacc
testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

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
