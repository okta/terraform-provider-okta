resource "okta_app_three_field" "test" {
  label                = "testAcc_replace_with_uuid"
  status               = "INACTIVE"
  button_selector      = "btn1"
  username_selector    = "user1"
  password_selector    = "pass1"
  url                  = "http://example.com"
  extra_field_selector = "mfa"
  extra_field_value    = "mfa"
  credentials_scheme   = "SHARED_USERNAME_AND_PASSWORD"
  shared_username      = "testAcc_replace_with_uuid@example.com"
  shared_password      = "PA11LVFDaNdulkNsKLeb"
}
