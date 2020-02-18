resource "okta_app_oauth" "test" {
  label          = "testAcc_replace_with_uuid"
  type           = "service"
  response_types = ["token"]
  grant_types    = ["client_credentials"]
  token_endpoint_auth_method = "private_key_jwt"

  jwks {
    kty = "RSA"
    kid = "SIGNING_KEY"
    e   = "AQAB"
    n   = "xyz"
  }
}
