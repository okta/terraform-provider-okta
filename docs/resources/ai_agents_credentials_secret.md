---
page_title: "Resource: okta_ai_agents_credentials_secret"
description: |-
  Creates a client secret for an AI agent.
---

# Resource: okta_ai_agents_credentials_secret

Creates a client secret for an AI agent. Okta generates the secret unless you supply one.

On destroy the secret is deactivated and then deleted. Okta rejects deactivating an agent's only active secret, so keep at least one other active secret, or destroy the agent instead.

~> **Note:** A secret can't be changed after it is created, and that includes `status`. Changing `agent_id` replaces the secret; changing `status` or `client_secret` isn't supported and fails the apply with "Provider produced inconsistent result". To rotate, add a new secret resource.

~> **Note:** Okta returns the plaintext secret only when it is created. It is kept in Terraform state from then on, so protect your state. After `terraform import`, `client_secret` is empty.

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

resource "okta_ai_agents_credentials_secret" "example" {
  agent_id = okta_ai_agent.example.id
  status   = "ACTIVE"
}

output "agent_client_secret" {
  value     = okta_ai_agents_credentials_secret.example.client_secret
  sensitive = true
}
```

## Schema

### Required

- `agent_id` (String) ID of the AI agent in Okta. Forces replacement if changed.

### Optional

- `client_secret` (String) The client secret. Leave unset to have Okta generate one.
- `status` (String) Status of the client secret: `ACTIVE` or `INACTIVE`. Set at creation; it can't be changed afterwards.

### Read-Only

- `id` (String) ID of the client secret.
- `secret_hash` (String) Hash of the client secret.

## Import

Import using `{agent_id}/{id}`:

```shell
terraform import okta_ai_agents_credentials_secret.example <agent-id>/<secret-id>
```
