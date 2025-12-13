# Grant a user access to a vault
resource "inferadb_vault_user_grant" "alice_admin" {
  organization_id = inferadb_organization.example.id
  vault_id        = inferadb_vault.production.id
  user_id         = "123456789012345678"
  role            = "admin" # Options: reader, writer, manager, admin
}
