resource "okta_app_three_field" "test" {
  label                = "testAcc_replace_with_uuid"
  button_selector      = "btn"
  username_selector    = "user"
  password_selector    = "pass"
  url                  = "http://example.com"
  extra_field_selector = "third"
  extra_field_value    = "third"
}
