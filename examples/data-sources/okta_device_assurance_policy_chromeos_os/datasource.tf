resource "okta_policy_device_assurance_chromeos" "test" {
  name     = "testAcc_replace_with_uuid"
  platform = "CHROMEOS"
  display_remediation_mode = "SHOW"
  third_party_signal_providers{
    device_posture_id_p{
      compliant ="true"
      managed = "true"
    }
  }
}

data "okta_policy_device_assurance_chromeos" "test" {
  id = okta_policy_device_assurance_chromeos.test.id
}
