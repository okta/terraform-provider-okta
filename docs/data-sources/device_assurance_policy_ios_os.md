---
page_title: "Data Source: okta_device_assurance_ios"
description: |-
  Get an iOS Device Assurance Policy data source.
---

# Data Source: okta_device_assurance_ios

Get an iOS Device Assurance Policy data source.

## Example Usage

```terraform
data "okta_device_assurance_ios" "example" {
  id = "device_assurance_policy_id"
}
```

## Argument Reference

- `id` - (Required) The unique identifier of the device assurance policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `name` - Display name of the device assurance policy.
- `platform` - The platform type (always "IOS").
- `created_by` - User ID who created the policy.
- `created_date` - When the policy was created.
- `last_updated_by` - User ID who last updated the policy.
- `last_update` - When the policy was last updated.
- `display_remediation_mode` - Remediation mode for non-compliant devices.
- `jailbreak` - Whether jailbreak detection is enabled.
- `os_version` - Operating system version requirements.
  - `minimum` - Minimum OS version.
- `third_party_signal_providers` - Third-party signal provider configuration.
  - `device_posture_id_p` - Device Posture IdP provider settings.
    - `compliant` - Whether device must be compliant.
    - `managed` - Whether device must be managed.
- `grace_period` - Grace period configuration.
  - `type` - Type of grace period (e.g., "DAYS").
  - `expiry` - Grace period duration.

## Import

This data source can be imported by specifying the policy ID:

```shell
terraform import okta_device_assurance_ios.example device_assurance_policy_id
```

Note: Data sources cannot be imported directly. Instead, import the corresponding resource and then reference it with a data source as shown in the example above.
