resource "okta_network_zone" "default_blocklist" {
  name     = "BlockedIpZone"
  type     = "IP"
  status   = "ACTIVE"
  gateways = []
}
