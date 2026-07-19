---
page_title: "okta_user_factor_token Resource - terraform-provider-okta"
description: |-
  Manages an Okta User Factor Token resource.
---

# okta_user_factor_token (Resource)

Manages an Okta User Factor Token.

## Example Usage

```hcl
resource "okta_user_factor_token" "example" {
  user_id = "<user-id>"
  factor_type = "<factor_type>"

  # Optional fields
  # profile = "<profile>"
  # provider = "<provider>"
  # verify = "<verify>"
}
```

## Schema

### Required

- `user_id` (String) ID of the user
- `factor_type` (String) Discriminator field identifying the variant type. Must be set to \

### Optional

- `profile` (String) Specific attributes related to the factor
- `provider` (String) Provider for the factor. Each provider can support a subset of factor types.
- `verify` (String) Verify

### Read-Only

- `id` (String) The unique identifier for the resource.
- `embedded` (String) Embedded
- `created` (String) Timestamp when the factor was enrolled
- `last_updated` (String) Timestamp when the factor was last updated
- `status` (String) Status of the factor
- `vendor_name` (String) Name of the factor vendor. This is usually the same as the provider except for On-Prem MFA, which depends on admin settings.

## Import

Import using `{user_id}/{id}`:

```shell
terraform import okta_user_factor_token.example <user_id> <id>
```
