---
layout: "okta"
page_title: "Okta: okta_inline_hook"
sidebar_current: "docs-okta-resource-inline-hook"
description: |-
  Creates an inline hook.
---

# okta_inline_hook

Creates an inline hook.

This resource allows you to create and configure an inline hook.

## Example Usage

```hcl
resource "okta_inline_hook" "example" {
  name    = "example"
  version = "1.0.1"
  type    = "com.okta.oauth2.tokens.transform"

  channel = {
    type    = "HTTP"
    version = "1.0.0"
    uri     = "https://example.com/test"
    method  = "POST"
  }

  auth = {
    key   = "Authorization"
    type  = "HEADER"
    value = "secret"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The inline hook display name.

* `version` - (Required) The version of the hook.

* `type` - (Required) The type of hook to create. [See here for supported types](https://developer.okta.com/docs/reference/api/inline-hooks/#supported-inline-hook-types).

## Attributes Reference

* `id` - The ID of the inline hooks.

## Import

An inline hook can be imported via the Okta ID.

```
$ terraform import okta_inline_hook.example <hook id>
```
