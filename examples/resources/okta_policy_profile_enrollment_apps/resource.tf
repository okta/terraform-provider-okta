data "okta_policy" "example" {
  name = "My Policy"
  type = "PROFILE_ENROLLMENT"
}

data "okta_app" "test" {
  label = "My App"
}

resource "okta_policy_profile_enrollment_apps" "example" {
  policy_id = okta_policy.example.id
  apps      = [data.okta_app.id]
}
