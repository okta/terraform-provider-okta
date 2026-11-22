resource "okta_ai_agent" "test" {
  profile {
    name        = "testAcc-jwk-rsa-replace_with_uuid"
    description = "Terraform acceptance test agent for okta_ai_agents_credentials_jwk_rsa"
  }

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

# A second active key, so the agent still has one active key after this one is destroyed.
# Okta won't deactivate an agent's last active key, so a key can only be deleted while another
# stays active.
resource "okta_ai_agents_credentials_jwk_rsa" "test_second" {
  agent_id = okta_ai_agent.test.id

  kty    = "RSA"
  alg    = "RS256"
  kid    = "testacc-rsa-second-replace_with_uuid"
  use    = "sig"
  status = "ACTIVE"

  n = "u0RK0ouLyfd8SfiIuEX1Ypet16uPbOZDZpiqfmOHVBhz8iHUQK54U2RSTkzn5iWcJXGVMrvb7EpRuHRIdB4oIN1hVpBM68FcIcLjwnof8U-i8oIq-22WsL3I7R3HJH6wdhsJ4_AxTmXWiMF1fzNwnXK6HMfp9SqZxVlZ4Ulaxou69NpxTCmg_8DOgjapZagxt1HHy-qeFDpn_Ui4b_jv9_3zgxkY54ZzItM2jARWxM6a-cry8a0pPFmSTusCUh4Vig4etIchzJCTyvVXMNDat1lwUgUexysyTAwgOPrm3Al3VF2qE_URplLAqOksj2-VN6_cjkW3_DVS37lGKRg9kw"
  e = "AQAB"
}
