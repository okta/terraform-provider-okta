---
page_title: "okta_device_assurance_chromeos Resource - terraform-provider-okta"
description: |-
  Manages an Okta Device Assurance Chromeos resource.
---

# okta_device_assurance_chromeos (Resource)

Manages an Okta Device Assurance Chromeos.

## Example Usage

```hcl
resource "okta_device_assurance_chromeos" "example" {
  platform = "<platform>"

  # Optional fields
  # value = "<value>"
  # variable_name = "Example Variable Name"
  # display_remediation_mode = "<display_remediation_mode>"
  # expiry = "<expiry>"
  # type = "<type>"
}
```

## Schema

### Required

- `platform` (String) Discriminator field identifying the variant type. Must be set to \

### Optional

- `value` (String) The device posture check value
- `variable_name` (String) The device posture check key
- `display_remediation_mode` (String) Represents the remediation mode of this device assurance policy when users are denied access due to device noncompliance
- `expiry` (String) Expiry
- `type` (String) Represents the type of Grace Period configured for the device assurance policy
- `name` (String) Display name of the device assurance policy
- `compliant` (Boolean) Indicates whether the device is compliant according to the custom IDP
- `managed` (Boolean) Indicates whether the device is managed according to the custom IDP
- `allow_screen_lock` (Boolean) Indicates whether the AllowScreenLock enterprise policy is enabled
- `minimum` (String) Minimum
- `built_in_dns_client_enabled` (Boolean) Indicates if a software stack is used to communicate with the DNS server
- `chrome_remote_desktop_app_blocked` (Boolean) Indicates whether access to the Chrome Remote Desktop application is blocked through a policy
- `device_enrollment_domain` (String) Enrollment domain of the customer that is currently managing the device
- `disk_encrypted` (Boolean) Indicates whether the main disk is encrypted
- `key_trust_level` (String) Represents the attestation strength used by the Chrome Verified Access API
- `managed_device` (Boolean) Indicates whether the device is enrolled in ChromeOS device management
- `os_firewall` (Boolean) Indicates whether a firewall is enabled at the OS-level on the device
- `minimum` (String) Minimum
- `password_protection_warning_trigger` (String) Indicates whether the Password Protection Warning feature is enabled
- `realtime_url_check_mode` (Boolean) Indicates whether enterprise-grade (custom) unsafe URL scanning is enabled
- `safe_browsing_protection_level` (String) Represents the current value of the Safe Browsing protection level
- `screen_lock_secured` (Boolean) Indicates whether the device is password-protected
- `site_isolation_enabled` (Boolean) Indicates whether the Site Isolation (also known as **Site Per Process**) setting is enabled

### Read-Only

- `id` (String) The unique identifier for the resource.
- `created_by` (String) CreatedBy
- `created_date` (String) CreatedDate
- `last_update` (String) LastUpdate
- `last_updated_by` (String) LastUpdatedBy
