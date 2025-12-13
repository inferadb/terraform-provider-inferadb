# Read an existing team
data "inferadb_team" "existing" {
  organization_id = "123456789012345678"
  id              = "789012345678901234"
}

output "team_name" {
  value = data.inferadb_team.existing.name
}

output "team_description" {
  value = data.inferadb_team.existing.description
}
