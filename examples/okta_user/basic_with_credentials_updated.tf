resource "okta_user" "test" {
  first_name        = "TestAcc"
  last_name         = "Smith"
  login             = "testAcc-replace_with_uuid@example.com"
  email             = "testAcc-replace_with_uuid@example.com"
  password          = "SuperSecret007"
  recovery_question = "Which symbol has the ASCII code of Forty Two?"
  recovery_answer   = "Asterisk"
}
