---
layout: 'okta'
page_title: 'Okta: okta_app_bookmark'
sidebar_current: 'docs-okta-resource-app-bookmark'
description: |-
  Creates a Bookmark Application.
---

# okta_app_bookmark

This resource allows you to create and configure a Bookmark Application.

## Example Usage

```hcl
resource "okta_app_bookmark" "example" {
  label  = "Example"
  url    = "https://example.com"
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

- `enduser_note` - (Optional) Application notes for end users.

- `groups` - (Optional) Groups associated with the application.
  - `DEPRECATED`: Please replace usage with the `okta_app_group_assignments` (or `okta_app_group_assignment`) resource.

- `hide_ios` - (Optional) Do not display application icon on mobile app.

- `hide_web` - (Optional) Do not display application icon to users.

- `label` - (Required) The Application's display name.

- `logo` - (Optional) Local file path to the logo. The file must be in PNG, JPG, or GIF format, and less than 1 MB in size.

- `request_integration` - (Optional) Would you like Okta to add an integration for this app?

- `skip_groups` - (Optional) Indicator that allows the app to skip `groups` sync (it's also can be provided during import). Default is `false`.

- `skip_users` - (Optional) Indicator that allows the app to skip `users` sync (it's also can be provided during import). Default is `false`.

- `status` - (Optional) Status of application. (`"ACTIVE"` or `"INACTIVE"`).

- `url` - (Optional) The URL of the bookmark.

- `users` - (Optional) Users associated with the application.
  - `DEPRECATED`: Please replace usage with the `okta_app_user` resource.

## Attributes Reference

- `id` - ID of the Application.

- `label` - The Application's display name.

- `url` - The URL of the bookmark.

- `logo_url` - Direct link of application logo.

## Timeouts

-> See [here](https://developer.okta.com/todo) for Considerations when Syncing Users/Groups

The `timeouts` block allows you to specify timeouts for certain actions: 

- `create` - (Defaults to no timeout) Used when creating the App with synced Users/Groups.

- `update` - (Defaults to no timeout) Used when updating the App with synced Users/Groups.

- `read` - (Defaults to no timeout) Used when reading the App with synced Users/Groups.

## Import

A Bookmark App can be imported via the Okta ID.

```
$ terraform import okta_app_bookmark.example &#60;app id&#62;
```

It's also possible to import app without groups or/and users. In this case ID may look like this:

```
$ terraform import okta_app_basic_auth.example &#60;app id&#62;/skip_users

$ terraform import okta_app_basic_auth.example &#60;app id&#62;/skip_users/skip_groups

$ terraform import okta_app_basic_auth.example &#60;app id&#62;/skip_groups
```
