---
layout: 'okta'
page_title: 'Okta: okta_app_bookmark'
sidebar_current: 'docs-okta-resource-app-bookmark'
description: |-
  Creates a Bookmark Application.
---

# okta_app_bookmark

Creates a Bookmark Application.

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

- `label` - (Required) The Application's display name.

- `url` - (Optional) The URL of the bookmark.

- `request_integration` - (Optional) Would you like Okta to add an integration for this app?

- `users` - (Optional) Users associated with the application.

- `groups` - (Optional) Groups associated with the application.

- `status` - (Optional) Status of application. (`"ACTIVE"` or `"INACTIVE"`).

- `hide_web` - (Optional) Do not display application icon to users.

- `hide_ios` - (Optional) Do not display application icon on mobile app.

- `auto_submit_toolbar` - (Optional) Display auto submit toolbar.

## Attributes Reference

- `id` - ID of the Application.

- `label` - The Application's display name.

- `url` - The URL of the bookmark.

## Import

A Bookmark App can be imported via the Okta ID.

```
$ terraform import okta_app_bookmark.example <app id>
```
