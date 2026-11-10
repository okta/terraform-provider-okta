---
page_title: "Resource: okta_ai_agents_schema_default"
description: |-
  Manages the org's AI agent schema, including its custom attributes.
---

# Resource: okta_ai_agents_schema_default

Manages the org's AI agent schema: the base attributes every AI agent has, and the custom attributes you add to it.

Each org has exactly one AI agent schema, so use at most one of these resources per org.
- **Create** updates the existing schema. It does not make a new one.
- **Destroy** only removes the resource from Terraform state. The schema and its custom attributes stay in Okta.

## Example Usage

```hcl
resource "okta_ai_agents_schema_default" "example" {
  definitions {
    custom {
      properties = jsonencode({
        department = {
          title       = "Department"
          type        = "string"
          description = "The department that owns this AI agent"
          mutability  = "READ_WRITE"
          required    = false
        }
      })
    }
  }
}
```

## Schema

### Optional

- `description` (String) Description for the schema.
- `title` (String) Display name for the schema.
- `properties` (Map of String) AI agent schema properties.
- `definitions` (Block, Optional) The schema's attribute definitions. See [below](#nested-schema-for-definitions).

### Read-Only

- `id` (String) Always `default`.
- `name` (String) Name of the schema.
- `schema` (String) JSON schema version identifier.
- `type` (String) Type of the root schema.

<a id="nested-schema-for-definitions"></a>
### Nested Schema for `definitions`

- `custom` (Block, Optional) The `#custom` subschema: attributes your org adds to AI agents.
  - `properties` (String) JSON object of custom attributes, keyed by attribute name. Each value is an attribute definition with fields such as `title`, `type`, `description`, `mutability` and `required`. Build it with `jsonencode()`.
  - `id` (String, Read-Only) The subschema name.
  - `type` (String, Read-Only) The object type.
  - `required` (List of String, Read-Only) Names of the custom attributes that are required.
- `base` (Block, Optional) The `#base` subschema: Okta's built-in attributes (`name`, `description`, `external_id`, `platform`). You normally leave this out of your configuration. Okta defines the base attributes, and they are only read back when you configure this block or import the resource.

~> **Note:** Custom attributes in `definitions.custom.properties` are sent to Okta but not read back. If someone changes them outside Terraform, `terraform plan` does not show a difference, and importing the resource does not fill them in.

## Import

Import using the fixed ID `default`. Any other value is accepted and ignored, because there is only one schema per org:

```shell
terraform import okta_ai_agents_schema_default.example default
```

After importing, add your `definitions.custom.properties` to the configuration: import cannot restore them.
