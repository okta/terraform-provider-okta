---
page_title: "okta_group Resource - terraform-provider-okta"
description: |-
  Manages an Okta Group resource.
---

# okta_group (Resource)

Manages an Okta Group.

## Example Usage

```hcl
resource "okta_group" "example" {

  # Optional fields
  # description = "Example description"
  # name = "Example Name"
  # type = "<type>"
}
```

## Schema

### Optional

- `description` (String) Description of the group
- `name` (String) Name of the group
- `type` (String) Determines how a group's profile and memberships are managed

### Read-Only

- `id` (String) The unique identifier for the resource.
- `embedded` (String) Embedded resources related to the group
- `created` (String) Timestamp when the group was created
- `last_membership_updated` (String) Timestamp when the groups memberships were last updated
- `last_updated` (String) Timestamp when the group's profile was last updated
- `object_class` (List) Determines the group's `profile`
