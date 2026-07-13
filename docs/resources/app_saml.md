---
page_title: "Resource: okta_app_saml"
description: |-
  This resource allows you to create and configure a SAML Application.
  -> During an apply if there is change in 'status' the app will first be
  activated or deactivated in accordance with the 'status' change. Then, all
  other arguments that changed will be applied.
  -> If you receive the error 'You do not have permission to access the feature
  you are requesting' contact support mailto:dev-inquiries@okta.com and
  request feature flag 'ADVANCED_SSO' be applied to your org.
---

# Resource: okta_app_saml

This resource allows you to create and configure a SAML Application.
-> During an apply if there is change in 'status' the app will first be
activated or deactivated in accordance with the 'status' change. Then, all
other arguments that changed will be applied.
		
-> If you receive the error 'You do not have permission to access the feature
you are requesting' [contact support](mailto:dev-inquiries@okta.com) and
request feature flag 'ADVANCED_SSO' be applied to your org.

## Example Usage

```terraform
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

### Example to demonstrate usage of acs_endpoints_indices
resource "okta_app_saml" "test" {
  label           = "dwdef"
  sso_url         = "https://example.com"
  recipient       = "https://example.com"
  destination     = "https://example.com"
  audience        = "example"
  response_signed = true
  acs_endpoints_indices {
    url   = "https://example5.com"
    index = 2999
  }
  acs_endpoints_indices {
    url   = "https://example1.com"
    index = 102
  }
  acs_endpoints_indices {
    url   = "https://example5.com"
    index = 29
  }
  acs_endpoints_indices {
    url   = "https://example4.com"
    index = 19
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `label` (String) The Application's display name.

### Optional

- `accessibility_error_redirect_url` (String) Custom error page URL
- `accessibility_login_redirect_url` (String) Custom login page URL
- `accessibility_self_service` (Boolean) Enable self service. Default is `false`
- `acs_endpoints` (Set of String) An array of ACS endpoints. You can configure a maximum of 100 endpoints.
- `admin_note` (String) Application notes for admins.
- `app_links_json` (String) Displays specific appLinks for the app. The value for each application link should be boolean.
- `app_settings_json` (String) Application settings in JSON format
- `assertion_signed` (Boolean) Determines whether the SAML assertion is digitally signed
- `attribute_statements` (Block List) (see [below for nested schema](#nestedblock--attribute_statements))
- `audience` (String) Audience Restriction
- `authentication_policy` (String) The ID of the associated `app_signon_policy`. If this property is removed from the application the `default` sign-on-policy will be associated with this application.y
- `authn_context_class_ref` (String) Identifies the SAML authentication context class for the assertion’s authentication statement
- `auto_submit_toolbar` (Boolean) Display auto submit toolbar. Default is: `false`
- `default_relay_state` (String) Identifies a specific application resource in an IDP initiated SSO scenario.
- `destination` (String) Identifies the location where the SAML response is intended to be sent inside of the SAML assertion
- `digest_algorithm` (String) Determines the digest algorithm used to digitally sign the SAML assertion and response
- `enduser_note` (String) Application notes for end users.
- `hide_ios` (Boolean) Do not display application icon on mobile app
- `hide_web` (Boolean) Do not display application icon to users
- `honor_force_authn` (Boolean) Prompt user to re-authenticate if SP asks for it. Default is: `false`
- `idp_issuer` (String) SAML issuer ID
- `implicit_assignment` (Boolean) *Early Access Property*. Enable Federation Broker Mode.
- `inline_hook_id` (String) Saml Inline Hook setting
- `key_name` (String) Certificate name. This modulates the rotation of keys. New name == new key. Required to be set with `key_years_valid`
- `key_years_valid` (Number) Number of years the certificate is valid (2 - 10 years).
- `logo` (String) Local file path to the logo. The file must be in PNG, JPG, or GIF format, and less than 1 MB in size.
- `preconfigured_app` (String) Name of application from the Okta Integration Network. For instance 'slack'. If not included a custom app will be created.  If not provided the following arguments are required:
'sso_url'
'recipient'
'destination'
'audience'
'subject_name_id_template'
'subject_name_id_format'
'signature_algorithm'
'digest_algorithm'
'authn_context_class_ref'
- `recipient` (String) The location where the app may present the SAML assertion
- `request_compressed` (Boolean) Denotes whether the request is compressed or not.
- `response_signed` (Boolean) Determines whether the SAML auth response message is digitally signed
- `saml_signed_request_enabled` (Boolean) SAML Signed Request enabled
- `saml_version` (String) SAML version for the app's sign-on mode. Valid values are: `2.0` or `1.1`. Default is `2.0`
- `signature_algorithm` (String) Signature algorithm used to digitally sign the assertion and response
- `single_logout_certificate` (String) x509 encoded certificate that the Service Provider uses to sign Single Logout requests. Note: should be provided without `-----BEGIN CERTIFICATE-----` and `-----END CERTIFICATE-----`, see [official documentation](https://developer.okta.com/docs/reference/api/apps/#service-provider-certificate).
- `single_logout_issuer` (String) The issuer of the Service Provider that generates the Single Logout request
- `single_logout_url` (String) The location where the logout response is sent
- `skip_metadata` (Boolean) Skip reading SAML metadata during read operations, this can improve performance for large numbers of applications. When enabled, the following computed attributes will be omitted from state: metadata, metadata_url, http_post_binding, http_redirect_binding, entity_key, entity_url, certificate. Default is `false`.
- `skip_keys` (Boolean) Skip reading/writing key credentials during read/write operations, this can improve performance for large numbers of applications. When enabled, the following computed attributes will be omitted from state: key_id, keys. Default is `false`.
- `sp_issuer` (String) SAML SP issuer ID
- `sso_url` (String) Single Sign On URL
- `status` (String) Status of application. By default, it is `ACTIVE`
- `subject_name_id_format` (String) Identifies the SAML processing rules.
- `subject_name_id_template` (String) Template for app user's username when a user is assigned to the app
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `user_name_template` (String) Username template. Default: `${source.login}`
- `user_name_template_push_status` (String) Push username on update. Valid values: `PUSH` and `DONT_PUSH`
- `user_name_template_suffix` (String) Username template suffix
- `user_name_template_type` (String) Username template type. Default: `BUILT_IN`
- `acs_endpoints_indices` (Set) ACS endpoints along with custom index as a set of maps called `acs_endpoints_indices` in JSON format.

### Read-Only

- `certificate` (String) cert from SAML XML metadata payload
- `embed_url` (String) The url that can be used to embed this application in other portals.
- `entity_key` (String) Entity ID, the ID portion of the entity_url
- `entity_url` (String) Entity URL for instance http://www.okta.com/exk1fcia6d6EMsf331d8
- `features` (Set of String) features to enable
- `http_post_binding` (String) urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Post location from the SAML metadata.
- `http_redirect_binding` (String) urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect location from the SAML metadata.
- `id` (String) The ID of this resource.
- `key_id` (String) Certificate ID
- `keys` (List of Object) Application keys (see [below for nested schema](#nestedatt--keys))
- `logo_url` (String) URL of the application's logo
- `metadata` (String) SAML xml metadata payload
- `metadata_url` (String) SAML xml metadata URL
- `name` (String) Name of the app.
- `sign_on_mode` (String) Sign on mode of application.

<a id="nestedblock--attribute_statements"></a>
### Nested Schema for `attribute_statements`

Required:

- `name` (String) The reference name of the attribute statement

Optional:

- `filter_type` (String) Type of group attribute filter. Valid values are: `STARTS_WITH`, `EQUALS`, `CONTAINS`, or `REGEX`
- `filter_value` (String) Filter value to use
- `namespace` (String) The attribute namespace. It can be set to `urn:oasis:names:tc:SAML:2.0:attrname-format:unspecified`, `urn:oasis:names:tc:SAML:2.0:attrname-format:uri`, or `urn:oasis:names:tc:SAML:2.0:attrname-format:basic`
- `type` (String) The type of attribute statements object
- `values` (List of String)


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `read` (String)
- `update` (String)


<a id="nestedatt--keys"></a>
### Nested Schema for `keys`

Read-Only:

- `created` (String)
- `e` (String)
- `expires_at` (String)
- `kid` (String)
- `kty` (String)
- `last_updated` (String)
- `n` (String)
- `use` (String)
- `x5c` (List of String)
- `x5t_s256` (String)

## Import

Import is supported using the following syntax:

```shell
terraform import okta_app_saml.example <app_id>
```
