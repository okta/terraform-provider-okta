# Complete Examples for okta_ai_agents_credentials_jwk_rsa

terraform {
  required_providers {
    okta = {
      source  = "okta/okta"
      version = ">= 4.0"
    }
  }
}

provider "okta" {
  org_name  = "your-org"
  api_token = var.okta_api_token
}

variable "okta_api_token" {
  type      = string
  sensitive = true
}

# ============================================================================
# Example 1: Basic RSA-2048 Signing Key (RS256)
# ============================================================================

resource "okta_ai_agents_credentials_jwk_rsa" "basic_signing" {
  agent_id = "wlp119dzl3aF4KKOc1d8"  # Replace with your AI agent ID

  kty = "RSA"
  alg = "RS256"
  kid = "basic-rsa-signing-key"
  use = "sig"

  n = "0vx7agoebGcQSuuPiLJXZptN9nndrQmbXEps2aiAFbWhM78LhWx4cbbfAAtVT86zwu1RK7aPFFxuhDR1L6tSoc_BJECPebWKRXjBZCiFV4n3oknjhMstn64tZ_2W-5JsGY4Hc5n9yBXArwl93lqt7_RN5w6Cf0h4QyQ5v-65YGjQR0_FDW2QvzqY368QQMicAtaSqzs8KJZgnYb9c7d0zgdAZHzu6qMQvRL5hajrn1n91CbOpbISD08qNLyrdkt-bFTWhAI4vMQFh6WeZu0fM4lFd2NcRwr3XPksINHaQ-G_xBniIqbw0Ls1jF44-csFCur-kEgU8awapJzKnqDKgw"
  e = "AQAB"

  status = "ACTIVE"
}

# ============================================================================
# Example 2: RSA-2048 Encryption Key (RSA-OAEP)
# ============================================================================

resource "okta_ai_agents_credentials_jwk_rsa" "encryption_key" {
  agent_id = "wlp119dzl3aF4KKOc1d8"

  kty = "RSA"
  alg = "RSA-OAEP"
  kid = "rsa-encryption-key"
  use = "enc"

  n = "0vx7agoebGcQSuuPiLJXZptN9nndrQmbXEps2aiAFbWhM78LhWx4cbbfAAtVT86zwu1RK7aPFFxuhDR1L6tSoc_BJECPebWKRXjBZCiFV4n3oknjhMstn64tZ_2W-5JsGY4Hc5n9yBXArwl93lqt7_RN5w6Cf0h4QyQ5v-65YGjQR0_FDW2QvzqY368QQMicAtaSqzs8KJZgnYb9c7d0zgdAZHzu6qMQvRL5hajrn1n91CbOpbISD08qNLyrdkt-bFTWhAI4vMQFh6WeZu0fM4lFd2NcRwr3XPksINHaQ-G_xBniIqbw0Ls1jF44-csFCur-kEgU8awapJzKnqDKgw"
  e = "AQAB"

  status = "ACTIVE"
}

# ============================================================================
# Example 3: RSA-4096 with PS256 (PSS padding - Recommended)
# ============================================================================

resource "okta_ai_agents_credentials_jwk_rsa" "pss_signing" {
  agent_id = "wlp119dzl3aF4KKOc1d8"

  kty = "RSA"
  alg = "PS256"      # PSS padding with SHA-256 (more secure than PKCS#1 v1.5)
  kid = "rsa-4096-ps256"
  use = "sig"

  # 4096-bit RSA modulus (base64url encoded, no padding)
  # Replace with your actual 4096-bit public key modulus
  n = "xjlCRBqkQRKkTr1wKKKs7v8xJyRTt8E7ItLk4R8vK5jQ0vK4mN1pR2sT3uV4wX5yZ6aB7cD8eF9gH0iJ1kL2mN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oO4pP5qQ6rR7sS8tT9uU0vV1wW2xX3yY4zZ5aA6bB7cC8dD9eE0fF1gG2hH3iI4jJ5kK6lL7mM8nN9oO0pP5qQ6rR7sS8tT9uU0vV1wW2xX3yY4zZ5aA6bB7cC8dD9eE0fF1gG2hH3iI4jJ5kK6lL7mM8nN9oO0pP5qQ6rR7sS8tT9uU0vV1wW2xX3yY4zZ5aA6bB7cC8dD9eE0fF1gG2hH3iI4jJ5kK6lL7mM8nN9oO0pP5qQ6rR7sS8tT9uU0vV1wW2xX3yY4zZ5aA6bB7cC8dD9eE0fF1gG2hH3iI4jJ5kK6lL7mM8nN9oO0pP5qQ6rR7sS8tT9uU0vV1wW2xX3yY4zZ5aA6bB7cC8dD9eE0fF1gG2hH3iI4jJ5kK6lL7mM8nN9oO0pP5qQ6rR7sS8tT9uU0vV1wW2xX3yY4zZ5aA6bB7cC8dD9eE0fF1gG2hH3iI4jJ5kK6lL7mM8nN9oO0pP5qQ6rR7sS8tT9uU0vV1wW2xX3yY"
  e = "AQAB"

  status = "ACTIVE"
}

# ============================================================================
# Example 4: Read back existing key using data source
# ============================================================================

data "okta_ai_agents_credentials_jwk_rsa" "existing" {
  agent_id = "wlp119dzl3aF4KKOc1d8"
  id       = okta_ai_agents_credentials_jwk_rsa.basic_signing.id
}

# ============================================================================
# Outputs
# ============================================================================

output "signing_key_id" {
  description = "ID of the RSA signing key"
  value       = okta_ai_agents_credentials_jwk_rsa.basic_signing.id
}

output "signing_key_kid" {
  description = "Key ID of the RSA signing key"
  value       = okta_ai_agents_credentials_jwk_rsa.basic_signing.kid
}

output "signing_key_created" {
  description = "Creation timestamp of signing key"
  value       = okta_ai_agents_credentials_jwk_rsa.basic_signing.created
}

output "encryption_key_id" {
  description = "ID of the RSA encryption key"
  value       = okta_ai_agents_credentials_jwk_rsa.encryption_key.id
}

output "encryption_key_status" {
  description = "Status of the encryption key"
  value       = okta_ai_agents_credentials_jwk_rsa.encryption_key.status
}

output "pss_key_algorithm" {
  description = "Algorithm used by the PSS key"
  value       = okta_ai_agents_credentials_jwk_rsa.pss_signing.alg
}

output "read_key_data" {
  description = "Data read from existing key"
  value = {
    id      = data.okta_ai_agents_credentials_jwk_rsa.existing.id
    kid     = data.okta_ai_agents_credentials_jwk_rsa.existing.kid
    alg     = data.okta_ai_agents_credentials_jwk_rsa.existing.alg
    created = data.okta_ai_agents_credentials_jwk_rsa.existing.created
  }
}

# ============================================================================
# Usage Instructions
# ============================================================================
#
# 1. Replace "wlp119dzl3aF4KKOc1d8" with your AI Agent ID
# 2. Replace modulus and exponent values with your RSA public key components
# 3. Set your Okta API token:
#    export TF_VAR_okta_api_token="your-api-token"
# 4. Initialize Terraform:
#    terraform init
# 5. Plan and apply:
#    terraform plan
#    terraform apply
# 6. Import existing key:
#    terraform import okta_ai_agents_credentials_jwk_rsa.basic_signing {agent_id}/{key_id}
#
# ============================================================================
