---
layout: 'okta'
page_title: 'Okta: okta_behavior'
sidebar_current: 'docs-okta-resource-behavior'
description: |-
  Creates different types of behavior.
---

# okta_behavior

This resource allows you to create and configure a behavior.

## Example Usage

```hcl
resource "okta_behavior" "my_location" {
  name = "My Location"
  type = "ANOMALOUS_LOCATION"
  number_of_authentications = 50
  location_granularity_type = "LAT_LONG"
  radius_from_location = 20
}

resource "okta_behavior" "my_city" {
  name = "My City"
  type = "ANOMALOUS_LOCATION"
  number_of_authentications = 50
  location_granularity_type = "CITY"
}

resource "okta_behavior" "my_device" {
  name = "My Device"
  type = "ANOMALOUS_DEVICE"
  number_of_authentications = 50
}

resource "okta_behavior" "my_ip" {
  name = "My IP"
  type = "ANOMALOUS_IP"
  number_of_authentications = 50
}

resource "okta_behavior" "my_velocity" {
  name = "My Velocity"
  type = "VELOCITY"
  velocity = 25
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) Name of the behavior.

- `type` - (Required) Type of the behavior. Can be set to `"ANOMALOUS_LOCATION"`, `"ANOMALOUS_DEVICE"`, `"ANOMALOUS_IP"`
  or `"VELOCITY"`. Resource will be recreated when the type changes.

- `status` - (Optional) The status of the behavior. By default, it is`"ACTIVE"`.

- `location_granularity_type` - (Optional) Determines the method and level of detail used to evaluate the behavior.
  Required for `"ANOMALOUS_LOCATION"` behavior type. Can be set to `"LAT_LONG"`, `"CITY"`, `"COUNTRY"`
  or `"SUBDIVISION"`.

- `radius_from_location` - (Optional) Radius from location (in kilometers). Should be at least 5. Required
  when `location_granularity_type` is set to `"LAT_LONG"`.

- `number_of_authentications` - (Optional) The number of recent authentications used to evaluate the behavior. Required
  for `"ANOMALOUS_LOCATION"`, `"ANOMALOUS_DEVICE"` and `"ANOMALOUS_IP"` behavior types.

- `velocity` - (Optional) Velocity (in kilometers per hour). Should be at least 1. Required for `"VELOCITY"` behavior
  type.

## Attributes Reference

- `id` - ID of the behavior.

## Import

Behavior can be imported via the Okta ID.

```
$ terraform import okta_behavior.example &#60;behavior id&#62;
```
