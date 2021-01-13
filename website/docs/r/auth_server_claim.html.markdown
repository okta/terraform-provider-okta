---
layout: 'okta'
page_title: 'Okta: okta_auth_server_claim'
sidebar_current: 'docs-okta-resource-auth-server-claim'
description: |-
  Creates an Authorization Server Claim.
---

# okta_auth_server_claim

Creates an Authorization Server Claim.

This resource allows you to create and configure an Authorization Server Claim.

## Example Usage

```hcl
resource "okta_auth_server_claim" "example" {
  auth_server_id = "<auth server id>"
  name           = "staff"
  value          = "String.substringAfter(user.email, \"@\") == \"example.com\""
  scopes         = ["${okta_auth_server_scope.example.name}"]
  claim_type     = "IDENTITY"
}
```

## Argument Reference

The following arguments are supported:

- `auth_server_id` - (Required) The Application's display name.

- `name` - (Required) The name of the claim.

- `value` - (Required) The value of the claim.

- `scopes` - (Optional) The list of scopes the auth server claim is tied to.

- `status` - (Optional) The status of the application. It defaults to `"ACTIVE"`.

- `value_type` - (Optional) The type of value of the claim. It can be set to `"EXPRESSION"` or `"GROUPS"`. It defaults to `"EXPRESSION"`.

- `claim_type` - (Required) Specifies whether the claim is for an access token `"RESOURCE"` or ID token `"IDENTITY"`.

- `always_include_in_token` - (Optional) Specifies whether to include claims in token, by default it is set to `true`.

- `group_filter_type` - (Optional) Specifies the type of group filter if `value_type` is `"GROUPS"`. Can be set to one of the following `"STARTS_WITH"`, `"EQUALS"`, `"CONTAINS"`, `"REGEX"`.

## Attributes Reference

- `id` - The ID for the auth server claim.

- `name` - The name of the claim.

## Import

Authorization Server Claim can be imported via the Auth Server ID and Claim ID.

```
$ terraform import okta_auth_server_claim.example <auth server id>/<claim id>
```
