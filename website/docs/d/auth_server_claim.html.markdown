---
layout: 'okta'
page_title: 'Okta: okta_auth_server_claim'
sidebar_current: 'docs-okta-datasource-auth-server-claim'
description: |-
  Get authorization server claim from Okta.
---

# okta_auth_server_claim

Use this data source to retrieve authorization server claim from Okta.

## Example Usage

```hcl
data "okta_auth_server_claim" "test" {
  auth_server_id = "default"
  name           = "birthdate"
}
```

## Arguments Reference

- `auth_server_id` - (Required) Auth server ID.

- `name` (Optional) Name of the claim. Conflicts with `id`.

- `id` (Optional) ID of the claim. Conflicts with `name`.

## Attributes Reference

- `id` - ID of the claim.

- `name` - Name of the claim.

- `scopes` - Specifies the scopes for this Claim.

- `status` - Status of the claim.

- `value` - Value of the claim

- `value_type` - Specifies whether the Claim is an Okta EL expression (`"EXPRESSION"`), a set of groups (`"GROUPS"`), or a system claim (`"SYSTEM"`)

- `claim_type` - Specifies whether the Claim is for an access token (`"RESOURCE"`) or ID token (`"IDENTITY"`).

- `always_include_in_token` - Specifies whether to include Claims in the token.
