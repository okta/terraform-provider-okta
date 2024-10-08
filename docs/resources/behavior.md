---
page_title: "Resource: okta_behavior"
description: |-
  This resource allows you to create and configure a behavior.
---

# Resource: okta_behavior

This resource allows you to create and configure a behavior.

## Example Usage

```terraform
resource "okta_behavior" "my_location" {
  name                      = "My Location"
  type                      = "ANOMALOUS_LOCATION"
  number_of_authentications = 50
  location_granularity_type = "LAT_LONG"
  radius_from_location      = 20
}

resource "okta_behavior" "my_city" {
  name                      = "My City"
  type                      = "ANOMALOUS_LOCATION"
  number_of_authentications = 50
  location_granularity_type = "CITY"
}

resource "okta_behavior" "my_device" {
  name                      = "My Device"
  type                      = "ANOMALOUS_DEVICE"
  number_of_authentications = 50
}

resource "okta_behavior" "my_ip" {
  name                      = "My IP"
  type                      = "ANOMALOUS_IP"
  number_of_authentications = 50
}

resource "okta_behavior" "my_velocity" {
  name     = "My Velocity"
  type     = "VELOCITY"
  velocity = 25
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the behavior
- `type` (String) Type of the behavior. Can be set to `ANOMALOUS_LOCATION`, `ANOMALOUS_DEVICE`, `ANOMALOUS_IP` or `VELOCITY`. Resource will be recreated when the type changes.e

### Optional

- `location_granularity_type` (String) Determines the method and level of detail used to evaluate the behavior. Required for `ANOMALOUS_LOCATION` behavior type. Can be set to `LAT_LONG`, `CITY`, `COUNTRY` or `SUBDIVISION`.
- `number_of_authentications` (Number) The number of recent authentications used to evaluate the behavior. Required for `ANOMALOUS_LOCATION`, `ANOMALOUS_DEVICE` and `ANOMALOUS_IP` behavior types.
- `radius_from_location` (Number) Radius from location (in kilometers). Should be at least 5. Required when `location_granularity_type` is set to `LAT_LONG`.
- `status` (String) Behavior status: ACTIVE or INACTIVE. Default: `ACTIVE`
- `velocity` (Number) Velocity (in kilometers per hour). Should be at least 1. Required for `VELOCITY` behavior

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import okta_behavior.example <behavior_id>
```
