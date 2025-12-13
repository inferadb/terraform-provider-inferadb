# Grant a team access to a vault
resource "inferadb_vault_team_grant" "engineering_writer" {
  organization_id = inferadb_organization.example.id
  vault_id        = inferadb_vault.production.id
  team_id         = inferadb_team.engineering.id
  role            = "writer" # Options: reader, writer, manager, admin
}
