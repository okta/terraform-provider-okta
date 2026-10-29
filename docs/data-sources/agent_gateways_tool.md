---
page_title: "Data Source: okta_agent_gateways_tool"
description: |-
  Retrieves a tool selected for an agent gateway
---

# Data Source: okta_agent_gateways_tool

Retrieves a specific tool that has been selected for an agent gateway by ID. Tools represent capabilities that agents can access through the gateway.

## Example Usage

```terraform
data "okta_agent_gateways_tool" "example" {
  id               = "tool_abc123def456"
  agent_gateway_id = "agw_xyz789"
}

output "tool_info" {
  value = {
    id     = data.okta_agent_gateways_tool.example.id
    status = data.okta_agent_gateways_tool.example.status
    name   = try(data.okta_agent_gateways_tool.example.profile[0].name, "")
  }
}
```

## Argument Reference

### Required

- `id` - (String) Unique identifier of the tool
- `agent_gateway_id` - (String) ID of the agent gateway

## Attribute Reference

- `connection_orn` - (String) The [ORN](https://developer.okta.com/docs/concepts/orn/) of the MCP server connection
- `source_orn` - (String) The [ORN](https://developer.okta.com/docs/concepts/orn/) of the source tool
- `status` - (String) Validity status of the tool relative to the upstream MCP server
- `profile` - (Block) Gateway tool profile information
  - `alias` - (String) Admin-defined alias that overrides the tool's default name
  - `description` - (String) Description of the tool from the upstream MCP server
  - `name` - (String) Name of the tool as reported by the upstream MCP server

## Example: Retrieve Tool with Alias

```terraform
data "okta_agent_gateways_tool" "github_tool" {
  id               = "tool_github_integration"
  agent_gateway_id = "agw_prod"
}

output "tool_details" {
  value = {
    id           = data.okta_agent_gateways_tool.github_tool.id
    status       = data.okta_agent_gateways_tool.github_tool.status
    display_name = try(data.okta_agent_gateways_tool.github_tool.profile[0].alias, data.okta_agent_gateways_tool.github_tool.profile[0].name)
    description  = try(data.okta_agent_gateways_tool.github_tool.profile[0].description, "")
    source_orn   = data.okta_agent_gateways_tool.github_tool.source_orn
  }
}
```

## Example: Query Multiple Tools

```terraform
variable "gateway_id" {
  type = string
}

variable "tool_ids" {
  type    = list(string)
  default = ["tool_create", "tool_read", "tool_update", "tool_delete"]
}

data "okta_agent_gateways_tool" "gateway_tools" {
  for_each = toset(var.tool_ids)
  
  id               = each.value
  agent_gateway_id = var.gateway_id
}

output "available_gateway_tools" {
  value = {
    for id, tool in data.okta_agent_gateways_tool.gateway_tools :
    id => {
      name        = try(tool.profile[0].alias, tool.profile[0].name)
      description = try(tool.profile[0].description, "")
      status      = tool.status
      source_orn  = tool.source_orn
    }
  }
}
```

## Example: Check Tool Validity

```terraform
data "okta_agent_gateways_tool" "check" {
  id               = "tool_integration"
  agent_gateway_id = "agw_main"
}

locals {
  tool_valid = data.okta_agent_gateways_tool.check.status == "VALID"
}

output "tool_validation" {
  value = {
    id      = data.okta_agent_gateways_tool.check.id
    valid   = local.tool_valid
    status  = data.okta_agent_gateways_tool.check.status
  }
}
```

## Notes

- This data source retrieves tools selected for a specific agent gateway
- Tools represent MCP server capabilities exposed through the gateway
- The `alias` field allows administrators to customize how the tool appears to agents
- Tool status indicates whether the tool is still valid relative to the upstream MCP server
- The `connection_orn` and `source_orn` fields track the connection and tool source in Okta's resource namespace