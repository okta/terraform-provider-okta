resource "okta_api_server" "test" {
  resource_url = "https://testacc-api-server-as-replace_with_uuid.example.com"

  metadata {
    display_name     = "testAcc-api-server-as-replace_with_uuid"
    description      = "Terraform acceptance test API server for okta_api_servers_authorization_server"
    scopes_supported = ["read:items"]
  }
}

resource "okta_api_servers_authorization_server" "test" {
  api_server_id = okta_api_server.test.id
  orn           = "orn:okta:directory:00otxk5i5T0Rk59Z3ao5:external-authorization-servers:easu7ekmjWxppv5PFao5"
}
