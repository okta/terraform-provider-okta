resource "okta_network_zone" "ip_network_zone_example" {
  name     = "testAcc_replace_with_uuid"
  type     = "IP"
  status   = "ACTIVE"
  gateways = ["1.2.3.4/24", "2.3.4.5-2.3.4.15"]
  proxies  = ["2.2.3.4/24", "3.3.4.5-3.3.4.15"]
}