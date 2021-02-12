---
layout: 'okta'
page_title: 'Okta: okta_admin_role_targets'
sidebar_current: 'docs-okta-resource-okta-admin-role-targets'
description: |-
  Manages targets for administrator roles.
---

# okta_admin_role_targets

Manages targets for administrator roles.

This resource allows you to define permissions for admin roles into a smaller subset of Groups or Apps within your org.
You can define admin roles to target Groups, Applications, and Application Instances.

```
Note 1: you have to assign a role to a user before creating this resource.

Note 2: You can target a mixture of both App and App Instance targets, but can't assign permissions to manage all 
        instances of an App and then a subset of that same App. For example, you can't specify that an admin has access 
        to manage all instances of a Salesforce app and then also specific configurations of the Salesforce app.
```

## Example Usage

```hcl
resource "okta_admin_role_targets" "example" {
  user_id   = "<user_id>"
  role_type = "APP_ADMIN"
  apps      = ["oidc_client.<app_id>", "facebook"]
}
```

## Argument Reference

The following arguments are supported:

- `user_id` - (Required) ID of the user.

- `role_type` - (Required) Name of the role associated with the user.

- `apps` - (Optional) List of app names (name represents set of app instances) or a combination of app name and app instance ID (like 'salesforce' or 'facebook.0oapsqQ6dv19pqyEo0g3').

- `groups` (Optional) List of group IDs. Conflicts with `apps`.

- `role_id` (Computed) Role ID.

## Attributes Reference

- `role_id` - Role ID

## Import

Okta Admin Role Targets can be imported via the Okta ID.

```
$ terraform import okta_admin_role_targets.example <user id>/<role type>
```
