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

# Okta refuses to deactivate an agent's last active key, and only deletes inactive keys.
# Forget the key without destroying it; deleting the agent at test cleanup removes it.
removed {
  from = okta_ai_agents_credentials_jwk_ec.test

  lifecycle {
    destroy = false
  }
}
