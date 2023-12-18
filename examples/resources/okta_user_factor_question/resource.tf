data "okta_user_security_questions" "example" {
  user_id = okta_user.example.id
}

resource "okta_user" "example" {
  first_name = "John"
  last_name  = "Smith"
  login      = "john.smith@example.com"
  email      = "john.smith@example.com"
}

resource "okta_factor" "example" {
  provider_id = "okta_question"
  active      = true
}

resource "okta_user_factor_question" "example" {
  user_id = okta_user.example.id
  key     = data.okta_user_security_questions.example.questions[0].key
  answer  = "meatball"
  depends_on = [
  okta_factor.example]
}
