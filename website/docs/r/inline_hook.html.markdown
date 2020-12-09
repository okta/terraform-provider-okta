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
  version = "1.0.0"
  type    = "com.okta.oauth2.tokens.transform"

  channel = {
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

- `name` - (Required) The inline hook display name.

- `version` - (Required) The version of the hook. The currently-supported version is `"1.0.0"`.

- `type` - (Required) The type of hook to create. [See here for supported types](https://developer.okta.com/docs/reference/api/inline-hooks/#supported-inline-hook-types).

- `headers` - (Optional) Map of headers to send along in inline hook request.

- `auth` - (Optional) Authentication required for inline hook request.

  - `key` - (Required) Key to use for authentication, usually the header name, for example `"Authorization"`.
  - `value` - (Required) Authentication secret.
  - `type` - (Optional) Auth type. Currently, the only supported type is `"HEADER"`.

- `channel` - (Required) Details of the endpoint the inline hook will hit.
  - `version` - (Required) Version of the channel. The currently-supported version is `"1.0.0"`.
  - `uri` - (Required) The URI the hook will hit.
  - `type` - (Optional) The type of hook to trigger. Currently, the only supported type is `"HTTP"`.
  - `method` - (Optional) The request method to use. Default is `"POST"`.

## Attributes Reference

- `id` - The ID of the inline hooks.

## Import

An inline hook can be imported via the Okta ID.

```
$ terraform import okta_inline_hook.example <hook id>
```
