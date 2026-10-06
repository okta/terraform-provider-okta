resource "okta_policy_device_assurance_chromeos" "test" {
  name                     = "testAcc_replace_with_uuid"
  platform                 = "CHROMEOS"
  display_remediation_mode = "HIDE"

  third_party_signal_providers {
    device_posture_id_p {
      compliant = true
    }
  }
}
