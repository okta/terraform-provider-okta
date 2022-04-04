resource "okta_app_saml" "test" {
  app_settings_json = <<JSON
{
    "groupFilter": "app1.*",
    "siteURL": "http://www.okta.com"
}
JSON
  label             = "testAcc_replace_with_uuid"
  preconfigured_app = "sharepoint_onpremise"
  saml_version      = "1.1"
}
