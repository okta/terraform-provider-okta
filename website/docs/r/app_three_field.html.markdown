---
layout: 'okta'
page_title: 'Okta: okta_app_three_field'
sidebar_current: 'docs-okta-resource-app-three-field'
description: |-
  Creates a Three Field Application.
---

# okta_app_three_field

Creates a Three Field Application.

This resource allows you to create and configure a Three Field Application.

## Example Usage

```hcl
resource "okta_app_three_field" "example" {
  label                = "Example App"
  sign_on_url          = "https://example.com/login.html"
  sign_on_redirect_url = "https://example.com"
  reveal_password      = true
  credentials_scheme   = "EDIT_USERNAME_AND_PASSWORD"
}
```

## Argument Reference

The following arguments are supported:

- `label` - (Required) The display name of the Application.

- `button_selector` - (Required) Login button field CSS selector.

- `password_selector` - (Required) Login password field CSS selector.

- `username_selector` - (Required) Login username field CSS selector.

- `extra_field_selector` - (Required) Extra field CSS selector.

- `extra_field_value` - (Required) Value for extra form field.

- `url` - (Required) Login URL.

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

A Three Field App can be imported via the Okta ID.

```
$ terraform import okta_app_three_field.example <app id>
```
