# Suggested Commands

## Build & Install
```bash
make build          # Build the provider
make install        # Build and install locally
make dev-install    # Install to ~/.terraform.d/plugins for manual testing
```

## Testing
```bash
make test           # Run unit tests (fast, no API needed)
make testacc        # Run acceptance tests (requires TF_ACC=1 and real API)
go test -v -run TestName ./internal/provider/  # Run specific test
```

## Code Quality
```bash
make fmt            # Format code with gofmt
make lint           # Run golangci-lint
make check          # Run fmt, lint, and test together
```

## Documentation
```bash
make docs           # Generate Terraform docs from code
```

## Dependencies
```bash
make deps           # Download and tidy go modules
make verify         # Verify module checksums
```

## Utility Commands (Darwin/macOS)
```bash
git status          # Check working tree status
git diff            # View uncommitted changes
ls -la              # List directory contents
grep -r "pattern" . # Search for pattern (prefer rg if available)
```

## Environment Variables for Testing
```bash
export TF_ACC=1                           # Enable acceptance tests
export INFERADB_ENDPOINT="https://..."    # API endpoint
export INFERADB_SESSION_TOKEN="..."       # Auth token
```
