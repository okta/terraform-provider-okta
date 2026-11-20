resource "okta_client_auth_setting" "test" {
  purposes      = ["AGENT_IMPORT"]
  protocol     = "OAUTH2"
  resource_orn = "orn:okta:idp:00otxk5i5T0Rk59Z3ao5:apps:office365:0oautxdhNWOdr5J9qao5"
  name                     = "testAcc-provider-cred-replace_with_uuid"
  oauth_client_configuration{
    grant_types = ["client_credentials"]
    client_credentials {
      client_id     = "1b2af17d-98de-408c-a4d3-d73f475bd222"
      client_secret = "replace_with_client_secret"
      credential_type = "CLIENT_SECRET"
    }
    scopes = ["https://graph.microsoft.com/.default"]
  }
}

resource "okta_ai_agent_provider" "test"{
  source_orn = "orn:okta:idp:00otxk5i5T0Rk59Z3ao5:apps:office365:0oautxdhNWOdr5J9qao5"
  configuration = {
    client_auth_id = okta_client_auth_setting.test.id
  }
  de_dup_expression = "target.profile.externalId eq source.profile.externalId"

  schedule {
    cron     = "0 0 4 ? * MON"
    timezone = "America/Los_Angeles"
  }
}
