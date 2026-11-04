resource "okta_device_posture_check" "test" {
  name          = "testAcc_replace_with_uuid"
  description   = "Query macOS devices to check if firewall is enabled"
  variable_name = "testvar"
  platform      = "MACOS"
  type          = "CUSTOM"
  mapping_type  = "CHECKBOX"
  query         = "SELECT CASE WHEN global_state = 0 THEN 0 ELSE 1 END AS firewall_enabled FROM alf;"

  remediation_settings {
    link {
      custom_url = "https://example.myremediationurl.com/docs"
    }
    message {
      custom_text = "It is our company policy that all devices must have a firewall enabled."
    }
  }
}

data "okta_device_posture_check" "test" {
  id = okta_device_posture_check.test.id
}