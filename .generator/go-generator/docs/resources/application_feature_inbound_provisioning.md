---
page_title: "okta_application_feature_inbound_provisioning Resource - terraform-provider-okta"
description: |-
  Manages an Okta Application Feature Inbound Provisioning resource.
---

# okta_application_feature_inbound_provisioning (Resource)

Manages an Okta Application Feature Inbound Provisioning.

## Example Usage

```hcl
resource "okta_application_feature_inbound_provisioning" "example" {
  app_id = "<app-id>"
  name = "Example Name"
  expression = "<expression>"
  expression = "<expression>"
  username_format = "Example Username Format"

  # Optional fields
  # allow_partial_match = true
  # auto_activate_new_users = true
  # auto_confirm_exact_match = true
  # auto_confirm_new_users = true
  # auto_confirm_partial_match = true
}
```

## Schema

### Required

- `app_id` (String) ID of the parent application
- `name` (String) Discriminator field identifying the variant type. Must be set to \
- `expression` (String) The import schedule in UNIX cron format
- `expression` (String) The import schedule in UNIX cron format
- `username_format` (String) Determines the username format when users sign in to Okta

### Optional

- `allow_partial_match` (Boolean) Allows user import upon partial matching. Partial matching occurs when the first and last names of an imported user match those of an existing Okta user, even if the username or email attributes do...
- `auto_activate_new_users` (Boolean) If set to `true`, imported new users are automatically activated.
- `auto_confirm_exact_match` (Boolean) If set to `true`, exact-matched users are automatically confirmed on activation. If set to `false`, exact-matched users need to be confirmed manually.
- `auto_confirm_new_users` (Boolean) If set to `true`, imported new users are automatically confirmed on activation. This doesn't apply to imported users that already exist in Okta.
- `auto_confirm_partial_match` (Boolean) If set to `true`, partially matched users are automatically confirmed on activation. If set to `false`, partially matched users need to be confirmed manually.
- `exact_match_criteria` (String) Determines the attribute to match users
- `timezone` (String) The import schedule time zone in Internet Assigned Numbers Authority (IANA) time zone name format
- `timezone` (String) The import schedule time zone in Internet Assigned Numbers Authority (IANA) time zone name format
- `status` (String) Setting status
- `user_name_expression` (String) For `usernameFormat=CUSTOM`, specifies the Okta Expression Language statement for a username format that imported users use to sign in to Okta
- `status` (String) Status

### Read-Only

- `id` (String) The unique identifier for the resource.
- `description` (String) Description of the feature

## Import

Import using `{app_id}/{id}`:

```shell
terraform import okta_application_feature_inbound_provisioning.example <app_id> <id>
```
