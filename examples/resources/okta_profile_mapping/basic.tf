resource "okta_profile_mapping" "test" {
  source_id          = okta_idp_social.google.id
  target_id          = data.okta_user_profile_mapping_source.user.id
  delete_when_absent = true

  mappings {
    id         = "firstName"
    expression = "appuser.firstName"
  }

  mappings {
    id         = "lastName"
    expression = "appuser.lastName"
  }

  mappings {
    id         = "email"
    expression = "appuser.email"
  }

  mappings {
    id         = "login"
    expression = "appuser.email"
  }
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

data "okta_user_profile_mapping_source" "user" {
  depends_on = [okta_idp_social.google]
}

