---
page_title: "okta_inline_hook Resource - terraform-provider-okta"
description: |-
  Manages an Okta Inline Hook resource.
---

# okta_inline_hook (Resource)

Manages an Okta Inline Hook.

## Example Usage

```hcl
resource "okta_inline_hook" "example" {

  # Optional fields
  # type = "<type>"
  # version = "<version>"
  # name = "Example Name"
  # status = "ACTIVE"
  # type = "<type>"
}
```

## Schema

### Optional

- `type` (String) Type
- `version` (String) Version of the inline hook type. The currently supported version is `1.0.0`.
- `name` (String) The display name of the inline hook
- `status` (String) Status
- `type` (String) One of the inline hook types

### Read-Only

- `id` (String) The unique identifier for the resource.
- `created` (String) Date of the inline hook creation
- `last_updated` (String) Date of the last inline hook update
- `version` (String) Version of the inline hook type. The currently supported version is `1.0.0`.
