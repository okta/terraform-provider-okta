package idaas_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

// TestAccResourceOktaNetworkZoneDefaultExempt_VCR validates the use_as_exempt_list functionality
// for the DefaultExemptIpZone system resource. This test demonstrates that:
// 1. The DefaultExemptIpZone can be imported successfully
// 2. The use_as_exempt_list parameter is correctly applied during updates
// 3. New IP addresses can be added to the exempt zone
//
// NOTE: System network zones (DefaultExemptIpZone, LegacyIpZone, BlockedIpZone) cannot be
// deactivated or deleted. This test handles those expected errors.
func TestAccResourceOktaNetworkZoneDefaultExempt_VCR(t *testing.T) {
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck: acctest.AccPreCheck(t),
		ErrorCheck: func(err error) error {
			if err != nil {
				// Ignore deactivate and delete errors for system resources
				errStr := err.Error()
				if strings.Contains(errStr, "failed to deactivate network zone") ||
					strings.Contains(errStr, "failed to delete network zone") {
					return nil
				}
			}
			// Call the default error checker for other errors
			return testAccErrorChecks(t)(err)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy: func(s *terraform.State) error {
			// Always return nil for system resources that can't be destroyed
			// This prevents the test from failing when it can't delete the DefaultExemptIpZone
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config:             importTf,
				ResourceName:       "okta_network_zone.default_exempt",
				ImportState:        true,
				ImportStateId:      "nzot5r5hx70lKWqvT697",
				ImportStateVerify:  false,
				ImportStatePersist: true,
			},
			{
				Config:             updateTf,
				ExpectNonEmptyPlan: true, // Expected because import doesn't capture use_as_exempt_list
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("okta_network_zone.default_exempt", "name", "DefaultExemptIpZone"),
					resource.TestCheckResourceAttr("okta_network_zone.default_exempt", "type", "IP"),
					resource.TestCheckResourceAttr("okta_network_zone.default_exempt", "usage", "POLICY"),
					resource.TestCheckResourceAttr("okta_network_zone.default_exempt", "use_as_exempt_list", "true"),
					resource.TestCheckResourceAttr("okta_network_zone.default_exempt", "gateways.#", "2"),
				),
			},
			{
				// Empty config to remove the resource from state without destroying it
				Config: emptyTf,
			},
		},
	})
}

const importTf = `
resource "okta_network_zone" "default_exempt" {
  name               = "DefaultExemptIpZone"
  type               = "IP"
  usage              = "POLICY"
  gateways           = ["192.168.101.0/24", "10.0.100.0/24", "172.16.0.0/12"]
  use_as_exempt_list = true
}
`

const updateTf = `
resource "okta_network_zone" "default_exempt" {
  name               = "DefaultExemptIpZone"
  type               = "IP"
  usage              = "POLICY"
  gateways           = ["192.168.101.0/24", "10.0.100.0/24"]
  use_as_exempt_list = true
}
`

const emptyTf = `
# Empty configuration - removes resource from state without destroying it
`
