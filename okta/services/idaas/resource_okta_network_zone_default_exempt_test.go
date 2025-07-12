package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

func TestAccResourceOktaNetworkZoneDefaultExempt_VCR(t *testing.T) {
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config:            importTf,
				ResourceName:      "okta_network_zone.default_exempt",
				ImportState:       true,
				ImportStateId:     "nzot5r5hx70lKWqvT697",
				ImportStateVerify: false,
			},
			{
				Config: updateTf,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("okta_network_zone.default_exempt", "name", "DefaultExemptIpZone"),
					resource.TestCheckResourceAttr("okta_network_zone.default_exempt", "type", "IP"),
					resource.TestCheckResourceAttr("okta_network_zone.default_exempt", "usage", "POLICY"),
					resource.TestCheckResourceAttr("okta_network_zone.default_exempt", "use_as_exempt_list", "true"),
					resource.TestCheckResourceAttr("okta_network_zone.default_exempt", "gateways.#", "3"),
				),
			},
		},
	})
}

const importTf = `
resource "okta_network_zone" "default_exempt" {
  name               = "DefaultExemptIpZone"
  type               = "IP"
  usage              = "POLICY"
  gateways           = ["192.168.101.0/24", "10.0.100.0/24"]
  use_as_exempt_list = true
}
`

const updateTf = `
resource "okta_network_zone" "default_exempt" {
  name               = "DefaultExemptIpZone"
  type               = "IP"
  usage              = "POLICY"
  gateways           = ["192.168.101.0/24", "10.0.100.0/24", "172.16.0.0/12"]
  use_as_exempt_list = true
}
`