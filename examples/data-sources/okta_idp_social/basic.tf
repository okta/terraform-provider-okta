resource "okta_idp_social" "facebook" {
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

resource "okta_user_schema_property" "schema_custom" {
  index       = "customfieldId"
  title       = "Test Custom Field"
  type        = "string"
  description = "Test custom field"
  required    = false
  min_length  = 1
  max_length  = 50
  permissions = "READ_ONLY"
  master      = "PROFILE_MASTER"
}

resource "okta_idp_social" "google" {
  type          = "GOOGLE"
  protocol_type = "OIDC"
  name          = "testAcc_google_replace_with_uuid"

  scopes = [
    "profile",
    "email",
    "openid",
  ]

  client_id         = "abcd123"
  client_secret     = "abcd123"
  username_template = "idpuser.email"
}

resource "okta_idp_social" "microsoft" {
  type          = "MICROSOFT"
  protocol_type = "OIDC"
  name          = "testAcc_microsoft_replace_with_uuid"

  scopes = [
    "openid",
    "email",
    "profile",
    "https://graph.microsoft.com/User.Read",
  ]

  client_id         = "abcd123"
  client_secret     = "abcd123"
  username_template = "idpuser.userPrincipalName"
  groups_action     = "ASSIGN"
  groups_assignment = [okta_group.test.id]
}

resource "okta_group" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "testing, testing, testing"
}
