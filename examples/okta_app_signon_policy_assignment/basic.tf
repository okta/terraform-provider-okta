data "okta_policy" "dashboard" {
  name = "Okta Dashboard"
  type = "ACCESS_POLICY"
}

resource "okta_app_signon_policy_assignment" "my_app_assignment" {
  app_id    = "0oa8wqnqnhKERkM721d7"
  policy_id = data.okta_policy.dashboard.id
}
