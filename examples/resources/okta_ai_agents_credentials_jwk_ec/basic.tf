resource "okta_ai_agent" "test" {
  profile {
    name        = "testAcc-jwk-ec-replace_with_uuid"
    description = "Terraform acceptance test agent for okta_ai_agents_credentials_jwk_ec"
  }

  resource_url = "https://example.com/api"

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
