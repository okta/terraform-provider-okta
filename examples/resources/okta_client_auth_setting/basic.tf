resource "okta_client_auth_setting" "test" {
  resource_orn = "orn:okta:directory:00otxk5i5T0Rk59Z3ao5:resource-servers:mcp:emsu6mjxY5YyMHF0zao5"
  purposes     = ["STS"]
  protocol     = "OAUTH2"
    authorization_server_orn = "orn:okta:directory:00otxk5i5T0Rk59Z3ao5:external-authorization-servers:easu6mjzHxpNlehMAao5"
    name                     = "testAcc-cred-replace_with_uuid"
  oauth_client_configuration {
    grant_types = ["authorization_code"]
    scopes      = ["okta.users.read"]
    registration_method="STATIC"
    client_credentials {
      client_id     = "abcd"
      client_secret = "BlxZuD-BY93TVg8uZc1Jx8dbqKJ5H74Bwot2Gon3fovQg0Yh3NxpIDr1AJG2yvWI"
      credential_type = "CLIENT_SECRET"
    }
  }
}

