---
page_title: "Data Source: okta_agent_gateway"
description: |-
  Retrieves an agent gateway by ID
---

# Data Source: okta_agent_gateway

Retrieves an agent gateway by its unique identifier. Agent gateways are communication bridges used to facilitate secure connections between AI agents and Okta systems.

## Example Usage

```terraform
data "okta_agent_gateway" "example" {
  id = "agw_abc123def456"
}

output "gateway_status" {
  value = data.okta_agent_gateway.example.status
}
```

## Argument Reference

### Required

- `id` - (String) Unique identifier of the agent gateway

## Attribute Reference

- `created` - (String) Timestamp when the agent gateway was created (RFC3339 format)
- `last_updated` - (String) Timestamp when the agent gateway was last updated (RFC3339 format)
- `orn` - (String) The [ORN](https://developer.okta.com/docs/concepts/orn/) identifier for this agent gateway
- `status` - (String) Current lifecycle status of the agent gateway
- `profile` - (Block) Admin-editable profile for the agent gateway
  - `display_name` - (String) Human-readable display name for the agent gateway
  - `description` - (String) Description of the agent gateway

## Example with Profile Access

```terraform
data "okta_agent_gateway" "example" {
  id = "agw_abc123def456"
}

output "gateway_info" {
  value = {
    id             = data.okta_agent_gateway.example.id
    display_name   = try(data.okta_agent_gateway.example.profile[0].display_name, null)
    description    = try(data.okta_agent_gateway.example.profile[0].description, null)
    status         = data.okta_agent_gateway.example.status
    last_updated   = data.okta_agent_gateway.example.last_updated
  }
}
```

## Notes

- This data source is read-only
- Agent gateways are managed separately from their configuration
- The ORN can be used to reference the gateway in other resources
- Status values typically include: `ACTIVE`, `INACTIVE`, etc.