---
page_title: "Resource: okta_device_posture_check"
description: |-
  Manages an okta_device_posture_check resource.
---

# Resource: okta_device_posture_check

Manages an okta_device_posture_check resource.



## Argument Reference

The following arguments are supported:
- `remediation_settings` - (Optional) Represents the remediation instructions shown to the end user when the device posture check fails

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

- `id` - The ID of the resource.
- `created_by` - User who created the device posture check
- `description` - Description of the device posture check
- `last_updated_by` - User who updated the device posture check
- `mapping_type` - Represents how the device posture check is rendered in device assurance policies
- `name` - Display name of the device posture check
- `platform` - Platform
- `query` - OSQuery for the device posture check
- `type` - Type
- `variable_name` - Unique name of the device posture check


