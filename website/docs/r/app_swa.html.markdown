---
layout: 'okta'
page_title: 'Okta: okta_app_swa'
sidebar_current: 'docs-okta-resource-app-swa'
description: |-
  Creates an SWA Application.
---

# okta_app_swa

Creates an SWA Application.

This resource allows you to create and configure an SWA Application.

## Example Usage

```hcl
resource "okta_app_swa" "example" {
  label          = "example"
  button_field   = "btn-login"
  password_field = "txtbox-password"
  username_field = "txtbox-username"
  url            = "https://example.com/login.html"
}
```

## Argument Reference

The following arguments are supported:

- `label` - (Required) The display name of the Application.

- `button_field` - (Required) Login button field.

- `preconfigured_app` - (Optional) name of application from the Okta Integration Network, if not included a custom app will be created.

- `password_field` - (Optional) Login password field.

- `username_field` - (Optional) Login username field.

- `url` - (Optional) Login URL.

- `url_regex` - (Optional) A regex that further restricts URL to the specified regex.

- `users` - (Optional) The users assigned to the application. See `okta_app_user` for a more flexible approach.

- `groups` - (Optional) Groups associated with the application. See `okta_app_group_assignment` for a more flexible approach.

- `status` - (Optional) Status of application. By default, it is `"ACTIVE"`.

- `accessibility_self_service` - (Optional) Enable self-service. By default, it is `false`.

- `accessibility_error_redirect_url` - (Optional) Custom error page URL.

- `auto_submit_toolbar` - (Optional) Display auto submit toolbar.

- `hide_ios` - (Optional) Do not display application icon on mobile app.

- `hide_web` - (Optional) Do not display application icon to users.

## Attributes Reference

- `name` - Name assigned to the application by Okta.

- `sign_on_mode` - Sign-on mode of application.

- `user_name_template` - The default username assigned to each user.

- `user_name_template_type` - The Username template type.

## Import

Okta SWA App can be imported via the Okta ID.

```
$ terraform import okta_app_swa.example <app id>
```
