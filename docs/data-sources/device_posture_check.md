---
page_title: "Data Source: okta_device_posture_check"
description: |-
  Retrieves a device posture check by ID.
---

# Data Source: okta_device_posture_check

Retrieves a device posture check by ID.

## Argument Reference

The following arguments are supported:
- `id` - (Required) Unique identifier of the device posture check.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

- `created_by` - User who created the device posture check
- `created_date` - Time the device posture check was created
- `description` - Description of the device posture check
- `last_update` - Time the device posture check was updated
- `last_updated_by` - User who updated the device posture check
- `mapping_type` - Represents how the device posture check is rendered in device assurance policies
- `name` - Display name of the device posture check
- `platform` - Platform
- `query` - OSQuery for the device posture check
- `remediation_settings` - Represents the remediation instructions shown to the end user when the device posture check fails
- `type` - Type
- `variable_name` - Unique name of the device posture check
