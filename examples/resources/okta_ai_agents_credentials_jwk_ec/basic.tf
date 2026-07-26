resource "okta_ai_agent" "test_agent" {
  profile {
    name        = "test-ai-agent-ec"
    description = "Test AI Agent for EC Keys"
  }

  resource_url = "https://example.com/api"
}

resource "okta_ai_agents_credentials_jwk_ec" "test" {
  agent_id = okta_ai_agent.test_agent.id

  kty    = "EC"
  alg    = "ES256"
  kid    = "test-ec-key"
  crv    = "P-256"
  use    = "sig"
  status = "ACTIVE"

  x = "WKn-ZIGevcwGIyyrzFoZNBdaq9_TsqzGl96oc0CWuis"
  y = "y77t-RvAHRKTsSGdIYUfweuOvwrvDD-Q3Hv5J0fSKQE"
}
