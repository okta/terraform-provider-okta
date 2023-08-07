---
layout: 'okta'
page_title: 'Okta: okta_policy_device_assurance_android'
sidebar_current: 'docs-okta-device-assurance-policy-android'
description: |-
    Manages a device assurance policy for android.
---

# okta_policy_device_assurance_android

This resource allows you to create and configure an device assurance policy for android.

## Example Usage

```hcl
resource okta_policy_device_assurance_android example{
    name = "example"
    os_version = "12"
    disk_encryption_type = toset(["FULL", "USER"])
    jailbreak = false
    secure_hardware_present = true
    screenlock_type = toset(["BIOMETRIC"])
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) Name of the device assurance policy.

- `disk_encryption_type` - (Optional) List of disk encryption type of the device assurance policy.

- `jailbreak` - (Optional)  Is the device jailbroken in the device assurance policy.

- `os_version` - (Optional) Minimum os version of the device in the device assurance policy.

- `secure_hardware_present` - (Optional) Is the device secure with hardware in the device assurance policy.

- `screenlock_type` - (Optional) List of screen lock type of the device assurance policy.

## Attributes Reference

- `id` - ID of the device assurance policy.

- `platform` - Platform of the device assurance policy.

- `created_date` - Created date of the device assurance policy.

- `created_by` - Created by of the device assurance policy.

- `last_update` - Last update of the device assurance policy.

- `last_updated_by` - Last updated by of the device assurance policy.

## Import

Okta Device Assurance Android can be imported via the Okta ID.

```
$ terraform import okta_policy_device_assurance_android.example &#60;device assurance id&#62;
```