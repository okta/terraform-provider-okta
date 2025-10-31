resource "okta_idp_social" "test" {
  type          = "FACEBOOK"
  protocol_type = "OAUTH2"
  name          = "testAcc_replace_with_uuid"

  scopes = [
    "public_profile",
    "email",
  ]

  client_id         = "abcd123"
  client_secret_wo  = "secret_from_writeonly_attr"
  username_template = "idpuser.email"
}
