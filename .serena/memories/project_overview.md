# InferaDB Terraform Provider - Project Overview

## Purpose
This is an official Terraform provider for InferaDB, a service for managing authorization policies. The provider enables Infrastructure-as-Code management of InferaDB resources via the Terraform ecosystem.

## Tech Stack
- **Language**: Go 1.24+
- **Framework**: HashiCorp Terraform Plugin Framework v1.15+
- **Testing**: Terraform Plugin Testing v1.14+ with acceptance tests
- **Documentation**: terraform-plugin-docs for automated doc generation
- **Release**: GoReleaser for cross-platform builds

## Managed Resources
- `inferadb_organization` - Multi-tenant containers with tier configuration
- `inferadb_vault` - Authorization policy storage
- `inferadb_client` - Backend service identities
- `inferadb_client_certificate` - Ed25519 certificates for client auth
- `inferadb_team` - Group-based access control
- `inferadb_team_member` - Team membership management
- `inferadb_vault_user_grant` - User-level vault access
- `inferadb_vault_team_grant` - Team-level vault access

## Data Sources
- `inferadb_organization`, `inferadb_vault`, `inferadb_client`, `inferadb_team`

## Authentication
The provider uses session tokens via:
1. Provider config: `session_token = "..."`
2. Environment variable: `INFERADB_SESSION_TOKEN`

## Status
⚠️ Under active development - not production-ready.
