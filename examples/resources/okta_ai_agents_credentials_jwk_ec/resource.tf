# Example: Create an EC JSON Web Key for an AI Agent

resource "okta_ai_agents_credentials_jwk_ec" "example" {
  agent_id = "YOUR_AGENT_ID"  # ID of the AI agent in Okta

  # Required discriminator field - must be "EC" for EC keys
  kty = "EC"

  # Algorithm used in the JSON Web Key (e.g., "ES256", "ES384", "ES512")
  alg = "ES256"

  # The cryptographic curve that's used for the key pair
  # Common values: "P-256", "P-384", "P-521"
  crv = "P-256"

  # The unique identifier of the JSON Web Key in the JWKS
  kid = "my-key-id-1"

  # Acceptable use of the JSON Web Key (e.g., "sig" for signing, "enc" for encryption)
  use = "sig"

  # Public x coordinate for the elliptic curve point
  # This is a base64url-encoded value from your EC public key
  x = "YOUR_X_COORDINATE"

  # Public y coordinate for the elliptic curve point
  # This is a base64url-encoded value from your EC public key
  y = "YOUR_Y_COORDINATE"

  # Status of the JSON Web Key (optional)
  # status = "ACTIVE"
}

# Example with multiple EC keys for the same agent
resource "okta_ai_agents_credentials_jwk_ec" "backup_key" {
  agent_id = "YOUR_AGENT_ID"

  kty = "EC"
  alg = "ES256"
  crv = "P-256"
  kid = "backup-key-1"
  use = "sig"

  x = "BACKUP_X_COORDINATE"
  y = "BACKUP_Y_COORDINATE"
}

# Output the created key ID for reference
output "jwk_ec_key_id" {
  value       = okta_ai_agents_credentials_jwk_ec.example.id
  description = "The ID of the created EC JSON Web Key"
}

output "jwk_ec_kid" {
  value       = okta_ai_agents_credentials_jwk_ec.example.kid
  description = "The key ID (kid) of the EC JSON Web Key"
}

output "jwk_ec_created" {
  value       = okta_ai_agents_credentials_jwk_ec.example.created
  description = "Timestamp of when the key was created"
}

# Data source to retrieve an existing EC key
data "okta_ai_agents_credentials_jwk_ec" "example" {
  agent_id = "YOUR_AGENT_ID"
  id       = okta_ai_agents_credentials_jwk_ec.example.id
}
