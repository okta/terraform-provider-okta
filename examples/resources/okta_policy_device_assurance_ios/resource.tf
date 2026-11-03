resource "okta_policy_device_assurance_ios" "example" {
  name      = "My iOS Policy"
  platform  = "IOS"
  jailbreak = false
}

resource "okta_policy_device_assurance_ios" "example_with_providers" {
  name                     = "My iOS Policy with Providers"
  platform                 = "IOS"
  jailbreak                = false
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
