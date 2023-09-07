package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaNetworkZone_crud(t *testing.T) {
	mgr := newFixtureManager(networkZone, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.ip_network_zone_example", networkZone)
	dynamicResourceName := fmt.Sprintf("%s.dynamic_network_zone_example", networkZone)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(networkZone, doesNetworkZoneExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "type", "IP"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "proxies.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "gateways.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "usage", "POLICY"),
					resource.TestCheckResourceAttr(dynamicResourceName, "name", fmt.Sprintf("testAcc_%d Dynamic", mgr.Seed)),
					resource.TestCheckResourceAttr(dynamicResourceName, "type", "DYNAMIC"),
					resource.TestCheckResourceAttr(dynamicResourceName, "dynamic_locations.#", "2"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%d Updated", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "type", "IP"),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "proxies.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "gateways.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "usage", "BLOCKLIST"),
					resource.TestCheckResourceAttr(dynamicResourceName, "name", fmt.Sprintf("testAcc_%d Dynamic Updated", mgr.Seed)),
					resource.TestCheckResourceAttr(dynamicResourceName, "type", "DYNAMIC"),
					resource.TestCheckResourceAttr(dynamicResourceName, "dynamic_locations.#", "3"),
					resource.TestCheckResourceAttr(dynamicResourceName, "asns.#", "1"),
				),
			},
		},
	})
}

func doesNetworkZoneExist(id string) (bool, error) {
	client := sdkV2ClientForTest()
	_, response, err := client.NetworkZone.GetNetworkZone(context.Background(), id)
	return doesResourceExist(response, err)
}
