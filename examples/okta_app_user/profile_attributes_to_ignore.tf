resource "okta_app_saml" "test" {
  label             = "testAcc_replace_with_uuid"
  preconfigured_app = "google"

  lifecycle {
    ignore_changes = [users, groups]
  }
}

resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc_replace_with_uuid@example.com"
  email      = "testAcc_replace_with_uuid@example.com"
}

resource "okta_app_user" "test" {
  app_id   = okta_app_saml.test.id
  user_id  = okta_user.test.id
  username = okta_user.test.email

  profile = <<JSON
{"someAttribute":"testing"}
JSON

  profile_attributes_to_ignore = ["testCustom"]
}
