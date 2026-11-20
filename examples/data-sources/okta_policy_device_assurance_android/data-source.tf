resource "okta_policy_device_assurance_android" "example" {
  name                    = "My Android Policy"
  platform                = "ANDROID"
  secure_hardware_present = true
}

data "okta_policy_device_assurance_android" "example" {
  id = okta_policy_device_assurance_android.example.id
}
