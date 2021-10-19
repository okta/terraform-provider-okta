---
layout: 'okta'
page_title: 'Okta: okta_security_notification_emails'
sidebar_current: 'docs-okta-resource-security-notification-emails'
description: |-
  Manages Security Notification Emails
---

# okta_security_notification_emails

This resource allows you to configure Security Notification Emails.

**Important Note**: available only when using api token in the provider config.

## Example Usage

```hcl
resource "okta_security_notification_emails" "example" {
  report_suspicious_activity_enabled       = true
  send_email_for_factor_enrollment_enabled = true
  send_email_for_factor_reset_enabled      = true
  send_email_for_new_device_enabled        = true
  send_email_for_password_changed_enabled  = true
}
```

## Argument Reference

- `send_email_for_new_device_enabled` (Optional) - Notifies end users about new sign-on activity. Default is `true`.

- `send_email_for_factor_enrollment_enabled` (Optional) - Notifies end users of any activity on their account related to MFA factor enrollment. Default is `true`.

- `send_email_for_factor_reset_enabled` (Optional) - Notifies end users that one or more factors have been reset for their account. Default is `true`.

- `send_email_for_password_changed_enabled` (Optional) - Notifies end users that the password for their account has changed. Default is `true`.

- `report_suspicious_activity_enabled` (Optional) - Notifies end users about suspicious or unrecognized activity from their account. Default is `true`.

## Import

Security Notification Emails can be imported without any parameters.

```
$ terraform import okta_security_notification_emails.example _
```
