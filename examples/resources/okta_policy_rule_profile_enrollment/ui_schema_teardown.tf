# The catch-all rule is a system rule: it cannot be deleted on its own and is
# removed by Okta when its policy is deleted. Forget it from state rather than
# destroying it, then let the policy be destroyed.
removed {
  from = okta_policy_rule_profile_enrollment.test

  lifecycle {
    destroy = false
  }
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
