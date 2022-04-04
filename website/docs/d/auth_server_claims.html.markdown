---
layout: 'okta'
page_title: 'Okta: okta_auth_server_claims'
sidebar_current: 'docs-okta-datasource-auth-server-claims'
description: |-
  Get a list of authorization server claims from Okta.
---

# okta_auth_server_claims

Use this data source to retrieve a list of authorization server claims from Okta.

## Example Usage

```hcl
data "okta_auth_server_claims" "test" {
  auth_server_id = "default"
}
```

## Arguments Reference

- `auth_server_id` - (Required) Auth server ID.

## Attributes Reference

- `claims` - collection of authorization server claims retrieved from Okta with the following properties.

    - `id` - ID of the claim.

    - `name` - Name of the claim.
    
    - `scopes` - Specifies the scopes for this Claim.
    
    - `status` - Status of the claim.
    
    - `value` - Value of the claim
    
    - `value_type` - Specifies whether the Claim is an Okta EL expression (`"EXPRESSION"`), a set of groups (`"GROUPS"`), or a system claim (`"SYSTEM"`)
    
    - `claim_type` - Specifies whether the Claim is for an access token (`"RESOURCE"`) or ID token (`"IDENTITY"`).
    
    - `always_include_in_token` - Specifies whether to include Claims in the token.
