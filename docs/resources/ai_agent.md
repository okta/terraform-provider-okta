---
page_title: "Resource: okta_ai_agent"
description: |-
  Manages an Okta AI Agent resource.
---

# Resource: okta_ai_agent

Manages an Okta AI Agent. AI agents are registered principals that represent external AI systems within Okta's identity and access management infrastructure.

## Example Usage

```hcl
resource "okta_ai_agent" "example" {
  profile {
    name        = "my-ai-agent"
    description = "External AI system agent"
    platform    = "AWS"
  }

  oauth_client {
    client_id_match_pattern = "https://example.com/.well-known/agent-metadata"
    type                    = "SERVICE_APP"
  }
}
```

## Schema

### Required

- `profile` (Block, Required) AI agent profile
  - `name` (String, Required) Unique name of the AI agent

### Optional

- `profile` (Block, Optional) AI agent profile
  - `description` (String, Optional) Description of the AI agent
  - `external_id` (String, Optional) Unique external ID for an AI agent
  - `platform` (String, Optional) The hosting platform of the AI agent
- `oauth_client` (Block, Optional) Registers the AI agent as a CIMD OAuth client
  - `client_id_match_pattern` (String, Optional) The HTTPS URL of the agent's Client ID Metadata Document (CIMD)
  - `type` (String, Optional) OAuth client type
- `sign_on_provider` (Block, Optional) Sign-on provider to bind to the workload principal
  - `app_instance_id` (String, Optional) ID of the app instance to bind the workload principal to
  - `type` (String, Optional) Sign-on provider type
- `resource_url` (String, Optional) The resource URL for the agent-to-agent (A2A) server associated with this AI agent

### Read-Only

- `id` (String) The unique identifier for the resource.
- `status` (String) Lifecycle status of the AI agent (typically `STAGED` when first created).
- `oauth_client_id` (String) The workload principal's OAuth `clientId`.

## Import

Import using `{id}`:

```shell
terraform import okta_ai_agent.example aia_abc123def456
```

> **Note**: Fields not returned by the API `GET` endpoint (such as `resource_url`) will not be restored after import.
> Add them back manually or use `lifecycle { ignore_changes = [...] }`.

## Notes

- AI agents are created in `STAGED` status
- The `oauth_client_id` is automatically generated upon creation and represents the workload principal's OAuth client
- External systems use the agent's identity to authenticate and authorize within Okta