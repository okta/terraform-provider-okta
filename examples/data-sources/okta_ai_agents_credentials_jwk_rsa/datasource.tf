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

resource "okta_ai_agents_credentials_jwk_rsa" "test" {
  agent_id = okta_ai_agent.test.id

  kty    = "RSA"
  alg    = "RS256"
  kid    = "testacc-rsa-replace_with_uuid"
  use    = "sig"
  status = "ACTIVE"

  n = "0vx7agoebGcQSuuPiLJXZptN9nndrQmbXEps2aiAFbWhM78LhWx4cbbfAAtVT86zwu1RK7aPFFxuhDR1L6tSoc_BJECPebWKRXjBZCiFV4n3oknjhMstn64tZ_2W-5JsGY4Hc5n9yBXArwl93lqt7_RN5w6Cf0h4QyQ5v-65YGjQR0_FDW2QvzqY368QQMicAtaSqzs8KJZgnYb9c7d0zgdAZHzu6qMQvRL5hajrn1n91CbOpbISD08qNLyrdkt-bFTWhAI4vMQFh6WeZu0fM4lFd2NcRwr3XPksINHaQ-G_xBniIqbw0Ls1jF44-csFCur-kEgU8awapJzKnqDKgw"
  e = "AQAB"
}

data "okta_ai_agents_credentials_jwk_rsa" "test" {
  id       = okta_ai_agents_credentials_jwk_rsa.test.id
  agent_id = okta_ai_agent.test.id
}
