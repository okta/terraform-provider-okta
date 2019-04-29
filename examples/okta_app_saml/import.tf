resource "okta_app_saml" "testAcc_replace_with_uuid" {
  preconfigured_app = "pagerduty"
  label             = "testAcc_replace_with_uuid"

  app_settings_json = <<JSON
{
  "subdomain": "articulate"
}
JSON
}
