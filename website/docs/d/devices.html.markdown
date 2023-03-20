---
layout: 'okta'
page_title: 'Okta: okta_devices'
sidebar_current: 'docs-okta-devices'
description: |-
List okta devices
---


# okta_devices

Use this data source to list the okta devices, searchable by user_id. This allows you to retrieve device information for use within Terraform.

## Example Usage

```hcl
data "okta_devices" "example" {
  user_id = "00u22mtxlrJ8YkzXQ357"
}
```

## Argument Reference

- `user_id` - (Required) The ID of the Okta user you want to retrieve the list of devices for.


## Attribute Reference

- `id` - ID of the device.

- `status` - Current status of device. One of CREATED, ACTIVE, SUSPENDED or DEACTIVATED.

- `profile` - Device profile properties.

- `resource_type` - Device type.

- `resource_display_name` - Device Display name

