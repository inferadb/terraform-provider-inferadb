# Create a team for group-based access control
resource "inferadb_team" "engineering" {
  organization_id = inferadb_organization.example.id
  name            = "Engineering"
  description     = "Engineering team with production access"
}

output "team_id" {
  value = inferadb_team.engineering.id
}
