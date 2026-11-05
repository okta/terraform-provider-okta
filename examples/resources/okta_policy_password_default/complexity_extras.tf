resource "okta_policy_password_default" "test" {
  password_max_consecutive_repeating_characters = 3
  password_oel_statement                        = "password.value.contains('okta')"
}
