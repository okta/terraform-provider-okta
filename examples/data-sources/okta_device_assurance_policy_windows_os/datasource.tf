resource "okta_policy_device_assurance_windows" "test" {
  name     = "testAcc_replace_with_uuid"
  platform                = "WINDOWS"
  secure_hardware_present = true
}

data "okta_policy_device_assurance_windows" "test" {
  id = okta_policy_device_assurance_windows.test.id
}
