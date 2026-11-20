resource "okta_policy_device_assurance_chromeos" "example" {
  name     = "My ChromeOS Policy"
  platform = "CHROMEOS"

  third_party_signal_providers {
    device_posture_id_p {
      compliant = true
    }
  }
}

data "okta_policy_device_assurance_chromeos" "example" {
  id = okta_policy_device_assurance_chromeos.example.id
}
