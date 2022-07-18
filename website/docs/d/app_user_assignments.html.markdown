---
layout: 'okta'
page_title: 'Okta: okta_app_user_assignments'
sidebar_current: 'docs-okta-datasource-app-user-assignments'
description: |-
Get a set of users assigned to an Okta application.
---


# okta_app_user_assignments

Use this data source to retrieve the list of users assigned to the given Okta application (by ID).

## Example Usage

```hcl
data "okta_app_user_assignments" "test" {
  id         = okta_app_oauth.test.id
}
```

## Argument Reference

- `id` - (Required) The ID of the Okta application you want to retrieve the groups for.

## Attribute Reference

- `id` - ID of application.

- `users` - List of user IDs assigned to the application.

## Timeouts

-> See [here](https://developer.okta.com/todo) for Considerations when Syncing Users/Groups

The `timeouts` block allows you to specify timeouts for certain actions: 

- `read` - (Defaults to no timeout) Used when reading the App with synced Users/Groups.
