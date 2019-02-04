resource "okta_saml_app" "testAcc_%[1]d" {
  preconfigured_app = "pagerduty"
  label             = "testAcc_%[1]d"

  app_settings_json = <<JSON
{
  "subdomain": "articulate"
}
JSON
}
