---
page_title: "Resource: okta_mcp_server"
description: |-
  Manages an Okta MCP server registration.
---

# Resource: okta_mcp_server

Manages a third-party Model Context Protocol (MCP) server registration. The MCP Server Registrations V2 API provides lifecycle management for MCP servers.

## Example Usage

```hcl
resource "okta_mcp_server" "example" {
  registration_type = "DISCOVERABLE"
  resource_url      = "https://example.com"

  metadata {
    display_name = "Example MCP Server"
    description  = "Example MCP server registration"
  }
}
```

### OIN-sourced MCP server

```hcl
resource "okta_mcp_server" "oin" {
  registration_type         = "OIN"
  app_name                  = "example_app"
  resource_server_global_id = "example-resource-server-global-id"

  org_config_property_values {
    name  = "tenant"
    value = "example-tenant"
  }
}
```

## Schema

### Optional

- `registration_type` (String, Sensitive) How the MCP server's metadata is sourced. Changing this forces replacement. `DISCOVERABLE` and `MANUAL` require `metadata` and `resource_url`. `OIN` requires `app_name` and `resource_server_global_id`.
- `app_name` (String, Sensitive) Name of the OIN app that the resource server catalog entry belongs to. Required when `registration_type` is `OIN`.
- `resource_server_global_id` (String) ID of the resource server in the Okta Integration Network (OIN) catalog that the MCP server was registered from. Required when `registration_type` is `OIN`.
- `resource_url` (String) The URL of the MCP server resource. Required when `registration_type` is `DISCOVERABLE` or `MANUAL`.
- `metadata` (Block, Optional) Metadata for the MCP server. See [below](#nested-schema-for-metadata).
- `oauth_client` (Block, Optional) The MCP server's own OAuth client identity. See [below](#nested-schema-for-oauth_client).
- `org_config_property_values` (Block List) Values for the org-configured properties that the OIN catalog entry defines. See [below](#nested-schema-for-org_config_property_values).

### Read-Only

- `id` (String) The unique identifier for the resource.
- `authorization_server_count` (Number) Number of authorization servers linked to this MCP server.
- `orn` (String) The [ORN](https://developer.okta.com/docs/reference/okta-resource-names/) of the MCP server.
- `status` (String) Current status of the MCP server in its lifecycle.

<a id="nested-schema-for-metadata"></a>
### Nested Schema for `metadata`

Optional:

- `description` (String) Description of the MCP server.
- `display_name` (String) Human-readable display name for the MCP server.
- `scopes_supported` (List of String) OAuth scopes supported by the MCP server.

<a id="nested-schema-for-oauth_client"></a>
### Nested Schema for `oauth_client`

Read-Only:

- `client_id` (String) The OAuth client ID that Okta uses when issuing downstream tokens for this MCP server.
- `connection_count` (Number) Number of resource connections defined on the MCP server's client identity.
- `grant_types` (List of String) OAuth 2.0 grant types for the client identity.
- `orn` (String) The [ORN](https://developer.okta.com/docs/reference/okta-resource-names/) of the client identity.
- `status` (String) Status of the MCP server's client identity.
- `token_endpoint_auth_method` (String) The client authentication method that the MCP server uses at the Okta token endpoint.

<a id="nested-schema-for-org_config_property_values"></a>
### Nested Schema for `org_config_property_values`

Optional:

- `name` (String) Name of the org-configured property, as defined by the OIN catalog entry.
- `value` (String) Value for the org-configured property.

## Import

Import using the MCP server ID:

```shell
terraform import okta_mcp_server.example <mcp-server-id>
```

> **Note**: `registration_type` and `app_name` are request-only inputs and are not returned by the API, so they are not restored on import. Set them in configuration after importing.
