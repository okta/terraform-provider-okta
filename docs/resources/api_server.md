---
page_title: "Resource: okta_api_server"
description: |-
  Registers a third-party API server as an Okta resource server.
---

# Resource: okta_api_server

Registers a third-party API server as an Okta resource server. The API Server Registrations V2 API provides registration and management operations for API servers. The provider activates the API server after it's created, and deactivates it before deleting it.

Changes to `metadata` (`display_name`, `description`, and `scopes_supported`) are applied in place. The update API doesn't accept `resource_url`, `registration_type`, or `metadata.app`, so don't change them after creation. To change them, recreate the resource (for example, with `terraform apply -replace`).

## Example Usage

```hcl
resource "okta_api_server" "example" {
  resource_url = "https://api.example.com"

  metadata {
    display_name     = "Example API"
    description      = "Internal inventory API"
    scopes_supported = ["read:items", "write:items"]
  }
}
```

To link an authorization server to the API server, use [`okta_api_servers_authorization_server`](api_servers_authorization_server.md).

## Schema

### Required

- `resource_url` (String) The URL of the API server. Can't be changed after creation.

### Optional

- `registration_type` (String, Sensitive) Specify how the API server's authorization server metadata is sourced. Only sent on create.
- `metadata` (Block, Optional) Metadata for the API server.
  - `display_name` (String) Human-readable display name for the API server.
  - `description` (String) Description of the API server.
  - `scopes_supported` (List of String) The OAuth 2.0 scopes that the API server supports. Provide these when you manually register an API server.
  - `app` (Block, Optional) The app instance that represents the API resource server in the Cross App Access (XAA) flow. Only sent on create.
    - `orn` (String) The [ORN](https://developer.okta.com/docs/api/openapi/okta-management/guides/roles/#okta-resource-name-orn) of the app instance.

### Read-Only

- `id` (String) The unique identifier of the API server.
- `orn` (String) The [ORN](https://developer.okta.com/docs/api/openapi/okta-management/guides/roles/#okta-resource-name-orn) of the API server.
- `status` (String) Current status of the API server in its lifecycle.

## Import

Import using the API server ID:

```shell
terraform import okta_api_server.example <api-server-id>
```

> **Note**: `registration_type` is not returned by the read endpoint, so it is not restored on import. Set it in configuration after importing.
