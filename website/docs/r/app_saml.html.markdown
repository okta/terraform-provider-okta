---
layout: 'okta'
page_title: 'Okta: okta_app_saml'
sidebar_current: 'docs-okta-resource-app-saml'
description: |-
  Creates a SAML Application.
---

# okta_app_saml

This resource allows you to create and configure a SAML Application.

## Example Usage

```hcl
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
```

### With inline hook

```hcl
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
  label                     = "testAcc_replace_with_uuid"
  sso_url                   = "https://google.com"
  recipient                 = "https://here.com"
  destination               = "https://its-about-the-journey.com"
  audience                  = "https://audience.com"
  subject_name_id_template  = "$${user.userName}"
  subject_name_id_format    = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
  response_signed           = true
  signature_algorithm       = "RSA_SHA256"
  digest_algorithm          = "SHA256"
  honor_force_authn         = false
  authn_context_class_ref   = "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"
  inline_hook_id            = okta_inline_hook.test.id

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
```

### Pre-configured app with SAML 1.1 sign-on mode

```hcl
resource "okta_app_saml" "test" {
  app_settings_json = <<JSON
{
    "groupFilter": "app1.*",
    "siteURL": "https://www.okta.com"
}
JSON
  label = "SharePoint (On-Premise)"
  preconfigured_app = "sharepoint_onpremise"
  saml_version = "1.1"
  status = "ACTIVE"
  user_name_template = "$${source.login}"
  user_name_template_type = "BUILT_IN"
}
```

### Pre-configured app with SAML 1.1 sign-on mode, `app_settings_json` and `app_links_json`

```hcl
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
```

## Argument Reference

The following arguments are supported:

- `accessibility_error_redirect_url` - (Optional) Custom error page URL.

- `accessibility_login_redirect_url` - (Optional) Custom login page for this application.

- `accessibility_self_service` - (Optional) Enable self-service. Default is: `false`.

- `acs_endpoints` - (Optional) An array of ACS endpoints. You can configure a maximum of 100 endpoints.

- `admin_note` - (Optional) Application notes for admins.

- `app_links_json` - (Optional) Displays specific appLinks for the app. The value for each application link should be boolean.

- `app_settings_json` - (Optional) Application settings in JSON format.

- `assertion_signed` - (Optional) Determines whether the SAML assertion is digitally signed.

- `attribute_statements` - (Optional) List of SAML Attribute statements.
  - `name` - (Required) The name of the attribute statement.
  - `filter_type` - (Optional) Type of group attribute filter. Valid values are: `"STARTS_WITH"`, `"EQUALS"`, `"CONTAINS"`, or `"REGEX"`
  - `filter_value` - (Optional) Filter value to use.
  - `namespace` - (Optional) The attribute namespace. It can be set to `"urn:oasis:names:tc:SAML:2.0:attrname-format:unspecified"`, `"urn:oasis:names:tc:SAML:2.0:attrname-format:uri"`, or `"urn:oasis:names:tc:SAML:2.0:attrname-format:basic"`.
  - `type` - (Optional) The type of attribute statement value. Valid values are: `"EXPRESSION"` or `"GROUP"`. Default is `"EXPRESSION"`.
  - `values` - (Optional) Array of values to use.

- `audience` - (Optional) Audience restriction.

- `authn_context_class_ref` - (Optional) Identifies the SAML authentication context class for the assertion’s authentication statement.

- `auto_submit_toolbar` - (Optional) Display auto submit toolbar. Default is: `false`

- `default_relay_state` - (Optional) Identifies a specific application resource in an IDP initiated SSO scenario.

- `destination` - (Optional) Identifies the location where the SAML response is intended to be sent inside the SAML assertion.

- `digest_algorithm` - (Optional) Determines the digest algorithm used to digitally sign the SAML assertion and response.

- `enduser_note` - (Optional) Application notes for end users.

- `features` - (Optional) features enabled. Notice: you can't currently configure provisioning features via the API.

- `groups` - (Optional) Groups associated with the application.
  - `DEPRECATED`: Please replace usage with the `okta_app_group_assignments` (or `okta_app_group_assignment`) resource.

- `hide_ios` - (Optional) Do not display application icon on mobile app. Default is: `false`

- `hide_web` - (Optional) Do not display application icon to users. Default is: `false`

- `honor_force_authn` - (Optional) Prompt user to re-authenticate if SP asks for it. Default is: `false`

- `idp_issuer` - (Optional) SAML issuer ID.

- `implicit_assignment` - (Optional) *Early Access Property*. Enables [Federation Broker Mode]( https://help.okta.com/en/prod/Content/Topics/Apps/apps-fbm-enable.htm). When this mode is enabled, `users` and `groups` arguments are ignored.

- `inline_hook_id` - (Optional) Saml Inline Hook associated with the application.

- `key_name` - (Optional) Certificate name. This modulates the rotation of keys. New name == new key. Required to be set with `key_years_valid`.

- `key_years_valid` - (Optional) Number of years the certificate is valid (2 - 10 years).

- `label` - (Required) label of application.

- `logo` - (Optional) Local file path to the logo. The file must be in PNG, JPG, or GIF format, and less than 1 MB in size.

- `preconfigured_app` - (Optional) name of application from the Okta Integration Network, if not included a custom app will be created.  
  If not provided the following arguments are required:
  - `sso_url`
  - `recipient`
  - `destination`
  - `audience`
  - `subject_name_id_template`
  - `subject_name_id_format`
  - `signature_algorithm`
  - `digest_algorithm`
  - `authn_context_class_ref`

- `recipient` - (Optional) The location where the app may present the SAML assertion.

- `request_compressed` - (Optional) Denotes whether the request is compressed or not.

- `response_signed` - (Optional) Determines whether the SAML auth response message is digitally signed.

- `saml_version` - (Optional) SAML version for the app's sign-on mode. Valid values are: `"2.0"` or `"1.1"`. Default is `"2.0"`.

- `signature_algorithm` - (Optional) Signature algorithm used ot digitally sign the assertion and response.

- `single_logout_certificate` - (Optional) x509 encoded certificate that the Service Provider uses to sign Single Logout requests.
  Note: should be provided without `-----BEGIN CERTIFICATE-----` and `-----END CERTIFICATE-----`, see [official documentation](https://developer.okta.com/docs/reference/api/apps/#service-provider-certificate).

- `single_logout_issuer` - (Optional) The issuer of the Service Provider that generates the Single Logout request.

- `single_logout_url` - (Optional) The location where the logout response is sent.

- `skip_groups` - (Optional) Indicator that allows the app to skip `groups` sync (it can also be provided during import). Default is `false`.

- `skip_users` - (Optional) Indicator that allows the app to skip `users` sync (it can also be provided during import). Default is `false`.

- `sp_issuer` - (Optional) SAML service provider issuer.

- `sso_url` - (Optional) Single Sign-on Url.

- `status` - (Optional) status of application.

- `subject_name_id_format` - (Optional) Identifies the SAML processing rules.

- `subject_name_id_template` - (Optional) Template for app user's username when a user is assigned to the app.

- `user_name_template` - (Optional) Username template. Default is: `"${source.login}"`

- `user_name_template_push_status` - (Optional) Push username on update. Valid values: `"PUSH"` and `"DONT_PUSH"`.

- `user_name_template_suffix` - (Optional) Username template suffix.

- `user_name_template_type` - (Optional) Username template type. Default is: `"BUILT_IN"`.

- `users` - (Optional) Users associated with the application.
  - `DEPRECATED`: Please replace usage with the `okta_app_user` resource.

- `authentication_policy` - (Optional) The ID of the associated `app_signon_policy`. If this property is removed from the application the `default` sign-on-policy will be associated with this application.

## Attributes Reference

- `id` - id of application.

- `name` - Name assigned to the application by Okta.

- `sign_on_mode` - Sign-on mode of application.

- `key_id` - Certificate key ID.

- `key_name` - Certificate name. This modulates the rotation of keys. New name == new key.

- `certificate` - The raw signing certificate.

- `metadata` - The raw SAML metadata in XML.

- `metadata_url` - SAML xml metadata URL.

- `http_post_binding` - `urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Post` location from the SAML metadata.

- `http_redirect_binding` - `urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect` location from the SAML metadata.

- `entity_key` - Entity ID, the ID portion of the `entity_url`.

- `entity_url` - Entity URL for instance [http://www.okta.com/exk1fcia6d6EMsf331d8](http://www.okta.com/exk1fcia6d6EMsf331d8).

- `logo_url` - Direct link of application logo.

## Import

A SAML App can be imported via the Okta ID.

```
$ terraform import okta_app_saml.example &#60;app id&#62;
```

It's also possible to import app without groups or/and users. In this case ID may look like this:

```
$ terraform import okta_app_basic_auth.example &#60;app id&#62;/skip_users

$ terraform import okta_app_basic_auth.example &#60;app id&#62;/skip_users/skip_groups

$ terraform import okta_app_basic_auth.example &#60;app id&#62;/skip_groups
```
