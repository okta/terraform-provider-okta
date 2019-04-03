resource "okta_group" "group" {
  name = "testAcc_%[1]d"
}

resource "okta_group" "group1" {
  name = "testAcc_%[1]d_1"
}

resource "okta_group" "group2" {
  name = "testAcc_%[1]d_2"
}

data "okta_group" "all" {
  name = "Everyone"
}

resource "okta_user" "user" {
  admin_roles = ["APP_ADMIN", "USER_ADMIN"]
  first_name  = "TestAcc"
  last_name   = "blah"
  login       = "test-acc-%[1]d@testing.com"
  email       = "test-acc-%[1]d@testing.com"
  status      = "ACTIVE"
}

resource "okta_user" "user1" {
  first_name = "TestAcc1"
  last_name  = "blah"
  login      = "test-acc-1-%[1]d@testing.com"
  email      = "test-acc-1-%[1]d@testing.com"
  status     = "ACTIVE"
}

resource "okta_saml_app" "test" {
  preconfigured_app = "amazon_aws"
  label             = "testAcc_%[1]d"
  groups            = ["${data.okta_group.all.id}"]

  users = [
    {
      id       = "${okta_user.user1.id}"
      username = "${okta_user.user1.email}"
    },
  ]
}
