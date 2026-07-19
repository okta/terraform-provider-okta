---
page_title: "okta_group_role_assignment_access_requests_admin Resource - terraform-provider-okta"
description: |-
  Manages an Okta Group Role Assignment Access Requests Admin resource.
---

# okta_group_role_assignment_access_requests_admin (Resource)

Manages an Okta Group Role Assignment Access Requests Admin.

## Example Usage

```hcl
resource "okta_group_role_assignment_access_requests_admin" "example" {
  group_id = "<group-id>"
  type = "<type>"

  # Optional fields
  # assignment_type = "<assignment_type>"
  # status = "ACTIVE"
}
```

## Schema

### Required

- `group_id` (String) ID of the group
- `type` (String) Discriminator field identifying the variant type. Must be set to \

### Optional

- `assignment_type` (String) Role assignment type
- `status` (String) Status

### Read-Only

- `id` (String) The unique identifier for the resource.
- `created` (String) Timestamp when the object was created
- `label` (String) Label for the role assignment
- `last_updated` (String) Timestamp when the object was last updated
- `resource_set` (String) Resource set ID
- `role` (String) Role ID

## Import

Import using `{group_id}/{id}`:

```shell
terraform import okta_group_role_assignment_access_requests_admin.example <group_id> <id>
```
