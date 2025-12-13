# Read an existing organization
data "inferadb_organization" "existing" {
  id = "123456789012345678"
}

output "organization_name" {
  value = data.inferadb_organization.existing.name
}

output "organization_tier" {
  value = data.inferadb_organization.existing.tier
}
