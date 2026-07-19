---
page_title: "okta_network_zone_ip Resource - terraform-provider-okta"
description: |-
  Manages an Okta Network Zone Ip resource.
---

# okta_network_zone_ip (Resource)

Manages an Okta Network Zone Ip.

## Example Usage

```hcl
resource "okta_network_zone_ip" "example" {
  type = "<type>"
  name = "Example Name"

  # Optional fields
  # type = "<type>"
  # value = "<value>"
  # type = "<type>"
  # value = "<value>"
  # status = "ACTIVE"
}
```

## Schema

### Required

- `type` (String) Discriminator field identifying the variant type. Must be set to \
- `name` (String) Unique name for this Network Zone

### Optional

- `type` (String) Format of the IP addresses
- `value` (String) Value in CIDR/range form, depending on the `type` specified
- `type` (String) Format of the IP addresses
- `value` (String) Value in CIDR/range form, depending on the `type` specified
- `status` (String) Network Zone status
- `usage` (String) The usage of the Network Zone
- `use_as_exempt_list` (Boolean) You can **only** use this parameter when making a request to the Replace the network zone endpoint (`/api/v1/zones/{zoneId}`). Set this parameter to `true` in your request when you update the `Defa...

### Read-Only

- `id` (String) The unique identifier for the resource.
- `created` (String) Timestamp when the object was created
- `last_updated` (String) Timestamp when the object was last modified
- `system` (Boolean) Indicates a system Network Zone: * `true` for system Network Zones * `false` for custom Network Zones  The Okta org provides the following default system Network Zones: * `LegacyIpZone` * `BlockedI...
