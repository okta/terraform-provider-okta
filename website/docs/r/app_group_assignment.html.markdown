---
layout: 'okta'
page_title: 'Okta: okta_app_group_assignment'
sidebar_current: 'docs-okta-resource-app-group-assignment'
description: |-
  Assigns a group to an application.
---

# okta_app_group_assignment

Assigns a group to an application.

This resource allows you to create an App Group assignment.

## Example Usage

```hcl
resource "okta_app_group_assignment" "example" {
  app_id   = "<app id>"
  group_id = "<group id>"
  profile = <<JSON
{
  "<app_profile_field>": "<value>"
}
JSON
}

```

!> **NOTE** When using this resource in conjunction with other application resources (e.g. `okta_app_oauth`) it is advisable to add the following `lifecycle` argument to the associated `app_*` resources to prevent the groups being unassigned on subsequent runs:

```hcl
resource "okta_app_oauth" "app" {
  //...
  lifecycle {
     ignore_changes = [groups]
  }
}
```

~> **Important:** When the `app_group_assignment` is retained, by setting `retain_assignment` to `true`, it is no longer managed by Terraform after it is destroyed. To truly delete the assignment, you will need to remove it either through the Okta Console or API. This argument exists for the use case where the same group is assigned in multiple places in order to prevent a single destruction removing all of them.

## Argument Reference

The following arguments are supported:

- `app_id` - (Required) The ID of the application to assign a group to.

- `group_id` - (Required) The ID of the group to assign the app to.

- `profile` - (Optional) JSON document containing [application profile](https://developer.okta.com/docs/reference/api/apps/#profile-object)

- `retain_assignment` - (Optional) Retain the group assignment on destroy. If set to true, the resource will be removed from state but not from the Okta app.

## Attributes Reference

- `id` - ID of the group assignment.

## Import

An application group assignment can be imported via the `app_id` and the `group_id`.

```
$ terraform import okta_app_group_assignment.example <app_id>/<group_id>
```
