// Notice users are on added to the group and group_membership is left empty.
// It is generally advisable to pick a single method of tying users to groups.
// To remove all membership specify an empty list. This is the only way to catch config drift
// and support multiple ways to outline the same config.
resource "okta_group" "test" {
  name        = "testAcc"
  description = "testing, testing"

  users = [
    okta_user.test.id,
    okta_user.test1.id,
    okta_user.test2.id,
    okta_user.test3.id,
  ]
}

resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Jones"
  login      = "john_replace_with_uuid@ledzeppelin.com"
  email      = "john_replace_with_uuid@ledzeppelin.com"
}

resource "okta_user" "test1" {
  first_name = "TestAcc"
  last_name  = "Entwhistle"
  login      = "john_replace_with_uuid@thewho.com"
  email      = "john_replace_with_uuid@thewho.com"
}

resource "okta_user" "test2" {
  first_name = "TestAcc"
  last_name  = "Doe"
  login      = "john_replace_with_uuid@unknown.com"
  email      = "john_replace_with_uuid@unknown.com"
}

resource "okta_user" "test3" {
  first_name = "TestAcc"
  last_name  = "Astley"
  login      = "rick_astley_replace_with_uuid@rickrollin.com"
  email      = "rick_astley_replace_with_uuid@rickrollin.com"
}
