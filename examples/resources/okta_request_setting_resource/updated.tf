resource "okta_request_setting_resource" "test" {
  id = "0oaoum6j3cElINe1z1d7"
  risk_settings {
    default_setting {
      request_submission_type = "RESTRICTED"
    }
  }
}
