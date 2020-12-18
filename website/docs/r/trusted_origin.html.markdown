---
layout: 'okta'
page_title: 'Okta: okta_trusted_origin'
sidebar_current: 'docs-okta-resource-trusted-origin'
description: |-
  Creates a Trusted Origin.
---

# okta_trusted_origin

Creates a Trusted Origin.

This resource allows you to create and configure a Trusted Origin.

## Example Usage

```hcl
resource "okta_trusted_origin" "example" {
  name   = "example"
  origin = "https://example.com"
  scopes = ["CORS"]
}
```

## Argument Reference

The following arguments are supported:

- `active` - (Optional) Whether the Trusted Origin is active or not - can only be issued post-creation. By default it is 'true'.

- `name` - (Required) Unique name for this trusted origin.

- `origin` - (Required) Unique origin URL for this trusted origin.

- `scopes` - (Required) Scopes of the Trusted Origin - can be `"CORS"` and/or `"REDIRECT"`.

## Attributes Reference

- `id` - The ID of the Trusted Origin.

## Import

A Trusted Origin can be imported via the Okta ID.

```
$ terraform import okta_trusted_origin.example <trusted origin id>
```
