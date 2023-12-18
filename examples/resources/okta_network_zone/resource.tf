resource "okta_network_zone" "example" {
  name     = "example"
  type     = "IP"
  gateways = ["1.2.3.4/24", "2.3.4.5-2.3.4.15"]
  proxies  = ["2.2.3.4/24", "3.3.4.5-3.3.4.15"]
}

## Example Usage - Dynamic Tor Blocker

resource "okta_network_zone" "example" {
  name               = "TOR Blocker"
  type               = "DYNAMIC"
  usage              = "BLOCKLIST"
  dynamic_proxy_type = "TorAnonymizer"
}
