resource "okta_request_setting_resource" "test" {
  id = "0oa15zse2by0vGzyP1d8"
  risk_settings {
    default_setting {
      request_submission_type = "RESTRICTED"
    }
  }
}
