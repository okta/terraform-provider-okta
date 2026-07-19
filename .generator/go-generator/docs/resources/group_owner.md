---
page_title: "okta_group_owner Resource - terraform-provider-okta"
description: |-
  Manages an Okta Group Owner resource.
---

# okta_group_owner (Resource)

Manages an Okta Group Owner.

## Example Usage

```hcl
resource "okta_group_owner" "example" {
  group_id = "<group-id>"

  # Optional fields
  # origin_id = "<origin-id>"
  # origin_type = "<origin_type>"
  # resolved = true
  # type = "<type>"
}
```

## Schema

### Required

- `group_id` (String) ID of the group

### Optional

- `origin_id` (String) The ID of the app instance if the `originType` is `APPLICATION`. This value is `NULL` if `originType` is `OKTA_DIRECTORY`.
- `origin_type` (String) The source where group ownership is managed
- `resolved` (Boolean) If `originType`is APPLICATION, this parameter is set to `FALSE` until the owner's `originId` is reconciled with an associated Okta ID.
- `type` (String) The entity type of the owner

### Read-Only

- `id` (String) The unique identifier for the resource.
- `display_name` (String) The display name of the group owner
- `last_updated` (String) Timestamp when the group owner was last updated

## Import

Import using `{group_id}/{id}`:

```shell
terraform import okta_group_owner.example <group_id> <id>
```
