# ENTITLEMENT-BUNDLE grant example
resource "okta_user" "test" {
  first_name = "Test"
  last_name  = "User"
  login      = "test.user@example.com"
  email      = "test.user@example.com"
}

resource "okta_grant" "test" {
  grant_type              = "ENTITLEMENT-BUNDLE"
  target_principal_id     = okta_user.test.id
  target_principal_type   = "OKTA_USER"
  target_resource_orn     = "orn:okta:idp:00o123:apps:salesforce:0oa456"
  entitlement_bundle_id   = "enb789"
  actor                   = "ADMIN"
}
