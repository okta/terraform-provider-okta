data "okta_app_group_assignments" "test" {
  id = okta_app_oauth.test.id
}
