---
page_title: "Data Source: okta_ai_agents_credentials_jwks_uri"
description: |-
  Retrieves the stored external JWKS URI for an AI agent.
---

# Data Source: okta_ai_agents_credentials_jwks_uri

Retrieves the stored external JWKS URI for an AI agent. Each agent has at most one JWKS URI; it is returned as a one-element `items` list.

## Example Usage

```hcl
data "okta_ai_agents_credentials_jwks_uri" "example" {
  agent_id = okta_ai_agent.example.id
}

output "agent_jwks_uri" {
  value = data.okta_ai_agents_credentials_jwks_uri.example.items[0].jwks_uri
}
```

## Schema

### Required

- `agent_id` (String) ID of the AI agent in Okta.

### Optional

- `id` (String) Identifier of the data source. Set by the provider to `okta_ai_agents_credentials_jwks_uri_list`; you do not need to set it.

### Read-Only

- `items` (Block List) The agent's JWKS URI, as a single-element list.
  - `id` (String) Always `ai_agents_credentials_jwks_uri`.
  - `jwks_uri` (String) The stored JWKS URI.
