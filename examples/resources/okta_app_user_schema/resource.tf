resource "okta_app_user_schema" "example" {
  app_id = okta_app_saml.example.id

  custom_property {
    index       = "salesforceId"
    title       = "Salesforce ID"
    type        = "string"
    description = "User's Salesforce ID"
    required    = false
    scope       = "NONE"
    permissions = "READ_ONLY"
  }

  custom_property {
    index       = "department"
    title       = "Department"
    type        = "string"
    description = "User's department"
    master      = "OKTA"
    scope       = "SELF"
  }

  custom_property {
    index      = "roles"
    title      = "User Roles"
    type       = "array"
    array_type = "string"
    scope      = "NONE"
  }
}
