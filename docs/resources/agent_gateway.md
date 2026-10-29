---
page_title: "Resource: okta_agent_gateway"
description: |-
  Manages an Okta agent gateway resource.
---

# Resource: okta_agent_gateway

Manages an Okta agent gateway. Agent gateways serve as communication bridges facilitating secure connections between AI agents and Okta systems.

## Example Usage

```hcl
resource "okta_agent_gateway" "example" {
  resource_path = "my-gateway-path"

  profile {
    display_name = "Production Agent Gateway"
    description  = "Main gateway for production agents"
  }
}
```

## Schema

### Required

- `resource_path` (String) A URI path suffix that uniquely identifies this agent gateway. Forces replacement if changed.

### Optional

- `profile` (Block, Optional) Admin-editable profile for the agent gateway
  - `display_name` (String, Required) Human-readable display name for the agent gateway
  - `description` (String, Optional) Description of the agent gateway

### Read-Only

- `id` (String) The unique identifier for the resource.
- `orn` (String) The [ORN](https://developer.okta.com/docs/concepts/orn/) identifier for this agent gateway.
- `status` (String) Current lifecycle status of the agent gateway.

## Import

Import using `{id}`:

```shell
terraform import okta_agent_gateway.example agw_abc123def456
```

> **Note**: The `resource_path` field is write-only and will not be restored after import.
> Add it back manually or use `lifecycle { ignore_changes = [resource_path] }`.

## Notes

- The `resource_path` field is sensitive and will not be displayed in logs
- Agent gateways must be properly configured to establish secure connections with agents
- Status changes automatically based on gateway lifecycle events