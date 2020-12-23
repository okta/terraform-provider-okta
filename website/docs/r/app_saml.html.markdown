---
layout: 'okta'
page_title: 'Okta: okta_app_saml'
sidebar_current: 'docs-okta-resource-app-saml'
description: |-
  Creates an SAML Application.
---

# okta_app_saml

Creates an SAML Application.

This resource allows you to create and configure an SAML Application.

## Example Usage

```hcl
resource "okta_app_saml" "example" {
  label                    = "example"
  sso_url                  = "http://example.com"
  recipient                = "http://example.com"
  destination              = "http://example.com"
  audience                 = "http://example.com/audience"
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

## Argument Reference

The following arguments are supported:

- `label` - (Required) label of application.

- `preconfigured_app` - (Optional) name of application from the Okta Integration Network, if not included a custom app will be created.

- `description` - (Optional) description of application.

- `status` - (Optional) status of application.

- `auto_submit_toolbar` - (Optional) Display auto submit toolbar.

- `hide_ios` - (Optional) Do not display application icon on mobile app.

- `hide_web` - (Optional) Do not display application icon to users

- `default_relay_state` - (Optional) Identifies a specific application resource in an IDP initiated SSO scenario.

- `sso_url` - (Optional) Single Sign on Url.

- `recipient` - (Optional) The location where the app may present the SAML assertion.

- `destination` - (Optional) Identifies the location where the SAML response is intended to be sent inside of the SAML assertion.

- `audience` - (Optional) Audience restriction.

- `idp_issuer` - (Optional) SAML issuer ID.

- `sp_issuer` - (Optional) SAML service provider issuer.

- `subject_name_id_template` - (Optional) Template for app user's username when a user is assigned to the app.

- `subject_name_id_format` - (Optional) Identifies the SAML processing rules.

- `response_signed` - (Optional) Determines whether the SAML auth response message is digitally signed.

- `request_compressed` - (Optional) Denotes whether the request is compressed or not.

- `assertion_signed` - (Optional) Determines whether the SAML assertion is digitally signed.

- `signature_algorithm` - (Optional) Signature algorithm used ot digitally sign the assertion and response.

- `digest_algorithm` - (Optional) Determines the digest algorithm used to digitally sign the SAML assertion and response.

- `honor_force_authn` - (Optional) Prompt user to re-authenticate if SP asks for it.

- `authn_context_class_ref` - (Optional) Identifies the SAML authentication context class for the assertion’s authentication statement.

- `accessibility_self_service` - (Optional) Enable self service.

- `accessibility_error_redirect_url` - (Optional) Custom error page URL.

- `accessibility_login_redirect_url` - (Optional) Custom login page URL.

- `features` - (Optional) features enabled.

- `user_name_template` - (Optional) Username template.

- `user_name_template_suffix` - (Optional) Username template suffix.

- `user_name_template_type` - (Optional) Username template type.

- `app_settings_json` - (Optional) Application settings in JSON format.

- `acs_endpoints` - An array of ACS endpoints. You can configure a maximum of 100 endpoints.

- `attribute_statements` - (Optional) List of SAML Attribute statements.
  - `name` - (Required) The name of the attribute statement.
  - `filter_type` - (Optional) Type of group attribute filter.
  - `filter_value` - (Optional) Filter value to use.
  - `namespace` - (Optional) The attribute namespace. It can be set to `"urn:oasis:names:tc:SAML:2.0:attrname-format:unspecified"`, `"urn:oasis:names:tc:SAML:2.0:attrname-format:uri"`, or `"urn:oasis:names:tc:SAML:2.0:attrname-format:basic"`.
  - `type` - (Optional) The type of attribute statement value. Can be `"EXPRESSION"` or `"GROUP"`.
  - `values` - (Optional) Array of values to use.

## Attributes Reference

- `id` - id of application.

- `name` - Name assigned to the application by Okta.

- `sign_on_mode` - Sign on mode of application.

- `key_id` - Certificate key ID.

- `key_name` - Certificate name. This modulates the rotation of keys. New name == new key.

- `certificate` - The raw signing certificate.

- `metadata` - The raw SAML metadata in XML.

- `http_post_binding` - `urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Post` location from the SAML metadata.

- `http_redirect_binding` - `urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect` location from the SAML metadata.

- `entity_key` - Entity ID, the ID portion of the `entity_url`.

- `entity_url` - Entity URL for instance [http://www.okta.com/exk1fcia6d6EMsf331d8](http://www.okta.com/exk1fcia6d6EMsf331d8).

## Import

A SAML App can be imported via the Okta ID.

```
$ terraform import okta_app_saml.example <app id>
```
