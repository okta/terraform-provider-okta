resource "okta_device_assurance_android" "test" {
  name                     = "testAcc_replace_with_uuid"
  platform                 = "ANDROID"
  secure_hardware_present  = true
  display_remediation_mode = "HIDE"
}
