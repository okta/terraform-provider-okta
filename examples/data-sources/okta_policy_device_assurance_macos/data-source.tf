resource "okta_policy_device_assurance_macos" "example" {
  name                    = "My macOS Policy"
  platform                = "MACOS"
  secure_hardware_present = true
}

data "okta_policy_device_assurance_macos" "example" {
  id = okta_policy_device_assurance_macos.example.id
}
