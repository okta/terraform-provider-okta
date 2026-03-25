resource "okta_app_oauth" "test" {
  label                      = "testAcc_replace_with_uuid"
  type                       = "service"
  response_types             = ["token"]
  grant_types                = ["client_credentials"]
  token_endpoint_auth_method = "private_key_jwt"

  jwks_uri = "https://example.com/.well-known/jwks.json"
}

resource "okta_app_oauth" "test_ec" {
  label                      = "test_ecAcc_replace_with_uuid"
  type                       = "service"
  response_types             = ["token"]
  grant_types                = ["client_credentials"]
  token_endpoint_auth_method = "private_key_jwt"

  jwks_uri = "https://example.com/.well-known/jwks-ec.json"
}

# Test EC Key
# {
#     "kty": "EC",
#     "use": "sig",
#     "crv": "P-256",
#     "kid": "testing",
#     "x": "K37X78mXJHHldZYMzrwipjKR-YZUS2SMye0KindHp6I",
#     "y": "8IfvsvXWzbFWOZoVOMwgF5p46mUj3kbOVf9Fk0vVVHo",
#     "alg": "ES256"
# }
