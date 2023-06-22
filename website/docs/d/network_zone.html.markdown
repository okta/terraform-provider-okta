---
layout: 'okta'
page_title: 'Okta: okta_network_zone'
sidebar_current: 'docs-okta-datasource-network-zone'
description: |-
  Gets Okta Network Zone.
---

# okta_network_zone

Use this data source to retrieve a network zone from Okta.

## Example Usage

```hcl
data "okta_network_zone" "example" {
  name = "Block Antarctica"
}
```

## Argument Reference

- `id` - (Optional) ID of the network zone to retrieve, conflicts with `name`.

- `name` - (Optional) Name of the network zone to retrieve, conflicts with `id`.

## Attributes Reference

- `id` - ID of the network zone.

- `name` - Name of the network zone.

- `type` - Type of the Network Zone.

- `status` - Network Status - can either be ACTIVE or INACTIVE only.

- `dynamic_locations` - Array of locations.

- `dynamic_proxy_type` - Type of proxy being controlled by this dynamic network zone.

- `gateways` -  Array of values in CIDR/range form.

- `proxies` -  Array of values in CIDR/range form.

- `usage` - Usage of the Network Zone.

- `asns` - Array of Autonomous System Numbers.
