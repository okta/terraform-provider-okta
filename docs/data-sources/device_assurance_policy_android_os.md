---
page_title: "Data Source: okta_policy_device_assurance_android"
description: |-
  Retrieves an Android Device Assurance Policy.
---

# Data Source: okta_policy_device_assurance_android

Retrieves a device assurance policy by `deviceAssuranceId`.

!> **BREAKING CHANGE**: This data source has been migrated to use Okta SDK v7, which introduces schema changes. Migration guide below.

## Breaking Changes

### 1. Data Source Replaces the Generic `okta_device_assurance_policy` Data Source
The previous generic `okta_device_assurance_policy` data source (queryable by `id` or `name`) has been replaced by platform-specific data sources, one per OS. This data source, `okta_policy_device_assurance_android`, retrieves Android device assurance policies.

### 2. Argument Changes
- `id` is now **required**; querying by `name` is no longer supported
- `platform` is no longer exposed (the platform is implied by the data source itself)

### 3. Schema Changes
- `screenlock_type` has been renamed to `screen_lock_type` and is now a nested block with an `include` attribute
- `disk_encryption_type` is now a nested block with an `include` attribute
- `os_version` is now a nested block with a `minimum` attribute and an optional `dynamic_version_requirement` sub-block
- New computed attributes: `display_remediation_mode`, `grace_period` (`type`, `expiry`), and `third_party_signal_providers` (`android_device_trust`, `device_posture_id_p`)

## Example Usage

```terraform
resource "okta_policy_device_assurance_android" "example" {
  name                    = "My Android Policy"
  platform                = "ANDROID"
  secure_hardware_present = true
}

data "okta_policy_device_assurance_android" "example" {
  id = okta_policy_device_assurance_android.example.id
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
- `jailbreak` - Whether jailbreak detection is enabled.
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
  - `android_device_trust` - Android Device Trust integration settings.
    - `device_integrity_level` - Device integrity attestation level.
    - `network_proxy_disabled` - Whether network proxy must be disabled.
    - `play_protect_verdict` - Play Protect verdict requirement.
    - `require_major_version_update` - Whether major version update is required.
    - `screen_lock_complexity` - Screen lock complexity requirement.
    - `usb_debugging_disabled` - Whether USB debugging must be disabled.
    - `wifi_secured` - Whether WiFi must be secured.
  - `device_posture_id_p` - Device Posture IdP provider settings.
    - `compliant` - Whether device must be compliant.
    - `managed` - Whether device must be managed.
- `grace_period` - Grace period configuration.
  - `type` - Type of grace period (e.g., "DAYS").
  - `expiry` - Grace period duration.
