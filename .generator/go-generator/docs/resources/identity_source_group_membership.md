---
page_title: "okta_identity_source_group_membership Resource - terraform-provider-okta"
description: |-
  Manages an Okta Identity Source Group Membership resource.
---

# okta_identity_source_group_membership (Resource)

Manages an Okta Identity Source Group Membership.

## Example Usage

```hcl
resource "okta_identity_source_group_membership" "example" {
  identity_source_id = "<identity-source-id>"
  group_or_external_id = "<group-or-external-id>"

  # Optional fields
  # member_external_ids = "<member_external_ids>"
  # member_external_id = "<member-external-id>"
}
```

## Schema

### Required

- `identity_source_id` (String) ID of the identity source
- `group_or_external_id` (String) External ID of the identity source group

### Optional

- `member_external_ids` (List) A list of app user external IDs that are members of the group in Okta
- `member_external_id` (String) The external ID of the user to be added as a member of the group in Okta

### Read-Only

- `id` (String) The unique identifier for the resource.

## Import

Import using `{identity_source_id}/{group_or_external_id}/{id}`:

```shell
terraform import okta_identity_source_group_membership.example <identity_source_id> <group_or_external_id> <id>
```
