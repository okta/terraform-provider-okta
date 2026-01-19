resource "okta_group" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "Test group for push mapping"
}

data "okta_app" "test" {
  label = "My Provisioning App"
}

resource "okta_app_group_push_mapping" "test" {
  app_id          = data.okta_app.test.id
  source_group_id = okta_group.test.id
  target_group_name = okta_group.test.name
}
