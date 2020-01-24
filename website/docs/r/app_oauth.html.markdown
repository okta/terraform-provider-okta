---
layout: "okta"
page_title: "Okta: okta_app_oauth"
sidebar_current: "docs-okta-resource-app-oauth"
description: |-
  Creates an OIDC Application.
---

# okta_app_oauth

Creates an OIDC Application.

This resource allows you to create and configure an OIDC Application.

## Example Usage

```hcl
resource "okta_app_oauth" "example" {
  label                      = "example"
  type                       = "web"
  grant_types                = ["authorization_code"]
  redirect_uris              = ["https://example.com/"]
  response_types             = ["code"]
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) The Application's display name.

* `status` - (Optional) The status of the application, by default it is `"ACTIVE"`.

* `type` - (Required) The type of OAuth application.

* `users` - (Optional) The users assigned to the application. It is recommended not to use this and instead use `okta_app_user`.

* `groups` - (Optional) The groups assigned to the application. It is recommended not to use this and instead use `okta_app_group_assignment`.

* `client_id` - (Optional) OAuth client ID. If set during creation, app is created with this id.

* `omit_secret` - (Optional) This tells the provider not to persist the application's secret to state. If this is ever changes from true => false your app will be recreated.

* `client_basic_secret` - (Optional) OAuth client secret key, this can be set when token_endpoint_auth_method is client_secret_basic.

* `token_endpoint_auth_method` - (Optional) Requested authentication method for the token endpoint. It can be set to `"none"`, `"client_secret_post"`, `"client_secret_basic"`, `"client_secret_jwt"`.

* `auto_key_rotation` - (Optional) Requested key rotation mode.

* `client_uri` - (Optional) URI to a web page providing information about the client.

* `logo_uri` - (Optional) URI that references a logo for the client.

* `login_uri` - (Optional) URI that initiates login.

* `redirect_uris` - (Optional) List of URIs for use in the redirect-based flow. This is required for all application types except service.

* `post_logout_redirect_uris` - (Optional) List of URIs for redirection after logout.

* `response_types` - (Optional) List of OAuth 2.0 response type strings.

* `grant_types` - (Optional) List of OAuth 2.0 grant types. Conditional validation params found here https://developer.okta.com/docs/api/resources/apps#credentials-settings-details. Defaults to minimum requirements per app type.

* `tos_uri` - (Optional) URI to web page providing client tos (terms of service).

* `policy_uri` - (Optional) URI to web page providing client policy document.

* `consent_method` - (Optional) Indicates whether user consent is required or implicit. Valid values: REQUIRED, TRUSTED. Default value is TRUSTED.

* `issuer_mode` - (Optional) Indicates whether the Okta Authorization Server uses the original Okta org domain URL or a custom domain URL as the issuer of ID token for this client.

* `auto_submit_toolbar` - (Optional) Display auto submit toolbar.

* `hide_ios` - (Optional) Do not display application icon on mobile app.

* `hide_web` - (Optional) Do not display application icon to users.

* `profile` - (Optional) Custom JSON that represents an OAuth application's profile.

## Attributes Reference

* `name` - Name assigned to the application by Okta.

* `sign_on_mode` - Sign on mode of application.

* `client_id` - The client ID of the application.

* `client_secret` - The client secret of the application.

## Import

An OIDC Application can be imported via the Okta ID.

```
$ terraform import okta_app_oauth.example <app id>
```
