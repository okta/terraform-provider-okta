resource "okta_idp_social" "google" {
  type                = "GOOGLE"
  protocol_type       = "OIDC"
  name                = "testAcc_google_replace_with_uuid"
  provisioning_action = "DISABLED"

  scopes = [
    "profile",
    "email",
    "openid",
  ]

  client_id         = "abcd123"
  client_secret     = "abcd123"
  username_template = "idpuser.email"
}
