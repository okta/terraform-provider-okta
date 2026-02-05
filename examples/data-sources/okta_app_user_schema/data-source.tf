data "okta_app_user_schema" "example" {
  app_id = okta_app_saml.example.id
}
