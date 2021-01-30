---
layout: 'okta' page_title: 'Okta: okta_group_membership' sidebar_current: 'docs-okta-resource-group-membership'
description: |- Manage an individual instance of group membership.
---

# okta_group_membership

Manage an individual instance of group membership.

This resource allows you to manage group membership for a given user and group at an individual level. This allows you
to manage group membership in terraform without overriding other automatic membership operations performed by group
rules and other non-managed actions.

When using this with a `okta_user` resource, you should add a lifecycle ignore for group memberships to avoid conflicts
in desired state.

## Example Usage

```hcl
resource "okta_group_membership" "example" {
  group_id = "00g1mana0vCrxzQY84x7"
  user_id  = "00u1manxvp7QBAGgk4x7"
}
```

## Argument Reference

The following arguments are supported:

- `group_id` - (Required) The ID of the Okta Group.

- `user_id` - (Required) The ID of the Okta User.

## Attributes Reference

N/A
