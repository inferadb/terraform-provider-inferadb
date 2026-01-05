# Codebase Structure

```
terraform-provider-inferadb/
├── main.go                    # Entry point, starts provider server
├── go.mod                     # Go module definition (Go 1.24)
├── Makefile                   # Build/test/lint commands
├── tools.go                   # Tool dependencies (tfplugindocs)
│
├── internal/
│   ├── provider/              # Terraform plugin implementation
│   │   ├── provider.go        # Main provider (Configure, Resources, DataSources)
│   │   ├── resource_*.go      # Resource implementations
│   │   ├── datasource_*.go    # Data source implementations
│   │   └── *_test.go          # Test files
│   │
│   └── client/                # InferaDB API client
│       ├── client.go          # Base HTTP client with auth
│       ├── models.go          # Shared types (Snowflake IDs)
│       ├── organizations.go   # Organization CRUD
│       ├── vaults.go          # Vault CRUD
│       ├── clients.go         # Client CRUD
│       ├── certificates.go    # Certificate operations
│       ├── teams.go           # Team CRUD
│       └── grants.go          # Vault grant operations
│
├── examples/                  # Terraform example configs
│   ├── provider/              # Provider configuration
│   ├── resources/             # Per-resource examples
│   └── data-sources/          # Per-datasource examples
│
├── docs/                      # Generated documentation (via tfplugindocs)
│
└── .github/
    └── workflows/
        ├── test.yml           # CI pipeline (build, lint, test)
        └── codeql.yml         # Security scanning
```

## Key Files
- `internal/provider/provider.go` - Central provider setup, registers all resources/datasources
- `internal/client/client.go` - HTTP client base, handles auth via session cookie
- `internal/client/models.go` - Shared types like `SnowflakeID`

## CI Pipeline
The `.github/workflows/test.yml` runs:
1. **Build** - Compiles the provider
2. **Lint** - Runs golangci-lint
3. **Test** - Unit tests
4. **Acceptance** - Full acceptance tests (main branch only, requires secrets)
