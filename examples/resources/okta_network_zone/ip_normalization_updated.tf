// Single IP resource - should not trigger update despite API format
resource "okta_network_zone" "ip_network_zone_single" {
  name     = "testAcc_replace_with_uuid Single"
  type     = "IP"
  gateways = ["192.168.1.1"] // Should not trigger update even though API returns as "192.168.1.1-192.168.1.1"
  usage    = "POLICY"
  status   = "ACTIVE"
}

// CIDR resource
resource "okta_network_zone" "ip_network_zone_cidr" {
  name     = "testAcc_replace_with_uuid CIDR"
  type     = "IP"
  gateways = ["192.168.1.0/24"]
  usage    = "POLICY"
  status   = "ACTIVE"
}

// Range resource
resource "okta_network_zone" "ip_network_zone_range" {
  name     = "testAcc_replace_with_uuid Range"
  type     = "IP"
  gateways = ["192.168.1.1-192.168.1.10"]
  usage    = "POLICY"
  status   = "ACTIVE"
}

// Resource with single IP that changes
resource "okta_network_zone" "ip_network_zone_changing_single" {
  name     = "testAcc_replace_with_uuid Changing Single"
  type     = "IP"
  gateways = ["192.168.1.2"] // Changed from "192.168.1.1"
  usage    = "POLICY"
  status   = "ACTIVE"
}

// Resource with CIDR that changes
resource "okta_network_zone" "ip_network_zone_changing_cidr" {
  name     = "testAcc_replace_with_uuid Changing CIDR"
  type     = "IP"
  gateways = ["10.0.0.0/16"] // Changed from "192.168.0.0/24"
  usage    = "POLICY"
  status   = "ACTIVE"
}

// Resource with range that changes
resource "okta_network_zone" "ip_network_zone_changing_range" {
  name     = "testAcc_replace_with_uuid Changing Range"
  type     = "IP"
  gateways = ["172.16.0.1-172.16.0.20"] // Changed end IP from .10 to .20
  usage    = "POLICY"
  status   = "ACTIVE"
}

// Resource with mixed notation types
resource "okta_network_zone" "ip_network_zone_mixed" {
  name = "testAcc_replace_with_uuid Mixed"
  type = "IP"
  gateways = [
    "192.168.1.1",           // Single IP - should not trigger update despite API format
    "10.0.0.0/24",           // CIDR
    "172.16.0.1-172.16.0.10" // Range
  ]
  usage  = "POLICY"
  status = "ACTIVE"
}

// Resources that should remain unchanged
resource "okta_network_zone" "ip_network_zone_unchanged_single" {
  name     = "testAcc_replace_with_uuid Unchanged Single"
  type     = "IP"
  gateways = ["192.168.2.1"] // Should not trigger update despite API format
  usage    = "POLICY"
  status   = "ACTIVE"
}

resource "okta_network_zone" "ip_network_zone_unchanged_cidr" {
  name     = "testAcc_replace_with_uuid Unchanged CIDR"
  type     = "IP"
  gateways = ["10.1.0.0/24"]
  usage    = "POLICY"
  status   = "ACTIVE"
}

resource "okta_network_zone" "ip_network_zone_unchanged_range" {
  name     = "testAcc_replace_with_uuid Unchanged Range"
  type     = "IP"
  gateways = ["172.17.0.1-172.17.0.10"]
  usage    = "POLICY"
  status   = "ACTIVE"
}
