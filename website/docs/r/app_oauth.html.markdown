---
layout: 'okta'
page_title: 'Okta: okta_app_oauth'
sidebar_current: 'docs-okta-resource-app-oauth'
description: |-
  Creates an OIDC Application.
---

# okta_app_oauth

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

```hcl
resource "okta_app_oauth" "example" {
  label          = "example"
  type           = "service"
  response_types = ["token"]
  grant_types    = ["client_credentials"]
  token_endpoint_auth_method = "private_key_jwt"

  jwks {
    kty = "RSA"
    kid = "SIGNING_KEY"
    e   = "AQAB"
    n   = "xyz"
  }
}
```

## Argument Reference

The following arguments are supported:

- `label` - (Required) The Application's display name.

- `status` - (Optional) The status of the application, by default, it is `"ACTIVE"`.

- `type` - (Required) The type of OAuth application. Valid values: `"web"`, `"native"`, `"browser"`, `"service"`.

- `users` - (Optional) The users assigned to the application. It is recommended not to use this and instead use `okta_app_user`.
  - `DEPRECATED`: Please replace usage with the `okta_app_user` resource.

- `groups` - (Optional) The groups assigned to the application. It is recommended not to use this and instead use `okta_app_group_assignment`.
  - `DEPRECATED`: Please replace usage with the `okta_app_group_assignments` (or `okta_app_group_assignment`) resource.

- `client_id` - (Optional) OAuth client ID. If set during creation, app is created with this id.

- `omit_secret` - (Optional) This tells the provider not to persist the application's secret to state. Your app will be recreated if this ever changes from true => false.

- `client_basic_secret` - (Optional) OAuth client secret key, this can be set when token_endpoint_auth_method is client_secret_basic.

- `token_endpoint_auth_method` - (Optional) Requested authentication method for the token endpoint. It can be set to `"none"`, `"client_secret_post"`, `"client_secret_basic"`, `"client_secret_jwt"`, `"private_key_jwt"`. To enable PKCE, set this to `"none"`.

- `auto_key_rotation` - (Optional) Requested key rotation mode.

- `client_uri` - (Optional) URI to a web page providing information about the client.

- `logo_uri` - (Optional) URI that references a logo for the client.

- `login_uri` - (Optional) URI that initiates login. Required when `login_mode` is NOT `DISABLED`.

- `redirect_uris` - (Optional) List of URIs for use in the redirect-based flow. This is required for all application types except service.

- `wildcard_redirect` - (Optional) *Early Access Property*. Indicates if the client is allowed to use wildcard matching of `redirect_uris`. Valid values: `"DISABLED"`, `"SUBDOMAIN"`. Default value is `"DISABLED"`.

- `post_logout_redirect_uris` - (Optional) List of URIs for redirection after logout.

- `response_types` - (Optional) List of OAuth 2.0 response type strings.

- `grant_types` - (Optional) List of OAuth 2.0 grant types. Conditional validation params found [here](https://developer.okta.com/docs/api/resources/apps#credentials-settings-details). 
  Defaults to minimum requirements per app type. Valid values: `"authorization_code"`, `"implicit"`, `"password"`, `"refresh_token"`, `"client_credentials"`, 
  `"urn:ietf:params:oauth:grant-type:saml2-bearer"` (*Early Access Property*), `"urn:ietf:params:oauth:grant-type:token-exchange"` (*Early Access Property*)

- `tos_uri` - (Optional) URI to web page providing client tos (terms of service).

- `policy_uri` - (Optional) URI to web page providing client policy document.

- `consent_method` - (Optional) Indicates whether user consent is required or implicit. Valid values: `"REQUIRED"`, `"TRUSTED"`. Default value is `"TRUSTED"`.

- `issuer_mode` - (Optional) Indicates whether the Okta Authorization Server uses the original Okta org domain URL or a custom domain URL as the issuer of ID token for this client.

- `refresh_token_rotation` - (Optional) Refresh token rotation behavior. Valid values: `"STATIC"` or `"ROTATE"`.

- `refresh_token_leeway` - (Optional) Grace period for token rotation. Valid values: 0 to 60 seconds.

- `auto_submit_toolbar` - (Optional) Display auto submit toolbar.

- `accessibility_error_redirect_url` - (Optional) Custom error page URL.

- `accessibility_login_redirect_url` - (Optional) Custom login page for this application.

- `accessibility_self_service` - (Optional) Enable self-service. By default, it is `false`.

- `user_name_template` - (Optional) Username template. Default: `"${source.login}"`

- `user_name_template_suffix` - (Optional) Username template suffix.

- `user_name_template_type` - (Optional) Username template type. Default: `"BUILT_IN"`.

- `user_name_template_push_status` - (Optional) Push username on update. Valid values: `"PUSH"` and `"DONT_PUSH"`.

- `hide_ios` - (Optional) Do not display application icon on mobile app.

- `hide_web` - (Optional) Do not display application icon to users.

- `profile` - (Optional) Custom JSON that represents an OAuth application's profile.

- `implicit_assignment` - (Optional) *Early Access Property*. Enables [Federation Broker Mode]( https://help.okta.com/en/prod/Content/Topics/Apps/apps-fbm-enable.htm). When this mode is enabled, `users` and `groups` arguments are ignored.

- `login_mode` - (Optional) The type of Idp-Initiated login that the client supports, if any. Valid values: `"DISABLED"`, `"SPEC"`, `"OKTA"`. Default is `"DISABLED"`.

- `login_scopes` - (Optional) List of scopes to use for the request. Valid values: `"openid"`, `"profile"`, `"email"`, `"address"`, `"phone"`. Required when `login_mode` is NOT `DISABLED`.

- `logo` - (Optional) Local file path to the logo. The file must be in PNG, JPG, or GIF format, and less than 1 MB in size.

- `groups_claim` - (Optional) Groups claim for an OpenID Connect client application. **IMPORTANT**: this field is available only when using api token in the provider config.
  - `type` - (Required) Groups claim type. Valid values: `"FILTER"`, `"EXPRESSION"`.
  - `filter_type` - (Optional) Groups claim filter. Can only be set if type is `"FILTER"`. Valid values: `"EQUALS"`, `"STARTS_WITH"`, `"CONTAINS"`, `"REGEX"`.
  - `name` - (Required) Name of the claim that will be used in the token.
  - `value` - (Required) Value of the claim. Can be an Okta Expression Language statement that evaluates at the time the token is minted.

- `admin_note` - (Optional) Application notes for admins.

- `enduser_note` - (Optional) Application notes for end users.

- `app_settings_json` - (Optional) Application settings in JSON format.

- `skip_users` - (Optional) Indicator that allows the app to skip `users` sync (it's also can be provided during import). Default is `false`.

- `skip_groups` - (Optional) Indicator that allows the app to skip `groups` sync (it's also can be provided during import). Default is `false`.

## Attributes Reference

- `id` - ID of the application.

- `name` - Name assigned to the application by Okta.

- `sign_on_mode` - Sign-on mode of application.

- `client_id` - The client ID of the application.

- `client_secret` - The client secret of the application.

- `logo_url` - Direct link of application logo.

## Import

An OIDC Application can be imported via the Okta ID.

```
$ terraform import okta_app_oauth.example <app id>
```

It's also possible to import app without groups or/and users. In this case ID may look like this:

```
$ terraform import okta_app_basic_auth.example <app id>/skip_users

$ terraform import okta_app_basic_auth.example <app id>/skip_users/skip_groups

$ terraform import okta_app_basic_auth.example <app id>/skip_groups
```
