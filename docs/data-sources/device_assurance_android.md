---
page_title: "Data Source: okta_device_assurance_android"
description: |-
  Get an Android Device Assurance Policy data source.
---

# Data Source: okta_device_assurance_android

Get an Android Device Assurance Policy data source.

## Example Usage

```terraform
data "okta_device_assurance_android" "example" {
  id = "device_assurance_policy_id"
}
```

## Argument Reference

- `id` - (Required) The unique identifier of the device assurance policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `name` - Display name of the device assurance policy.
- `platform` - The platform type (always "ANDROID").
- `created_by` - User ID who created the policy.
- `created_date` - When the policy was created.
- `last_updated_by` - User ID who last updated the policy.
- `last_update` - When the policy was last updated.
- `display_remediation_mode` - Remediation mode for non-compliant devices.
- `jailbreak` - Whether jailbreak detection is enabled.
- `secure_hardware_present` - Whether secure hardware is required.
- `disk_encryption_type` - Disk encryption requirements.
- `os_version` - Operating system version requirements.
- `screen_lock_type` - Screen lock type requirements.
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

## Import

This data source can be imported by specifying the policy ID:

```shell
terraform import okta_device_assurance_android.example device_assurance_policy_id
```

Note: Data sources cannot be imported directly. Instead, import the corresponding resource and then reference it with a data source as shown in the example above.
