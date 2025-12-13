# InferaDB Terraform Provider

Terraform provider for InferaDB Control Plane resources using terraform-plugin-framework.

## Quick Commands

```bash
# Build
make build

# Install locally
make install

# Run tests
make test

# Run acceptance tests (requires API)
TF_ACC=1 INFERADB_ENDPOINT=... INFERADB_SESSION_TOKEN=... make testacc

# Format code
make fmt
go fmt ./...

# Generate docs
make docs
```

## Architecture

### Package Structure

| Package             | Purpose                                         |
| ------------------- | ----------------------------------------------- |
| `internal/client`   | HTTP client for InferaDB Control                |
| `internal/provider` | Terraform provider, resources, and data sources |

### Key Files

- `main.go` - Provider entry point
- `internal/client/client.go` - HTTP client with authentication
- `internal/client/models.go` - API request/response types
- `internal/provider/provider.go` - Provider configuration
- `internal/provider/resource_*.go` - Resource implementations
- `internal/provider/datasource_*.go` - Data source implementations

## Resource Patterns

All resources follow this pattern:

```go
type XxxResource struct {
    client *client.Client
}

type XxxResourceModel struct {
    ID   types.String `tfsdk:"id"`
    // ... other fields
}

func (r *XxxResource) Configure(...)  // Get client from provider
func (r *XxxResource) Create(...)     // POST to API
func (r *XxxResource) Read(...)       // GET from API
func (r *XxxResource) Update(...)     // PATCH to API
func (r *XxxResource) Delete(...)     // DELETE to API
func (r *XxxResource) ImportState(...) // Parse import ID
```

### Error Handling

Check for 404 errors to handle deleted resources:

```go
if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
    resp.State.RemoveResource(ctx)
    return
}
```

### Import ID Formats

| Resource                      | Format                           |
| ----------------------------- | -------------------------------- |
| `inferadb_organization`       | `<org_id>`                       |
| `inferadb_vault`              | `<org_id>/<vault_id>`            |
| `inferadb_client`             | `<org_id>/<client_id>`           |
| `inferadb_client_certificate` | `<org_id>/<client_id>/<cert_id>` |
| `inferadb_team`               | `<org_id>/<team_id>`             |
| `inferadb_team_member`        | `<org_id>/<team_id>/<member_id>` |
| `inferadb_vault_user_grant`   | `<org_id>/<vault_id>/<grant_id>` |
| `inferadb_vault_team_grant`   | `<org_id>/<vault_id>/<grant_id>` |

## API Client

The client uses session cookie authentication:

```go
client := client.New(client.Config{
    Endpoint:     "https://api.inferadb.com",
    SessionToken: "...",
})
```

All IDs are Snowflake IDs (64-bit, represented as strings).

## Adding a New Resource

1. Add API methods to `internal/client/` (e.g., `foos.go`)
2. Create resource file `internal/provider/resource_foo.go`
3. Implement: `NewFooResource`, `Metadata`, `Schema`, `Configure`, `Create`, `Read`, `Update`, `Delete`, `ImportState`
4. Register in `provider.go` `Resources()` method
5. Add example in `examples/resources/inferadb_foo/resource.tf`

## Testing

Acceptance tests require a running InferaDB Control:

```bash
export TF_ACC=1
export INFERADB_ENDPOINT=http://localhost:3000
export INFERADB_SESSION_TOKEN=your-session-token
go test -v ./internal/provider/
```

## Common Issues

1. **Certificate private key empty after import**: Private keys are only returned on creation. Cannot be retrieved later.

2. **Resource not found after deletion**: API uses soft deletes. 404 handling in Read removes from state.

3. **Organization ID required**: Most resources require `organization_id`. It's used in API paths.

## Dependencies

- `terraform-plugin-framework` v1.15+ (not SDKv2)
- `terraform-plugin-go` v0.25+
- `terraform-plugin-testing` v1.11+ (for acceptance tests)
