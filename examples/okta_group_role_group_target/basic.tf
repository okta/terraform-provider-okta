resource okta_group test {
  name        = "testAcc_replace_with_uuid"
  description = "testing"
}

resource okta_group_role test {
  group_id    = "${okta_group.test.id}"
  role_type   = "USER_ADMIN"
}

resource okta_group_role_group_target test {
  role_id     = "${okta_group_role.test.id}"
  group_id    = "${okta_group.test.id}"
}

resource okta_user test {
  first_name        = "TestAcc"
  last_name         = "Smith"
  login             = "test-acc-replace_with_uuid@example.com"
  email             = "test-acc-replace_with_uuid@example.com"
  group_memberships = ["${okta_group.test.id}"]
}
