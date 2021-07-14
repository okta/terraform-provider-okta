resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc-replace_with_uuid@example.com"
  email      = "testAcc-replace_with_uuid@example.com"

  lifecycle {
    ignore_changes = [group_memberships]
  }
}

resource "okta_group" "test_1" {
  name        = "testAcc_1_replace_with_uuid"
  description = "testing, testing"
  lifecycle {
    ignore_changes = [users]
  }
}

resource "okta_group" "test_2" {
  name        = "testAcc_2_replace_with_uuid"
  description = "testing, testing"
  lifecycle {
    ignore_changes = [users]
  }
}

resource "okta_group" "test_3" {
  name        = "testAcc_3_replace_with_uuid"
  description = "testing, testing"
  lifecycle {
    ignore_changes = [users]
  }
}

resource "okta_group" "test_4" {
  name        = "testAcc_4_replace_with_uuid"
  description = "testing, testing"
  lifecycle {
    ignore_changes = [users]
  }
}

resource "okta_user_group_memberships" "test" {
  user_id = okta_user.test.id
  groups = [
    okta_group.test_1.id,
    okta_group.test_2.id,
  ]
}
