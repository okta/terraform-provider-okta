resource "okta_policy_device_assurance_ios" "example" {
  name            = "testAcc-replace_with_uuid"
  os_version      = "12.4.5"
  jailbreak       = false
  screenlock_type = toset(["BIOMETRIC"])
}
