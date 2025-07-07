resource "okta_group" "test_group" {
  name        = "testAcc_replace_with_uuid_group"
  description = "Test group for data source testing"
}

resource "okta_app_oauth" "test_app" {
  label          = "testAcc_replace_with_uuid_app"
  type           = "web"
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://d.com/"]
  response_types = ["code"]
}

resource "okta_resource_set" "test" {
  label       = "testAcc_replace_with_uuid"
  description = "A test resource set with multiple valid resources"
  resources = [
    "https://${var.hostname}/api/v1/groups/${okta_group.test_group.id}",
    "https://${var.hostname}/api/v1/apps/${okta_app_oauth.test_app.id}",
    "https://${var.hostname}/api/v1/users"
  ]
}

data "okta_resource_set" "test" {
  id = okta_resource_set.test.id
}
