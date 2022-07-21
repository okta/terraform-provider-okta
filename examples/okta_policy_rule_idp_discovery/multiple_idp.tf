data "okta_policy" "test" {
  type = "IDP_DISCOVERY"
  name = "Idp Discovery Policy"

}

resource "okta_policy_rule_idp_discovery" "test" {
  name      = "Test Rule"
  policy_id = data.okta_policy.test.id

  idp {
    id   = okta_idp_social.google.id
    type = "GOOGLE"
  }
  idp {
    id   = okta_idp_social.facebook.id
    type = "FACEBOOK"
  }
}

resource "okta_idp_social" "facebook" {
  name          = "Facebook"
  type          = "FACEBOOK"
  protocol_type = "OAUTH2"
  client_id     = "var.idp_social_facebook_client_id"
  client_secret = "var.idp_social_facebook_client_secret"
  scopes        = ["public_profile", "email"]

}
resource "okta_idp_social" "google" {
  name          = "Google"
  type          = "GOOGLE"
  protocol_type = "OIDC"
  client_id     = "var.idp_social_google_client_id"
  client_secret = "var.idp_social_google_client_secret"
  scopes        = ["openid", "profile", "email"]
}
