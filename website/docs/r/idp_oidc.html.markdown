---
layout: 'okta'
page_title: 'Okta: okta_idp_oidc'
sidebar_current: 'docs-okta-resource-idp-oidc'
description: |-
  Creates an OIDC Identity Provider.
---

# okta_idp_oidc

Creates an OIDC Identity Provider.

This resource allows you to create and configure an OIDC Identity Provider.

## Example Usage

```hcl
resource "okta_idp_oidc" "example" {
  name                  = "example"
  acs_type              = "INSTANCE"
  acs_binding           = "HTTP-POST"
  authorization_url     = "https://idp.example.com/authorize"
  authorization_binding = "HTTP-REDIRECT"
  token_url             = "https://idp.example.com/token"
  token_binding         = "HTTP-POST"
  user_info_url         = "https://idp.example.com/userinfo"
  user_info_binding     = "HTTP-REDIRECT"
  jwks_url              = "https://idp.example.com/keys"
  jwks_binding          = "HTTP-REDIRECT"
  scopes                = ["openid"]
  client_id             = "efg456"
  client_secret         = "efg456"
  issuer_url            = "https://id.example.com"
  username_template     = "idpuser.email"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The Application's display name.

- `scopes` - (Required) The scopes of the IdP.

- `authorization_url` - (Required) IdP Authorization Server (AS) endpoint to request consent from the user and obtain an authorization code grant.

- `authorization_binding` - (Required) The method of making an authorization request. It can be set to `"HTTP-POST"` or `"HTTP-REDIRECT"`.

- `token_url` - (Required) IdP Authorization Server (AS) endpoint to exchange the authorization code grant for an access token.

- `token_binding` - (Required) The method of making a token request. It can be set to `"HTTP-POST"` or `"HTTP-REDIRECT"`.

- `jwks_url` - (Required) Endpoint where the signer of the keys publishes its keys in a JWK Set.

- `jwks_binding` - (Required) The method of making a request for the OIDC JWKS. It can be set to `"HTTP-POST"` or `"HTTP-REDIRECT"`.

- `acs_binding` - (Required) The method of making an ACS request. It can be set to `"HTTP-POST"` or `"HTTP-REDIRECT"`.

- `client_id` - (Required) Unique identifier issued by AS for the Okta IdP instance.

- `client_secret` - (Required) Client secret issued by AS for the Okta IdP instance.

- `issuer_url` - (Required) URI that identifies the issuer.

- `status` - (Optional) Status of the IdP.

- `user_info_url` - (Optional) Protected resource endpoint that returns claims about the authenticated user.

- `user_info_binding` - (Optional)

- `acs_type` - (Optional) The type of ACS. Default is `"INSTANCE"`.

- `protocol_type` - (Optional) The type of protocol to use. It can be `"OIDC"` or `"OAUTH2"`.

- `issuer_mode` - (Optional) Indicates whether Okta uses the original Okta org domain URL, or a custom domain URL. It can be `"ORG_URL"` or `"CUSTOM_URL"`.

- `max_clock_skew` - (Optional) Maximum allowable clock-skew when processing messages from the IdP.

- `account_link_action` - (Optional) Specifies the account linking action for an IdP user.

- `account_link_group_include` - (Optional) Group memberships to determine link candidates.

- `provisioning_action` - (Optional) Provisioning action for an IdP user during authentication.

- `deprovisioned_action` - (Optional) Action for a previously deprovisioned IdP user during authentication. Can be `"NONE"` or `"REACTIVATE"`.

- `suspended_action` - (Optional) Action for a previously suspended IdP user during authentication. Can be set to `"NONE"` or `"UNSUSPEND"`

- `groups_action` - (Optional) Provisioning action for IdP user's group memberships. It can be `"NONE"`, `"SYNC"`, `"APPEND"`, or `"ASSIGN"`.

- `groups_attribute` - (Optional) IdP user profile attribute name (case-insensitive) for an array value that contains group memberships.

- `groups_assignment` - (Optional) List of Okta Group IDs to add an IdP user as a member with the `"ASSIGN"` `groups_action`.

- `groups_filter` - (Optional) Whitelist of Okta Group identifiers that are allowed for the `"APPEND"` or `"SYNC"` `groups_action`.

- `username_template` - (Optional) Okta EL Expression to generate or transform a unique username for the IdP user.

- `subject_match_type` - (Optional) Determines the Okta user profile attribute match conditions for account linking and authentication of the transformed IdP username. By default it is set to `"USERNAME"`. It can be set to `"USERNAME"`, `"EMAIL"`, `"USERNAME_OR_EMAIL"` or `"CUSTOM_ATTRIBUTE"`.

- `subject_match_attribute` - (Optional) Okta user profile attribute for matching transformed IdP username. Only for matchType `"CUSTOM_ATTRIBUTE"`.

- `profile_master` - (Optional) Determines if the IdP should act as a source of truth for user profile attributes.

## Attributes Reference

- `id` - ID of the IdP.

- `type` - Type of OIDC IdP.

## Import

An OIDC IdP can be imported via the Okta ID.

```
$ terraform import okta_idp_oidc.example <idp id>
```
