terraform {
  required_providers {
    okta = {
      source = "okta/okta"
    }
  }
}

# Credentials come from the environment:
# OKTA_ORG_NAME, OKTA_BASE_URL, OKTA_API_TOKEN (or OKTA_API_CLIENT_ID/OKTA_API_PRIVATE_KEY)
resource "okta_group" "example_target_group" {
  name                      = "OKTA-1129669 Target Group"
  description               = "Group new users are added to on enrollment"
  custom_profile_attributes = "{}"
  skip_users                = false
}

resource "okta_inline_hook" "example_registration_hook" {
  name    = "OKTA-1129669 Registration Hook"
  status  = "ACTIVE"
  type    = "com.okta.user.pre-registration"
  version = "1.0.3"

  channel = {
    type    = "HTTP"
    version = "1.0.0"
    uri     = "https://example.com/okta-1129669"
    method  = "POST"
  }
}

resource "okta_policy_profile_enrollment" "example_profile_enrollment_policy" {
  name   = "OKTA-1129669 Enrollment Policy"
  status = "ACTIVE"
}

resource "okta_policy_rule_profile_enrollment" "example_profile_enrollment_rule" {
  conditions                  = null
  system                      = true
  priority                    = 99
  name = "Catch-all Rule"
  status = "ACTIVE"
  type = "PROFILE_ENROLLMENT"
  policy_id                    = okta_policy_profile_enrollment.example_profile_enrollment_policy.id
  actions {
    profile_enrollment {
      access                       = "ALLOW"
      target_group_ids             = [okta_group.example_target_group.id]
      unknown_user_action          = "REGISTER"
      pre_registration_inline_hooks {
        inline_hook_id               = okta_inline_hook.example_registration_hook.id
        scopes = ["PROGRESSIVE_PROFILING"] # "SELF_SERVICE_REGISTRATION"
      }
      activation_requirements {
        email_verification           = false
      }
      
      progressive_profiling_action = "ENABLED"
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
          required = false
        }
        # Backend-generated value; not managed by any resource in this state.
        ui_schema_id = "uis14w6973lbErFYJ1d8"
    }
  }
}
