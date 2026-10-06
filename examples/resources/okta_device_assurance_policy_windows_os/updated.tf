resource "okta_policy_device_assurance_windows" "test" {
  name                     = "testAcc_replace_with_uuid"
  platform                 = "WINDOWS"
  secure_hardware_present  = true
  display_remediation_mode = "HIDE"
}
