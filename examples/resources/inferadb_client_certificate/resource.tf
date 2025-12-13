# Generate a certificate for client authentication
resource "inferadb_client_certificate" "cert_2025" {
  organization_id = inferadb_organization.example.id
  client_id       = inferadb_client.api_server.id
  name            = "Production Certificate 2025"
}

# IMPORTANT: The private key is only available after initial creation
# Store it securely - it cannot be retrieved again!
output "private_key" {
  value     = inferadb_client_certificate.cert_2025.private_key_pem
  sensitive = true
}

output "public_key" {
  value = inferadb_client_certificate.cert_2025.public_key_pem
}

output "key_id" {
  value = inferadb_client_certificate.cert_2025.kid
}
