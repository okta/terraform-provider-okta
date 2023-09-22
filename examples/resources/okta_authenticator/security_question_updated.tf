resource "okta_authenticator" "test" {
  name = "Security Question"
  key  = "security_question"
  settings = jsonencode(
    {
      "allowedFor" : "any"
    }
  )
}