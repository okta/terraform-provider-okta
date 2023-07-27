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

### With JWKS value

See also [Advanced PEM secrets and JWKS example](#advanced-pem-and-jwks-example).

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

- `accessibility_error_redirect_url` - (Optional) Custom error page URL.

- `accessibility_login_redirect_url` - (Optional) Custom login page for this application.

- `accessibility_self_service` - (Optional) Enable self-service. By default, it is `false`.

- `authentication_policy` - (Optional) The ID of the associated `app_signon_policy`. If this property is removed from the application the `default` sign-on-policy will be associated with this application.

- `admin_note` - (Optional) Application notes for admins.

- `app_links_json` - (Optional) Displays specific appLinks for the app. The value for each application link should be boolean.

- `app_settings_json` - (Optional) Application settings in JSON format.

- `auto_key_rotation` - (Optional) Requested key rotation mode.  If
    `auto_key_rotation` isn't specified, the client automatically opts in for Okta's
    key rotation. You can update this property via the API or via the administrator
    UI.
    See: https://developer.okta.com/docs/reference/api/apps/#oauth-credential-object

- `auto_submit_toolbar` - (Optional) Display auto submit toolbar.

- `client_basic_secret` - (Optional) OAuth client secret key, this can be set when `token_endpoint_auth_method` is `"client_secret_basic"`.

- `client_id` - (Optional) OAuth client ID. If set during creation, app is created with this id. See: https://developer.okta.com/docs/reference/api/apps/#oauth-credential-object

- `client_uri` - (Optional) URI to a web page providing information about the client.

- `consent_method` - (Optional) Indicates whether user consent is required or implicit. Valid values: `"REQUIRED"`, `"TRUSTED"`. Default value is `"TRUSTED"`.

- `custom_client_id` - (Optional) This property allows you to set your client_id during creation. NOTE: updating after creation will be a no-op, use client_id for that behavior instead.
  - `DEPRECATED`: This field is being replaced by `client_id`. Please use that field instead.",

- `enduser_note` - (Optional) Application notes for end users.

- `grant_types` - (Optional) List of OAuth 2.0 grant types. Conditional validation params found [here](https://developer.okta.com/docs/api/resources/apps#credentials-settings-details).
  Defaults to minimum requirements per app type. Valid values: `"authorization_code"`, `"implicit"`, `"password"`, `"refresh_token"`, `"client_credentials"`,
  `"urn:ietf:params:oauth:grant-type:saml2-bearer"` (*Early Access Property*), `"urn:ietf:params:oauth:grant-type:token-exchange"` (*Early Access Property*),
  `"interaction_code"` (*OIE only*).

- `groups_claim` - (Optional) Groups claim for an OpenID Connect client application. **IMPORTANT**: this field is available only when using api token in the provider config.
  - `type` - (Required) Groups claim type. Valid values: `"FILTER"`, `"EXPRESSION"`.
  - `filter_type` - (Optional) Groups claim filter. Can only be set if type is `"FILTER"`. Valid values: `"EQUALS"`, `"STARTS_WITH"`, `"CONTAINS"`, `"REGEX"`.
  - `name` - (Required) Name of the claim that will be used in the token.
  - `value` - (Required) Value of the claim. Can be an Okta Expression Language statement that evaluates at the time the token is minted.
  - `issuer_mode` - (Read-Only) Issuer Mode is inherited from the Issuer Mode on the OAuth app itself.

- `hide_ios` - (Optional) Do not display application icon on mobile app.

- `hide_web` - (Optional) Do not display application icon to users.

- `implicit_assignment` - (Optional) *Early Access Property*. Enables [Federation Broker Mode](https://help.okta.com/en/prod/Content/Topics/Apps/apps-fbm-enable.htm). When this mode is enabled, `users` and `groups` arguments are ignored.

- `issuer_mode` - (Optional) Indicates whether the Okta Authorization Server uses the original Okta org domain URL or a custom domain URL as the issuer of ID token for this client.
Valid values: `"CUSTOM_URL"`,`"ORG_URL"` or `"DYNAMIC"`. Default is `"ORG_URL"`.

- `jwks` - (Optional) JSON Web Key set. [Admin Console JWK Reference](https://developer.okta.com/docs/guides/implement-oauth-for-okta-serviceapp/main/#generate-the-jwk-in-the-admin-console)

- `jwks_uri` - (Optional) URL of the custom authorization server's JSON Web Key Set document.

- `label` - (Required) The Application's display name.

- `login_mode` - (Optional) The type of Idp-Initiated login that the client supports, if any. Valid values: `"DISABLED"`, `"SPEC"`, `"OKTA"`. Default is `"DISABLED"`.

- `login_scopes` - (Optional) List of scopes to use for the request. Valid values: `"openid"`, `"profile"`, `"email"`, `"address"`, `"phone"`. Required when `login_mode` is NOT `DISABLED`.

- `login_uri` - (Optional) URI that initiates login. Required when `login_mode` is NOT `DISABLED`.

- `logo` - (Optional) Local file path to the logo. The file must be in PNG, JPG, or GIF format, and less than 1 MB in size.

- `logo_uri` - (Optional) URI that references a logo for the client.

- `omit_secret` - (Optional) This tells the provider not to persist the application's secret to state. Your app's `client_secret` will be recreated if this ever changes from true => false.

- `pkce_required` - (Optional) Require Proof Key for Code Exchange (PKCE) for
    additional verification.  If `pkce_required` isn't specified when adding a new
    application, Okta sets it to `true` by default for `"browser"` and `"native"`
    application types.
    See https://developer.okta.com/docs/reference/api/apps/#oauth-credential-object

- `policy_uri` - (Optional) URI to web page providing client policy document.

- `post_logout_redirect_uris` - (Optional) List of URIs for redirection after logout.

- `profile` - (Optional) Custom JSON that represents an OAuth application's profile.

- `redirect_uris` - (Optional) List of URIs for use in the redirect-based flow. This is required for all application types except service.

- `refresh_token_leeway` - (Optional) Grace period for token rotation. Valid values: 0 to 60 seconds.

- `refresh_token_rotation` - (Optional) Refresh token rotation behavior. Valid values: `"STATIC"` or `"ROTATE"`.

- `response_types` - (Optional) List of OAuth 2.0 response type strings. Array
    values of `"code"`, `"token"`, `"id_token"`. The `grant_types` and `response_types`
    values described are partially orthogonal, as they refer to arguments
    passed to different endpoints in the OAuth 2.0 protocol (opens new window).
    However, they are related in that the `grant_types` available to a client
    influence the `response_types` that the client is allowed to use, and vice versa.
    For instance, a grant_types value that includes authorization_code implies a
    `response_types` value that includes code, as both values are defined as part of
    the OAuth 2.0 authorization code grant.
    See: https://developer.okta.com/docs/reference/api/apps/#add-oauth-2-0-client-application

- `status` - (Optional) The status of the application, by default, it is `"ACTIVE"`.

- `token_endpoint_auth_method` - (Optional) Requested authentication method for
    the token endpoint. It can be set to `"none"`, `"client_secret_post"`,
    `"client_secret_basic"`, `"client_secret_jwt"`, `"private_key_jwt"`.  Use
    `pkce_required` to require PKCE for your confidential clients using the
    Authorization Code flow. If `"token_endpoint_auth_method"` is `"none"`,
    `pkce_required` needs to be `true`. If `pkce_required` isn't specified when
    adding a new application, Okta sets it to `true` by default for `"browser"` and
    `"native"` application types.
    See https://developer.okta.com/docs/reference/api/apps/#oauth-credential-object

- `tos_uri` - (Optional) URI to web page providing client tos (terms of service).

- `type` - (Required) The type of OAuth application. Valid values: `"web"`, `"native"`, `"browser"`, `"service"`. For SPA apps use `browser`.

- `user_name_template` - (Optional) Username template. Default: `"${source.login}"`

- `user_name_template_push_status` - (Optional) Push username on update. Valid values: `"PUSH"` and `"DONT_PUSH"`.

- `user_name_template_suffix` - (Optional) Username template suffix.

- `user_name_template_type` - (Optional) Username template type. Default: `"BUILT_IN"`.

- `wildcard_redirect` - (Optional) *Early Access Property*. Indicates if the client is allowed to use wildcard matching of `redirect_uris`. Valid values: `"DISABLED"`, `"SUBDOMAIN"`. Default value is `"DISABLED"`.

## Attributes Reference

- `client_id` - The client ID of the application.

- `client_secret` - The client secret of the application. See: https://developer.okta.com/docs/reference/api/apps/#oauth-credential-object

- `id` - ID of the application.

- `logo_url` - Direct link of application logo.

- `name` - Name assigned to the application by Okta.

- `sign_on_mode` - Sign-on mode of application.

## Timeouts

The `timeouts` block allows you to specify custom [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

- `create` - Create timeout (default 1 hour).

- `update` - Update timeout (default 1 hour).

- `read` - Read timeout (default 1 hour).

## Import

An OIDC Application can be imported via the Okta ID.

```
$ terraform import okta_app_oauth.example &#60;app id&#62;
```

## Etc.

### Resetting client secret

If the client secret needs to be reset run an apply with `omit_secret` set to
true in the resource. This causes `client_secret` to be set to blank. Remove
`omit_secret` and run apply again. The resource will set a new `client_secret`
for the app.

### Advanced PEM and JWKS example


```hcl
# This example config illustrates how Terraform can be used to generate a
# private keys PEM file and valid JWKS to be used as part of an Okta OAuth
# app's creation.
#
terraform {
  required_providers {
    okta = {
      source = "okta/okta"
    }
    tls = {
      source = "hashicorp/tls"
    }
    jwks = {
      source = "iwarapter/jwks"
    }
  }
}

# NOTE: Example to generate a PEM easily as a tool. These secrets will be saved
# to the state file and shouldn't be persisted. Instead, save the secrets into
# a secrets manager to be reused.
# https://registry.terraform.io/providers/hashicorp/tls/latest/docs/resources/private_key
#
# NOTE: Even though tls is a Hashicorp provider you should still audit its code
# to be satisfied with its security.
# https://github.com/hashicorp/terraform-provider-tls
#
resource "tls_private_key" "rsa" {
  algorithm = "RSA"
  rsa_bits  = 4096
}
#
# Pretty print a PEM with TF show and jq
# terraform show -json | jq -r '.values.root_module.resources[] | select(.address == "tls_private_key.rsa").values.private_key_pem'
#
# Delete the secrets explicitly or just remove them from the config and run
# apply again.
# terraform apply -destroy -auto-approve -target=tls_private_key.rsa

# NOTE: Even though the iwarapter/jwks is listed in the registry you should
# still audit its code to be satisfied with its security.
# https://registry.terraform.io/providers/iwarapter/jwks/latest/docs/data-sources/from_key
# https://github.com/iwarapter/terraform-provider-jwks
#
data "jwks_from_key" "jwks" {
  key = tls_private_key.rsa.private_key_pem
  kid = "my-kid"
}
#
# Pretty print the jwks
# terraform show -json | jq -r '.values.root_module.resources[] | select(.address == "data.jwks_from_key.jwks").values.jwks' | jq .

# Feed values into Okta OAuth app's jwks
locals {
  jwks = jsondecode(data.jwks_from_key.jwks.jwks)
}

# https://registry.terraform.io/providers/okta/okta/latest/docs/resources/app_oauth
resource "okta_app_oauth" "app" {
  label                      = "My OAuth App"
  type                       = "service"
  response_types             = ["token"]
  grant_types                = ["client_credentials"]
  token_endpoint_auth_method = "private_key_jwt"

  jwks {
    kty = local.jwks.kty
    kid = local.jwks.kid
    e   = local.jwks.e
    n   = local.jwks.n
  }
}
#
# Pretty print OAuth app Client ID
# terraform show -json | jq -r '.values.root_module.resources[] | select(.address == "okta_app_oauth.app").values.id'
```
