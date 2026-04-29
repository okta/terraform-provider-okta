---
page_title: "Resource: okta_identity_source_group"
description: |-
  Manages an Okta Identity Source Group resource.
---

# Resource: okta_identity_source_group

Manages a group within an Okta Identity Source. Groups created here are staged for import via a session; use [`okta_identity_source_session_import`](identity_source_session_import.md) to trigger the actual sync into Okta.

## Example Usage

```terraform
resource "okta_identity_source_group" "example" {
  identity_source_id = "<identity-source-id>"
  external_id        = "GRPEXT123456EXAMPLE"

  profile {
    display_name = "Engineering"
    description  = "Engineering team group"
  }
}
```

## Schema

### Required

- `identity_source_id` (String) ID of the identity source. Forces replacement when changed.

### Optional

- `external_id` (String) The external ID of the identity source group as defined in the upstream identity provider.
- `profile` (Block, Optional) Display attributes for the group. (see [below for nested schema](#nested-schema-for-profile))

### Read-Only

- `id` (String) The unique identifier for the resource (server-assigned on create).

### Nested Schema for `profile`

Optional:

- `display_name` (String) Name of the group.
- `description` (String) Description of the group.

## Import

Import using `{identity_source_id}/{id}`:

```shell
terraform import okta_identity_source_group.example <identity-source-id>/<group-id>
```
