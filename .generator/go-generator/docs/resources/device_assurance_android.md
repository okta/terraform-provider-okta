---
page_title: "okta_device_assurance_android Resource - terraform-provider-okta"
description: |-
  Manages an Okta Device Assurance Android resource.
---

# okta_device_assurance_android (Resource)

Manages an Okta Device Assurance Android.

## Example Usage

```hcl
resource "okta_device_assurance_android" "example" {
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
- `jailbreak` (Boolean) Jailbreak
- `name` (String) Display name of the device assurance policy
- `distance_from_latest_major` (Number) Indicates the distance from the latest major version
- `latest_security_patch` (Boolean) Indicates whether the device needs to be on the latest security patch
- `type` (String) Indicates the type of the dynamic OS version requirement
- `minimum` (String) The device version must be equal to or newer than the specified version string (maximum of three components for iOS and macOS, and maximum of four components for Android)
- `include` (List) Include
- `secure_hardware_present` (Boolean) SecureHardwarePresent
- `device_integrity_level` (String) Indicates how well a device can enforce app integrity
- `network_proxy_disabled` (Boolean) Indicates whether a device has a network proxy disabled
- `play_protect_verdict` (String) Indicates if Google Play Protect is enabled on the device and whether it has found known malware
- `require_major_version_update` (Boolean) Indicates whether the device needs to be on the latest major version available to the device  **Note:** This option requires an `osVersion.dynamicVersionRequirement` value to be supplied with the `...
- `screen_lock_complexity` (String) Indicates whether a device has a screen lock set, and the type or complexity of the screen lock  **Note:** This option requires a `screenLockType.include` value.
- `usb_debugging_disabled` (Boolean) Indicates whether Android Debug Bridge (adb) over USB is disabled
- `wifi_secured` (Boolean) Indicates whether a device is on a password-protected Wi-Fi network
- `compliant` (Boolean) Indicates whether the device is compliant according to the custom IDP
- `managed` (Boolean) Indicates whether the device is managed according to the custom IDP

### Read-Only

- `id` (String) The unique identifier for the resource.
- `created_by` (String) CreatedBy
- `created_date` (String) CreatedDate
- `last_update` (String) LastUpdate
- `last_updated_by` (String) LastUpdatedBy
