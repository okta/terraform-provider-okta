---
layout: 'okta'
page_title: 'Okta: okta_auth_server_scope'
sidebar_current: 'docs-okta-resource-auth-server-scope'
description: |-
  Creates an Authorization Server Scope.
---

# okta_auth_server_scope

Creates an Authorization Server Scope.

This resource allows you to create and configure an Authorization Server Scope.

## Example Usage

```hcl
resource "okta_auth_server_scope" "example" {
  auth_server_id   = "<auth server id>"
  metadata_publish = "NO_CLIENTS"
  name             = "example"
  consent          = "IMPLICIT"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) Auth Server scope name.

- `auth_server_id` - (Required) Auth Server ID.

- `description` - (Optional) Description of the Auth Server Scope.

- `consent` - (Optional) Indicates whether a consent dialog is needed for the scope. It can be set to `"REQUIRED"` or `"IMPLICIT"`.

- `metadata_publish` - (Optional) Whether to publish metadata or not. It can be set to `"ALL_CLIENTS"` or `"NO_CLIENTS"`.

- `default` - (Optional) A default scope will be returned in an access token when the client omits the scope parameter in a token request, provided this scope is allowed as part of the access policy rule.

## Attributes Reference

- `id` - ID of the Auth Server Scope.

- `auth_server_id` - The ID of the Auth Server.

## Import

Okta Auth Server Scope can be imported via the Auth Server ID and Scope ID.

```
$ terraform import okta_auth_server_scope.example <auth server id>/<scope id>
```
