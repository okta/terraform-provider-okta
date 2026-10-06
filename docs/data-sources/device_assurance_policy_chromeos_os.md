---
page_title: "Data Source: okta_device_assurance_chromeos"
description: |-
  Get a ChromeOS Device Assurance Policy data source.
---

# Data Source: okta_device_assurance_chromeos

Get a ChromeOS Device Assurance Policy data source.

## Example Usage

```terraform
data "okta_device_assurance_chromeos" "example" {
  id = "device_assurance_policy_id"
}
```

## Argument Reference

- `id` - (Required) The unique identifier of the device assurance policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `name` - Display name of the device assurance policy.
- `platform` - The platform type (always "CHROMEOS").
- `created_by` - User ID who created the policy.
- `created_date` - When the policy was created.
- `last_updated_by` - User ID who last updated the policy.
- `last_update` - When the policy was last updated.
- `display_remediation_mode` - Remediation mode for non-compliant devices.
- `third_party_signal_providers` - Third-party signal provider configuration.
  - `device_posture_id_p` - Device Posture IdP provider settings.
    - `compliant` - Whether device must be compliant.
    - `managed` - Whether device must be managed.
  - `dtc` - Google Chrome Device Trust Connector settings.
    - `allow_screen_lock` - Whether AllowScreenLock enterprise policy is enabled.
    - `built_in_dns_client_enabled` - Whether built-in DNS client is enabled.
    - `chrome_remote_desktop_app_blocked` - Whether Chrome Remote Desktop is blocked.
    - `device_enrollment_domain` - Device enrollment domain.
    - `disk_encrypted` - Whether main disk is encrypted.
    - `key_trust_level` - Chrome Verified Access API attestation strength.
    - `managed_device` - Whether device is enrolled in ChromeOS management.
    - `os_firewall` - Whether OS-level firewall is enabled.
    - `password_protection_warning_trigger` - Password Protection Warning feature status.
    - `realtime_url_check_mode` - Whether enterprise unsafe URL scanning is enabled.
    - `safe_browsing_protection_level` - Safe Browsing protection level.
    - `screen_lock_secured` - Whether device is password-protected.
    - `site_isolation_enabled` - Whether Site Isolation is enabled.
    - `browser_version` - Chrome browser version requirements.
      - `minimum` - Minimum browser version.
    - `os_version` - Operating system version requirements.
      - `minimum` - Minimum OS version.
- `grace_period` - Grace period configuration.
  - `type` - Type of grace period (e.g., "DAYS").
  - `expiry` - Grace period duration.

## Import

This data source can be imported by specifying the policy ID:

```shell
terraform import okta_device_assurance_chromeos.example device_assurance_policy_id
```

Note: Data sources cannot be imported directly. Instead, import the corresponding resource and then reference it with a data source as shown in the example above.
