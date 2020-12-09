---
layout: 'okta'
page_title: 'Okta: okta_network_zone'
sidebar_current: 'docs-okta-resource-network-zone'
description: |-
  Creates an Okta Network Zone.
---

# okta_network_zone

Creates an Okta Network Zone.

This resource allows you to create and configure an Okta Network Zone.

## Example Usage

```hcl
resource "okta_network_zone" "example" {
  name     = "example"
  type     = "IP"
  gateways = ["1.2.3.4/24", "2.3.4.5-2.3.4.15"]
  proxies  = ["2.2.3.4/24", "3.3.4.5-3.3.4.15"]
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) Name of the Network Zone Resource.

- `type` - (Required) Type of the Network Zone - can either be IP or DYNAMIC only.

- `dynamic_locations` - (Optional) Array of locations ISO-3166-1(2). Format code: countryCode OR countryCode-regionCode.

- `gateways` - (Optional) Array of values in CIDR/range form.

- `proxies` - (Optional) Array of values in CIDR/range form.

## Attributes Reference

- `id` - Network Zone ID.

## Import

Okta Network Zone can be imported via the Okta ID.

```
$ terraform import okta_network_zone.example <zone id>
```
