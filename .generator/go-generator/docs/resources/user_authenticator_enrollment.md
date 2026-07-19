---
page_title: "okta_user_authenticator_enrollment Resource - terraform-provider-okta"
description: |-
  Manages an Okta User Authenticator Enrollment resource.
---

# okta_user_authenticator_enrollment (Resource)

Manages an Okta User Authenticator Enrollment.

## Example Usage

```hcl
resource "okta_user_authenticator_enrollment" "example" {
  user_id = "<user-id>"
}
```

## Schema

### Required

- `user_id` (String) ID of the user

### Read-Only

- `id` (String) The unique identifier for the resource.

## Import

Import using `{user_id}/{id}`:

```shell
terraform import okta_user_authenticator_enrollment.example <user_id> <id>
```
