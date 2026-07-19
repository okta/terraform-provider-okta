---
page_title: "okta_identity_source_group Resource - terraform-provider-okta"
description: |-
  Manages an Okta Identity Source Group resource.
---

# okta_identity_source_group (Resource)

Manages an Okta Identity Source Group.

## Example Usage

```hcl
resource "okta_identity_source_group" "example" {
  identity_source_id = "<identity-source-id>"

  # Optional fields
  # external_id = "<external-id>"
  # description = "Example description"
  # display_name = "Example Display Name"
}
```

## Schema

### Required

- `identity_source_id` (String) ID of the identity source

### Optional

- `external_id` (String) The external ID of the identity source group
- `description` (String) Description of the group
- `display_name` (String) Name of the group

### Read-Only

- `id` (String) The unique identifier for the resource.

## Import

Import using `{identity_source_id}/{id}`:

```shell
terraform import okta_identity_source_group.example <identity_source_id> <id>
```
