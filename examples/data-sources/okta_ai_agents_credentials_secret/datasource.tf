resource "okta_ai_agent" "test" {
  profile {
    name        = "testAcc-secret-replace_with_uuid"
    description = "Terraform acceptance test agent for the secret data source"
  }

  resource_url = "https://example.com/api"

  sign_on_provider {
    type = "NEW_OIDC_APP"
  }
}

resource "okta_ai_agents_credentials_secret" "test" {
  agent_id = okta_ai_agent.test.id
  status   = "ACTIVE"
}

data "okta_ai_agents_credentials_secret" "test" {
  id       = okta_ai_agents_credentials_secret.test.id
  agent_id = okta_ai_agent.test.id
}
