---
layout: 'okta'
page_title: 'Okta: okta_auth_server_policy'
sidebar_current: 'docs-okta-resource-auth-server-policy'
description: |-
  Creates an Authorization Server Policy.
---

# okta_auth_server_policy

Creates an Authorization Server Policy.

This resource allows you to create and configure an Authorization Server Policy.

## Example Usage

```hcl
resource "okta_auth_server_policy" "example" {
  auth_server_id   = "<auth server id>"
  status           = "ACTIVE"
  name             = "example"
  description      = "example"
  priority         = 1
  client_whitelist = ["ALL_CLIENTS"]
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The name of the Auth Server Policy.

- `auth_server_id` - (Required) The ID of the Auth Server.

- `status` - (Optional) The status of the Auth Server Policy.

- `priority` - (Required) The priority of the Auth Server Policy.

- `description` - (Optional) The description of the Auth Server Policy.

- `client_whitelist` - (Required) The clients to whitelist the policy for. `["ALL_CLIENTS"]` is a special value that can be used to whitelist for all clients. Otherwise it is a list of client ids.

## Attributes Reference

- `id` - (Required) The ID of the authorization server policy.

- `auth_server_id` - (Required) The ID of the Auth Server.

- `type` - The type of the Auth Server Policy.

## Import

Authorization Server Policy can be imported via the Auth Server ID and Policy ID.

```
$ terraform import okta_auth_server_policy.example <auth server id>/<policy id>
```
