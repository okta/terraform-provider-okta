---
page_title: "okta_identity_source_session Resource - terraform-provider-okta"
description: |-
  Manages an Okta Identity Source Session resource.
---

# okta_identity_source_session (Resource)

Manages an Okta Identity Source Session.

## Example Usage

```hcl
resource "okta_identity_source_session" "example" {
  identity_source_id = "<identity-source-id>"

  # Optional fields
  # status = "ACTIVE"
}
```

## Schema

### Required

- `identity_source_id` (String) ID of the identity source

### Optional

- `status` (String) The current status of the identity source session

### Read-Only

- `id` (String) The unique identifier for the resource.
- `created` (String) The timestamp when the identity source session was created
- `import_type` (String) The type of import.  All imports are `INCREMENTAL` imports.
- `last_updated` (String) The timestamp when the identity source session was created

## Import

Import using `{identity_source_id}/{id}`:

```shell
terraform import okta_identity_source_session.example <identity_source_id> <id>
```
