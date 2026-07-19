data "okta_group" "all" {
  name = "Everyone"
}

resource "okta_idp_social" "google" {
  type                = "GOOGLE"
  protocol_type       = "OIDC"
  name                = "testAcc_replace_with_uuid_google"
  provisioning_action = "DISABLED"
  scopes = [
    "profile",
    "email",
    "openid",
  ]
  client_id         = "foo"
  client_secret     = "bar"
  username_template = "idpuser.email"
}

resource "okta_idp_social" "facebook" {
  type                = "FACEBOOK"
  protocol_type       = "OIDC"
  name                = "testAcc_replace_with_uuid_facebook"
  provisioning_action = "DISABLED"
  scopes = [
    "public_profile",
    "email",
  ]
  client_id         = "foo"
  client_secret     = "bar"
  username_template = "idpuser.email"
}

resource "okta_policy_signon" "test" {
  name            = "testAcc_replace_with_uuid"
  status          = "ACTIVE"
  description     = "Terraform Acceptance Test SignOn Policy GH-2793"
  groups_included = [data.okta_group.all.id]
}

resource "okta_policy_rule_signon" "test" {
  policy_id             = okta_policy_signon.test.id
  name                  = "testAcc_replace_with_uuid"
  status                = "ACTIVE"
  mfa_required          = false
  identity_provider     = "SPECIFIC_IDP"
  identity_provider_ids = [okta_idp_social.google.id, okta_idp_social.facebook.id]
}