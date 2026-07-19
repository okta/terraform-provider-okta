---
page_title: "okta_device_assurance_macos Resource - terraform-provider-okta"
description: |-
  Manages an Okta Device Assurance Macos resource.
---

# okta_device_assurance_macos (Resource)

Manages an Okta Device Assurance Macos.

## Example Usage

```hcl
resource "okta_device_assurance_macos" "example" {
  platform = "<platform>"

  # Optional fields
  # value = "<value>"
  # variable_name = "Example Variable Name"
  # include = "<include>"
  # display_remediation_mode = "<display_remediation_mode>"
  # expiry = "<expiry>"
}
```

## Schema

### Required

- `platform` (String) Discriminator field identifying the variant type. Must be set to \

### Optional

- `value` (String) The device posture check value
- `variable_name` (String) The device posture check key
- `include` (List) Include
- `display_remediation_mode` (String) Represents the remediation mode of this device assurance policy when users are denied access due to device noncompliance
- `expiry` (String) Expiry
- `type` (String) Represents the type of Grace Period configured for the device assurance policy
- `name` (String) Display name of the device assurance policy
- `distance_from_latest_major` (Number) Indicates the distance from the latest major version
- `latest_security_patch` (Boolean) Indicates whether the device needs to be on the latest security patch
- `type` (String) Indicates the type of the dynamic OS version requirement
- `minimum` (String) The device version must be equal to or newer than the specified version string (maximum of three components for iOS and macOS, and maximum of four components for Android)
- `include` (List) Include
- `secure_hardware_present` (Boolean) SecureHardwarePresent
- `compliant` (Boolean) Indicates whether the device is compliant according to the custom IDP
- `managed` (Boolean) Indicates whether the device is managed according to the custom IDP
- `minimum` (String) Minimum
- `built_in_dns_client_enabled` (Boolean) Indicates if a software stack is used to communicate with the DNS server
- `chrome_remote_desktop_app_blocked` (Boolean) Indicates whether access to the Chrome Remote Desktop application is blocked through a policy
- `device_enrollment_domain` (String) Enrollment domain of the customer that is currently managing the device
- `disk_encrypted` (Boolean) Indicates whether the main disk is encrypted
- `key_trust_level` (String) Represents the attestation strength used by the Chrome Verified Access API
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
