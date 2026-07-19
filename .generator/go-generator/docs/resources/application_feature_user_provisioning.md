---
page_title: "okta_application_feature_user_provisioning Resource - terraform-provider-okta"
description: |-
  Manages an Okta Application Feature User Provisioning resource.
---

# okta_application_feature_user_provisioning (Resource)

Manages an Okta Application Feature User Provisioning.

## Example Usage

```hcl
resource "okta_application_feature_user_provisioning" "example" {
  app_id = "<app-id>"
  name = "Example Name"

  # Optional fields
  # status = "ACTIVE"
  # status = "ACTIVE"
  # change = "<change>"
  # seed = "<seed>"
  # status = "ACTIVE"
}
```

## Schema

### Required

- `app_id` (String) ID of the parent application
- `name` (String) Discriminator field identifying the variant type. Must be set to \

### Optional

- `status` (String) Status
- `status` (String) Status
- `change` (String) Determines whether a change in a user's password also updates the user's password in the app
- `seed` (String) Determines whether the generated password is the user's Okta password or a randomly generated password
- `status` (String) Status
- `status` (String) Status
- `status` (String) Status

### Read-Only

- `id` (String) The unique identifier for the resource.
- `description` (String) Description of the feature

## Import

Import using `{app_id}/{id}`:

```shell
terraform import okta_application_feature_user_provisioning.example <app_id> <id>
```
