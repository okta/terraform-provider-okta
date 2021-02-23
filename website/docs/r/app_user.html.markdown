---
layout: 'okta'
page_title: 'Okta: okta_app_user'
sidebar_current: 'docs-okta-resource-app-user'
description: |-
  Creates an Application User.
---

# okta_app_user

Creates an Application User.

This resource allows you to create and configure an Application User.

**When using this resource, make sure to add the following `lifefycle` argument to the application resource you are assigning to:**

```hcl
lifecycle {
  ignore_changes = ["users"]
}
```

~> **Important:** When the `okta_app_user` is retained, by setting `retain_assignment` to `true`, it is no longer managed by Terraform after it is destroyed. To truly delete the assignment, you will need to remove it either through the Okta Console or API. This argument exists for the use case where the same user is assigned in multiple places in order to prevent a single destruction removing all of them.

## Example Usage

```hcl
resource "okta_app_user" "example" {
  app_id   = "<app_id>"
  user_id  = "<user id>"
  username = "example"
}
```

## Argument Reference

The following arguments are supported:

- `app_id` - (Required) App to associate user with.

- `user_id` - (Required) User to associate the application with.

- `username` - (Required) The username to use for the app user.

- `password` - (Optional) The password to use.

- `profile` - (Optional) The JSON profile of the App User.

- `retain_assignment` - (Optional) Retain the user association on destroy. If set to true, the resource will be removed from state but not from the Okta app.

## Attributes Reference

- `id` - The ID of the app user.

## Import

An Application User can be imported via the Okta ID.

```
$ terraform import okta_app_user.example <app id>/<user id>
```
