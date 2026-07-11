---
page_title: "Resource: okta_inline_hook"
description: |-
  Manages an okta_inline_hook resource.
---

# Resource: okta_inline_hook

Manages an okta_inline_hook resource.

## Example Usage

```terraform
resource "okta_inline_hook" "example" {
  name    = "example"
  type    = "com.okta.oauth2.tokens.transform"
  version = "1.0.0"

  channel {
    type    = "HTTP"
    version = "1.0.0"
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The display name of the inline hook
- `type` - (Optional) One of the inline hook types
- `version` - (Optional) Version of the inline hook type.
- `status` - (Optional) Status
- `channel` - (Optional) Channel
  - `type` - (Optional) Type
  - `version` - (Optional) Version of the inline hook type.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

- `id` - The ID of the resource.
- `created` - Date of the inline hook creation
- `last_updated` - Date of the last inline hook update

## Import

Import is supported using the following syntax:

```shell
terraform import okta_inline_hook.example <hook_id>
```
