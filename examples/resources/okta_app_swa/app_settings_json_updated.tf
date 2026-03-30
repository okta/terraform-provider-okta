resource "okta_app_swa" "test" {
  preconfigured_app = "office365"
  label             = "testAcc_replace_with_uuid"
  status            = "INACTIVE"
  app_settings_json = <<JSON
    {
      "wsFedConfigureType": "AUTO",
      "windowsTransportEnabled": false,
      "domain": "example-updated.com",
      "msftTenant": "exampletenant",
      "domains": [],
      "requireAdminConsent": false
    }
JSON
  app_links_json = <<JSON
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
      "video": false,
      "word": false,
      "yammer": false,
      "login": true
    }
JSON
  user_name_template      = "user.login"
  user_name_template_type = "CUSTOM"
}
