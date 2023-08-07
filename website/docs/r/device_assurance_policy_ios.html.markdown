---
layout: 'okta'
page_title: 'Okta: okta_policy_device_assurance_ios'
sidebar_current: 'docs-okta-device-assurance-policy-ios'
description: |-
    Manages a device assurance policy for ios.
---

# okta_policy_device_assurance_ios

This resource allows you to create and configure an device assurance policy for ios.

## Example Usage

```hcl
resource okta_policy_device_assurance_ios example{
    name = "example"
    os_version = "12.4.5"
    jailbreak = false
    screenlock_type = toset(["BIOMETRIC"])
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) Name of the device assurance policy.

- `jailbreak` - (Optional)  Is the device jailbroken in the device assurance policy.

- `os_version` - (Optional) Minimum os version of the device in the device assurance policy.

- `screenlock_type` - (Optional) List of screen lock type of the device assurance policy.

## Attributes Reference

- `id` - ID of the device assurance policy.

- `platform` - Platform of the device assurance policy.

- `created_date` - Created date of the device assurance policy.

- `created_by` - Created by of the device assurance policy.

- `last_update` - Last update of the device assurance policy.

- `last_updated_by` - Last updated by of the device assurance policy.

## Import

Okta Device Assurance iOS can be imported via the Okta ID.

```
$ terraform import okta_policy_device_assurance_ios.example &#60;device assurance id&#62;
```