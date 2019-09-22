---
layout: "okta"
page_title: "Okta: okta_app_oauth"
sidebar_current: "docs-okta-resource-app-auto-login"
description: |-
  Creates an Auto Login Okta Application.
---

# okta_app_oauth

Creates an Auto Login Okta Application.

This resource allows you to create and configure an Auto Login Okta Application.

## Example Usage

```hcl
resource "okta_app_oauth" "example" {
  label                = "Example App"
  sign_on_url          = "https://example.com/login.html"
  sign_on_redirect_url = "https://example.com"
  reveal_password      = true
  credentials_scheme   = "EDIT_USERNAME_AND_PASSWORD"
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) The Application's display name.

* `status` - (Optional) The status of the application, by default it is `"ACTIVE"`.

* `type` - (Required) The type of OAuth application.

* `users` - (Optional) The users assigned to the application. It is recommended not to use this and instead use `okta_app_user`.

* `groups` - (Optional) The groups assigned to the application. It is recommended not to use this and instead use `okta_app_group_assignment`.

* `custom_client_id` - (Optional) This property allows you to set the application's client id.

## Attributes Reference

* `name` - Name assigned to the application by Okta.

* `sign_on_mode` - Sign on mode of application.

* `client_id` - The client ID of the application.

* `client_secret` - The client secret of the application.

## Import

Okta Auto Login App can be imported via the Okta ID.

```
$ terraform import okta_app_oauth.example <app id>
```
