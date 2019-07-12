resource okta_group test {
  name        = "testAcc_replace_with_uuid"
  description = "testing"
}

resource okta_group_roles test {
  group_id    = "${okta_group.test.id}"
  admin_roles = ["SUPER_ADMIN"]
}
