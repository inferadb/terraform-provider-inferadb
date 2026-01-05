<div align="center">
    <p><a href="https://inferadb.com"><img src=".github/inferadb.png" width="100" alt="InferaDB Logo" /></a></p>
    <h1>InferaDB Terraform Provider</h1>
    <p>Manage InferaDB resources including organizations, vaults, clients, teams, and access grants</p>
</div>

> [!IMPORTANT]
> Under active development. Not production-ready.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23 (for building from source)

## Installation

### From Source

```bash
git clone https://github.com/inferadb/terraform-provider-inferadb.git
cd terraform-provider-inferadb
make install
```

### Local Development

```bash
make dev-install
```

This installs the provider to `~/.terraform.d/plugins/` for local testing.

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

The provider requires a session token for authentication. You can provide it via:

1. Provider configuration: `session_token = "..."`
2. Environment variable: `INFERADB_SESSION_TOKEN`

To obtain a session token, log in via the InferaDB CLI:

```bash
inferadb login
```

## Resources

| Resource                      | Description                                              |
| ----------------------------- | -------------------------------------------------------- |
| `inferadb_organization`       | Manages organizations (multi-tenant containers)          |
| `inferadb_vault`              | Manages vaults (authorization policy storage)            |
| `inferadb_client`             | Manages clients (backend service identities)             |
| `inferadb_client_certificate` | Generates Ed25519 certificates for client authentication |
| `inferadb_team`               | Manages teams for group-based access control             |
| `inferadb_team_member`        | Manages team memberships                                 |
| `inferadb_vault_user_grant`   | Grants user access to vaults                             |
| `inferadb_vault_team_grant`   | Grants team access to vaults                             |

## Data Sources

| Data Source             | Description               |
| ----------------------- | ------------------------- |
| `inferadb_organization` | Read organization details |
| `inferadb_vault`        | Read vault details        |
| `inferadb_client`       | Read client details       |
| `inferadb_team`         | Read team details         |

## Development

```bash
# Build
make build

# Run tests
make test

# Run acceptance tests (requires TF_ACC=1 and API credentials)
make testacc

# Generate documentation
make docs

# Format code
make fmt

# Run linter
make lint
```

## Related Resources

- [InferaDB Documentation](https://inferadb.com/docs)
- [Control API OpenAPI Spec](../control/openapi.yaml)
- [Engine Terraform Modules](../engine/terraform/)

## License

Apache License 2.0 - see [LICENSE](LICENSE) for details.
