---
page_title: "Resource: okta_mcp_servers_authorization_server"
description: |-
  Links an authorization server to an Okta MCP (Model Context Protocol) server.
---

# Resource: okta_mcp_servers_authorization_server

Links an authorization server to an Okta MCP server. The MCP Server Registrations V2 API provides registration and lifecycle management operations for third-party MCP servers, including the authorization servers linked to them.

Updates are not supported in place. Changing `mcp_server_id` or `orn` replaces the link.

## Example Usage

```hcl
resource "okta_mcp_servers_authorization_server" "example" {
  mcp_server_id = "<mcp-server-id>"
  orn           = "<authorization-server-orn>"
}
```

### With protocol configuration

```hcl
resource "okta_mcp_servers_authorization_server" "example" {
  mcp_server_id = "<mcp-server-id>"
  orn           = "<authorization-server-orn>"

  protocol_configurations {
    identity_assertion {
      enabled   = true
      tenant_id = "<tenant-id>"
    }

    sts {
      enabled = true
    }
  }
}
```

## Schema

### Required

- `mcp_server_id` (String) ID of an MCP server. Forces replacement if changed.
- `orn` (String) The [ORN](https://developer.okta.com/docs/reference/okta-resource-names/) of the authorization server to link. Forces replacement if changed.

### Optional

- `protocol_configurations` (Block, Optional) Admin-controlled configuration of each protocol for this authorization server on this MCP server.
  - `identity_assertion` (Block, Optional) Configuration of the `IDENTITY_ASSERTION` protocol (cross-app access).
    - `enabled` (Boolean) Whether `IDENTITY_ASSERTION` is enabled for this authorization server on this MCP server.
    - `tenant_id` (String) Identifier of the tenant at the authorization server for this MCP server (for example, a Slack workspace ID).
  - `sts` (Block, Optional) Configuration of the `STS` protocol.
    - `enabled` (Boolean) Whether `STS` is enabled for this authorization server on this MCP server.

### Read-Only

- `id` (String) The unique identifier of the link, which is the authorization server's ID.
- `authorization_server` (Block) An authorization server linked to the MCP server. Attributes:
  - `id` (String) Unique identifier for the authorization server.
  - `issuer` (String) OAuth 2.0 issuer of the authorization server.
  - `orn` (String) The [ORN](https://developer.okta.com/docs/reference/okta-resource-names/) of the authorization server.
  - `status` (String) Current status of the authorization server in its lifecycle.
  - `subject_mapping_expression` (String) Okta Expression Language formula that computes the `sub` claim value for the ID-JAG token minted for this authorization server.
  - `type` (String) How the authorization server's metadata is sourced.
  - `metadata` (Block) Metadata about the authorization server. See the generated provider documentation for the full list of fields.

## Import

Import using `{mcp_server_id}/{id}`:

```shell
terraform import okta_mcp_servers_authorization_server.example <mcp-server-id>/<authorization-server-id>
```

> **Note**: `orn` is a configuration input and is not returned by the read endpoint, so it is not restored on import. Set it in configuration after importing.
