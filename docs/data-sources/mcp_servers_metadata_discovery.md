---
page_title: "Data Source: okta_mcp_servers_metadata_discovery"
description: |-
  Retrieves discovery metadata for a Model Context Protocol server
---

# Data Source: okta_mcp_servers_metadata_discovery

Retrieves the discovery metadata for a third-party Model Context Protocol (MCP) server. This includes information about available resources, tools, and capabilities that the MCP server exposes.

## Example Usage

```terraform
data "okta_mcp_servers_metadata_discovery" "example" {
  id            = "meta_abc123def456"
  mcp_server_id = "mcp_xyz789"
}

output "mcp_metadata" {
  value = {
    id                 = data.okta_mcp_servers_metadata_discovery.example.id
    status             = data.okta_mcp_servers_metadata_discovery.example.status
    last_discovered_at = data.okta_mcp_servers_metadata_discovery.example.last_discovered_at
    tools_count        = data.okta_mcp_servers_metadata_discovery.example.tools_count
  }
}
```

## Argument Reference

### Required

- `id` - (String) Unique identifier of the metadata discovery resource
- `mcp_server_id` - (String) ID of the MCP server

## Attribute Reference

- `error_reason` - (String) Reason for the most recent discovery failure, if any
- `last_discovered_at` - (String) Timestamp when discovery last completed successfully (RFC3339 format)
- `status` - (String) Current status of the metadata discovery (e.g., `PENDING`, `IN_PROGRESS`, `COMPLETED`, `FAILED`)
- `tools_count` - (Int64) Number of tools found by the most recent successful discovery
- `credential` - (Block) The credential context that produced this discovery
  - `auth_type` - (String) Authentication type used for discovery
  - `client_auth_settings_id` - (String) ID of the resource server client auth settings
  - `scopes` - (List of String) The scope set used at the resource authorization server

## Example: Inspect MCP Server Discovery Status

```terraform
data "okta_mcp_servers_metadata_discovery" "github_mcp" {
  id            = "meta_github_discovery"
  mcp_server_id = "mcp_github_server"
}

output "github_mcp_capabilities" {
  value = {
    status             = data.okta_mcp_servers_metadata_discovery.github_mcp.status
    last_discovered_at = data.okta_mcp_servers_metadata_discovery.github_mcp.last_discovered_at
    tools_count        = data.okta_mcp_servers_metadata_discovery.github_mcp.tools_count
    error_reason       = try(data.okta_mcp_servers_metadata_discovery.github_mcp.error_reason, "")
  }
}
```

## Example: Check Discovery Status and Retry

```terraform
variable "mcp_server_id" {
  type = string
}

data "okta_mcp_servers_metadata_discovery" "check" {
  id            = "meta_discovery"
  mcp_server_id = var.mcp_server_id
}

locals {
  discovery_complete = data.okta_mcp_servers_metadata_discovery.check.status == "COMPLETED"
  discovery_failed   = data.okta_mcp_servers_metadata_discovery.check.status == "FAILED"
  has_tools          = data.okta_mcp_servers_metadata_discovery.check.tools_count > 0
}

output "discovery_status" {
  value = {
    id           = data.okta_mcp_servers_metadata_discovery.check.id
    complete     = local.discovery_complete
    failed       = local.discovery_failed
    has_tools    = local.has_tools
    error_reason = try(data.okta_mcp_servers_metadata_discovery.check.error_reason, "none")
  }
}
```

## Notes

- This data source retrieves read-only metadata about MCP server discovery results
- Metadata discovery is typically performed when a server is registered or manually triggered
- The `status` field indicates the current discovery state
- `tools_count` indicates how many tools were discovered in the most recent successful discovery
- If discovery fails, `error_reason` provides diagnostic information about the failure
- Use this data source to verify that an MCP server's capabilities have been properly discovered