---
page_title: "okta_application_group_push_mapping Resource - terraform-provider-okta"
description: |-
  Manages an Okta Application Group Push Mapping resource.
---

# okta_application_group_push_mapping (Resource)

Manages an Okta Application Group Push Mapping.

## Example Usage

```hcl
resource "okta_application_group_push_mapping" "example" {
  app_id = "<app-id>"

  # Optional fields
  # target_group_name = "Example Target Group Name"
}
```

## Schema

### Required

- `app_id` (String) ID of the parent application

### Optional

- `target_group_name` (String) The name of the target group for the group push mapping. This is used when creating a new downstream group. If the group already exists, it links to the existing group. Required if `targetGroupId` ...

### Read-Only

- `id` (String) The unique identifier for the resource.
- `app_config` (String) Additional app configuration for group push mappings. Currently only required for Active Directory.
- `created` (String) Timestamp when the group push mapping was created
- `error_summary` (String) The error message summary if the latest push failed
- `last_push` (String) Timestamp when the group push mapping was pushed
- `last_updated` (String) Timestamp when the group push mapping was last updated
- `source_group_id` (String) The ID of the source group for the group push mapping
- `status` (String) The status of the group push mapping
- `target_group_id` (String) The ID of the target group for the group push mapping

## Import

Import using `{app_id}/{id}`:

```shell
terraform import okta_application_group_push_mapping.example <app_id> <id>
```
