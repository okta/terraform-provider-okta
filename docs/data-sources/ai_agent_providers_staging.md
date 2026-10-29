---
page_title: "Data Source: okta_ai_agent_providers_staging"
description: |-
  Retrieves a staging item for an AI agent provider
---

# Data Source: okta_ai_agent_providers_staging

Retrieves a staging item for the specified AI agent provider. Staging items represent agents that are being prepared for import or comparison during the AI agent provisioning workflow.

## Example Usage

```terraform
data "okta_ai_agent_providers_staging" "example" {
  id          = "staging-item-id"
  provider_id = okta_ai_agent_provider.example.id
}
```

## Argument Reference

### Required

- `id` - (String) Unique identifier of the staging item
- `provider_id` - (String) ID of the AI agent provider

## Attribute Reference

- `type` - (String) The staging item type

### Existing Agent Block

Returned only if an existing agent is present in the staging item:

- `existing_agent` - (Block) Existing AI agent information
  - `id` - (String) Unique ID for the AI agent in Okta
  - `created` - (String) Timestamp when the AI agent was created (RFC3339 format)
  - `last_updated` - (String) Timestamp when the AI agent was last updated (RFC3339 format)
  - `status` - (String) Current status of the AI agent (typically `STAGED` for newly created agents)
  - `oauth_client` - (Block) The workload principal's OAuth client identity
    - `client_id` - (String) The workload principal's OAuth client ID
    - `type` - (String) OAuth client type
  - `profile` - (Block) AI agent profile information
    - `name` - (String) Unique name of the AI agent
    - `description` - (String) Description of the AI agent
    - `platform` - (String) The hosting platform of the AI agent
    - `external_id` - (String) Unique external ID for the AI agent
  - `sign_on_provider` - (Block) The sign-on provider bound to the workload principal
    - `app_instance_id` - (String) ID of the app instance bound to the workload principal

### Incoming Agent Block

Returned only if an incoming agent is present in the staging item:

- `incoming_agent` - (Block) Incoming AI agent information (same structure as `existing_agent`)

## Example Usage with Local Values

```terraform
data "okta_ai_agent_providers_staging" "check_import" {
  id          = "staging-001"
  provider_id = okta_ai_agent_provider.example.id
}

locals {
  existing_agent_name = try(data.okta_ai_agent_providers_staging.check_import.existing_agent[0].profile[0].name, "None")
  incoming_agent_name = try(data.okta_ai_agent_providers_staging.check_import.incoming_agent[0].profile[0].name, "None")
}

output "staging_comparison" {
  value = {
    staging_type       = data.okta_ai_agent_providers_staging.check_import.type
    existing_agent     = local.existing_agent_name
    incoming_agent     = local.incoming_agent_name
    needs_reconciliation = local.existing_agent_name != local.incoming_agent_name && local.incoming_agent_name != "None"
  }
}
```

## Use Cases

This data source is useful for:

- **Agent Import Auditing**: Monitor what agents are staged for import from external platforms
- **Conflict Detection**: Compare existing agents with incoming agents to identify potential duplicates or conflicts
- **Reconciliation Workflows**: Build automation to approve or reject staged agent imports
- **Profile Validation**: Verify that agent profiles (name, description, platform, external ID) are correct before import

## Notes

- This is a read-only data source that retrieves staging information
- Staging items are typically temporary and represent agents in the import pipeline
- Both `existing_agent` and `incoming_agent` blocks may be present to show a comparison
- The `status` field for newly created agents is typically `STAGED`
- Use this data source to make decisions about agent import approval/rejection based on existing vs. incoming agent data