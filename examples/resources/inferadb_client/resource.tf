# Create a client (backend service identity)
resource "inferadb_client" "api_server" {
  organization_id = inferadb_organization.example.id
  vault_id        = inferadb_vault.production.id
  name            = "Production API Server"
}

output "client_id" {
  value = inferadb_client.api_server.id
}
