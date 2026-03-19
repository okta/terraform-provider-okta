resource "okta_app_saml" "example" {
  label                    = "example"
  sso_url                  = "https://example.com"
  recipient                = "https://example.com"
  destination              = "https://example.com"
  audience                 = "https://example.com/audience"
  subject_name_id_template = "$${user.userName}"
  subject_name_id_format   = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
  response_signed          = true
  signature_algorithm      = "RSA_SHA256"
  digest_algorithm         = "SHA256"
  honor_force_authn        = false
  authn_context_class_ref  = "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"

  attribute_statements {
    type         = "GROUP"
    name         = "groups"
    filter_type  = "REGEX"
    filter_value = ".*"
  }
}

### With inline hook
resource "okta_inline_hook" "test" {
  name    = "testAcc_replace_with_uuid"
  status  = "ACTIVE"
  type    = "com.okta.saml.tokens.transform"
  version = "1.0.2"

  channel = {
    type    = "HTTP"
    version = "1.0.0"
    uri     = "https://example.com/test1"
    method  = "POST"
  }
  auth = {
    key   = "Authorization"
    type  = "HEADER"
    value = "secret"
  }
}

resource "okta_app_saml" "test" {
  label                    = "testAcc_replace_with_uuid"
  sso_url                  = "https://google.com"
  recipient                = "https://here.com"
  destination              = "https://its-about-the-journey.com"
  audience                 = "https://audience.com"
  subject_name_id_template = "$${user.userName}"
  subject_name_id_format   = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
  response_signed          = true
  signature_algorithm      = "RSA_SHA256"
  digest_algorithm         = "SHA256"
  honor_force_authn        = false
  authn_context_class_ref  = "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"
  inline_hook_id           = okta_inline_hook.test.id

  depends_on = [
    okta_inline_hook.test
  ]
  attribute_statements {
    type         = "GROUP"
    name         = "groups"
    filter_type  = "REGEX"
    filter_value = ".*"
  }
}

### Pre-configured app with SAML 1.1 sign-on mode
resource "okta_app_saml" "test" {
  app_settings_json       = <<JSON
{
    "groupFilter": "app1.*",
    "siteURL": "https://www.okta.com"
}
JSON
  label                   = "SharePoint (On-Premise)"
  preconfigured_app       = "sharepoint_onpremise"
  saml_version            = "1.1"
  status                  = "ACTIVE"
  user_name_template      = "$${source.login}"
  user_name_template_type = "BUILT_IN"
}

### Pre-configured app with SAML 1.1 sign-on mode, `app_settings_json` and `app_links_json`
resource "okta_app_saml" "office365" {
  preconfigured_app = "office365"
  label             = "Microsoft Office 365"
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
      "video": false,
      "word": false,
      "yammer": false,
      "login": true
  }
JSON
}
