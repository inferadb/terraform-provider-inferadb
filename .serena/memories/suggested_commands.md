# Suggested Commands

## Build & Install
```bash
go build -v ./...                         # Build the provider
go install -v ./...                       # Build and install locally
```

## Testing
```bash
go test -v -cover -timeout=120s -parallel=4 ./...   # Run unit tests
TF_ACC=1 go test -v -cover -timeout 120m ./...      # Run acceptance tests
go test -v -run TestName ./internal/provider/       # Run specific test
```

## Code Quality
```bash
gofmt -s -w -e .                          # Format code
golangci-lint run                         # Run linter
```

## Documentation
```bash
go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs  # Generate docs
```

## Dependencies
```bash
go mod download                           # Download dependencies
go mod tidy                               # Clean up go.mod/go.sum
go mod verify                             # Verify module checksums
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
