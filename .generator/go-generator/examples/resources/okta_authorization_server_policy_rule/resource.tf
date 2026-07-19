resource "okta_authorization_server_policy_rule" "example" {
  auth_server_id = "<auth-server-id>"
  policy_id = "<policy-id>"

  # Optional fields
  # access_token_lifetime_minutes = 0
  # inline_hook = "<inline_hook>"
  # refresh_token_lifetime_minutes = 0
  # refresh_token_window_minutes = 0
  # include = "<include>"
}
