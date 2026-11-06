resource "okta_policy_profile_enrollment" "test" {
  name   = "testAcc_replace_with_uuid"
  status = "ACTIVE"
}

resource "okta_inline_hook" "test" {
  name    = "testAcc_replace_with_uuid"
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

resource "okta_group" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "testing, testing"
}

resource "okta_policy_rule_profile_enrollment" "test" {
  policy_id = okta_policy_profile_enrollment.test.id
  type      = "PROFILE_ENROLLMENT"
  name      = "Catch-all Rule"
  status    = "ACTIVE"
  priority  = 99
  system    = true

  actions {
    profile_enrollment {
      access                       = "ALLOW"
      allowed_identifiers          = ["login"]
      target_group_ids             = [okta_group.test.id]
      unknown_user_action          = "REGISTER"
      progressive_profiling_action = "ENABLED"

      pre_registration_inline_hooks {
        inline_hook_id = okta_inline_hook.test.id
        scopes         = ["PROGRESSIVE_PROFILING"]
      }

      activation_requirements {
        email_verification = false
      }

      profile_attributes {
        label    = "Email"
        name     = "email"
        required = true
      }
      profile_attributes {
        label    = "First Name"
        name     = "firstName"
        required = true
      }
      profile_attributes {
        label    = "Last Name"
        name     = "lastName"
        required = true
      }
      profile_attributes {
        label    = "Mobile Phone"
        name     = "mobilePhone"
        required = false
      }
    }
  }
}
