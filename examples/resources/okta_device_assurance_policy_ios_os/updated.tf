resource "okta_policy_device_assurance_ios" "test" {
  name                     = "testAcc_replace_with_uuid"
  platform                 = "IOS"
  jailbreak                = false
  display_remediation_mode = "HIDE"
}
