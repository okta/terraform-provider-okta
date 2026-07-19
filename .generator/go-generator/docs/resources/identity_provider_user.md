---
page_title: "okta_identity_provider_user Resource - terraform-provider-okta"
description: |-
  Manages an Okta Identity Provider User resource.
---

# okta_identity_provider_user (Resource)

Manages an Okta Identity Provider User.

## Example Usage

```hcl
resource "okta_identity_provider_user" "example" {
  idp_id = "<idp-id>"

  # Optional fields
  # profile = "<profile>"
}
```

## Schema

### Required

- `idp_id` (String) ID of the identity provider

### Optional

- `profile` (String) IdP-specific profile for the user.  IdP user profiles are IdP-specific but may be customized by the Profile Editor in the Admin Console.  > **Note:** Okta variable names have reserved characters th...

### Read-Only

- `id` (String) The unique identifier for the resource.
- `embedded` (String) Embedded resources related to the IdP user
- `created` (String) Timestamp when the object was created
- `external_id` (String) Unique IdP-specific identifier for the user
- `last_updated` (String) Timestamp when the object was last updated

## Import

Import using `{idp_id}/{id}`:

```shell
terraform import okta_identity_provider_user.example <idp_id> <id>
```
