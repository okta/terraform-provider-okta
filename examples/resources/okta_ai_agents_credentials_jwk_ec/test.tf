# Test configuration for okta_ai_agents_credentials_jwk_ec resource
# This shows various test scenarios for the EC JWK resource

# Scenario 1: Basic EC key creation with required fields
resource "okta_ai_agents_credentials_jwk_ec" "test_basic" {
  agent_id = "test-agent-001"
  kty      = "EC"
  alg      = "ES256"
  crv      = "P-256"
  kid      = "test-key-1"
  use      = "sig"
  x        = "gI0GAILBdu7T53-nP0_eAfGLBLO_5K4dHo0_yF2gSMvZD8o5uZvFn5_pZQRyBFe1"
  y        = "w5RiAv0tBZFz9bWGqxdZqP1i8tIDvBu1wnPAi4JNKZ4Ypxqv7Pq8mK9N3L2Q5R6S7"
  status   = "ACTIVE"
}

# Scenario 2: P-384 key
resource "okta_ai_agents_credentials_jwk_ec" "test_p384" {
  agent_id = "test-agent-001"
  kty      = "EC"
  alg      = "ES384"
  crv      = "P-384"
  kid      = "test-key-384"
  use      = "sig"
  x        = "gI0GAILBdu7T53-nP0_eAfGLBLO_5K4dHo0_yF2gSMvZD8o5uZvFn5_pZQRyBFe1gI0GAILBdu7T53"
  y        = "w5RiAv0tBZFz9bWGqxdZqP1i8tIDvBu1wnPAi4JNKZ4Ypxqv7Pq8mK9N3L2Q5R6S7w5RiAv0tBZFz9"
}

# Scenario 3: P-521 key
resource "okta_ai_agents_credentials_jwk_ec" "test_p521" {
  agent_id = "test-agent-001"
  kty      = "EC"
  alg      = "ES512"
  crv      = "P-521"
  kid      = "test-key-521"
  use      = "sig"
  x        = "AHKyLmZN1T4LcC0j8gKW1c5QY8v91R0DXrB9E0cqKoQ7lx4PtLt_4aEz0PeOu-cxEDzR8zGCFjH2-qYYjjJPf33J"
  y        = "AEpqHC1jZTqyrcL9HGg4ER_-B1cdHRDhXo3YGhDhXH_bXHYW9gjY4bDj_lmgJMVKMaH9eHaMf1VFkrXCc3c2J4bD"
}

# Data source to read back the created key
data "okta_ai_agents_credentials_jwk_ec" "test_read" {
  agent_id = "test-agent-001"
  id       = okta_ai_agents_credentials_jwk_ec.test_basic.id
}

# Output for verification
output "test_key_id" {
  value = okta_ai_agents_credentials_jwk_ec.test_basic.id
}

output "test_key_created" {
  value = okta_ai_agents_credentials_jwk_ec.test_basic.created
}

output "test_key_last_updated" {
  value = okta_ai_agents_credentials_jwk_ec.test_basic.last_updated
}

output "test_key_status" {
  value = okta_ai_agents_credentials_jwk_ec.test_basic.status
}

output "test_data_source_result" {
  value = data.okta_ai_agents_credentials_jwk_ec.test_read
}
