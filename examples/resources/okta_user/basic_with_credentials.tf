resource "okta_user" "test" {
  first_name        = "TestAcc"
  last_name         = "Smith"
  login             = "testAcc-replace_with_uuid@example.com"
  email             = "testAcc-replace_with_uuid@example.com"
  password          = "Abcd1234"
  recovery_question = "What is the answer to life, the universe, and everything?"
  recovery_answer   = "Forty Two"
}
