# Sample Terraform Config for okta_ai_agents_credentials_jwk_rsa
# With realistic RSA key values for testing

terraform {
  required_providers {
    okta = {
      source  = "okta/okta"
      version = ">= 4.0"
    }
  }
}

provider "okta" {
  org_name  = "dev-12345678"
  api_token = var.okta_api_token
}

variable "okta_api_token" {
  type      = string
  sensitive = true
}

# ============================================================================
# Example 1: RS256 Signing Key (RSA-2048)
# ============================================================================

resource "okta_ai_agents_credentials_jwk_rsa" "rs256_signing" {
  agent_id = "wlp119dzl3aF4KKOc1d8"

  kty    = "RSA"
  alg    = "RS256"
  kid    = "rsa-2048-rs256-sig"
  use    = "sig"
  status = "ACTIVE"

  # RSA-2048 Public Key Modulus (base64url encoded)
  n = "0vx7agoebGcQSuuPiLJXZptN9nndrQmbXEps2aiAFbWhM78LhWx4cbbfAAtVT86zwu1RK7aPFFxuhDR1L6tSoc_BJECPebWKRXjBZCiFV4n3oknjhMstn64tZ_2W-5JsGY4Hc5n9yBXArwl93lqt7_RN5w6Cf0h4QyQ5v-65YGjQR0_FDW2QvzqY368QQMicAtaSqzs8KJZgnYb9c7d0zgdAZHzu6qMQvRL5hajrn1n91CbOpbISD08qNLyrdkt-bFTWhAI4vMQFh6WeZu0fM4lFd2NcRwr3XPksINHaQ-G_xBniIqbw0Ls1jF44-csFCur-kEgU8awapJzKnqDKgw"

  # Standard RSA Exponent (65537 = 0x10001 in base64url)
  e = "AQAB"
}

# ============================================================================
# Example 2: RS384 Signing Key (RSA-2048)
# ============================================================================

resource "okta_ai_agents_credentials_jwk_rsa" "rs384_signing" {
  agent_id = "wlp119dzl3aF4KKOc1d8"

  kty    = "RSA"
  alg    = "RS384"
  kid    = "rsa-2048-rs384-sig"
  use    = "sig"
  status = "ACTIVE"

  n = "0vx7agoebGcQSuuPiLJXZptN9nndrQmbXEps2aiAFbWhM78LhWx4cbbfAAtVT86zwu1RK7aPFFxuhDR1L6tSoc_BJECPebWKRXjBZCiFV4n3oknjhMstn64tZ_2W-5JsGY4Hc5n9yBXArwl93lqt7_RN5w6Cf0h4QyQ5v-65YGjQR0_FDW2QvzqY368QQMicAtaSqzs8KJZgnYb9c7d0zgdAZHzu6qMQvRL5hajrn1n91CbOpbISD08qNLyrdkt-bFTWhAI4vMQFh6WeZu0fM4lFd2NcRwr3XPksINHaQ-G_xBniIqbw0Ls1jF44-csFCur-kEgU8awapJzKnqDKgw"
  e = "AQAB"
}

# ============================================================================
# Example 3: RS512 Signing Key (RSA-2048)
# ============================================================================

resource "okta_ai_agents_credentials_jwk_rsa" "rs512_signing" {
  agent_id = "wlp119dzl3aF4KKOc1d8"

  kty    = "RSA"
  alg    = "RS512"
  kid    = "rsa-2048-rs512-sig"
  use    = "sig"
  status = "ACTIVE"

  n = "0vx7agoebGcQSuuPiLJXZptN9nndrQmbXEps2aiAFbWhM78LhWx4cbbfAAtVT86zwu1RK7aPFFxuhDR1L6tSoc_BJECPebWKRXjBZCiFV4n3oknjhMstn64tZ_2W-5JsGY4Hc5n9yBXArwl93lqt7_RN5w6Cf0h4QyQ5v-65YGjQR0_FDW2QvzqY368QQMicAtaSqzs8KJZgnYb9c7d0zgdAZHzu6qMQvRL5hajrn1n91CbOpbISD08qNLyrdkt-bFTWhAI4vMQFh6WeZu0fM4lFd2NcRwr3XPksINHaQ-G_xBniIqbw0Ls1jF44-csFCur-kEgU8awapJzKnqDKgw"
  e = "AQAB"
}

# ============================================================================
# Example 4: PS256 Signing Key (PSS Padding - Recommended)
# ============================================================================

resource "okta_ai_agents_credentials_jwk_rsa" "ps256_signing" {
  agent_id = "wlp119dzl3aF4KKOc1d8"

  kty    = "RSA"
  alg    = "PS256"
  kid    = "rsa-4096-ps256-sig"
  use    = "sig"
  status = "ACTIVE"

  # 4096-bit RSA Public Key Modulus (base64url encoded)
  n = "xjlCRBqkQRKkTr1wKKKs7v8xJyRTt8E7ItLk4R8vK5jQ0vK4mN1pR2sT3uV4wX5yZ6aB7cD8eF9gH0iJ1kL2mN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9"
  e = "AQAB"
}

# ============================================================================
# Example 5: PS384 Signing Key (PSS Padding)
# ============================================================================

resource "okta_ai_agents_credentials_jwk_rsa" "ps384_signing" {
  agent_id = "wlp119dzl3aF4KKOc1d8"

  kty    = "RSA"
  alg    = "PS384"
  kid    = "rsa-4096-ps384-sig"
  use    = "sig"
  status = "ACTIVE"

  n = "xjlCRBqkQRKkTr1wKKKs7v8xJyRTt8E7ItLk4R8vK5jQ0vK4mN1pR2sT3uV4wX5yZ6aB7cD8eF9gH0iJ1kL2mN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9"
  e = "AQAB"
}

# ============================================================================
# Example 6: PS512 Signing Key (PSS Padding)
# ============================================================================

resource "okta_ai_agents_credentials_jwk_rsa" "ps512_signing" {
  agent_id = "wlp119dzl3aF4KKOc1d8"

  kty    = "RSA"
  alg    = "PS512"
  kid    = "rsa-4096-ps512-sig"
  use    = "sig"
  status = "ACTIVE"

  n = "xjlCRBqkQRKkTr1wKKKs7v8xJyRTt8E7ItLk4R8vK5jQ0vK4mN1pR2sT3uV4wX5yZ6aB7cD8eF9gH0iJ1kL2mN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9aA0bB1cC2dD3eE4fF5gG6hH7iI8jJ9kK0lL1mM2nN3oP4qR5sT6uV7wX8yZ9"
  e = "AQAB"
}

# ============================================================================
# Example 7: RSA-OAEP Encryption Key
# ============================================================================

resource "okta_ai_agents_credentials_jwk_rsa" "rsa_oaep_encryption" {
  agent_id = "wlp119dzl3aF4KKOc1d8"

  kty    = "RSA"
  alg    = "RSA-OAEP"
  kid    = "rsa-2048-oaep-enc"
  use    = "enc"
  status = "ACTIVE"

  n = "0vx7agoebGcQSuuPiLJXZptN9nndrQmbXEps2aiAFbWhM78LhWx4cbbfAAtVT86zwu1RK7aPFFxuhDR1L6tSoc_BJECPebWKRXjBZCiFV4n3oknjhMstn64tZ_2W-5JsGY4Hc5n9yBXArwl93lqt7_RN5w6Cf0h4QyQ5v-65YGjQR0_FDW2QvzqY368QQMicAtaSqzs8KJZgnYb9c7d0zgdAZHzu6qMQvRL5hajrn1n91CbOpbISD08qNLyrdkt-bFTWhAI4vMQFh6WeZu0fM4lFd2NcRwr3XPksINHaQ-G_xBniIqbw0Ls1jF44-csFCur-kEgU8awapJzKnqDKgw"
  e = "AQAB"
}

# ============================================================================
# Example 8: Read existing key
# ============================================================================

data "okta_ai_agents_credentials_jwk_rsa" "read_signing_key" {
  agent_id = "wlp119dzl3aF4KKOc1d8"
  id       = okta_ai_agents_credentials_jwk_rsa.rs256_signing.id
}

# ============================================================================
# Outputs
# ============================================================================

output "rs256_key" {
  value = {
    id      = okta_ai_agents_credentials_jwk_rsa.rs256_signing.id
    kid     = okta_ai_agents_credentials_jwk_rsa.rs256_signing.kid
    alg     = okta_ai_agents_credentials_jwk_rsa.rs256_signing.alg
    use     = okta_ai_agents_credentials_jwk_rsa.rs256_signing.use
    created = okta_ai_agents_credentials_jwk_rsa.rs256_signing.created
    status  = okta_ai_agents_credentials_jwk_rsa.rs256_signing.status
  }
}

output "ps256_key" {
  value = {
    id  = okta_ai_agents_credentials_jwk_rsa.ps256_signing.id
    kid = okta_ai_agents_credentials_jwk_rsa.ps256_signing.kid
    alg = okta_ai_agents_credentials_jwk_rsa.ps256_signing.alg
  }
}

output "encryption_key" {
  value = {
    id  = okta_ai_agents_credentials_jwk_rsa.rsa_oaep_encryption.id
    kid = okta_ai_agents_credentials_jwk_rsa.rsa_oaep_encryption.kid
    alg = okta_ai_agents_credentials_jwk_rsa.rsa_oaep_encryption.alg
    use = okta_ai_agents_credentials_jwk_rsa.rsa_oaep_encryption.use
  }
}

output "read_key_data" {
  value = {
    id      = data.okta_ai_agents_credentials_jwk_rsa.read_signing_key.id
    kid     = data.okta_ai_agents_credentials_jwk_rsa.read_signing_key.kid
    created = data.okta_ai_agents_credentials_jwk_rsa.read_signing_key.created
  }
}
