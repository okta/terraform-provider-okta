resource "okta_api_server" "test" {
  resource_url = "https://testacc-api-server-replace_with_uuid.example.com"

  metadata {
    display_name     = "testAcc-api-server-replace_with_uuid"
    description      = "Terraform acceptance test API server"
    scopes_supported = ["read:items"]
  }
}
