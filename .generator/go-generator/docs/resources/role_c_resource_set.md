---
page_title: "okta_role_c_resource_set Resource - terraform-provider-okta"
description: |-
  Manages an Okta Role C Resource Set resource.
---

# okta_role_c_resource_set (Resource)

Manages an Okta Role C Resource Set.

## Example Usage

```hcl
resource "okta_role_c_resource_set" "example" {
  resources = "<resources>"

  # Optional fields
  # description = "Example description"
  # label = "Example Label"
}
```

## Schema

### Required

- `resources` (List) The endpoint (URL) that references all resource objects included in the resource set. Resources are identified by either an Okta Resource Name (ORN) or by a REST URL format. See [Okta Resource Name...

### Optional

- `description` (String) Description of the resource set
- `label` (String) Unique label for the resource set

### Read-Only

- `id` (String) The unique identifier for the resource.
- `created` (String) Timestamp when the role was created
- `last_updated` (String) Timestamp when the role was last updated
