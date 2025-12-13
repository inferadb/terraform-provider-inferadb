# Create a vault within an organization
resource "inferadb_vault" "production" {
  organization_id = inferadb_organization.example.id
  name            = "Production Policies"
  description     = "Authorization policies for production environment"
}

output "vault_id" {
  value = inferadb_vault.production.id
}
