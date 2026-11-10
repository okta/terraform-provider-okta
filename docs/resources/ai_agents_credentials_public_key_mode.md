---
page_title: "Resource: okta_ai_agents_credentials_public_key_mode"
description: |-
  Sets the active key source mode for an AI agent.
---

# Resource: okta_ai_agents_credentials_public_key_mode

Sets where Okta looks for an AI agent's public keys: the JSON Web Keys stored on the agent (`JWKS`), or the agent's external JWKS URI (`JWKS_URI`).

~> **Note:** Every agent always has a key source mode, so there is nothing to delete. Destroying this resource only removes it from Terraform state and shows a warning; the agent keeps its current mode.

## Example Usage

```hcl
resource "okta_ai_agent" "example" {
  profile {
    name = "example-agent"
  }

  resource_url = "https://example.com/api"

  sign_on_provider {
    type = "NEW_OIDC_APP"
  }
}

resource "okta_ai_agents_credentials_public_key_mode" "example" {
  agent_id = okta_ai_agent.example.id
  mode     = "JWKS"
}
```

## Schema

### Required

- `agent_id` (String) ID of the AI agent in Okta. Forces replacement if changed.
- `mode` (String) The active key source mode: `JWKS` or `JWKS_URI`. Can be changed in place.

### Read-Only

- `id` (String) Same as `agent_id`.

## Import

Import using the agent ID:

```shell
terraform import okta_ai_agents_credentials_public_key_mode.example <agent-id>
```
