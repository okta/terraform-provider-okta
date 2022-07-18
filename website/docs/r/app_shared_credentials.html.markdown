---
layout: 'okta'
page_title: 'Okta: okta_app_shared_credentials'
sidebar_current: 'docs-okta-resource-app-shared-credentials'
description: |-
    Creates a SWA shared credentials app.
---

# okta_app_shared_credentials

This resource allows you to create and configure SWA shared credentials app.

## Example Usage

```hcl
resource "okta_app_shared_credentials" "example" {
  label                            = "Example App"
  status                           = "ACTIVE"
  button_field                     = "btn-login"
  username_field                   = "txtbox-username"
  password_field                   = "txtbox-password"
  url                              = "https://example.com/login.html"
  redirect_url                     = "https://example.com/redirect_url"
  checkbox                         = "checkbox_red"
  user_name_template               = "user.firstName"
  user_name_template_type          = "CUSTOM"
  user_name_template_suffix        = "hello"
  shared_password                  = "sharedpass"
  shared_username                  = "sharedusername"
  accessibility_self_service       = true
  accessibility_error_redirect_url = "https://example.com/redirect_url_1"
  accessibility_login_redirect_url = "https://example.com/redirect_url_2"
  auto_submit_toolbar              = true
  hide_ios                         = true
}
```

## Argument Reference

The following arguments are supported:

- `accessibility_error_redirect_url` - (Optional) Custom error page URL.

- `accessibility_login_redirect_url` - (Optional) Custom login page for this application.

- `accessibility_self_service` - (Optional) Enable self-service. By default, it is `false`.

- `admin_note` - (Optional) Application notes for admins.

- `app_links_json` - (Optional) Displays specific appLinks for the app. The value for each application link should be boolean.

- `auto_submit_toolbar` - (Optional) Display auto submit toolbar.

- `button_field` - (Optional) CSS selector for the Sign-In button in the sign-in form.

- `checkbox` - (Optional) CSS selector for the checkbox.

- `enduser_note` - (Optional) Application notes for end users.

- `groups` - (Optional) Groups associated with the application. See `okta_app_group_assignment` for a more flexible approach.
  - `DEPRECATED`: Please replace usage with the `okta_app_group_assignments` (or `okta_app_group_assignment`) resource.

- `hide_ios` - (Optional) Do not display application icon on mobile app.

- `hide_web` - (Optional) Do not display application icon to users.

- `label` - (Required) The Application's display name.

- `logo` - (Optional) Local file path to the logo. The file must be in PNG, JPG, or GIF format, and less than 1 MB in size.

- `password_field` - (Optional) CSS selector for the Password field in the sign-in form.

- `preconfigured_app` - (Optional) name of application from the Okta Integration Network, if not included a custom app will be created.

- `redirect_url` - (Optional) Redirect URL. If going to the login page URL redirects to another page, then enter that URL here.

- `shared_password` - (Optional) Shared password, required for certain schemes.

- `shared_username` - (Optional) Shared username, required for certain schemes.

- `skip_groups` - (Optional) Indicator that allows the app to skip `groups` sync (it's also can be provided during import). Default is `false`.

- `skip_users` - (Optional) Indicator that allows the app to skip `users` sync (it's also can be provided during import). Default is `false`.

- `status` - (Optional) The status of the application, by default, it is `"ACTIVE"`.

- `url` - (Optional) The URL of the sign-in page for this app.

- `url_regex` - (Optional) A regular expression that further restricts url to the specified regular expression.

- `user_name_template` - (Optional) Username template. Default: `"${source.login}"`

- `user_name_template_push_status` - (Optional) Push username on update. Valid values: `"PUSH"` and `"DONT_PUSH"`.

- `user_name_template_suffix` - (Optional) Username template suffix.

- `user_name_template_type` - (Optional) Username template type. Default: `"BUILT_IN"`.

- `username_field` - (Optional) CSS selector for the username field.

- `users` - (Optional) The users assigned to the application. See `okta_app_user` for a more flexible approach.
  - `DEPRECATED`: Please replace usage with the `okta_app_user` resource.

## Attributes Reference

- `id` - ID of an app.

- `name` - Name assigned to the application by Okta.

- `sign_on_mode` - Sign-on mode of the application.

- `logo_url` - Direct link of application logo.

- `sign_on_mode` - Authentication mode of app.

## Timeouts

-> See [here](https://developer.okta.com/todo) for Considerations when Syncing Users/Groups

The `timeouts` block allows you to specify timeouts for certain actions: 

- `create` - (Defaults to no timeout) Used when creating the App with synced Users/Groups.

- `update` - (Defaults to no timeout) Used when updating the App with synced Users/Groups.

- `read` - (Defaults to no timeout) Used when reading the App with synced Users/Groups.

## Import

Okta SWA Shared Credentials App can be imported via the Okta ID.

```
$ terraform import okta_app_shared_credentials.example &#60;app id&#62;
```

It's also possible to import app without groups or/and users. In this case ID may look like this:

```
$ terraform import okta_app_basic_auth.example &#60;app id&#62;/skip_users

$ terraform import okta_app_basic_auth.example &#60;app id&#62;/skip_users/skip_groups

$ terraform import okta_app_basic_auth.example &#60;app id&#62;/skip_groups
```

