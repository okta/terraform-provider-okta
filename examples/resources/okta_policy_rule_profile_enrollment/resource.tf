resource "okta_policy_profile_enrollment" "example" {
  name = "My Enrollment Policy"
}

resource "okta_inline_hook" "example" {
  name    = "My Inline Hook"
  status  = "ACTIVE"
  type    = "com.okta.user.pre-registration"
  version = "1.0.3"

  channel = {
    type    = "HTTP"
    version = "1.0.0"
    uri     = "https://example.com/test2"
    method  = "POST"
  }
}

resource "okta_group" "example" {
  name        = "My Group"
  description = "Group of some users"
}

resource "okta_policy_rule_profile_enrollment" "example" {
  policy_id           = okta_policy_profile_enrollment.example.id
  inline_hook_id      = okta_inline_hook.example.id
  target_group_id     = okta_group.example.id
  unknown_user_action = "REGISTER"
  email_verification  = true
  access              = "ALLOW"
  profile_attributes {
    name     = "email"
    label    = "Email"
    required = true
  }
  profile_attributes {
    name     = "name"
    label    = "Name"
    required = true
  }
  profile_attributes {
    name     = "t-shirt"
    label    = "T-Shirt Size"
    required = false
  }
}
