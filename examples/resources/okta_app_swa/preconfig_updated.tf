resource "okta_app_swa" "test" {
  preconfigured_app = "aws_console"
  label             = "testAcc_replace_with_uuid"
  status            = "INACTIVE"
}
