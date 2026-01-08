<div align="center">
    <p><a href="https://inferadb.com"><img src=".github/inferadb.png" width="100" alt="InferaDB Logo" /></a></p>
    <h1>InferaDB Terraform Provider</h1>
    <p>
        <a href="https://discord.gg/inferadb"><img src="https://img.shields.io/badge/Discord-Join%20us-5865F2?logo=discord&logoColor=white" alt="Discord" /></a>
        <a href="#license"><img src="https://img.shields.io/badge/license-MIT%2FApache--2.0-blue.svg" alt="License" /></a>
    </p>
    <p>Manage InferaDB organizations, vaults, clients, teams, and access grants</p>
</div>

> [!IMPORTANT]
> Under active development. Not production-ready.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23 (to build from source)

## Installation

### From Source

```bash
git clone https://github.com/inferadb/terraform-provider-inferadb.git
cd terraform-provider-inferadb
go build -v ./...
go install -v ./...
```

## Usage

```hcl
terraform {
  required_providers {
    inferadb = {
      source = "inferadb/inferadb"
    }
  }
}

provider "inferadb" {
  endpoint      = "https://api.inferadb.com"
  session_token = var.inferadb_session_token
}

# Create an organization
resource "inferadb_organization" "example" {
  name = "My Organization"
  tier = "dev"
}

# Create a vault
resource "inferadb_vault" "production" {
  organization_id = inferadb_organization.example.id
  name            = "Production Policies"
  description     = "Authorization policies for production"
}

# Create a client (backend service identity)
resource "inferadb_client" "api" {
  organization_id = inferadb_organization.example.id
  vault_id        = inferadb_vault.production.id
  name            = "API Server"
}

# Generate a certificate for the client
resource "inferadb_client_certificate" "api_cert" {
  organization_id = inferadb_organization.example.id
  client_id       = inferadb_client.api.id
  name            = "API Certificate 2025"
}

# Output the private key (only available on creation!)
output "api_private_key" {
  value     = inferadb_client_certificate.api_cert.private_key_pem
  sensitive = true
}
```

## Authentication

The provider requires a session token. Provide it via:

1. Provider configuration: `session_token = "..."`
2. Environment variable: `INFERADB_SESSION_TOKEN`

To obtain a session token, log in via the InferaDB CLI:

```bash
inferadb login
```

## Resources

| Resource                      | Description                                     |
| ----------------------------- | ----------------------------------------------- |
| `inferadb_organization`       | Manages organizations                           |
| `inferadb_vault`              | Manages vaults for authorization policies       |
| `inferadb_client`             | Manages client identities for backend services  |
| `inferadb_client_certificate` | Generates client authentication certificates    |
| `inferadb_team`               | Manages teams for group access control          |
| `inferadb_team_member`        | Manages team memberships                        |
| `inferadb_vault_user_grant`   | Grants users vault access                       |
| `inferadb_vault_team_grant`   | Grants teams vault access                       |

## Data Sources

| Data Source             | Description              |
| ----------------------- | ------------------------ |
| `inferadb_organization` | Reads organization data  |
| `inferadb_vault`        | Reads vault data         |
| `inferadb_client`       | Reads client data        |
| `inferadb_team`         | Reads team data          |

## Development

```bash
# Build
go build -v ./...

# Run tests
go test -v -cover -timeout=120s -parallel=4 ./...

# Run acceptance tests (requires TF_ACC=1 and API credentials)
TF_ACC=1 go test -v -cover -timeout 120m ./...

# Generate documentation
go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

# Format code
gofmt -s -w -e .

# Run linter
golangci-lint run
```

## Related Resources

- [InferaDB Documentation](https://inferadb.com/docs)
- [Control API OpenAPI Spec](../control/openapi.yaml)
- [Engine Terraform Modules](../engine/terraform/)

## Community

Join our [Discord](https://discord.gg/inferadb) to discuss InferaDB, get help, and connect with other developers.

## License

Dual-licensed under [Apache 2.0](LICENSE-APACHE) and [MIT](LICENSE-MIT).
