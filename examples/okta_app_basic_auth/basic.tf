resource "okta_group" "group" {
  name = "testAcc_replace_with_uuid"
}

resource "okta_app_basic_auth" "test" {
  label    = "testAcc_replace_with_uuid"
  url      = "https://example.com/login.html"
  auth_url = "https://example.com/auth.html"
  groups   = ["${okta_group.group.id}"]
}
