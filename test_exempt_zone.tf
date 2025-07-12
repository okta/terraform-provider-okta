terraform {
  required_providers {
    okta = {
      source = "github.com/your-username/terraform-provider-okta"
      # Or use your fork's URL
    }
  }
}

provider "okta" {
  org_name  = "your-org-name"
  base_url  = "okta.com"
  api_token = "your-api-token"
}

resource "okta_network_zone" "test_exempt_zone" {
  name               = "TestExemptZone"
  type               = "IP"
  gateways           = ["192.168.1.0/24"]
  usage              = "POLICY"
  use_as_exempt_list = true
}