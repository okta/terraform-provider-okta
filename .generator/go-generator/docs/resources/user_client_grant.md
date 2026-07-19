---
page_title: "okta_user_client_grant Resource - terraform-provider-okta"
description: |-
  Manages an Okta User Client Grant resource.
---

# okta_user_client_grant (Resource)

Manages an Okta User Client Grant.

## Example Usage

```hcl
resource "okta_user_client_grant" "example" {
  user_id = "<user-id>"
  client_id = "<client-id>"
}
```

## Schema

### Required

- `user_id` (String) ID of the user
- `client_id` (String) ID of the OAuth client

### Read-Only

- `id` (String) The unique identifier for the resource.

## Import

Import using `{user_id}/{client_id}/{id}`:

```shell
terraform import okta_user_client_grant.example <user_id> <client_id> <id>
```
