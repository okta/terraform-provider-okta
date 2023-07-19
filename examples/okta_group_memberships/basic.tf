resource "okta_group" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "testing, testing"
}

resource "okta_user" "test1" {
  first_name = "TestAcc1"
  last_name  = "Smith"
  login      = "testAcc1-replace_with_uuid@example.com"
  email      = "testAcc1-replace_with_uuid@example.com"
}

resource "okta_user" "test2" {
  first_name = "TestAcc2"
  last_name  = "Brando"
  login      = "testAcc2-replace_with_uuid@example.com"
  email      = "testAcc2-replace_with_uuid@example.com"
}

resource "okta_user" "test3" {
  first_name = "TestAcc3"
  last_name  = "Python"
  login      = "testAcc3-replace_with_uuid@example.com"
  email      = "testAcc3-replace_with_uuid@example.com"
}

resource "okta_user" "test4" {
  first_name = "TestAcc4"
  last_name  = "Jenkins"
  login      = "testAcc4-replace_with_uuid@example.com"
  email      = "testAcc4-replace_with_uuid@example.com"
}


resource "okta_group_memberships" "test" {
  group_id = okta_group.test.id
  users    = [
    okta_user.test1.id,
    okta_user.test2.id,
  ]
}
