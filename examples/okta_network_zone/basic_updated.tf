resource "okta_network_zone" "ip_network_zone_example" {
  name     = "testAcc_replace_with_uuid Updated"
  type     = "IP"
  gateways = ["1.2.3.4/24", "2.3.4.5-2.3.4.10"]
  usage    = "BLOCKLIST"
}

resource "okta_network_zone" "dynamic_network_zone_example" {
  name              = "testAcc_replace_with_uuid Dynamic Updated"
  type              = "DYNAMIC"
  dynamic_locations = ["US", "AF-BGL", "UA-26"]
}
