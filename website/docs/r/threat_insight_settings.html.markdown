---
layout: 'okta'
page_title: 'Okta: okta_threat_insight_settings'
sidebar_current: 'docs-okta-resource-threat-insight-settings'
description: |-
  Manages Okta Threat Insight Settings
---

# okta_threat_insight_settings

This resource allows you to configure Threat Insight Settings.

## Example Usage

```hcl
resource "okta_network_zone" "ip_network_zone_example" {
  name     = "example"
  type     = "IP"
  gateways = ["1.2.3.4/24", "2.3.4.5-2.3.4.15"]
  proxies  = ["2.2.3.4/24", "3.3.4.5-3.3.4.15"]
}

resource "okta_threat_insight_settings" "example" {
  action           = "block"
  network_excludes = [okta_network_zone.ip_network_zone_example.id]
}
```

## Argument Reference

The following arguments are supported:

- `action` - (Required) Specifies how Okta responds to authentication requests from suspicious IPs. Valid values 
are `"none"`, `"audit"`, or `"block"`. A value of `"none"` indicates that ThreatInsight is disabled. A value of `"audit"` 
indicates that Okta logs suspicious requests in the System Log. A value of `"block"` indicates that Okta logs suspicious 
requests in the System Log and blocks the requests.

- `network_excludes` - (Optional) Accepts a list of Network Zone IDs. Can only accept zones of `"IP"` type. 
IPs in the excluded Network Zones aren't logged or blocked by Okta ThreatInsight and proceed to Sign On rules evaluation. 
This ensures that traffic from known, trusted IPs isn't accidentally logged or blocked.

## Import

Threat Insight Settings can be imported without any parameters.

```
$ terraform import okta_threat_insight_settings.example _
```
