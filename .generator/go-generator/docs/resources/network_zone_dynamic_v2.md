---
page_title: "okta_network_zone_dynamic_v2 Resource - terraform-provider-okta"
description: |-
  Manages an Okta Network Zone Dynamic V2 resource.
---

# okta_network_zone_dynamic_v2 (Resource)

Manages an Okta Network Zone Dynamic V2.

## Example Usage

```hcl
resource "okta_network_zone_dynamic_v2" "example" {
  type = "<type>"
  name = "Example Name"

  # Optional fields
  # exclude = "<exclude>"
  # include = "<include>"
  # exclude = "<exclude>"
  # include = "<include>"
  # exclude = "<exclude>"
}
```

## Schema

### Required

- `type` (String) Discriminator field identifying the variant type. Must be set to \
- `name` (String) Unique name for this Network Zone

### Optional

- `exclude` (String) Exclude
- `include` (String) Include
- `exclude` (List) IP services to exclude for an Enhanced Dynamic Network Zone
- `include` (List) IP services to include for an Enhanced Dynamic Network Zone
- `exclude` (String) Exclude
- `include` (String) Include
- `status` (String) Network Zone status
- `usage` (String) The usage of the Network Zone

### Read-Only

- `id` (String) The unique identifier for the resource.
- `created` (String) Timestamp when the object was created
- `last_updated` (String) Timestamp when the object was last modified
- `system` (Boolean) Indicates a system Network Zone: * `true` for system Network Zones * `false` for custom Network Zones  The Okta org provides the following default system Network Zones: * `LegacyIpZone` * `BlockedI...
