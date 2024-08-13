resource "okta_network_zone" "ip_network_zone_example" {
  name     = "testAcc_replace_with_uuid Updated"
  type     = "IP"
  status   = "INACTIVE"
  gateways = ["1.2.3.4/24", "2.3.4.5-2.3.4.10"]
  usage    = "BLOCKLIST"
}

resource "okta_network_zone" "dynamic_network_zone_example" {
  name              = "testAcc_replace_with_uuid Dynamic Updated"
  type              = "DYNAMIC"
  status            = "INACTIVE"
  dynamic_locations = ["US", "AF-BGL", "UA-26"]
  asns              = ["2232"]
}

resource "okta_network_zone" "dynamic_proxy_example" {
  name               = "testAcc_replace_with_uuid Dynamic Proxy Updated"
  type               = "DYNAMIC"
  status             = "INACTIVE"
  usage              = "POLICY"
  dynamic_proxy_type = "NotTorAnonymizer"
}

resource "okta_network_zone" "dynamic_v2_network_zone_example" {
  name                          = "testAcc_replace_with_uuid Dynamic V2 Updated"
  type                          = "DYNAMIC_V2"
  status                        = "ACTIVE"
  dynamic_locations_exclude     = ["BE-VAN", "CN-BJ"]
  ip_service_categories_include = ["SYMANTEC_VPN", "TRENDMICRO_VPN", "GOOGLE_VPN"]
}
