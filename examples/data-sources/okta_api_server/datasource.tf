resource "okta_api_server" "test" {
  resource_url = "https://testacc-api-server-ds-replace_with_uuid.example.com"

  metadata {
    display_name     = "testAcc-api-server-ds-replace_with_uuid"
    description      = "Terraform acceptance test API server for the okta_api_server data source"
    scopes_supported = ["read:items", "write:items"]
  }
}

data "okta_api_server" "test" {
  id = okta_api_server.test.id
}
