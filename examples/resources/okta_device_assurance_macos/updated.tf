resource "okta_device_assurance_macos" "test" {
  name                     = "testAcc_replace_with_uuid"
  platform                 = "MACOS"
  secure_hardware_present  = true
  display_remediation_mode = "HIDE"
}
