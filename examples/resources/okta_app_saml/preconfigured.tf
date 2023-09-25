resource "okta_app_saml" "test" {
  preconfigured_app = "office365"
  label             = "testAcc_replace_with_uuid"
  status            = "ACTIVE"
  saml_version      = "1.1"
  app_settings_json = <<JSON
    {
       "wsFedConfigureType": "AUTO",
       "windowsTransportEnabled": false,
       "domain": "okta.com",
       "msftTenant": "okta",
       "domains": [],
       "requireAdminConsent": false
    }
JSON
  app_links_json    = <<JSON
  {
      "calendar": false,
      "crm": false,
      "delve": false,
      "excel": false,
      "forms": false,
      "mail": false,
      "newsfeed": false,
      "onedrive": false,
      "people": false,
      "planner": false,
      "powerbi": false,
      "powerpoint": false,
      "sites": false,
      "sway": false,
      "tasks": false,
      "teams": false,
      "word": false,
      "yammer": false,
      "login": true
  }
JSON
}
