resource "okta_app_user_schema_property" "example" {
  app_id      = "<app id>"
  index       = "customPropertyName"
  title       = "customPropertyName"
  type        = "string"
  description = "My custom property name"
  master      = "OKTA"
  scope       = "SELF"
}
