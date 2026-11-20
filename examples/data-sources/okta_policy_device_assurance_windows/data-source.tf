resource "okta_policy_device_assurance_windows" "example" {
  name                    = "My Windows Policy"
  platform                = "WINDOWS"
  secure_hardware_present = true
}

data "okta_policy_device_assurance_windows" "example" {
  id = okta_policy_device_assurance_windows.example.id
}
