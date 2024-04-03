resource "okta_policy_device_assurance_android" "example" {
  name                    = "example"
  os_version              = "12"
  disk_encryption_type    = toset(["FULL", "USER"])
  jailbreak               = false
  secure_hardware_present = true
  screenlock_type         = toset(["BIOMETRIC"])
}
