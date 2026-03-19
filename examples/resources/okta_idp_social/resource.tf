resource "okta_idp_social" "example" {
  type          = "FACEBOOK"
  protocol_type = "OAUTH2"
  name          = "testAcc_facebook_replace_with_uuid"

  scopes = [
    "public_profile",
    "email",
  ]

  client_id         = "abcd123"
  client_secret     = "abcd123"
  username_template = "idpuser.email"
}
