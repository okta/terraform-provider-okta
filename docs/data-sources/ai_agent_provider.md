---
page_title: "Data Source: okta_ai_agent_provider"
description: |-
  Retrieves an AI agent provider
---

# Data Source: okta_ai_agent_provider

Retrieves an AI agent provider by ID. An AI agent provider connects a source app, such as AWS IAM Identity Center or Microsoft Office 365, so Okta can import the AI agents registered there.

## Example Usage

```terraform
data "okta_ai_agent_provider" "example" {
  id = "<ai-agent-provider-id>"
}

output "provider_source" {
  value = {
    name   = data.okta_ai_agent_provider.example.source_name
    status = data.okta_ai_agent_provider.example.status
  }
}
```

## Argument Reference

### Required

- `id` (String) Unique identifier of the AI agent provider.

## Attributes Reference

- `source_id` (String) The ID of the source app.
- `source_name` (String) The name of the source app, for example `amazon_aws_sso` or `office365`.
- `source_orn` (String) The [ORN](https://developer.okta.com/docs/api/openapi/okta-management/guides/roles/#okta-resource-name-orn) of the source app.
- `source_type` (String) The type of the source, for example `apps`.
- `status` (String) The status of the AI agent provider.
- `stage` (String) Preview option to unify identity for imported agents.
- `de_dup_expression` (String) Expression to uniquely identify imported agents.
- `import_in_progress` (Boolean) Whether an import job for this provider is currently queued or running.
- `integration_date` (String) The date and time when the AI agent provider was integrated.
- `last_import_date` (String) The date and time when the last import from the AI agent provider occurred.
- `next_scheduled_import_date` (String) The date and time of the next scheduled import from the AI agent provider.
- `configuration` (Map of String) Provider-specific configuration. The keys depend on the source type. Values that aren't strings in the API, such as lists, are returned in Go's default format rather than as JSON, for example `[us-east-1]`.
- `owners` (Block) Owners of the AI agent provider.
  - `group_orn` (String) The ORN of the owning group.
  - `user_orns` (List of String) The ORNs of the owning users.
- `schedule` (Block) The import schedule.
  - `cron` (String) Cron expression for scheduled imports.
  - `timezone` (String) The timezone the cron expression is evaluated in.
