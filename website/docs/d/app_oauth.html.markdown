---
layout: 'okta'
page_title: 'Okta: okta_app_oauth'
sidebar_current: 'docs-okta-datasource-app-oauth'
description: |-
Get a OIDC application from Okta.
---

# okta_app_oauth

Use this data source to retrieve an OIDC application from Okta.

## Example Usage

```hcl
data "okta_app_oauth" "test" {
  label = "Example App"
}
```

## Argument Reference

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

- `id` - ID of application.

- `label` - Label of application.

- `name` - Name of application.

- `status` - Status of application.

- `type` - The type of OAuth application.

- `auto_submit_toolbar` - Display auto submit toolbar.

- `hide_ios` - Do not display application icon on mobile app.

- `hide_web` - Do not display application icon to users.

- `grant_types` - List of OAuth 2.0 grant types.

- `response_types` - List of OAuth 2.0 response type strings.

- `redirect_uris` - List of URIs for use in the redirect-based flow.

- `post_logout_redirect_uris` - List of URIs for redirection after logout.

- `logo_uri` - URI that references a logo for the client.

- `login_uri` - URI that initiates login.

- `login_mode` - The type of Idp-Initiated login that the client supports, if any.

- `login_scopes` - List of scopes to use for the request.

- `client_id` - OAuth client ID. If set during creation, app is created with this id.

- `client_uri` - URI to a web page providing information about the client.

- `policy_uri` - URI to web page providing client policy document.

- `links` - generic JSON containing discoverable resources related to the app

- `users` - List of users IDs assigned to the application.
  - `DEPRECATED`: Please replace all usage of this field with the data source `okta_app_user_assignments`.

- `groups` - List of groups IDs assigned to the application.
  - `DEPRECATED`: Please replace all usage of this field with the data source `okta_app_group_assignments`.
