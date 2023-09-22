resource "okta_user" "example" {
  first_name = "John"
  last_name  = "Smith"
  login      = "john.smith@example.com"
  email      = "john.smith@example.com"
}

data "okta_user_security_questions" "example" {
  user_id = okta_user.example.id
}

// TODU
