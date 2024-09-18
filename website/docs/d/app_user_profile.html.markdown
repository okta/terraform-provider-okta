---
layout: 'okta'
page_title: 'Okta: okta_app_user_profile'
sidebar_current: 'docs-okta-datasource-app-user-profile'
description: |-
Get the application profile for a user assigned to an Okta application.
---


# okta_app_user_profile

Use this data source to get the app profile for a user assigned to the given Okta application (by ID). This allows you to retrieve
custom app profile attributes for use within Terraform.

## Example Usage

```hcl
data "okta_app" "example" {
  label = "Example App"

  skip_groups = true
  skip_users  = true
}

data "okta_user" "example" {
  user_id = "00u22mtxlrJ8YkzXQ357"
}

data "okta_app_user_profile" "test" {
  app_id  = data.okta_app.example.id
  user_id = data.okta_user.example.id
}
```

## Argument Reference

- `app_id` - (Required) The ID of the Okta application you want to retrieve the user profile for.

- `user_id` - (Required) The ID of the user you want to retrieve the user profile for.

## Attribute Reference

- `app_id` - ID of the application.

- `user_id` - ID of the user.

- `profile` - The profile for the user.

