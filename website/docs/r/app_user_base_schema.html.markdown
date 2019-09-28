---
layout: "okta"
page_title: "Okta: okta_app_user_base_schema"
sidebar_current: "docs-okta-resource-app-user-base-schema"
description: |-
  Manages a base Application User Schema property.
---

# okta_app_user_base_schema

Manages a base Application User Schema property.

This resource allows you to configure a base Application User Schema property.

## Example Usage

```hcl
resource "okta_app_user_base_schema" "example" {
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
$ terraform import okta_app_user_base_schema.example <app id>
```
