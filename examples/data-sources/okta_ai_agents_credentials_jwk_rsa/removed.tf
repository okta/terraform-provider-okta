resource "okta_ai_agent" "test" {
  profile {
    name        = "testAcc-jwk-rsa-replace_with_uuid"
    description = "Terraform acceptance test agent for the jwk-rsa data source"
  }

  resource_url = "https://testacc-jwk-rsa-ds-replace_with_uuid.example.com/api"

  sign_on_provider {
    type = "NEW_OIDC_APP"
  }
}

# Okta refuses to deactivate an agent's last active key, and only deletes inactive keys.
# Forget the key without destroying it; deleting the agent at test cleanup removes it.
removed {
  from = okta_ai_agents_credentials_jwk_rsa.test

  lifecycle {
    destroy = false
  }
}
