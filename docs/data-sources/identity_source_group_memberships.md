---
page_title: "Data Source: okta_identity_source_group_memberships"
description: |-
  Retrieves the member list for an Okta Identity Source group.
---

# Data Source: okta_identity_source_group_memberships

Retrieves the list of member external IDs for a group in an Okta Identity Source.

## Example Usage

```terraform
data "okta_identity_source_group_memberships" "example" {
  identity_source_id = "<identity-source-id>"
  group_external_id  = "GROUPEXT123456EXAMPLE"
}

output "member_ids" {
  value = data.okta_identity_source_group_memberships.example.member_external_ids
}
```

## Schema

### Required

- `identity_source_id` (String) ID of the identity source.
- `group_external_id` (String) The external ID (or Okta-assigned ID) of the group whose memberships to retrieve.

### Read-Only

- `id` (String) Composite identifier (`{identity_source_id}/{group_external_id}`).
- `member_external_ids` (List of String) External IDs of members belonging to the group in the identity source.
