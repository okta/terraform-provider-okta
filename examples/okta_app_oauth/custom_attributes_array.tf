resource "okta_app_oauth" "test" {
  label                      = "testAcc_replace_with_uuid"
  type                       = "web"
  grant_types                = ["authorization_code"]
  redirect_uris              = ["http://d.com/"]
  response_types             = ["code"]
  client_basic_secret        = "something_from_somewhere"
  custom_client_id           = "something_from_somewhere"
  token_endpoint_auth_method = "client_secret_basic"
  
  custom_profile_attributes = <<JSON
  {
    "array123": ["test"],
    "number123": 1
  }
JSON
}
