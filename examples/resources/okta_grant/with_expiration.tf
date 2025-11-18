# Grant with expiration
resource "okta_group" "test" {
  name = "Test Group"
}

resource "okta_grant" "test" {
  grant_type              = "CUSTOM"
  target_principal_id     = okta_group.test.id
  target_principal_type   = "OKTA_GROUP"
  target_resource_orn     = "orn:okta:idp:00o123:apps:jira:0oa789"
  expiration_date         = "2025-06-30T23:59:59Z"
  time_zone               = "America/New_York"
  
  entitlements {
    id = "ent123"
    values {
      id = "val456"
    }
  }
}
