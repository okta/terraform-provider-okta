---
layout: 'okta'
page_title: 'Okta: okta_app_saml'
sidebar_current: 'docs-okta-datasource-app-saml'
description: |-
  Get a SAML application from Okta.
---

# okta_app_saml

Use this data source to retrieve an SAML application from Okta.

## Example Usage

```hcl
data "okta_app_saml" "example" {
  label = "Example App"
}
```

## Arguments Reference

- `active_only` - (Optional) tells the provider to query for only `ACTIVE` applications.

- `id` - (Optional) `id` of application to retrieve, conflicts with `label` and `label_prefix`.

- `label` - (Optional) The label of the app to retrieve, conflicts with `label_prefix` and `id`. Label uses
  the `?q=<label>` query parameter exposed by Okta's API. It should be noted that at this time this searches both `name`
  and `label`. This is used to avoid paginating through all applications.

- `label_prefix` - (Optional) Label prefix of the app to retrieve, conflicts with `label` and `id`. This will tell the
  provider to do a `starts with` query as opposed to an `equals` query.

## Attributes Reference

- `accessibility_error_redirect_url` - Custom error page URL.

- `accessibility_login_redirect_url` - Custom login page URL.

- `accessibility_self_service` - Enable self-service.

- `acs_endpoints` - An array of ACS endpoints. You can configure a maximum of 100 endpoints.

- `app_settings_json` - Application settings in JSON format.

- `assertion_signed` - Determines whether the SAML assertion is digitally signed.

- `attribute_statements` - List of SAML Attribute statements.
    - `name` - The name of the attribute statement.
    - `filter_type` - Type of group attribute filter.
    - `filter_value` - Filter value to use.
    - `namespace` - The attribute namespace.
    - `type` - The type of attribute statement value.
    - `values` - Array of values to use.

- `audience` - Audience restriction.

- `authn_context_class_ref` - Identifies the SAML authentication context class for the assertion’s authentication statement.

- `auto_submit_toolbar` - Display auto submit toolbar.

- `default_relay_state` - Identifies a specific application resource in an IDP initiated SSO scenario.

- `destination` - Identifies the location where the SAML response is intended to be sent inside the SAML assertion.

- `digest_algorithm` - Determines the digest algorithm used to digitally sign the SAML assertion and response.

- `features` - features enabled.

- `groups` - List of groups IDs assigned to the application.
  - `DEPRECATED`: Please replace all usage of this field with the data source `okta_app_group_assignments`.

- `hide_ios` - Do not display application icon on mobile app.

- `hide_web` - Do not display application icon to users

- `honor_force_authn` - Prompt user to re-authenticate if SP asks for it.

- `id` - id of application.

- `idp_issuer` - SAML issuer ID.

- `inline_hook_id` - Saml Inline Hook associated with the application.

- `key_id` - Certificate key ID.

- `label` - label of application.

- `links` - Generic JSON containing discoverable resources related to the app.

- `name` - name of application.

- `recipient` - The location where the app may present the SAML assertion.

- `request_compressed` - Denotes whether the request is compressed or not.

- `response_signed` - Determines whether the SAML auth response message is digitally signed.

- `saml_signed_request_enabled` - SAML Signed Request enabled

- `signature_algorithm` - Signature algorithm used to digitally sign the assertion and response.

- `single_logout_certificate` - x509 encoded certificate that the Service Provider uses to sign Single Logout requests.

- `single_logout_issuer` - The issuer of the Service Provider that generates the Single Logout request.

- `single_logout_url` - The location where the logout response is sent.

- `sp_issuer` - SAML service provider issuer.

- `sso_url` - Single Sign-on Url.

- `status` - status of application.

- `subject_name_id_format` - Identifies the SAML processing rules.

- `subject_name_id_template` - Template for app user's username when a user is assigned to the app.

- `user_name_template_push_status` - Push username on update.

- `user_name_template_suffix` - Username template suffix.

- `user_name_template_type` - Username template type.

- `user_name_template` - Username template.