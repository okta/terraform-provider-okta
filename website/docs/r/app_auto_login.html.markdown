---
layout: 'okta'
page_title: 'Okta: okta_app_auto_login'
sidebar_current: 'docs-okta-resource-app-auto-login'
description: |-
  Creates an Auto Login Okta Application.
---

# okta_app_auto_login

This resource allows you to create and configure an Auto Login Okta Application.

## Example Usage

```hcl
resource "okta_app_auto_login" "example" {
  label                = "Example App"
  sign_on_url          = "https://example.com/login.html"
  sign_on_redirect_url = "https://example.com"
  reveal_password      = true
  credentials_scheme   = "EDIT_USERNAME_AND_PASSWORD"
}
```

#### Pre-configured application
```hcl
resource "okta_app_auto_login" "example" {
  label             = "Google Example App"
  status            = "ACTIVE"
  preconfigured_app = "google"
  app_settings_json = <<JSON
{
    "domain": "okta",
    "afwOnly": false
}
JSON
}
```

## Argument Reference

The following arguments are supported:

- `label` - (Required) The Application's display name.

- `status` - (Optional) The status of the application, by default, it is `"ACTIVE"`.

- `preconfigured_app` - (Optional) Tells Okta to use an existing application in their application catalog, as opposed to a custom application.

- `sign_on` - (Required) App login page URL

- `sign_on_redirect_url` - (Optional) Redirect URL; if going to the login page URL redirects to another page, then enter that URL here

- `credentials_scheme` - (Optional) One of: `"EDIT_USERNAME_AND_PASSWORD"`, `"ADMIN_SETS_CREDENTIALS"`, `"EDIT_PASSWORD_ONLY"`, `"EXTERNAL_PASSWORD_SYNC"`, or `"SHARED_USERNAME_AND_PASSWORD"`.

- `shared_username` - (Optional) Shared username, required for certain schemes

- `shared_password` - (Optional) Shared password, required for certain schemes

- `user_name_template` - (Optional) Username template. Default: `"${source.login}"`

- `user_name_template_suffix` - (Optional) Username template suffix.

- `user_name_template_type` - (Optional) Username template type. Default: `"BUILT_IN"`.

- `user_name_template_push_status` - (Optional) Push username on update. Valid values: `"PUSH"` and `"DONT_PUSH"`. Default: `"DONT_PUSH"`.

- `hide_web` - (Optional) Do not display application icon to users.

- `hide_ios` - (Optional) Do not display application icon on mobile app.

- `auto_submit_toolbar` - (Optional) Display auto submit toolbar.

- `accessibility_error_redirect_url` - (Optional) Custom error page URL.

- `accessibility_login_redirect_url` - (Optional) Custom login page for this application.

- `accessibility_self_service` - (Optional) Enable self-service. By default, it is `false`.

- `users` - (Optional) The users assigned to the application. See `okta_app_user` for a more flexible approach.
  - `DEPRECATED`: Please replace usage with the `okta_app_user` resource.

- `groups` - (Optional) Groups associated with the application. See `okta_app_group_assignment` for a more flexible approach.
  - `DEPRECATED`: Please replace usage with the `okta_app_group_assignments` (or `okta_app_group_assignment`) resource.

- `logo` - (Optional) Local file path to the logo. The file must be in PNG, JPG, or GIF format, and less than 1 MB in size.

- `admin_note` - (Optional) Application notes for admins.

- `enduser_note` - (Optional) Application notes for end users.

- `app_settings_json` - (Optional) Application settings in JSON format.

- `skip_users` - (Optional) Indicator that allows the app to skip `users` sync (it's also can be provided during import). Default is `false`.

- `skip_groups` - (Optional) Indicator that allows the app to skip `groups` sync (it's also can be provided during import). Default is `false`.

## Attributes Reference

- `name` - Name assigned to the application by Okta.
  
- `sign_on_mode` - Sign-on mode of the application.

- `logo_url` - Direct link of application logo.

## Import

Okta Auto Login App can be imported via the Okta ID.

```
$ terraform import okta_app_auto_login.example <app id>
```

It's also possible to import app without groups or/and users. In this case ID may look like this:

```
$ terraform import okta_app_basic_auth.example <app id>/skip_users

$ terraform import okta_app_basic_auth.example <app id>/skip_users/skip_groups

$ terraform import okta_app_basic_auth.example <app id>/skip_groups
```
