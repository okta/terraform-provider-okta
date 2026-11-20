resource "okta_app_oauth" "test" {
  preconfigured_app = "moveworks"
  label             = "testAcc_replace_with_uuid"
  hide_ios          = false
  hide_web          = false
}
