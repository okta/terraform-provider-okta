data "okta_app_user_assignments" "test" {
  id = okta_app_oauth.test.id
}
