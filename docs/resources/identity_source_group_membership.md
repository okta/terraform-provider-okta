---
page_title: "Resource: okta_identity_source_group_membership"
description: |-
  Manages a single group membership within an Okta Identity Source.
---

# Resource: okta_identity_source_group_membership

Manages a single member association between a user and a group in an Okta Identity Source. Each resource instance adds one user (identified by `member_external_id`) to the specified identity source group. All three key attributes force resource replacement when changed.

## Example Usage

```terraform
resource "okta_identity_source_group_membership" "example" {
  identity_source_id   = "<identity-source-id>"
  group_or_external_id = "GRPEXT123456EXAMPLE"
  member_external_id   = "USEREXT123456EXAMPLE"
}
```

## Schema

### Required

- `identity_source_id` (String) ID of the identity source. Forces replacement when changed.
- `group_or_external_id` (String) External ID of the identity source group. Forces replacement when changed.
- `member_external_id` (String) The external ID of the user to be added as a member of the group in Okta. Forces replacement when changed.

### Read-Only

- `id` (String) The unique identifier for the resource (set to the member external ID).
- `member_external_ids` (List of String) A list of all member external IDs currently in the group (populated after creation).

## Import

Import using `{identity_source_id}/{group_or_external_id}/{id}`:

```shell
terraform import okta_identity_source_group_membership.example <identity-source-id>/<group-external-id>/<member-external-id>
```
