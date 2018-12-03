resource "okta_group" "group-%[1]d" {
  name = "testAcc-%[1]d"
}

resource "okta_user" "user-%[1]d" {
  admin_roles = ["APP_ADMIN", "USER_ADMIN"]
  first_name  = "TestAcc"
  last_name   = "blah"
  login       = "test-acc-%[1]d@testing.com"
  email       = "test-acc-%[1]d@testing.com"
  status      = "ACTIVE"
}

resource "okta_saml_app" "testAcc-%[1]d" {
  preconfigured_app = "amazon_aws"
  label             = "testAcc-%[1]d"
}
