data "okta_policy" "access" {
  name = "Any two factors"
  type = "ACCESS_POLICY"
}

data "okta_app" "example" {
  label = "Example App"
}

resource "okta_app_access_policy_assignment" "assignment" {
  app_id    = data.okta_app.example.id
  policy_id = data.okta_policy.access.id
}
