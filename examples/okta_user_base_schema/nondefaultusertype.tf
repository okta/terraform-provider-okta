resource "okta_user_type" userbaseschematest_usertype {
  name   = "testAcc_replace_with_uuid"
  display_name = "testAcc_replace_with_uuid"
  description = "Terraform Acceptance Test Schema User Type"
}

resource "okta_user_base_schema" "firstName" {
  index       = "firstName"
  master      = "PROFILE_MASTER"
  permissions = "READ_ONLY"
  title       = "First name"
  type        = "string"
  user_type   = "${okta_user_type.userbaseschematest_usertype.id}"
}
