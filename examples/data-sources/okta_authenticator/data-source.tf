data "okta_authenticator" "test" {
  key = "security_question"
}

data "okta_authenticator" "test_1" {
  name = "Okta Verify"
}
