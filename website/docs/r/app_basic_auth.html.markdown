---
layout: "okta"
page_title: "Okta: okta_app_basic_auth"
sidebar_current: "docs-okta-resource-app-basic-auth"
description: |-
  Creates a Basic Auth Application.
---

# okta_app_basic_auth

Creates a Bsaic Auth Application.

This resource allows you to create and configure a Basic Auth Application.

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

* `label` - (Required) The Application's display name.

* `url` - (Optional) The URL of the sign-in page for this app.

* `auth_url` - (Optional) The URL of the authenticating site for this app.

## Attributes Reference

* `id` - ID of the Application.

* `label` - The Application's display name.

* `url` - The URL of the sign-in page for basic auth app.

* `auth_url` - The URL of the authenticating site for basic auth app.

## Import

A Basic Auth App can be imported via the Okta ID.

```
$ terraform import okta_app_basic_auth.example <app id>
```
