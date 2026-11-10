---
page_title: "Data Source: okta_ai_agents_schema_default"
description: |-
  Retrieves the AI agent schema, which defines the base attributes and org-defined custom attributes for AI agents.
---

# Data Source: okta_ai_agents_schema_default

Retrieves the org's AI agent schema, which defines the base attributes and org-defined custom attributes for AI agents.

Each org has exactly one AI agent schema. It is returned as a one-element `items` list.

## Example Usage

```hcl
data "okta_ai_agents_schema_default" "example" {}

output "ai_agent_schema_title" {
  value = data.okta_ai_agents_schema_default.example.items[0].title
}
```

## Schema

### Optional

- `id` (String) Unique identifier of the data source. Set by the provider to `okta_ai_agents_schema_default_list`; you do not need to set it.

### Read-Only

- `items` (Block List) The org's AI agent schema, as a single-element list. See [below](#nested-schema-for-items).

<a id="nested-schema-for-items"></a>
### Nested Schema for `items`

Read-Only:

- `id` (String) Identifier of the schema. Always `ai_agents_schema_default`.
- `schema` (String) JSON schema version identifier.
- `name` (String) Name of the schema.
- `title` (String) Display name for the schema.
- `description` (String) Description for the schema.
- `type` (String) Type of the root schema.
- `created` (String) Timestamp when the schema was created.
- `last_updated` (String) Timestamp when the schema was last updated.
- `properties` (Map of String) AI agent schema properties, as returned by the API. Each value is a string; nested objects are rendered with Go's default formatting, not JSON.

-> **Note:** This data source does not return the schema's `definitions` (the base attributes and the `#custom` properties). To manage custom properties, use the `okta_ai_agents_schema_default` resource.
