---
page_title: "Resource: okta_ai_agents_credentials_jwks_uri"
description: |-
  Sets the external JWKS URI for an AI agent.
---

# Resource: okta_ai_agents_credentials_jwks_uri

Sets the external JWKS URI that Okta uses to fetch an AI agent's public keys. Each agent has at most one JWKS URI.

To make Okta use the URI rather than keys stored on the agent, also set the agent's key source mode to `JWKS_URI` with `okta_ai_agents_credentials_public_key_mode`.

~> **Note:** Okta has no API to remove a JWKS URI. Destroying this resource only removes it from Terraform state and shows a warning; the URI stays set on the agent until you change it or delete the agent.

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

resource "okta_ai_agents_credentials_jwks_uri" "example" {
  agent_id = okta_ai_agent.example.id
  jwks_uri = "https://example.com/.well-known/jwks.json"
}
```

## Schema

### Required

- `agent_id` (String) ID of the AI agent in Okta. Forces replacement if changed.
- `jwks_uri` (String) The JWKS URI. Can be changed in place.

### Read-Only

- `id` (String) Same as `agent_id`.

## Import

Import using the agent ID:

```shell
terraform import okta_ai_agents_credentials_jwks_uri.example <agent-id>
```
