resource "okta_ui_schema" "test" {
  ui_schema {
    button_label = "submit"
    label        = "Sign in"
    type         = "Group"

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

data "okta_ui_schema" "example" {
  id = okta_ui_schema.test.id
}