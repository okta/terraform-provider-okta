resource "okta_api_service_integration" "example" {
  type = "anzennaapiservice"
  granted_scopes {
    scope = "okta.users.read"
  }
  granted_scopes {
    scope = "okta.groups.read"
  }
  granted_scopes {
    scope = "okta.logs.read"
  }
}
