package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceOktaNetworkZone_crud(t *testing.T) {
	mgr := newFixtureManager("resources", networkZone, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.ip_network_zone_example", networkZone)
	dynamicResourceName := fmt.Sprintf("%s.dynamic_network_zone_example", networkZone)
	dynamicV2ResourceName := fmt.Sprintf("%s.dynamic_v2_network_zone_example", networkZone)

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
					resource.TestCheckResourceAttr(dynamicV2ResourceName, "name", fmt.Sprintf("testAcc_%d Dynamic V2", mgr.Seed)),
					resource.TestCheckResourceAttr(dynamicV2ResourceName, "type", "DYNAMIC_V2"),
					resource.TestCheckResourceAttr(dynamicV2ResourceName, "dynamic_locations.#", "2"),
					resource.TestCheckNoResourceAttr(dynamicV2ResourceName, "dynamic_locations_exclude.#"),
					resource.TestCheckResourceAttr(dynamicV2ResourceName, "ip_service_categories_include.#", "1"),
					resource.TestCheckResourceAttr(dynamicV2ResourceName, "ip_service_categories_exclude.#", "1"),
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
					resource.TestCheckResourceAttr(dynamicV2ResourceName, "name", fmt.Sprintf("testAcc_%d Dynamic V2 Updated", mgr.Seed)),
					resource.TestCheckResourceAttr(dynamicV2ResourceName, "type", "DYNAMIC_V2"),
					resource.TestCheckNoResourceAttr(dynamicV2ResourceName, "dynamic_locations.#"),
					resource.TestCheckResourceAttr(dynamicV2ResourceName, "dynamic_locations_exclude.#", "2"),
					resource.TestCheckResourceAttr(dynamicV2ResourceName, "ip_service_categories_include.#", "3"),
					resource.TestCheckNoResourceAttr(dynamicV2ResourceName, "ip_service_categories_exclude.#"),
				),
			},
		},
	})
}

func TestAccResourceOktaNetworkZone_ActivateDefaultNetworkZone(t *testing.T) {
	resourceName := fmt.Sprintf("%s.default_enhanced_network_zone_example", networkZone)

	mgr := newFixtureManager("resources", networkZone, t.Name())
	importConfig := mgr.ConfigReplace(`
resource "okta_network_zone" "default_enhanced_dynamic_zone_example" {
  name                          = "DefaultEnhancedDynamicZone"
  type                          = "DYNAMIC"
  usage                         = "BLOCKLIST"
  ip_service_categories_include = ["ALL_ANONYMIZERS"]
  status                        = "INACTIVE"
}`)

	updateConfig := mgr.ConfigReplace(`
resource "okta_network_zone" "default_enhanced_dynamic_zone_example" {
  name                          = "DefaultEnhancedDynamicZone"
  type                          = "DYNAMIC"
  usage                         = "BLOCKLIST"
  ip_service_categories_include = ["ALL_ANONYMIZERS"]
  status                        = "ACTIVE"
}`)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(networkZone, doesNetworkZoneExist),
		Steps: []resource.TestStep{
			{
				ResourceName:      resourceName,
				Config:            importConfig,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					fmt.Println(s.RootModule())
					if !ok {
						return "", fmt.Errorf("failed to find %s", resourceName)
					}

					return rs.Primary.ID, nil
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOktaNetworkZoneExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
				),
			},
			{
				ResourceName: resourceName,
				Config:       updateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOktaNetworkZoneExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
		},
	})
}

func testAccCheckOktaNetworkZoneExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		return nil
	}
}

func doesNetworkZoneExist(id string) (bool, error) {
	client := sdkV2ClientForTest()
	_, response, err := client.NetworkZone.GetNetworkZone(context.Background(), id)
	return doesResourceExist(response, err)
}
