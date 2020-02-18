---
layout: "okta"
page_title: "Okta: okta_app_user"
sidebar_current: "docs-okta-resource-app-user"
description: |-
  Creates an Application User.
---

# okta_app_user

Creates an Application User.

This resource allows you to create and configure an Application User.

__When using this resource, make sure to add the following `lifefycle` argument to the application resource you are assigning to:__

```hcl
lifecycle {
  ignore_changes = ["users"]
}
```

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

* `app_id` - (Required) App to associate user with.

* `user_id` - (Required) User to associate the application with.

* `username` - (Required) The username to use for the app user.

* `password` - (Optional) The password to use.

* `profile` - (Optional) The JSON profile of the App User.

## Attributes Reference

* `id` - The ID of the app user.

## Import

An Application User can be imported via the Okta ID.

```
$ terraform import okta_app_user.example <app id>/<user id>
```
