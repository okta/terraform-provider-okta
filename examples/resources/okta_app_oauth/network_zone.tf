### Example: restricting an app to specific network zones

data "okta_network_zone" "example" {
  name = "LegacyIpZone"
}

resource "okta_app_oauth" "network_zone_example" {
  label          = "example_zone_restricted"
  type           = "web"
  grant_types    = ["authorization_code"]
  redirect_uris  = ["https://example.com/callback"]
  response_types = ["code"]

  network {
    connection = "ZONE"
    include    = [data.okta_network_zone.example.id]
    exclude    = []
  }
}
