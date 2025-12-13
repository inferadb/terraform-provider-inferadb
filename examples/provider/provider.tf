# Configure the InferaDB provider
provider "inferadb" {
  # The API endpoint (defaults to https://api.inferadb.com)
  # Can also be set via INFERADB_ENDPOINT environment variable
  endpoint = "https://api.inferadb.com"

  # Session token for authentication
  # Can also be set via INFERADB_SESSION_TOKEN environment variable
  # Obtain this by logging in via the InferaDB CLI: inferadb login
  session_token = var.inferadb_session_token
}

variable "inferadb_session_token" {
  type        = string
  description = "InferaDB session token for authentication"
  sensitive   = true
}
