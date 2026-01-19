resource "okta_idp_oidc" "test" {
  name                  = "testAcc_replace_with_uuid"
  authorization_url     = "https://idp.example.com/authorize"
  authorization_binding = "HTTP-REDIRECT"
  token_url             = "https://idp.example.com/token"
  token_binding         = "HTTP-POST"
  user_info_url         = "https://idp.example.com/userinfo"
  user_info_binding     = "HTTP-REDIRECT"
  jwks_url              = "https://idp.example.com/keys"
  jwks_binding          = "HTTP-REDIRECT"
  scopes                = ["openid"]
  client_id             = "efg456"
  client_secret         = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  issuer_url            = "https://id.example.com"
  username_template     = "idpuser.email"
  filter                = "abc"
}
