---
layout: 'okta'
page_title: 'Okta: okta_idp_social'
sidebar_current: 'docs-okta-datasource-idp-social'
description: |-
  Get a social IdP from Okta.
---

# okta_idp_social

Use this data source to retrieve a social IdP from Okta, namely `APPLE`, `FACEBOOK`, `LINKEDIN`, `MICROSOFT`, or  `GOOGLE`.

## Example Usage

```hcl
data "okta_idp_social" "example" {
  name = "My Facebook IdP"
}
```

## Arguments Reference

- `name` - (Optional) The name of the social idp to retrieve, conflicts with `id`.

- `id` - (Optional) The id of the social idp to retrieve, conflicts with `name`.

## Attributes Reference
  
- `status` - Status of the IdP.
  
- `account_link_action` - Specifies the account linking action for an IdP user.
  
- `account_link_group_include` - Group memberships to determine link candidates.
  
- `provisioning_action` - Provisioning action for an IdP user during authentication.
  
- `deprovisioned_action` -  Action for a previously deprovisioned IdP user during authentication.
  
- `suspended_action` - Action for a previously suspended IdP user during authentication.
  
- `groups_action` - Provisioning action for IdP user's group memberships.
  
- `groups_attribute` - IdP user profile attribute name for an array value that contains group memberships.
  
- `groups_assignment` - List of Okta Group IDs.
  
- `groups_filter` - Whitelist of Okta Group identifiers.
  
- `username_template` - Okta EL Expression to generate or transform a unique username for the IdP user.
  
- `subject_match_type` - Determines the Okta user profile attribute match conditions for account linking and authentication of the transformed IdP username.
  
- `subject_match_attribute` - Okta user profile attribute for matching transformed IdP username.
  
- `profile_master` - Determines if the IdP should act as a source of truth for user profile attributes.
  
- `authorization_url` - IdP Authorization Server (AS) endpoint to request consent from the user and obtain an authorization code grant.
  
- `authorization_binding` - The method of making an authorization request.
  
- `token_url` - IdP Authorization Server (AS) endpoint to exchange the authorization code grant for an access token.
  
- `token_binding` - The method of making a token request.
  
- `type` - The type of Social IdP.
  
- `scopes` - The scopes of the IdP.
  
- `protocol_type` - The type of protocol to use.
  
- `client_id` - Unique identifier issued by AS for the Okta IdP instance.
  
- `client_secret` - Client secret issued by AS for the Okta IdP instance.
  
- `max_clock_skew` - Maximum allowable clock-skew when processing messages from the IdP.
  
- `issuer_mode` - Indicates whether Okta uses the original Okta org domain URL, or a custom domain URL.
