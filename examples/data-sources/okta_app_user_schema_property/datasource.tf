# Example: Query an existing app user schema property
# This is useful for referencing properties that were auto-created
# when provisioning was enabled, or created outside of Terraform

data "okta_app_saml" "example" {
  label = "My SAML App"
}

# Query a property that may have been auto-created
data "okta_app_user_schema_property" "username" {
  app_id = data.okta_app_saml.example.id
  index  = "userName"
}

# Use the property details in other resources
output "username_property_type" {
  value = data.okta_app_user_schema_property.username.type
}

output "username_property_title" {
  value = data.okta_app_user_schema_property.username.title
}
