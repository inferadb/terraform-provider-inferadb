# Create an organization
resource "inferadb_organization" "example" {
  name = "Acme Corp"
  tier = "dev" # Options: dev, pro, max
}

output "organization_id" {
  value = inferadb_organization.example.id
}
