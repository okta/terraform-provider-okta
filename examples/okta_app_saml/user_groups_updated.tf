resource "okta_group" "group" {
  name = "testAcc_replace_with_uuid"
}

resource "okta_group" "group1" {
  name = "testAcc_replace_with_uuid_1"
}

resource "okta_group" "group2" {
  name = "testAcc_replace_with_uuid_2"
}

data "okta_group" "all" {
  name = "Everyone"
}

resource "okta_user" "user" {
  admin_roles = ["APP_ADMIN", "USER_ADMIN"]
  first_name  = "TestAcc"
  last_name   = "blah"
  login       = "testAcc-replace_with_uuid@example.com"
  email       = "testAcc-replace_with_uuid@example.com"
  status      = "ACTIVE"
}

resource "okta_user" "user1" {
  first_name = "TestAcc1"
  last_name  = "blah"
  login      = "testAcc-1-replace_with_uuid@example.com"
  email      = "testAcc-1-replace_with_uuid@example.com"
  status     = "ACTIVE"
}

resource "okta_app_saml" "test" {
  preconfigured_app = "amazon_aws"
  label             = "testAcc_replace_with_uuid"
  groups            = [data.okta_group.all.id]

  users {
    id       = okta_user.user1.id
    username = okta_user.user1.email
  }
}
