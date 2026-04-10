resource "okta_policy_device_assurance_ios" "example" {
  name            = "testAcc-replace_with_uuid"
  os_version      = "12.4.5"
  jailbreak       = false
  screenlock_type = toset(["BIOMETRIC"])

  grace_period {
    type   = "BY_DATE_TIME"
    expiry = "2026-12-01T00:00:00.000Z"
  }

  display_remediation_mode = "HIDE"
}
