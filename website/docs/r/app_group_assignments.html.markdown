---
layout: 'okta'
page_title: 'Okta: okta_app_group_assignments'
sidebar_current: 'docs-okta-resource-app-group-assignment'
description: |-
  Assigns groups to an application.
---

# okta_app_group_assignments

Assigns groups to an application.

This resource allows you to create multiple App Group assignments.

## Example Usage

```hcl
resource "okta_app_group_assignments" "example" {
  app_id   = "<app id>"
  group {
    id = "<group id>"
    priority = 1
  }
  group {
    id = "<another group id>"
    priority = 2
    profile = jsonencode({"application profile field": "application profile value"})
  }
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

~> **IMPORTANT:** When using `okta_app_group_assignments` it is expected to manage ALL group assignments for the target application.

## Argument Reference

The following arguments are supported:

- `app_id` - (Required) The ID of the application to assign a group to.

- `group` - (Required) A group to assign the app to.

    - `id` - ID of the group to assign.

    - `profile` - (Optional) JSON document containing [application profile](https://developer.okta.com/docs/reference/api/apps/#profile-object)

    - `priority` - (Optional) Priority of group assignment



## Attributes Reference


## Import

An application's group assignments can be imported via `app_id`.

```
$ terraform import okta_app_group_assignments.example <app_id>
```
