---
layout: "okta"
page_title: "Okta: okta_app_group_assignment"
sidebar_current: "docs-okta-resource-app-group-assignment"
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

## Argument Reference

The following arguments are supported:

* `app_id` - (Required) The ID of the application to assign a group to.

* `group_id` - (Required) The ID of the group to assign the app to.

* `profile` - (Optional) JSON document containing [application profile](https://developer.okta.com/docs/reference/api/apps/#profile-object)

## Attributes Reference

* `id` - ID of the group assignment.

## Import

An application group assignment can be imported via assignment ID.

```
$ terraform import okta_app_group_assignment.example <id>
```
