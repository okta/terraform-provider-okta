---
page_title: "Data Source: okta_policy_device_assurance_macos"
description: |-
  Retrieves a macOS Device Assurance Policy.
---

# Data Source: okta_policy_device_assurance_macos

Retrieves a device assurance policy by `deviceAssuranceId`.

!> **BREAKING CHANGE**: This data source has been migrated to use Okta SDK v7, which introduces schema changes. Migration guide below.

## Breaking Changes

### 1. Data Source Replaces the Generic `okta_device_assurance_policy` Data Source
The previous generic `okta_device_assurance_policy` data source (queryable by `id` or `name`) has been replaced by platform-specific data sources, one per OS. This data source, `okta_policy_device_assurance_macos`, retrieves macOS device assurance policies.

### 2. Argument Changes
- `id` is now **required**; querying by `name` is no longer supported
- `platform` is no longer exposed (the platform is implied by the data source itself)

### 3. Schema Changes
- `screenlock_type` has been renamed to `screen_lock_type` and is now a nested block with an `include` attribute
- `disk_encryption_type` is now a nested block with an `include` attribute
- `os_version` is now a nested block with a `minimum` attribute and an optional `dynamic_version_requirement` sub-block
- Third-party signal provider attributes previously flattened with a `tpsp_` prefix have been restructured into nested blocks under `third_party_signal_providers.dtc` (the Google Chrome Device Trust Connector provider), using names without the `tpsp_` prefix
- New computed attributes: `display_remediation_mode`, `grace_period` (`type`, `expiry`), and `third_party_signal_providers.device_posture_id_p` (`compliant`, `managed`)

## Example Usage

```terraform
resource "okta_policy_device_assurance_macos" "example" {
  name                    = "My macOS Policy"
  platform                = "MACOS"
  secure_hardware_present = true
}

data "okta_policy_device_assurance_macos" "example" {
  id = okta_policy_device_assurance_macos.example.id
}
```

## Argument Reference

- `id` - (Required) The unique identifier of the device assurance policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `name` - Display name of the device assurance policy.
- `created_by` - User ID who created the policy.
- `created_date` - When the policy was created.
- `last_updated_by` - User ID who last updated the policy.
- `last_update` - When the policy was last updated.
- `display_remediation_mode` - Remediation mode for non-compliant devices.
- `secure_hardware_present` - Whether secure hardware is required.
- `disk_encryption_type` - Disk encryption requirements.
  - `include` - List of disk encryption types required.
- `os_version` - Operating system version requirements.
  - `minimum` - Minimum OS version required.
  - `dynamic_version_requirement` - Dynamic OS version requirement settings.
    - `distance_from_latest_major` - Distance from the latest major version.
    - `latest_security_patch` - Whether the latest security patch is required.
    - `type` - Type of the dynamic OS version requirement.
- `screen_lock_type` - Screen lock type requirements.
  - `include` - List of screen lock types required.
- `third_party_signal_providers` - Third-party signal provider configuration.
  - `device_posture_id_p` - Device Posture IdP provider settings.
    - `compliant` - Whether device must be compliant.
    - `managed` - Whether device must be managed.
  - `dtc` - Google Chrome Device Trust Connector provider settings.
    - `built_in_dns_client_enabled` - Whether a software DNS client stack is used.
    - `chrome_remote_desktop_app_blocked` - Whether Chrome Remote Desktop access is blocked.
    - `device_enrollment_domain` - Enrollment domain of the managing customer.
    - `disk_encrypted` - Whether the main disk is encrypted.
    - `key_trust_level` - Attestation strength used by the Chrome Verified Access API.
    - `os_firewall` - Whether an OS-level firewall is enabled.
    - `password_protection_warning_trigger` - Password Protection Warning feature state.
    - `realtime_url_check_mode` - Whether enterprise-grade unsafe URL scanning is enabled.
    - `safe_browsing_protection_level` - Current Safe Browsing protection level.
    - `screen_lock_secured` - Whether the device is password-protected.
    - `site_isolation_enabled` - Whether Site Isolation is enabled.
    - `browser_version` - Current version of the Chrome Browser.
      - `minimum` - Minimum browser version.
    - `os_version` - Current version of the operating system.
      - `minimum` - Minimum OS version.
- `grace_period` - Grace period configuration.
  - `type` - Type of grace period (e.g., "DAYS").
  - `expiry` - Grace period duration.
