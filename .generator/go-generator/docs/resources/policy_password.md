---
page_title: "okta_policy_password Resource - terraform-provider-okta"
description: |-
  Manages an Okta Policy Password resource.
---

# okta_policy_password (Resource)

Manages an Okta Policy Password.

## Example Usage

```hcl
resource "okta_policy_password" "example" {
  type = "<type>"
  name = "Example Name"

  # Optional fields
  # include = "<include>"
  # provider = "<provider>"
  # include = "<include>"
  # description = "Example description"
  # priority = 0
}
```

## Schema

### Required

- `type` (String) Discriminator field identifying the variant type. Must be set to \
- `name` (String) Name of the policy

### Optional

- `include` (List) Include
- `provider` (String) Provider
- `include` (List) Groups to be included
- `description` (String) Description of the policy
- `priority` (Number) Specifies the order in which this policy is evaluated in relation to the other policies
- `skip_unlock` (Boolean) Indicates if, when performing an unlock operation on an Active Directory sourced User who is locked out of Okta, the system should also attempt to unlock the User's Windows account
- `expire_warn_days` (Number) Specifies the number of days prior to password expiration when a User is warned to reset their password: `0` indicates no warning
- `history_count` (Number) Specifies the number of distinct passwords that a User must create before they can reuse a previous password: `0` indicates none
- `max_age_days` (Number) Specifies how long (in days) a password remains valid before it expires: `0` indicates no limit
- `min_age_minutes` (Number) Specifies the minimum time interval (in minutes) between password changes: `0` indicates no limit
- `delegated_workflow_id` (String) The `id` of the workflow that runs when a breached password is found during a sign-in attempt.
- `expire_after_days` (Number) Specifies the number of days after a breached password is found during a sign-in attempt that the user's password should expire. Valid values: 0 through 10. If set to 0, it happens immediately.
- `logout_enabled` (Boolean) (Optional, default is false) If true, you must also specify a value for `expireAfterDays`. When enabled, the user's session(s) are terminated immediately the first time the user's credentials are d...
- `exclude` (Boolean) Indicates whether to check passwords against the common password dictionary
- `exclude_attributes` (List) The User profile attributes whose values must be excluded from the password: currently only supports `firstName` and `lastName`
- `exclude_username` (Boolean) Indicates if the Username must be excluded from the password
- `max_consecutive_characters` (Number) <x-lifecycle-container><x-lifecycle class='oie'></x-lifecycle></x-lifecycle-container>Specifies the maximum number of consecutive repeating characters that can be used in a password
- `min_length` (Number) Minimum password length
- `min_lower_case` (Number) Indicates if a password must contain at least one lower case letter: `0` indicates no, `1` indicates yes
- `min_number` (Number) Indicates if a password must contain at least one number: `0` indicates no, `1` indicates yes
- `min_symbol` (Number) Indicates if a password must contain at least one symbol (For example: !@#$%^&*): `0` indicates no, `1` indicates yes
- `min_upper_case` (Number) Indicates if a password must contain at least one upper case letter: `0` indicates no, `1` indicates yes
- `oel_statement` (String) <x-lifecycle-container><x-lifecycle class='oie'></x-lifecycle></x-lifecycle-container>Use an [Expression Language](https://developer.okta.com/docs/reference/okta-expression-language-in-identity-eng...
- `auto_unlock_minutes` (Number) Specifies the time interval (in minutes) a locked account remains locked before it is automatically unlocked: `0` indicates no limit
- `max_attempts` (Number) Specifies the number of times Users can attempt to sign in to their accounts with an invalid password before their accounts are locked: `0` indicates no limit
- `show_lockout_failures` (Boolean) Indicates if the User should be informed when their account is locked
- `user_lockout_notification_channels` (List) How the user is notified when their account becomes locked. The only acceptable values are `[]` and `['EMAIL']`.
- `status` (String) Status
- `token_lifetime_minutes` (Number) Lifetime (in minutes) of the recovery token
- `status` (String) Status
- `status` (String) Status
- `status` (String) Status
- `status` (String) Status
- `system` (Boolean) Specifies whether Okta created the policy

### Read-Only

- `id` (String) The unique identifier for the resource.
- `embedded` (String) Embedded
- `created` (String) Timestamp when the policy was created
- `last_updated` (String) Timestamp when the policy was last modified
- `min_length` (Number) Minimum length of the password recovery question answer
