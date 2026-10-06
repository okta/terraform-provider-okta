resource "okta_user_schema_property" "example" {
  index       = "customPropertyName"
  title       = "customPropertyName"
  type        = "string"
  description = "My custom property name"
  master      = "OKTA"
  scope       = "SELF"
  user_type   = data.okta_user_type.example.id
}

resource "okta_user_schema_property" "auth_type" {
  index       = "authType"
  title       = "Authentication Type"
  type        = "string"
  description = "Determines which authentication method the user should use"
  master      = "OKTA"
  scope       = "SELF"
  enum        = ["EMAIL", "PASSWORD"]
  default     = "EMAIL"

  one_of {
    const = "EMAIL"
    title = "Email"
  }

  one_of {
    const = "PASSWORD"
    title = "Password"
  }
}
