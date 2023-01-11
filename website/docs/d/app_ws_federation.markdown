---
layout: 'okta'
page_title: 'Okta: okta_app_ws_federation'
sidebar_current: 'docs-okta-datasource-app-ws-federation'
description: |-
  Get a WS Fed application from Okta.
---

# okta_app_ws_federation

Use this data source to retrieve an WS Fed application from Okta.

## Example Usage

```hcl
data "okta_app_ws_federation" "example" {
  label = "ExampleApp"
}
```

## Arguments Reference

- `label` - (Optional) The label of the app to retrieve, conflicts with `label_prefix` and `id`. Label uses
  the `?q=<label>` query parameter exposed by Okta's API. It should be noted that at this time this searches both `name`
  and `label`. This is used to avoid paginating through all applications.

- `label_prefix` - (Optional) Label prefix of the app to retrieve, conflicts with `label` and `id`. This will tell the
  provider to do a `starts with` query as opposed to an `equals` query.

- `id` - (Optional) `id` of application to retrieve, conflicts with `label` and `label_prefix`.

- `active_only` - (Optional) tells the provider to query for only `ACTIVE` applications.

- `skip_users` - (Optional) Indicator that allows the app to skip `users` sync. Default is `false`.

- `skip_groups` - (Optional) Indicator that allows the app to skip `groups` sync. Default is `false`.

## Attributes Reference

- `id` - id of application.

- `label` - label of application.

- `name` - name of application.

- `status` - status of application.

- `site_url` - Launch URL for the Web Application.

- `realm` - The trust realm for the Web Application.

- `reply_url` - The ReplyTo URL to which responses are directed.

- `reply_override` - Enable web application to override ReplyTo URL with reply param.

- `name_id_format` - Name ID Format.

- `audience_restriction` - The assertion containing a bearer subject confirmation MUST contain an Audience Restriction including the service provider's unique identifier as an Audience.

- `authn_context_class_ref` - Specifies the Authentication Context for the issued SAML Assertion.

- `group_filter` - An expression that will be used to filter groups. If the Okta group name matches the expression, the group name will be included in the SAML Assertion Attribute Statement

- `group_name` - Specifies the SAML attribute name for a user's group memberships.

- `group_value_format` - Specifies the SAML assertion attribute value for filtered groups.

- `username_attribute` - Specifies additional username attribute statements to include in the SAML Assertion.

- `attribute_statements` - Defines custom SAML attribute statements.

- `visibility` - Application icon visibility to users.