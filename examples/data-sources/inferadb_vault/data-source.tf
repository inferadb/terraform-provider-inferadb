# Read an existing vault
data "inferadb_vault" "existing" {
  organization_id = "123456789012345678"
  id              = "987654321098765432"
}

output "vault_name" {
  value = data.inferadb_vault.existing.name
}

output "vault_sync_status" {
  value = data.inferadb_vault.existing.sync_status
}
