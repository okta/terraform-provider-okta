
resource "okta_app_oauth" "example" {
  label          = "example"
  type           = "web"
  grant_types    = ["authorization_code"]
  redirect_uris  = ["https://example.com/"]
  response_types = ["code"]
}

### With JWKS value
### See also [Advanced PEM secrets and JWKS example](#advanced-pem-and-jwks-example).

resource "okta_app_oauth" "example" {
  label                      = "example"
  type                       = "service"
  response_types             = ["token"]
  grant_types                = ["client_credentials"]
  token_endpoint_auth_method = "private_key_jwt"

  jwks {
    kty = "RSA"
    kid = "SIGNING_KEY_RSA"
    e   = "AQAB"
    n   = "xyz"
  }

  jwks {
    kty = "EC"
    kid = "SIGNING_KEY_EC"
    x   = "K37X78mXJHHldZYMzrwipjKR-YZUS2SMye0KindHp6I"
    y   = "8IfvsvXWzbFWOZoVOMwgF5p46mUj3kbOVf9Fk0vVVHo"
  }
}

