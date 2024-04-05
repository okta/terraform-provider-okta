resource "okta_policy_device_assurance_ios" "example" {
  name            = "example"
  os_version      = "12.4.5"
  jailbreak       = false
  screenlock_type = toset(["BIOMETRIC"])
}
