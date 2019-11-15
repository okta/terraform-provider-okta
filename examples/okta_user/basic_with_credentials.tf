resource "okta_user" "test" {
  first_name        = "TestAcc"
  last_name         = "Smith"
  login             = "test-acc-replace_with_uuid@example.com"
  email             = "test-acc-replace_with_uuid@example.com"
  password          = "SuperSecret007"
  recovery_question = "What is the answer to life, the universe, and everything?"
  recovery_answer   = "Forty Two"
}
