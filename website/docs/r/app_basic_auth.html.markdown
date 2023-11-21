---
layout: "okta"
page_title: "Okta: okta_app_basic_auth"
sidebar_current: "docs-okta-resource-app-basic-auth"
description: |-
  Creates a Basic Auth Application.
---

# okta_app_basic_auth

This resource allows you to create and configure a Basic Auth Application.

-> During an apply if there is change in `status` the app will first be
activated or deactivated in accordance with the `status` change. Then, all
other arguments that changed will be applied.

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

- `accessibility_error_redirect_url` - (Optional) Custom error page URL.

- `accessibility_login_redirect_url` - (Optional) Custom login page for this application.

- `accessibility_self_service` - (Optional) Enable self-service. By default, it is `false`.

- `admin_note` - (Optional) Application notes for admins.

- `app_links_json` - (Optional) Displays specific appLinks for the app. The value for each application link should be boolean.

- `auth_url` - (Required) The URL of the authenticating site for this app.

- `auto_submit_toolbar` - (Optional) Display auto submit toolbar.

- `enduser_note` - (Optional) Application notes for end users.

- `hide_ios` - (Optional) Do not display application icon on mobile app.

- `hide_web` - (Optional) Do not display application icon to users.

- `label` - (Required) The Application's display name.

- `logo` - (Optional) Local path to the logo. The file must be in PNG, JPG, or GIF format, and less than 1 MB in size.

- `status` - (Optional) Status of application. (`"ACTIVE"` or `"INACTIVE"`).

- `url` - (Required) The URL of the sign-in page for this app.

## Attributes Reference

- `id` - ID of the Application.

- `label` - The Application's display name.

- `url` - The URL of the sign-in page for basic auth app.

- `auth_url` - The URL of the authenticating site for basic auth app.

- `logo_url` - Direct link of application logo.

## Timeouts

The `timeouts` block allows you to specify custom [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions: 

- `create` - Create timeout (default 1 hour).

- `update` - Update timeout (default 1 hour).

- `read` - Read timeout (default 1 hour).

## Import

A Basic Auth App can be imported via the Okta ID.

```
$ terraform import okta_app_basic_auth.example &#60;app id&#62;
```
