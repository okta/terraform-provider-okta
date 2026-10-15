resource "okta_policy_device_assurance_android" "test" {
  name                    = "testAcc_replace_with_uuid"
  platform                = "ANDROID"
  secure_hardware_present = true
}

data "okta_policy_device_assurance_android" "test" {
  id = okta_policy_device_assurance_android.test.id
}
