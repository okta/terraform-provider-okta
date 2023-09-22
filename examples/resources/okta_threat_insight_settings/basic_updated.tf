resource "okta_network_zone" "test" {
  name     = "testAcc_replace_with_uuid"
  type     = "IP"
  gateways = ["1.2.3.4/24", "2.3.4.5-2.3.4.15"]
  proxies  = ["2.2.3.4/24", "3.3.4.5-3.3.4.15"]
  status   = "ACTIVE"
}

resource "okta_threat_insight_settings" "test" {
  action           = "block"
  network_excludes = [okta_network_zone.test.id]
}
