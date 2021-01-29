---
layout: 'okta'
page_title: 'Okta: okta_app_saml'
sidebar_current: 'docs-okta-datasource-app-saml'
description: |-
  Get a SAML application from Okta.
---

# okta_app_saml

Use this data source to retrieve the collaborators for a given repository.

## Example Usage

```hcl
data "okta_app_saml" "example" {
  label = "Example App"
}
```

## Arguments Reference

- `label` - (Optional) The label of the app to retrieve, conflicts with `label_prefix` and `id`.

- `label_prefix` - (Optional) Label prefix of the app to retrieve, conflicts with `label` and `id`. This will tell the provider to do a `starts with` query as opposed to an `equals` query.

- `id` - (Optional) `id` of application to retrieve, conflicts with `label` and `label_prefix`.

- `active_only` - (Optional) tells the provider to query for only `ACTIVE` applications.

## Attributes Reference

- `id` - id of application.

- `label` - label of application.

- `name` - name of application.

- `status` - status of application.

- `key_id` - Certificate key ID.

- `auto_submit_toolbar` - Display auto submit toolbar.

- `hide_ios` - Do not display application icon on mobile app.

- `hide_web` - Do not display application icon to users

- `default_relay_state` - Identifies a specific application resource in an IDP initiated SSO scenario.

- `sso_url` - Single Sign-on Url.

- `recipient` - The location where the app may present the SAML assertion.

- `destination` - Identifies the location where the SAML response is intended to be sent inside the SAML assertion.

- `audience` - Audience restriction.

- `idp_issuer` - SAML issuer ID.

- `sp_issuer` - SAML service provider issuer.

- `subject_name_id_template` - Template for app user's username when a user is assigned to the app.

- `subject_name_id_format` - Identifies the SAML processing rules.

- `response_signed` - Determines whether the SAML auth response message is digitally signed.

- `request_compressed` - Denotes whether the request is compressed or not.

- `assertion_signed` - Determines whether the SAML assertion is digitally signed.

- `signature_algorithm` - Signature algorithm used ot digitally sign the assertion and response.

- `digest_algorithm` - Determines the digest algorithm used to digitally sign the SAML assertion and response.

- `honor_force_authn` - Prompt user to re-authenticate if SP asks for it.

- `authn_context_class_ref` - Identifies the SAML authentication context class for the assertion’s authentication statement.

- `accessibility_self_service` - Enable self-service.

- `accessibility_error_redirect_url` - Custom error page URL.

- `accessibility_login_redirect_url` - Custom login page URL.

- `features` - features enabled.

- `user_name_template` - Username template.

- `user_name_template_suffix` - Username template suffix.

- `user_name_template_type` - Username template type.

- `app_settings_json` - Application settings in JSON format.

- `acs_endpoints` - An array of ACS endpoints. You can configure a maximum of 100 endpoints.

- `attribute_statements` - List of SAML Attribute statements.
  - `name` - The name of the attribute statement.
  - `filter_type` - Type of group attribute filter.
  - `filter_value` - Filter value to use.
  - `namespace` - The attribute namespace.
  - `type` - The type of attribute statement value.
  - `values` - Array of values to use.

- `single_logout_issuer` - The issuer of the Service Provider that generates the Single Logout request.

- `single_logout_url` - The location where the logout response is sent.

- `single_logout_certificate` - x509 encoded certificate that the Service Provider uses to sign Single Logout requests.
