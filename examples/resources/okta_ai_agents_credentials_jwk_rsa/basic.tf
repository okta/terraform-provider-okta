resource "okta_ai_agent" "test_agent" {
  profile {
    name        = "test-ai-agent-rsa"
    description = "Test AI Agent for RSA Keys"
  }

  resource_url = "https://example.com/api"
}

resource "okta_ai_agents_credentials_jwk_rsa" "test" {
  agent_id = okta_ai_agent.test_agent.id

  kty    = "RSA"
  alg    = "RS256"
  kid    = "test-rsa-key"
  use    = "sig"
  status = "ACTIVE"

  n = "0vx7agoebGcQSuuPiLJXZptN9nndrQmbXEps2aiAFbWhM78LhWx4cbbfAAtVT86zwu1RK7aPFFxuhDR1L6tSoc_BJECPebWKRXjBZCiFV4n3oknjhMstn64tZ_2W-5JsGY4Hc5n9yBXArwl93lqt7_RN5w6Cf0h4QyQ5v-65YGjQR0_FDW2QvzqY368QQMicAtaSqzs8KJZgnYb9c7d0zgdAZHzu6qMQvRL5hajrn1n91CbOpbISD08qNLyrdkt-bFTWhAI4vMQFh6WeZu0fM4lFd2NcRwr3XPksINHaQ-G_xBniIqbw0Ls1jF44-csFCur-kEgU8awapJzKnqDKgw"
  e = "AQAB"
}
