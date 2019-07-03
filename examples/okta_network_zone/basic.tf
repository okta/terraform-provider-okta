resource "okta_network_zone" "ip_network_zone_example" {
  name		    = "IP Network Zone Example"
  type		    = "IP"
  gateways      = [ "1.2.3.4/24", "2.3.4.5-2.3.4.15" ]
  proxies       = [ "2.2.3.4/24", "3.3.4.5-3.3.4.15" ]
}

resource "okta_network_zone" "dynamic_network_zone_example" {
  name		        = "Dynamic Network Zone Example"
  type		        = "DYNAMIC"
  dynamic_locations = [ "US", "AF-BGL" ]
}