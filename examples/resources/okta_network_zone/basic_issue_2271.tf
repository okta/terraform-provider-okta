resource "okta_network_zone" "default" {
  name                          = "DefaultExemptIpZone"
  type                          = "IP"
  set_usage_as_exempt_list = true
  usage                         = "POLICY"
  status                        = "ACTIVE"
  gateways = [
            "1.2.3.4/32",    # Cloudflare
            "4.5.6.7/32",    # Cloudflare
          ]
  proxies = []
}