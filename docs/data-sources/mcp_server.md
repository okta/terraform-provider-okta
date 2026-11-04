---
page_title: "Data Source: okta_mcp_server"
description: |-
  Retrieves a third-party Model Context Protocol (MCP) server by ID.
---

# Data Source: okta_mcp_server

Retrieves a third-party Model Context Protocol (MCP) server by ID.

## Example Usage

```hcl
data "okta_mcp_server" "example" {
  id = "<mcp-server-id>"
}
```

## Schema

### Required

- `id` (String) Unique identifier of the MCP server.

### Read-Only

- `authorization_server_count` (Number) Number of authorization servers linked to this MCP server.
- `created` (String) Timestamp when the MCP server was created.
- `last_updated` (String) Timestamp when the MCP server was last updated.
- `org_config_property_values` (List of String) Values for the org-configured properties that the OIN catalog entry defines. Each element is a JSON-encoded object.
- `orn` (String) The [ORN](https://developer.okta.com/docs/reference/okta-resource-names/) of the MCP server.
- `resource_server_global_id` (String) ID of the resource server in the Okta Integration Network (OIN) catalog that the MCP server was registered from.
- `resource_url` (String) The URL of the MCP server resource.
- `status` (String) Current status of the MCP server in its lifecycle.
- `metadata` (Block) Metadata about the MCP server. See [below](#nested-schema-for-metadata).
- `oauth_client` (Block) The MCP server's own OAuth client identity. See [below](#nested-schema-for-oauth_client).

<a id="nested-schema-for-metadata"></a>
### Nested Schema for `metadata`

Read-Only:

- `description` (String) Description of the MCP server.
- `display_name` (String) Human-readable display name for the MCP server.
- `last_refreshed_at` (String) Timestamp when the metadata was last detected from the MCP server.
- `resource_name` (String) Canonical resource name of the MCP server as reported by the server itself (auto-detected).
- `scopes_supported` (List of String) List of OAuth scopes supported by the MCP server.

<a id="nested-schema-for-oauth_client"></a>
### Nested Schema for `oauth_client`

Read-Only:

- `client_id` (String) The OAuth client ID that Okta uses when issuing downstream tokens for this MCP server.
- `connection_count` (Number) Number of resource connections defined on the MCP server's client identity.
- `grant_types` (List of String) OAuth 2.0 grant types for the client identity.
- `orn` (String) The [ORN](https://developer.okta.com/docs/reference/okta-resource-names/) of the client identity.
- `status` (String) Status of the MCP server's client identity.
- `token_endpoint_auth_method` (String) The client authentication method that the MCP server uses at the Okta token endpoint.
