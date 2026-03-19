resource "okta_app_oauth_api_scope" "example" {
  app_id = "<application_id>"
  issuer = "<your org domain>"
  scopes = ["okta.users.read", "okta.users.manage"]
}
