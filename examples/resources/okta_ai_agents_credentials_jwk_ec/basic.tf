resource "okta_ai_agent" "test" {
  profile {
    name        = "testAcc-jwk-ec-replace_with_uuid"
    description = "Terraform acceptance test agent for okta_ai_agents_credentials_jwk_ec"
  }

  resource_url = "https://testacc-jwk-ec-replace_with_uuid.example.com/api"

  sign_on_provider {
    type = "NEW_OIDC_APP"
  }
}

resource "okta_ai_agents_credentials_jwk_ec" "test" {
  agent_id = okta_ai_agent.test.id

  kty    = "EC"
  alg    = "ES256"
  kid    = "testacc-ec-replace_with_uuid"
  crv    = "P-256"
  use    = "sig"
  status = "ACTIVE"

  x = "WKn-ZIGevcwGIyyrzFoZNBdaq9_TsqzGl96oc0CWuis"
  y = "y77t-RvAHRKTsSGdIYUfweuOvwrvDD-Q3Hv5J0fSKQE"
}

# A second active key, so the agent still has one active key after this one is destroyed.
# Okta won't deactivate an agent's last active key, so a key can only be deleted while another
# stays active.
resource "okta_ai_agents_credentials_jwk_ec" "test_second" {
  agent_id = okta_ai_agent.test.id

  kty    = "EC"
  alg    = "ES256"
  kid    = "testacc-ec-second-replace_with_uuid"
  crv    = "P-256"
  use    = "sig"
  status = "ACTIVE"

  x = "6a8IsP1QXnW191GI_z3vkO4rYl4Ak2kiMytZv7OL34k"
  y = "dmUCk20O_3WmkyIhh0h2b9bUWSR4eqwuCgGL9_8T7f0"
}
