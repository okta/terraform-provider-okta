---
layout: "okta"
page_title: "Okta: okta_app_auto_login"
sidebar_current: "docs-okta-resource-app-auto-login"
description: |-
  Creates an Auto Login Okta Application.
---

# okta\_app\_auto\_login

Creates an Auto Login Okta Application.

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

## Argument Reference

The following arguments are supported:

* `label` - (Required) The Application's display name.
* `status` - (Optional) The status of the application, by default it is `"ACTIVE"`.
* `preconfigured_app` - (Optional) Tells Okta to use an existing application in their application catalog, as opposed to a custom application.

## Attributes Reference

* `name` - Name assigned to the application by Okta.
* `sign_on_mode` - Sign on mode of application.

## Import

Okta Auto Login App can be imported via the Okta ID.

```
$ terraform import okta_app_auto_login.example <app id>
```
