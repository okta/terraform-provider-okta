resource "okta_user_type" "test" {
  name         = "testAcc_replace_with_uuid"
  display_name = "testAcc_replace_with_uuid Name"
  description  = "testAcc_replace_with_uuid Description"
}

data "okta_user_type" "test-find-by-id" {
  id         = okta_user_type.test.id
  depends_on = [okta_user_type.test]
}

data "okta_user_type" "test-find-by-name" {
  name       = okta_user_type.test.name
  depends_on = [okta_user_type.test]
}
