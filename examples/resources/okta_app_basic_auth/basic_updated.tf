resource "okta_app_basic_auth" "test" {
  label    = "testAcc_replace_with_uuid"
  url      = "https://example.com/login.html"
  auth_url = "https://example.org/auth.html"
  logo     = "../examples/resources/okta_app_basic_auth/terraform_icon.png"
}
