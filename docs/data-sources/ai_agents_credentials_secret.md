---
page_title: "Data Source: okta_ai_agents_credentials_secret"
description: |-
  Retrieves the metadata for an AI agent client secret by secret ID.
---

# Data Source: okta_ai_agents_credentials_secret

Retrieves the metadata for an AI agent client secret by secret ID.

## Example Usage

```hcl
data "okta_ai_agents_credentials_secret" "example" {
  id       = okta_ai_agents_credentials_secret.example.id
  agent_id = okta_ai_agent.example.id
}
```

## Schema

### Required

- `id` (String) ID of the client secret.
- `agent_id` (String) ID of the AI agent in Okta.

### Read-Only

- `status` (String) Status of the AI agent client secret.
- `secret_hash` (String) Hash of the client secret.
- `client_secret` (String) The client secret string. Okta returns it in plaintext only when the secret is created, so this data source always reads it as empty. Use the `okta_ai_agents_credentials_secret` resource's output at creation time if you need the value.
- `created` (String) Timestamp when the client secret was created.
- `last_updated` (String) Timestamp when the client secret was last updated.
