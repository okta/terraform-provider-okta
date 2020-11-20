resource okta_group test {
  name        = "testAcc_replace_with_uuid"
  description = "testing"
}

resource okta_group_roles test {
  group_id = okta_group.test.id

  admin_roles = [
    "SUPER_ADMIN",
    "ORG_ADMIN",
    "API_ACCESS_MANAGEMENT_ADMIN",
    "APP_ADMIN",
    "USER_ADMIN",
    "MOBILE_ADMIN",
    "READ_ONLY_ADMIN",
    "HELP_DESK_ADMIN",
    "REPORT_ADMIN",
    "GROUP_MEMBERSHIP_ADMIN"
  ]
}

resource okta_user test {
  first_name        = "TestAcc"
  last_name         = "Smith"
  login             = "test-acc-replace_with_uuid@example.com"
  email             = "test-acc-replace_with_uuid@example.com"
  group_memberships = [okta_group.test.id]
}
