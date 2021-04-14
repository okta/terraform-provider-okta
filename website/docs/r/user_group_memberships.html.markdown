---
layout: "okta"
page_title: "Okta: okta_user_group_membership"
sidebar_current: "docs-okta-resource-user-group-memberships"
description: |-
  Resource to manage a set of group memberships for a specific user.
---

# okta_user_group_memberships

Resource to manage a set of group memberships for a specific user.

This resource allows you to bulk manage groups for a single user, independent of the user schema itself. This allows you
to manage group membership in terraform without overriding other automatic membership operations performed by group
rules and other non-managed actions.

When using this with a `okta_user` resource, you should add a lifecycle ignore for group memberships to avoid conflicts
in desired state.

## Example Usage

```hcl
resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc-replace_with_uuid@example.com"
  email      = "testAcc-replace_with_uuid@example.com"

  lifecycle {
    ignore_changes = [group_memberships]
  }
}

resource "okta_user_group_memberships" "test" {
  user_id = okta_user.test.id
  groups = [
    okta_group.test_1.id,
    okta_group.test_2.id,
  ]
}
```

## Argument Reference

The following arguments are supported:

- `user_id` - (Required) ID of a Okta User.
- `groups` - (Required) The list of Okta group IDs which the user should have membership managed for.

## Attributes Reference

N/A
