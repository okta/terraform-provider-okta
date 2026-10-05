resource "okta_device_assurance_macos" "test" {
  name     = "testAcc_replace_with_uuid"
  platform                = "MACOS"
  secure_hardware_present = true
}

data "okta_device_assurance_macos" "test" {
  id = okta_device_assurance_macos.test.id
}
