---
page_title: "Data Source: okta_ai_agent_providers_staging_stat"
description: |-
  Retrieves the counts of staging items by type for an AI agent provider
---

# Data Source: okta_ai_agent_providers_staging_stat

Retrieves the counts of staging items broken down by type (`NEW`, `DUPLICATE`, `IGNORED`) for the specified AI agent provider.

## Example Usage

```terraform
data "okta_ai_agent_providers_staging_stat" "example" {
  provider_id = okta_ai_agent_provider.example.id
}

output "staging_counts" {
  value = {
    new       = data.okta_ai_agent_providers_staging_stat.example.items[0].new_count
    duplicate = data.okta_ai_agent_providers_staging_stat.example.items[0].duplicate_count
    ignored   = data.okta_ai_agent_providers_staging_stat.example.items[0].ignored_count
  }
}
```

## Argument Reference

### Required

- `provider_id` (String) ID of the AI agent provider

### Optional

- `id` (String) Unique identifier of the ai_agent_providers_staging_stat.

## Attribute Reference

- `items` (Block List) List of ai_agent_providers_staging_stat items. Contains one entry with the counts below.
  - `id` (String) Identifier of the staging stat item.
  - `new_count` (Number) The number of staging items of type `NEW`.
  - `duplicate_count` (Number) The number of staging items of type `DUPLICATE`.
  - `ignored_count` (Number) The number of staging items of type `IGNORED`.