---
layout: 'okta'
page_title: 'Okta: okta_auth_server_policy'
sidebar_current: 'docs-okta-datasource-auth-server-policy'
description: |-
  Get an authorization server policy from Okta.
---

# okta_auth_server_policy

Use this data source to retrieve a authorization server policy from Okta.

## Example Usage

```hcl
data "okta_auth_server_policy" "example" {
  auth_server_id = "<auth server id>"
  name           = "staff"
}
```

## Arguments Reference

- `auth_server_id` - (Required) The ID of the Auth Server.

- `name` - (Required) Name of policy to retrieve.

## Attributes Reference

- `id` - id of authorization server policy.

- `description` - description of authorization server policy.

- `assigned_clients` - list of clients this policy is assigned to. `["ALL_CLIENTS"]` is a special value when policy is assigned to all clients.



