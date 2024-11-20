resource "okta_user_type" "test" {
  name         = "testAcc_replace_with_uuid"
  display_name = "Terraform Acceptance Test User Type"
  description  = "Terraform Acceptance Test User Type"
}

data "okta_user_type" "test" {
  name = okta_user_type.test.name
}
