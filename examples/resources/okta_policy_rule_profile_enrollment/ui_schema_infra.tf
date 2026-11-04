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
