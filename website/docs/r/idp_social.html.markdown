---
layout: 'okta'
page_title: 'Okta: okta_idp_social'
sidebar_current: 'docs-okta-resource-idp-social'
description: |-
  Creates an Social Identity Provider.
---

# okta_idp_social

Creates a Social Identity Provider.

This resource allows you to create and configure a Social Identity Provider.

## Example Usage

```hcl
resource "okta_idp_social" "example" {
  type          = "FACEBOOK"
  protocol_type = "OAUTH2"
  name          = "testAcc_facebook_replace_with_uuid"

  scopes = [
    "public_profile",
    "email",
  ]

  client_id         = "abcd123"
  client_secret     = "abcd123"
  username_template = "idpuser.email"
  match_type        = "CUSTOM_ATTRIBUTE"
  match_attribute   = "customfieldId"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The Application's display name.

- `type` - (Required) The type of Social IdP. It can be `"FACEBOOK"`, `"LINKEDIN"`, `"MICROSOFT"`, or `"GOOGLE"`.

- `scopes` - (Required) The scopes of the IdP.

- `authorization_url` - (Optional) IdP Authorization Server (AS) endpoint to request consent from the user and obtain an authorization code grant.

- `authorization_binding` - (Optional) The method of making an authorization request. It can be set to `"HTTP-POST"` or `"HTTP-REDIRECT"`.

- `token_url` - (Optional) IdP Authorization Server (AS) endpoint to exchange the authorization code grant for an access token.

- `token_binding` - (Optional) The method of making a token request. It can be set to `"HTTP-POST"` or `"HTTP-REDIRECT"`.

- `status` - (Optional) Status of the IdP.

- `client_id` - (Optional) Unique identifier issued by AS for the Okta IdP instance.

- `client_secret` - (Optional) Client secret issued by AS for the Okta IdP instance.

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

- `subject_match_type` - (Optional) Determines the Okta user profile attribute match conditions for account linking and authentication of the transformed IdP username. By default, it is set to `"USERNAME"`. It can be set to `"USERNAME"`, `"EMAIL"`, `"USERNAME_OR_EMAIL"` or `"CUSTOM_ATTRIBUTE"`.

- `subject_match_attribute` - (Optional) Okta user profile attribute for matching transformed IdP username. Only for matchType `"CUSTOM_ATTRIBUTE"`.

- `profile_master` - (Optional) Determines if the IdP should act as a source of truth for user profile attributes.

- `request_signature_algorithm` - (Optional) The XML digital signature algorithm used when signing an AuthnRequest message.

- `request_signature_scope` - (Optional) Specifies whether to digitally sign an AuthnRequest messages to the IdP. It can be `"REQUEST"` or `"NONE"`.

- `response_signature_algorithm` - (Optional) The minimum XML digital signature algorithm allowed when verifying a SAMLResponse message or Assertion element.

- `response_signature_scope` - (Optional) Specifies whether to verify a SAMLResponse message or Assertion element XML digital signature. It can be `"RESPONSE"`, `"ASSERTION"`, or `"ANY"`.

## Attributes Reference

- `id` - ID of the IdP.

## Import

A Social IdP can be imported via the Okta ID.

```
$ terraform import okta_idp_social.example <idp id>
```
