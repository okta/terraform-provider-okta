resource "okta_oauth_app" "test" {
  label = "testAcc_%[1]d"
  type  = "service"
}
