resource "okta_request_setting_resource" "test" {
  id = "0oa15zse2by0vGzyP1d8"
  risk_settings {
    default_setting {
      request_submission_type = "ALLOWED_WITH_OVERRIDES"
      approval_sequence_id    = "69cbcb93a55844b4f3090b11"
    }
  }
  request_on_behalf_of_settings {
    allowed = true
  }
}
