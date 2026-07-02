resource "okta_group" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "testing import"
}

resource "okta_user" "test1" {
  first_name = "ImportTestAcc1"
  last_name  = "User"
  login      = "import-testAcc1-replace_with_uuid@example.com"
  email      = "import-testAcc1-replace_with_uuid@example.com"
}

resource "okta_user" "test2" {
  first_name = "ImportTestAcc2"
  last_name  = "User"
  login      = "import-testAcc2-replace_with_uuid@example.com"
  email      = "import-testAcc2-replace_with_uuid@example.com"
}

resource "okta_group_memberships" "test" {
  group_id        = okta_group.test.id
  track_all_users = true
  users = [
    okta_user.test1.id,
    okta_user.test2.id,
  ]
}
