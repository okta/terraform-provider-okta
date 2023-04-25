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

- `authentication_policy` - (Optional) The ID of the associated `app_signon_policy`. If this property is removed from the application the `default` sign-on-policy will be associated with this application.

- `auto_submit_toolbar` - (Optional) Display auto submit toolbar.

- `enduser_note` - (Optional) Application notes for end users.

- `hide_ios` - (Optional) Do not display application icon on mobile app.

- `hide_web` - (Optional) Do not display application icon to users.

- `label` - (Required) The Application's display name.

- `logo` - (Optional) Local file path to the logo. The file must be in PNG, JPG, or GIF format, and less than 1 MB in size.

- `request_integration` - (Optional) Would you like Okta to add an integration for this app?

- `status` - (Optional) Status of application. (`"ACTIVE"` or `"INACTIVE"`).

- `url` - (Optional) The URL of the bookmark.

## Attributes Reference

- `id` - ID of the Application.

- `label` - The Application's display name.

- `url` - The URL of the bookmark.

- `logo_url` - Direct link of application logo.

## Timeouts

The `timeouts` block allows you to specify custom [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions: 

- `create` - Create timeout (default 1 hour).

- `update` - Update timeout (default 1 hour).

- `read` - Read timeout (default 1 hour).

## Import

A Bookmark App can be imported via the Okta ID.

```
$ terraform import okta_app_bookmark.example &#60;app id&#62;
```
