resource "okta_app_user_base_schema_property" "example" {
  app_id = "<app id>"
  index  = "customPropertyName"
  title  = "customPropertyName"
  type   = "string"
  master = "OKTA"
}
