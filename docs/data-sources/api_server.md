---
page_title: "Data Source: okta_api_server"
description: |-
  Retrieves an API server registered in Okta by ID.
---

# Data Source: okta_api_server

Retrieves an API server registered in Okta by ID. An API server is an API resource server that AI agents can be granted access to.

## Example Usage

```hcl
data "okta_api_server" "example" {
  id = "<api-server-id>"
}

output "api_server_scopes" {
  value = data.okta_api_server.example.metadata.scopes_supported
}
```

## Schema

### Required

- `id` (String) ID of the API server.

### Read-Only

- `resource_url` (String) The URL of the API server.
- `orn` (String) The [ORN](https://developer.okta.com/docs/api/openapi/okta-management/guides/roles/#okta-resource-name-orn) of the API server.
- `status` (String) Current status of the API server: `ACTIVE`, `INACTIVE` or `INVALID`.
- `created` (String) Timestamp when the API server was created.
- `last_updated` (String) Timestamp when the API server was last updated.
- `metadata` (Block) Metadata about the API server. See [below](#nested-schema-for-metadata).

<a id="nested-schema-for-metadata"></a>
### Nested Schema for `metadata`

Read-Only:

- `display_name` (String) Human-readable display name for the API server.
- `description` (String) Description of the API server.
- `scopes_supported` (List of String) The OAuth 2.0 scopes that the API server supports.
- `app` (Block) The app instance that represents the API resource server in the Cross App Access (XAA) flow.
  - `orn` (String) The ORN of the app instance.
