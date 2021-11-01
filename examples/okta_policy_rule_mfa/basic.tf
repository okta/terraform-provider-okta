data "okta_default_policy" "test" {
  type = "MFA_ENROLL"
}

resource "okta_policy_rule_mfa" "test" {
  policy_id = data.okta_default_policy.test.id
  name      = "testAcc_replace_with_uuid"
  status    = "ACTIVE"
  app_include {
    id   = okta_app_oauth.test.id
    type = "APP"
  }
}

resource "okta_app_oauth" "test" {
  label          = "testAcc_replace_with_uuid"
  type           = "web"
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://localhost:8000"]
  response_types = ["code"]
  skip_groups    = true
}
