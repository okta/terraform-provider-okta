resource "okta_user" "test" {
  first_name        = "TestAcc"
  last_name         = "Smith"
  login             = "testAcc-replace_with_uuid@example.com"
  email             = "testAcc-replace_with_uuid@example.com"
  password          = "Super#Secret@007"
  old_password      = "SuperSecret007"
  recovery_question = "0011 & 1010"
  recovery_answer   = "0010"
}
