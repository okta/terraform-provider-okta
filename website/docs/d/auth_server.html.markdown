---
layout: 'okta'
page_title: 'Okta: okta_auth_server'
sidebar_current: 'docs-okta-datasource-auth-server'
description: |-
  Get an auth server from Okta.
---

# okta_auth_server

Use this data source to retrieve an auth server from Okta.

## Example Usage

```hcl
data "okta_auth_server" "example" {
  name = "Example Auth"
}
```

## Arguments Reference

- `name` - (Required) The name of the auth server to retrieve.

## Attributes Reference

- `id` - Authorization server id.

- `name` - The name of the auth server.

- `description` - description of Authorization server.

- `audiences` - array of audiences,

- `kid` - auth server key id.

- `credentials_last_rotated` - last time credentials were rotated.

- `credentials_next_rotation` - next time credentials will be rotated

- `credentials_rotation_mode` - mode of credential rotation, auto or manual.

- `status` - the activation status of the authorization server.
