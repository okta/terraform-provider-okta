resource "okta_ai_agent" "test" {
  profile {
    name        = "testAcc-jwks-uri-replace_with_uuid"
    description = "Terraform acceptance test agent for okta_ai_agents_credentials_jwks_uri"
  }

  sign_on_provider {
    type = "NEW_OIDC_APP"
  }
}

resource "okta_ai_agents_credentials_jwks_uri" "test" {
  agent_id = okta_ai_agent.test.id
  jwks_uri = "https://example.com/.well-known/jwks.json"
}
