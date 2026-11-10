---
page_title: "Data Source: okta_ai_agents_credentials_public_key_mode"
description: |-
  Retrieves the active key source mode for an AI agent.
---

# Data Source: okta_ai_agents_credentials_public_key_mode

Retrieves the active key source mode for an AI agent: whether Okta verifies the agent's tokens against its stored JSON Web Keys (`JWKS`) or against an external JWKS URI (`JWKS_URI`). It is returned as a one-element `items` list.

## Example Usage

```hcl
data "okta_ai_agents_credentials_public_key_mode" "example" {
  agent_id = okta_ai_agent.example.id
}

output "agent_key_mode" {
  value = data.okta_ai_agents_credentials_public_key_mode.example.items[0].mode
}
```

## Schema

### Required

- `agent_id` (String) ID of the AI agent in Okta.

### Optional

- `id` (String) Identifier of the data source. Set by the provider to `okta_ai_agents_credentials_public_key_mode_list`; you do not need to set it.

### Read-Only

- `items` (Block List) The agent's key source mode, as a single-element list.
  - `id` (String) Always `ai_agents_credentials_public_key_mode`.
  - `mode` (String) The active key source mode: `JWKS` or `JWKS_URI`.
