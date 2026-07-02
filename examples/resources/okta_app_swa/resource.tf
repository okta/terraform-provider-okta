resource "okta_app_swa" "example" {
  label          = "example"
  button_field   = "btn-login"
  password_field = "txtbox-password"
  username_field = "txtbox-username"
  url            = "https://example.com/login.html"
}

resource "okta_app_swa" "example_with_app_settings_json" {
  preconfigured_app       = "office365"
  label                   = "Microsoft Office 365 SWA"
  status                  = "ACTIVE"
  app_settings_json       = <<JSON
    {
      "wsFedConfigureType": "AUTO",
      "windowsTransportEnabled": false,
      "domain": "example.com",
      "msftTenant": "exampletenant",
      "domains": [],
      "requireAdminConsent": false
    }
JSON
  app_links_json          = <<JSON
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
