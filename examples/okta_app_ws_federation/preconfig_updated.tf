resource "okta_app_ws_federation" "test" {
  preconfigured_app = "aws_console"
  label             = "testAcc_replace_with_uuid"
  visibility        = false
}
