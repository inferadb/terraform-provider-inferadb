# Read an existing client
data "inferadb_client" "existing" {
  organization_id = "123456789012345678"
  id              = "456789012345678901"
}

output "client_name" {
  value = data.inferadb_client.existing.name
}

output "client_is_active" {
  value = data.inferadb_client.existing.is_active
}
