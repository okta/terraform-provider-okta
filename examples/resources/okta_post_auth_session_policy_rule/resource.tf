data "okta_post_auth_session_policy" "example" {
}

resource "okta_post_auth_session_policy_rule" "example" {
  policy_id         = data.okta_post_auth_session_policy.example.id
  name              = "Session Protection Rule"
  terminate_session = true
}
