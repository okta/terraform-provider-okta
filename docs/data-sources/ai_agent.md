---
page_title: "Data Source: okta_ai_agent"
description: |-
  Retrieves a registered AI agent by ID
---

# Data Source: okta_ai_agent

Retrieves a registered AI agent by its unique identifier from Okta's AI Agents registry in Universal Directory.

## Example Usage

```terraform
data "okta_ai_agent" "example" {
  id = "aiagent_abc123def456"
}

output "agent_profile" {
  value = {
    name        = try(data.okta_ai_agent.example.profile[0].name, null)
    description = try(data.okta_ai_agent.example.profile[0].description, null)
    platform    = try(data.okta_ai_agent.example.profile[0].platform, null)
    status      = data.okta_ai_agent.example.status
  }
}
```

## Argument Reference

### Required

- `id` - (String) Unique identifier of the AI agent

## Attribute Reference

- `created` - (String) Timestamp when the AI agent was created (RFC3339 format)
- `last_updated` - (String) Timestamp when the AI agent was last updated (RFC3339 format)
- `orn` - (String) The [ORN](https://developer.okta.com/docs/concepts/orn/) identifier for this agent
- `status` - (String) Current lifecycle status of the AI agent
- `profile` - (Block) AI agent profile information
  - `name` - (String) Unique name of the AI agent
  - `description` - (String) Description of the AI agent
  - `platform` - (String) The hosting platform of the AI agent (e.g., OpenAI, Anthropic)
  - `external_id` - (String) External identifier from the source platform
- `oauth_client` - (Block) OAuth client credentials for the agent
  - `client_id` - (String) The OAuth client ID for the agent
  - `type` - (String) Type of OAuth client
- `sign_on_provider` - (Block) Sign-on provider configuration
  - `app_instance_id` - (String) ID of the application instance providing authentication

## Example: Agent Integration Lookup

```terraform
data "okta_ai_agent" "copilot" {
  id = "aiagent_copilot_01"
}

locals {
  agent_name = try(data.okta_ai_agent.copilot.profile[0].name, "Unknown")
  is_active  = data.okta_ai_agent.copilot.status == "ACTIVE"
}

output "agent_ready" {
  value = {
    name        = local.agent_name
    active      = local.is_active
    created_at  = data.okta_ai_agent.copilot.created
    platform    = try(data.okta_ai_agent.copilot.profile[0].platform, null)
    client_id   = try(data.okta_ai_agent.copilot.oauth_client[0].client_id, null)
  }
}
```

## Example: Conditional Usage

```terraform
variable "required_agent_id" {
  type        = string
  description = "ID of the required AI agent"
}

data "okta_ai_agent" "required" {
  id = var.required_agent_id
}

locals {
  agent_available = data.okta_ai_agent.required.status == "ACTIVE"
  agent_details   = data.okta_ai_agent.required.profile[0]
}

resource "null_resource" "validate_agent" {
  count = local.agent_available ? 1 : 0
  
  triggers = {
    agent_id   = data.okta_ai_agent.required.id
    agent_name = local.agent_details.name
  }
}
```

## Notes

- This data source is read-only
- AI agents must be registered in Universal Directory before they can be queried
- Status values typically include: `ACTIVE`, `INACTIVE`, `STAGED`, etc.
- The ORN can be used to reference the agent in other contexts
- OAuth client credentials should be handled securely if exposed