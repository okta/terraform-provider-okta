resource "okta_request_setting_resource" "test" {
  id = "0oaoum6j3cElINe1z1d7"
  risk_settings {
    default_setting {
      request_submission_type = "ALLOWED_WITH_OVERRIDES"
      approval_sequence_id    = "68920b41386747a673869356"
    }
  }
  request_on_behalf_of_settings {
    allowed = true
  }
}
