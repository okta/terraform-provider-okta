---
page_title: "okta_group_role_assignment_user_admin Resource - terraform-provider-okta"
description: |-
  Manages an Okta Group Role Assignment User Admin resource.
---

# okta_group_role_assignment_user_admin (Resource)

Manages an Okta Group Role Assignment User Admin.

## Example Usage

```hcl
resource "okta_group_role_assignment_user_admin" "example" {
  group_id = "<group-id>"
  type = "<type>"

  # Optional fields
  # catalog = "<catalog>"
  # profile = "<profile>"
  # type = "<type>"
  # assignment_type = "<assignment_type>"
  # status = "ACTIVE"
}
```

## Schema

### Required

- `group_id` (String) ID of the group
- `type` (String) Discriminator field identifying the variant type. Must be set to \

### Optional

- `catalog` (String) App targets
- `profile` (String) Specifies required and optional properties for a group. The `objectClass` of a group determines which additional properties are available.  You can extend group profiles with custom properties, but...
- `type` (String) Determines how a group's profile and memberships are managed
- `assignment_type` (String) Role assignment type
- `status` (String) Status

### Read-Only

- `id` (String) The unique identifier for the resource.
- `embedded` (String) Embedded resources related to the group
- `created` (String) Timestamp when the group was created
- `last_membership_updated` (String) Timestamp when the groups memberships were last updated
- `last_updated` (String) Timestamp when the group's profile was last updated
- `object_class` (List) Determines the group's `profile`
- `created` (String) Timestamp when the object was created
- `label` (String) Label for the role assignment
- `last_updated` (String) Timestamp when the object was last updated

## Import

Import using `{group_id}/{id}`:

```shell
terraform import okta_group_role_assignment_user_admin.example <group_id> <id>
```
