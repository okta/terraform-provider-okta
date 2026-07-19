# Sample Terraform Configuration for okta_ai_agents_credentials_jwk_rsa

terraform {
  required_providers {
    okta = {
      source  = "okta/okta"
      version = ">= 4.0"
    }
  }
}

provider "okta" {
  org_name  = "your-org"           # Replace with your org name
  api_token = var.okta_api_token   # Set via environment: export TF_VAR_okta_api_token="..."
}

variable "okta_api_token" {
  description = "Okta API token"
  type        = string
  sensitive   = true
}

# ============================================================================
# Example 1: Basic RSA-2048 Key with RS256 Algorithm
# ============================================================================

resource "okta_ai_agents_credentials_jwk_rsa" "example_2048" {
  agent_id = "00u1a2b3c4d5e6f7g8h9i"  # Replace with your AI agent ID

  kty = "RSA"        # Key type (required, must be "RSA")
  alg = "RS256"      # Algorithm (RS256, RS384, RS512, PS256, PS384, PS512)
  kid = "rsa-key-2048"
  use = "sig"        # Usage (sig = signing, enc = encryption)

  # RSA public key components (base64url-encoded, no padding)
  # n = modulus (the public exponent part of the key)
  # e = exponent (typically "AQAB" for RSA-2048)

  n = "0vx7agoebGcQSuuPiLJXZptN9nndrQmbXEps2aiAFbWhM78LhWx4cbbfAAtVT86zwu1RK7aPFFxuhDR1L6tSoc_BJECPebWKRXjBZCiFV4n3oknjhMstn64tZ_2W-5JsGY4Hc5n9yBXArwl93lqt7_RN5w6Cf0h4QyQ5v-65YGjQR0_FDW2QvzqY368QQMicAtaSqzs8KJZgnYb9c7d0zgdAZHzu6qMQvRL5hajrn1n91CbOpbISD08qNLyrdkt-bFTWhAI4vMQFh6WeZu0fM4lFd2NcRwr3XPksINHaQ-G_xBniIqbw0Ls1jF44-csFCur-kEgU8awapJzKnqDKgw"
  e = "AQAB"

  status = "ACTIVE"
}

# ============================================================================
# Example 2: RSA-4096 Key with PS256 Algorithm (Recommended for new keys)
# ============================================================================

resource "okta_ai_agents_credentials_jwk_rsa" "example_4096" {
  agent_id = "00u1a2b3c4d5e6f7g8h9i"

  kty = "RSA"
  alg = "PS256"      # PSS padding with SHA-256 (more secure than PKCS#1 v1.5)
  kid = "rsa-key-4096-ps256"
  use = "sig"

  # Longer modulus for 4096-bit key
  n = "xjlCRBqkQRKkTr1wKKKs...VERY_LONG_BASE64_STRING...Zf7gY"  # Replace with actual 4096-bit modulus
  e = "AQAB"

  status = "ACTIVE"
}

# ============================================================================
# Example 3: RSA Key for Encryption
# ============================================================================

resource "okta_ai_agents_credentials_jwk_rsa" "encryption_key" {
  agent_id = "00u1a2b3c4d5e6f7g8h9i"

  kty = "RSA"
  alg = "RSA-OAEP"   # RSA with OAEP padding for encryption
  kid = "rsa-encryption-key"
  use = "enc"        # For encryption instead of signing

  n = "0vx7agoebGcQSuuPiLJXZptN9nndrQmbXEps2aiAFbWhM78LhWx4cbbfAAtVT86zwu1RK7aPFFxuhDR1L6tSoc_BJECPebWKRXjBZCiFV4n3oknjhMstn64tZ_2W-5JsGY4Hc5n9yBXArwl93lqt7_RN5w6Cf0h4QyQ5v-65YGjQR0_FDW2QvzqY368QQMicAtaSqzs8KJZgnYb9c7d0zgdAZHzu6qMQvRL5hajrn1n91CbOpbISD08qNLyrdkt-bFTWhAI4vMQFh6WeZu0fM4lFd2NcRwr3XPksINHaQ-G_xBniIqbw0Ls1jF44-csFCur-kEgU8awapJzKnqDKgw"
  e = "AQAB"

  status = "ACTIVE"
}

# ============================================================================
# Data Source: Read an existing RSA key
# ============================================================================

data "okta_ai_agents_credentials_jwk_rsa" "existing_key" {
  agent_id = "00u1a2b3c4d5e6f7g8h9i"
  id       = okta_ai_agents_credentials_jwk_rsa.example_2048.id
}

# ============================================================================
# Outputs
# ============================================================================

output "rsa_key_id" {
  description = "The ID of the created RSA JSON Web Key"
  value       = okta_ai_agents_credentials_jwk_rsa.example_2048.id
}

output "rsa_key_kid" {
  description = "The key ID (kid) of the RSA key"
  value       = okta_ai_agents_credentials_jwk_rsa.example_2048.kid
}

output "rsa_key_created" {
  description = "Timestamp of when the key was created"
  value       = okta_ai_agents_credentials_jwk_rsa.example_2048.created
}

output "rsa_key_algorithm" {
  description = "The algorithm used by the key"
  value       = okta_ai_agents_credentials_jwk_rsa.example_2048.alg
}

output "rsa_key_status" {
  description = "Current status of the key"
  value       = okta_ai_agents_credentials_jwk_rsa.example_2048.status
}

# ============================================================================
# Usage Instructions:
# ============================================================================
#
# 1. Set your Okta API token:
#    export TF_VAR_okta_api_token="your-api-token"
#
# 2. Initialize Terraform:
#    terraform init
#
# 3. Plan the deployment:
#    terraform plan
#
# 4. Apply the configuration:
#    terraform apply
#
# 5. Verify the keys in Okta Console:
#    Admin > Applications > AI Agents > [Your Agent] > Keys
#
# 6. Import an existing key:
#    terraform import okta_ai_agents_credentials_jwk_rsa.example_2048 00u1a2b3c4d5e6f7g8h9i/00c1a2b3c4d5e6f7g8h9i
#
# ============================================================================
