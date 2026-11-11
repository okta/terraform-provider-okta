resource "okta_policy_profile_enrollment" "test" {
  name   = "testAcc_replace_with_uuid"
  status = "ACTIVE"
}

resource "okta_group" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "terraform group"
}

resource "okta_ui_schema" "test" {
  ui_schema {
    button_label = "Register"
    label        = "Sign in"
    type         = "Group"

    elements {
      type  = "Control"
      scope = "#/properties/email"
      label = "Primary Email"
      options {
        format = "text"
      }
    }
    elements {
      type  = "Control"
      scope = "#/properties/firstName"
      label = "First Name"
      options {
        format = "text"
      }
    }
    elements {
      type  = "Control"
      scope = "#/properties/lastName"
      label = "Last Name"
      options {
        format = "text"
      }
    }
  }
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
      unknown_user_action          = "REGISTER"
      progressive_profiling_action = "DISABLED"
      target_group_ids             = [okta_group.test.id]
      ui_schema_id                 = okta_ui_schema.test.id

      activation_requirements {
        email_verification = false
      }

      profile_attributes {
        name     = "email"
        label    = "Primary Email"
        required = true
      }
      profile_attributes {
        name     = "firstName"
        label    = "First Name"
        required = true
      }
      profile_attributes {
        name     = "lastName"
        label    = "Last Name"
        required = true
      }
    }
  }
}
