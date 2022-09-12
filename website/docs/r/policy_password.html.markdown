---
layout: 'okta'
page_title: 'Okta: okta_policy_password'
sidebar_current: 'docs-okta-resource-policy-password'
description: |-
  Creates a Password Policy.
---

# okta_policy_password

Creates a Password Policy.

This resource allows you to create and configure a Password Policy.

## Example Usage

```hcl
resource "okta_policy_password" "example" {
  name                   = "example"
  status                 = "ACTIVE"
  description            = "Example"
  password_history_count = 4
  groups_included        = ["${data.okta_group.everyone.id}"]
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) Policy Name. Type `"string"`

- `description` - (Optional) Policy Description. Type `"string"`

- `priority` - (Optional) Priority of the policy. Type `"number"`

- `status` - (Optional) Policy Status: `"ACTIVE"` or `"INACTIVE"`. Type `"string"`

- `groups_included` - (Optional) List of Group IDs to Include. Type `"list(string)"`

- `auth_provider` - (Optional) Authentication Provider: `"OKTA"`, `"ACTIVE_DIRECTORY"` or `"LDAP"`. Default is `"OKTA"`. Type `"string"`

- `password_min_length` - (Optional) Minimum password length. Default is 8. Type `"number"`

- `password_min_lowercase` - (Optional) Minimum number of lower case characters in a password. Type `"number"`

- `password_min_uppercase` - (Optional) Minimum number of upper case characters in a password. Type `"number"`

- `password_min_number` - (Optional) Minimum number of numbers in a password. Type `"number"`

- `password_min_symbol` - (Optional) Minimum number of symbols in a password. Type `"number"`

- `password_exclude_username` - (Optional) If the username must be excluded from the password. Type `"bool"`

- `password_exclude_first_name` - (Optional) User firstName attribute must be excluded from the password. Type `"bool"`

- `password_exclude_last_name` - (Optional) User lastName attribute must be excluded from the password. Type `"bool"`

- `password_dictionary_lookup` - (Optional) Check Passwords Against Common Password Dictionary. Type `"bool"`

- `password_max_age_days` - (Optional) Length in days a password is valid before expiry: 0 = no limit.  Type `"number"`

- `password_expire_warn_days` - (Optional) Length in days a user will be warned before password expiry: 0 = no warning.  Type `"number"`

- `password_min_age_minutes` - (Optional) Minimum time interval in minutes between password changes: 0 = no limit.  Type `"number"`

- `password_history_count` - (Optional) Number of distinct passwords that must be created before they can be reused: 0 = none.  Type `"number"`

- `password_max_lockout_attempts` - (Optional) Number of unsuccessful login attempts allowed before lockout: 0 = no limit.  Type `"number"`

- `password_auto_unlock_minutes` - (Optional) Number of minutes before a locked account is unlocked: 0 = no limit.  Type `"number"`

- `password_show_lockout_failures` - (Optional) If a user should be informed when their account is locked.  Type `"bool"`

- `password_lockout_notification_channels` - (Optional) Notification channels to use to notify a user when their account has been locked. Type `"set(string)"`

- `question_min_length` - (Optional) Min length of the password recovery question answer.  Type `"number"`

- `email_recovery` - (Optional) Enable or disable email password recovery: ACTIVE or INACTIVE. Type `"string"`

- `recovery_email_token` - (Optional) Lifetime in minutes of the recovery email token.  Type `"number"`

- `sms_recovery` - (Optional) Enable or disable SMS password recovery: ACTIVE or INACTIVE. Type `"string"`

- `call_recovery` - (Optional) Enable or disable voice call password recovery: ACTIVE or INACTIVE.  Type `"string"`

- `question_recovery` - (Optional) Enable or disable security question password recovery: ACTIVE or INACTIVE. Type `"string"`

- `skip_unlock` - (Optional) When an Active Directory user is locked out of Okta, the Okta unlock operation should also attempt to unlock the user's Windows account. Type `"bool"`

## Attributes Reference

- `id` - ID of the Policy.

## Import

A Password Policy can be imported via the Okta ID.

```
$ terraform import okta_policy_password.example &#60;policy id&#62;
```
