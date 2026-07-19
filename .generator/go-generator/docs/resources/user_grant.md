---
page_title: "okta_user_grant Resource - terraform-provider-okta"
description: |-
  Manages an Okta User Grant resource.
---

# okta_user_grant (Resource)

Manages an Okta User Grant.

## Example Usage

```hcl
resource "okta_user_grant" "example" {
  user_id = "<user-id>"
  issuer = "<issuer>"
  scope_id = "<scope-id>"
}
```

## Schema

### Required

- `user_id` (String) ID of the user
- `issuer` (String) The issuer of your org authorization server. This is typically your Okta domain.
- `scope_id` (String) The name of the [Okta scope](https://developer.okta.com/docs/api/oauth2/#oauth-20-scopes) for which consent is granted

### Read-Only

- `id` (String) The unique identifier for the resource.
- `embedded` (String) Embedded resources related to the Grant
- `client_id` (String) Client ID of the app integration
- `created` (String) Timestamp when the object was created
- `created_by` (String) User that created the object
- `last_updated` (String) Timestamp when the object was last updated
- `source` (String) User type source that granted consent
- `status` (String) Status

## Import

Import using `{user_id}/{id}`:

```shell
terraform import okta_user_grant.example <user_id> <id>
```
