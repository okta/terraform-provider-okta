---
page_title: "Data Source: okta_ai_agent_providers_profile_mapping"
description: |-
  Retrieves profile mappings for an AI agent provider
---

# Data Source: okta_ai_agent_providers_profile_mapping

Retrieves the AI agent provider's profile mappings that define how source app attributes are mapped to AI agent profile attributes during import.

## Example Usage

```terraform
data "okta_ai_agent_providers_profile_mapping" "example" {
  provider_id = okta_ai_agent_provider.example.id
}
```

## Argument Reference

### Required

- `provider_id` - (String) ID of the AI agent provider

### Optional

- `id` - (String) Unique identifier of the profile mapping. If not specified, the data source will return all mappings for the provider.

## Attribute Reference

- `items` - (List of Objects) List of profile mappings for the provider, each containing:
  - `id` - (String) Unique identifier for the mapping
  - `source_expression` - (String) OEL (Okta Expression Language) expression that maps the source attribute value to the target field
  - `target_field` - (String) Target AI agent profile attribute name
  - `per_platform` - (List of Strings) Optional per-platform sources for this attribute, useful for integrations covering several platforms whose source fields differ

## Example Usage with Filters

```terraform
data "okta_ai_agent_providers_profile_mapping" "mappings" {
  provider_id = "00a1gjdxyz0000000000"
}

output "mapping_details" {
  value = [
    for mapping in data.okta_ai_agent_providers_profile_mapping.mappings.items : {
      target_field       = mapping.target_field
      source_expression  = mapping.source_expression
      per_platform_count = length(mapping.per_platform)
    }
  ]
}
```

## Notes

- This data source is read-only and retrieves existing profile mappings
- Profile mappings are managed through the AI Agent Providers API
- The `source_expression` field contains OEL expressions used to transform source attributes during the import process
- Use this data source to audit or reference existing mappings in your Terraform configuration