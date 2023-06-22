resource "okta_network_zone" "ip_network_zone_example" {
  name     = "testAcc_replace_with_uuid"
  type     = "IP"
  status   = "ACTIVE"
  gateways = ["1.2.3.4/24", "2.3.4.5-2.3.4.15"]
  proxies  = ["2.2.3.4/24", "3.3.4.5-3.3.4.15"]
}

# resource "okta_network_zone" "dynamic_network_zone_example" {
#   name              = "testAcc_replace_with_uuid Dynamic"
#   type              = "DYNAMIC"
#   status            = "ACTIVE"
#   dynamic_locations = ["US", "AF-BGL"]
# }

# resource "okta_network_zone" "dynamic_proxy_example" {
#   name               = "testAcc_replace_with_uuid Dynamic Proxy"
#   type               = "DYNAMIC"
#   status             = "ACTIVE"
#   usage              = "BLOCKLIST"
#   dynamic_proxy_type = "TorAnonymizer"
# }
