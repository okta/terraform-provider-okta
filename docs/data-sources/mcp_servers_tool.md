---
page_title: "Data Source: okta_mcp_servers_tool"
description: |-
  Retrieves a tool discovered from a Model Context Protocol server
---

# Data Source: okta_mcp_servers_tool

Retrieves a specific tool that has been discovered from a third-party Model Context Protocol (MCP) server. Tools represent capabilities or functions that agents can invoke through the MCP server.

## Example Usage

```terraform
data "okta_mcp_servers_tool" "example" {
  id            = "tool_abc123def456"
  mcp_server_id = "mcp_xyz789"
}

output "tool_details" {
  value = {
    id           = data.okta_mcp_servers_tool.example.id
    name         = data.okta_mcp_servers_tool.example.name
    description  = data.okta_mcp_servers_tool.example.description
    created      = data.okta_mcp_servers_tool.example.created
    last_updated = data.okta_mcp_servers_tool.example.last_updated
  }
}
```

## Argument Reference

### Required

- `id` - (String) Unique identifier of the tool
- `mcp_server_id` - (String) ID of the MCP server

## Attribute Reference

- `created` - (String) Timestamp when the tool was first discovered (RFC3339 format)
- `description` - (String) Human-readable description of what the tool does
- `input_schema` - (Map of String) JSON Schema describing the tool's input parameters
- `last_updated` - (String) Timestamp when the tool was last updated by a discovery run (RFC3339 format)
- `name` - (String) Name of the tool as reported by the MCP server
- `credential` - (Block) The credential context that produced this discovery
  - `auth_type` - (String) Authentication type used for discovery
  - `client_auth_settings_id` - (String) ID of the resource server client auth settings
  - `scopes` - (List of String) The scope set used at the resource authorization server

## Example: Retrieve Tool Information

```terraform
data "okta_mcp_servers_tool" "github_list_repos" {
  id            = "tool_github_list_repos"
  mcp_server_id = "mcp_github"
}

output "github_tool_info" {
  value = {
    name         = data.okta_mcp_servers_tool.github_list_repos.name
    description  = data.okta_mcp_servers_tool.github_list_repos.description
    created      = data.okta_mcp_servers_tool.github_list_repos.created
    last_updated = data.okta_mcp_servers_tool.github_list_repos.last_updated
  }
}
```

## Example: Query Multiple Tools

```terraform
variable "mcp_server_id" {
  type = string
}

variable "tool_ids" {
  type    = list(string)
  default = ["tool_query", "tool_create", "tool_delete"]
}

data "okta_mcp_servers_tool" "tools" {
  for_each = toset(var.tool_ids)
  
  id            = each.value
  mcp_server_id = var.mcp_server_id
}

output "discovered_tools" {
  value = {
    for id, tool in data.okta_mcp_servers_tool.tools :
    id => {
      name        = tool.name
      description = tool.description
      created     = tool.created
    }
  }
}
```

## Example: Inspect Tool Input Schema

```terraform
data "okta_mcp_servers_tool" "create_issue" {
  id            = "tool_create_github_issue"
  mcp_server_id = "mcp_github"
}

output "tool_schema" {
  value = {
    name         = data.okta_mcp_servers_tool.create_issue.name
    description  = data.okta_mcp_servers_tool.create_issue.description
    input_schema = data.okta_mcp_servers_tool.create_issue.input_schema
  }
}
```

## Notes

- This data source is read-only and retrieves metadata about discovered tools
- Tools are discovered from MCP servers during the metadata discovery process
- The `input_schema` provides JSON Schema information about the tool's parameters
- Use this data source to understand what capabilities are available from an MCP server