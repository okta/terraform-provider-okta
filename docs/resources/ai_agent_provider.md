---
page_title: "Resource: okta_ai_agent_provider"
description: |-
  Manages an AI Agent Provider in Okta's AI Agents registry.
---

# Resource: okta_ai_agent_provider

Manages an AI Agent Provider, allowing integrations that enable AI agents from external platforms to be imported and registered in Okta's AI Agents registry in Universal Directory.

## Example Usage

### Basic Configuration

```terraform
resource "okta_ai_agent_provider" "example" {
  source_orn = "orn:okta:ai-agents-v2:us-east-1:00o1gjdxyz:provider/example-provider"

  configuration = {
    api_endpoint = "https://api.example.com"
    api_key      = sensitive("sk-proj-xxxxx")
  }
}
```

### With Schedule and Owners

```terraform
resource "okta_ai_agent_provider" "example_with_schedule" {
  source_orn = "orn:okta:ai-agents-v2:us-east-1:00o1gjdxyz:provider/example-provider"

  configuration = {
    api_endpoint = "https://api.example.com"
    api_key      = sensitive("sk-proj-xxxxx")
  }

  de_dup_expression = "id"
  stage             = "PREVIEW"

  schedule {
    cron     = "0 0 * * *"
    timezone = "America/New_York"
  }

  owners {
    group_orn = "orn:okta:directory:us-east-1:00o1gjdxyz:group/admin-group-id"
    user_orns = [
      "orn:okta:directory:us-east-1:00o1gjdxyz:user/user-id-1",
      "orn:okta:directory:us-east-1:00o1gjdxyz:user/user-id-2"
    ]
  }
}
```

## Argument Reference

### Required

- `source_orn` - (String) The [ORN](https://developer.okta.com/docs/concepts/orn/) of the source representing the AI agent provider in the external platform.

### Optional

- `configuration` - (Map of String) Configuration details for connecting to the AI agent provider. This field is sensitive and will not be displayed in logs or state files.
- `de_dup_expression` - (String) Expression to uniquely identify imported agents, used to prevent duplicates during import.
- `stage` - (String) Preview option to unify identity for imported agents.
- `schedule` - (Block) Defines the schedule for importing agents from the provider.
  - `cron` - (String) The import schedule in cron format (e.g., `"0 0 * * *"` for daily at midnight).
  - `timezone` - (String) Timezone where the scheduled import takes place (e.g., `"America/New_York"`).
- `owners` - (Block) Defines ownership and administrative access to this provider.
  - `group_orn` - (String) The resource-owner group [ORN](https://developer.okta.com/docs/concepts/orn/).
  - `user_orns` - (List of String) List of resource-owner user [ORNs](https://developer.okta.com/docs/concepts/orn/).

### Read-Only

- `id` - (String) The unique identifier for the provider.
- `source_id` - (String) The ID of the source.
- `source_name` - (String) The name of the source.
- `source_type` - (String) The type of the source.
- `status` - (String) The current status of the AI agent provider.
- `import_in_progress` - (Boolean) Whether an import job for this provider is currently queued or running.
- `integration_date` - (String) The ISO 8601 date and time when the AI agent provider was integrated.
- `last_import_date` - (String) The ISO 8601 date and time when the last import from the AI agent provider occurred.
- `next_scheduled_import_date` - (String) The ISO 8601 date and time of the next scheduled import.

## Import

Import an AI Agent Provider using its ID:

```shell
terraform import okta_ai_agent_provider.example <provider_id>
```

## Notes

- The `configuration` field is sensitive and should use the `sensitive()` function for sensitive values like API keys.
- The resource automatically deactivates the provider before deletion to ensure proper cleanup.
- Import operations are performed asynchronously; monitor the `import_in_progress` and `last_import_date` attributes to track import status.
- The `de_dup_expression` field helps prevent duplicate agents during import by identifying unique agents.