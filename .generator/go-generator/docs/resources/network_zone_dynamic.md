---
page_title: "okta_network_zone_dynamic Resource - terraform-provider-okta"
description: |-
  Manages an Okta Network Zone Dynamic resource.
---

# okta_network_zone_dynamic (Resource)

Manages an Okta Network Zone Dynamic.

## Example Usage

```hcl
resource "okta_network_zone_dynamic" "example" {
  type = "<type>"
  name = "Example Name"

  # Optional fields
  # asns = "<asns>"
  # locations = "<locations>"
  # proxy_type = "<proxy_type>"
  # status = "ACTIVE"
  # usage = "<usage>"
}
```

## Schema

### Required

- `type` (String) Discriminator field identifying the variant type. Must be set to \
- `name` (String) Unique name for this Network Zone

### Optional

- `asns` (String) Asns
- `locations` (String) Locations
- `proxy_type` (String) The proxy type used for a Dynamic Network Zone
- `status` (String) Network Zone status
- `usage` (String) The usage of the Network Zone

### Read-Only

- `id` (String) The unique identifier for the resource.
- `created` (String) Timestamp when the object was created
- `last_updated` (String) Timestamp when the object was last modified
- `system` (Boolean) Indicates a system Network Zone: * `true` for system Network Zones * `false` for custom Network Zones  The Okta org provides the following default system Network Zones: * `LegacyIpZone` * `BlockedI...
