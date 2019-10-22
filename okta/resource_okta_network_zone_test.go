package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOktaNetworkZone_crud(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(networkZone)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.ip_network_zone_example", networkZone)
	dynamicResourceName := fmt.Sprintf("%s.dynamic_network_zone_example", networkZone)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(networkZone, doesNetworkZoneExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "type", "IP"),
					resource.TestCheckResourceAttr(resourceName, "proxies.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "gateways.#", "2"),

					resource.TestCheckResourceAttr(dynamicResourceName, "name", fmt.Sprintf("testAcc_%d Dynamic", ri)),
					resource.TestCheckResourceAttr(dynamicResourceName, "type", "DYNAMIC"),
					resource.TestCheckResourceAttr(dynamicResourceName, "dynamic_locations.#", "2"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%d Updated", ri)),
					resource.TestCheckResourceAttr(resourceName, "type", "IP"),
					resource.TestCheckResourceAttr(resourceName, "proxies.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "gateways.#", "2"),

					resource.TestCheckResourceAttr(dynamicResourceName, "name", fmt.Sprintf("testAcc_%d Dynamic Updated", ri)),
					resource.TestCheckResourceAttr(dynamicResourceName, "type", "DYNAMIC"),
					resource.TestCheckResourceAttr(dynamicResourceName, "dynamic_locations.#", "2"),
				),
			},
		},
	})
}

func doesNetworkZoneExist(id string) (bool, error) {
	client := getSupplementFromMetadata(testAccProvider.Meta())
	_, response, err := client.GetNetworkZone(id)

	return doesResourceExist(response, err)
}
