---
page_title: "okta_linked_object Resource - terraform-provider-okta"
description: |-
  Manages an Okta Linked Object resource.
---

# okta_linked_object (Resource)

Manages an Okta Linked Object.

## Example Usage

```hcl
resource "okta_linked_object" "example" {
  name = "Example Name"
  title = "Example Title"
  type = "<type>"
  name = "Example Name"
  title = "Example Title"
  type = "<type>"

  # Optional fields
  # description = "Example description"
  # description = "Example description"
}
```

## Schema

### Required

- `name` (String) API name of the `primary` or the `associated` link. The `name` parameter can't start with a number and can only contain the following characters: `a-z`, `A-Z`,` 0-9`, and `_`.
- `title` (String) Display name of the `primary` or the `associated` link
- `type` (String) The object type for this relationship
- `name` (String) API name of the `primary` or the `associated` link. The `name` parameter can't start with a number and can only contain the following characters: `a-z`, `A-Z`,` 0-9`, and `_`.
- `title` (String) Display name of the `primary` or the `associated` link
- `type` (String) The object type for this relationship

### Optional

- `description` (String) Description of the `primary` or the `associated` relationship
- `description` (String) Description of the `primary` or the `associated` relationship

### Read-Only

- `id` (String) The unique identifier for the resource.
