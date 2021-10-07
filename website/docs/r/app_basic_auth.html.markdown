---
layout: "okta"
page_title: "Okta: okta_app_basic_auth"
sidebar_current: "docs-okta-resource-app-basic-auth"
description: |-
  Creates a Basic Auth Application.
---

# okta_app_basic_auth

This resource allows you to create and configure a Basic Auth Application.

## Example Usage

```hcl
resource "okta_app_basic_auth" "example" {
  label    = "Example"
  url      = "https://example.com/login.html"
  auth_url = "https://example.com/auth.html"
}
```

## Argument Reference

The following arguments are supported:

- `label` - (Required) The Application's display name.

- `url` - (Required) The URL of the sign-in page for this app.

- `auth_url` - (Required) The URL of the authenticating site for this app.

- `users` - (Optional) Users associated with the application.
  - `DEPRECATED`: Please replace usage with the `okta_app_user` resource.

- `groups` - (Optional) Groups associated with the application.
  - `DEPRECATED`: Please replace usage with the `okta_app_group_assignments` (or `okta_app_group_assignment`) resource.

- `status` - (Optional) Status of application. (`"ACTIVE"` or `"INACTIVE"`).

- `hide_web` - (Optional) Do not display application icon to users.

- `hide_ios` - (Optional) Do not display application icon on mobile app.

- `auto_submit_toolbar` - (Optional) Display auto submit toolbar.

- `accessibility_error_redirect_url` - (Optional) Custom error page URL.

- `accessibility_login_redirect_url` - (Optional) Custom login page for this application.

- `accessibility_self_service` - (Optional) Enable self-service. By default, it is `false`.

- `logo` - (Optional) Local path to the logo. The file must be in PNG, JPG, or GIF format, and less than 1 MB in size.

- `admin_note` - (Optional) Application notes for admins.

- `enduser_note` - (Optional) Application notes for end users.

- `skip_users` - (Optional) Indicator that allows the app to skip `users` sync (it's also can be provided during import). Default is `false`. 

- `skip_groups` - (Optional) Indicator that allows the app to skip `groups` sync (it's also can be provided during import). Default is `false`.

## Attributes Reference

- `id` - ID of the Application.

- `label` - The Application's display name.

- `url` - The URL of the sign-in page for basic auth app.

- `auth_url` - The URL of the authenticating site for basic auth app.

- `logo_url` - Direct link of application logo.

## Import

A Basic Auth App can be imported via the Okta ID.

```
$ terraform import okta_app_basic_auth.example <app id>
```

It's also possible to import app without groups or/and users. In this case ID may look like this:

```
$ terraform import okta_app_basic_auth.example <app id>/skip_users

$ terraform import okta_app_basic_auth.example <app id>/skip_users/skip_groups

$ terraform import okta_app_basic_auth.example <app id>/skip_groups
```
