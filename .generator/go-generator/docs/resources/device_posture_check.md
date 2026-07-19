---
page_title: "okta_device_posture_check Resource - terraform-provider-okta"
description: |-
  Manages an Okta Device Posture Check resource.
---

# okta_device_posture_check (Resource)

Manages an Okta Device Posture Check.

## Example Usage

```hcl
resource "okta_device_posture_check" "example" {

  # Optional fields
  # description = "Example description"
  # mapping_type = "<mapping_type>"
  # name = "Example Name"
  # platform = "<platform>"
  # query = "<query>"
}
```

## Schema

### Optional

- `description` (String) Description of the device posture check
- `mapping_type` (String) Represents how the device posture check is rendered in device assurance policies
- `name` (String) Display name of the device posture check
- `platform` (String) Platform
- `query` (String) OSQuery for the device posture check
- `custom_url` (String) Custom URL for the link
- `default_url` (String) Default URL for the link. This property is only relevant if type is set to `BUILTIN`. If type is set to `CUSTOM`, this field is ignored.
- `custom_text` (String) Custom text for the message
- `default_i18n_key` (String) Default i18n key for the message. This property is only relevant if type is set to `BUILTIN`. If type is set to `CUSTOM`, this field is ignored.
- `type` (String) Type
- `variable_name` (String) Unique name of the device posture check

### Read-Only

- `id` (String) The unique identifier for the resource.
- `created_by` (String) User who created the device posture check
- `created_date` (String) Time the device posture check was created
- `last_update` (String) Time the device posture check was updated
- `last_updated_by` (String) User who updated the device posture check
