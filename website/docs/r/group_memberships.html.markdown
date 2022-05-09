---
layout: "okta"
page_title: "Okta: okta_group_memberships"
sidebar_current: "docs-okta-resource-group-memberships"
description: |-
  Resource to manage a set of memberships for a specific group.
---

# okta_group_memberships

Resource to manage a set of memberships for a specific group.

This resource will allow you to bulk manage group membership in Okta for a given group. This offers an interface to pass multiple users into a single resource call, for better API resource usage. Effectively this is the same as using the `okta_group_membership` resource several times with a single group and different users. If you need a relationship of a single user to many groups, please use the `okta_user_group_memberships` resource.

When using this with a `okta_user` resource, you should add a lifecycle ignore for group memberships to avoid conflicts in desired state.

## Example Usage

```hcl
resource "okta_group" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "testing, testing"
}

resource "okta_group_memberships" "test" {
  group_id = okta_group.test.id
  users = [
    okta_user.test1.id,
    okta_user.test2.id,
  ]
}
```

## Argument Reference

The following arguments are supported:

- `group_id` - (Required) Okta group ID.
- `users` - (Required) The list of Okta user IDs which the group should have membership managed for.

## Attributes Reference

N/A

## Import

an Okta Group's memberships can be imported via the Okta group ID.

```
$ terraform import okta_group_memberships.test &#60;group id&#62;
```
