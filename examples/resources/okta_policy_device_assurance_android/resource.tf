resource "okta_policy_device_assurance_android" "example" {
  name                    = "My Android Policy"
  platform                = "ANDROID"
  secure_hardware_present = true
}

resource "okta_policy_device_assurance_android" "example_with_remediation" {
  name                     = "My Android Policy with Remediation"
  platform                 = "ANDROID"
  display_remediation_mode = "HIDE"

  third_party_signal_providers {
    device_posture_id_p {
      compliant = true
      managed   = true
    }
  }

  grace_period {
    type   = "DAYS"
    expiry = "30"
  }
}
