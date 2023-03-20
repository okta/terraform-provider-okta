---
layout: 'okta'
page_title: 'Okta: okta_device'
sidebar_current: 'docs-okta-device'
description: |-
Get an okta device by its id
---


# okta_device

Use this data source to get the okta device (by ID). This allows you to retrieve device information for use within Terraform.

## Example Usage

```hcl
data "okta_device" "example" {
  id = "guo4a5u7YAHhjXrMN0g4"
}
```

## Argument Reference

- `id` - (Required) The ID of the Okta device you want to retrieve.


## Attribute Reference

- `id` - ID of the device.

- `status` - Current status of device. One of CREATED, ACTIVE, SUSPENDED or DEACTIVATED.

- `profile` - Device profile properties.

- `resource_type` - Device type.

- `resource_display_name` - Device Display name

