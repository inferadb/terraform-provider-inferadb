# Add a user to a team
resource "inferadb_team_member" "alice" {
  organization_id = inferadb_organization.example.id
  team_id         = inferadb_team.engineering.id
  user_id         = "123456789012345678" # User must already exist
  role            = "maintainer"         # Options: maintainer, member
}
