resource okta_group test {
  name        = "testAcc_replace_with_uuid"
  description = "testing"
}

resource okta_group_roles test {
  group_id = "${okta_group.test.id}"

  admin_roles = [
    "SUPER_ADMIN",
    "ORG_ADMIN",
    "API_ACCESS_MANAGEMENT_ADMIN",
    "APP_ADMIN",
    "USER_ADMIN",
    "MOBILE_ADMIN",
    "READ_ONLY_ADMIN",
    "HELP_DESK_ADMIN",
  ]
}
