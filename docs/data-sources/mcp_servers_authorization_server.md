---
page_title: "Data Source: okta_mcp_servers_authorization_server"
description: |-
  Retrieves the link between an MCP server and an authorization server, by the authorization server's ID.
---

# Data Source: okta_mcp_servers_authorization_server

Retrieves the link between an Okta MCP server and an authorization server, by the authorization server's ID.

## Example Usage

```hcl
data "okta_mcp_servers_authorization_server" "example" {
  id            = "<authorization-server-id>"
  mcp_server_id = "<mcp-server-id>"
}
```

## Schema

### Required

- `id` (String) Unique identifier of the authorization server linked to the MCP server.
- `mcp_server_id` (String) ID of an MCP server.

### Read-Only

- `created` (String) Timestamp when the authorization server was linked to the MCP server.
- `last_updated` (String) Timestamp when the link was last updated.
- `authorization_server` (Block) An authorization server linked to the MCP server. See [below](#nested-schema-for-authorization_server).

<a id="nested-schema-for-authorization_server"></a>
### Nested Schema for `authorization_server`

Read-Only:

- `created` (String) Timestamp when the authorization server was created.
- `id` (String) Unique identifier for the authorization server.
- `issuer` (String) OAuth 2.0 issuer of the authorization server.
- `last_updated` (String) Timestamp when the authorization server was last updated.
- `orn` (String) The [ORN](https://developer.okta.com/docs/reference/okta-resource-names/) of the authorization server.
- `status` (String) Current status of the authorization server in its lifecycle.
- `subject_mapping_expression` (String) Okta Expression Language formula that computes the `sub` claim value for the ID-JAG token minted for this authorization server.
- `type` (String) How the authorization server's metadata is sourced.
- `metadata` (Block) Metadata about the authorization server. See the schema in the generated provider documentation for the full list of fields.
