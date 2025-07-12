package idaas_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func TestAccResourceOktaNetworkZone_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSNetworkZone, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.ip_network_zone_example", resources.OktaIDaaSNetworkZone)
	dynamicResourceName := fmt.Sprintf("%s.dynamic_network_zone_example", resources.OktaIDaaSNetworkZone)
	dynamicV2ResourceName := fmt.Sprintf("%s.dynamic_v2_network_zone_example", resources.OktaIDaaSNetworkZone)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSNetworkZone, doesNetworkZoneExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
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

func doesNetworkZoneExist(id string) (bool, error) {
	client := iDaaSAPIClientForTestUtil.OktaSDKClientV2()
	_, response, err := client.NetworkZone.GetNetworkZone(context.Background(), id)
	return utils.DoesResourceExist(response, err)
}

func TestAccResourceOktaNetworkZone_exempt_zone(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSNetworkZone, t.Name())
	resourceName := fmt.Sprintf("%s.exempt_zone_example", resources.OktaIDaaSNetworkZone)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSNetworkZone, doesNetworkZoneExist),
		Steps: []resource.TestStep{
			{
				Config: testOktaNetworkZoneConfig_exempt(mgr.Seed),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%d Exempt Zone", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "type", "IP"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "gateways.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "usage", "POLICY"),
					resource.TestCheckResourceAttr(resourceName, "use_as_exempt_list", "true"),
				),
			},
		},
	})
}

func testOktaNetworkZoneConfig_exempt(rInt int) string {
	return fmt.Sprintf(`
resource "okta_network_zone" "exempt_zone_example" {
  name               = "testAcc_%d Exempt Zone"
  type               = "IP"
  gateways           = ["1.2.3.4/32"]
  usage              = "POLICY"
  status             = "ACTIVE"
  use_as_exempt_list = true
}
`, rInt)
}

// TestAccResourceOktaNetworkZone_exempt_zone_update - Test update operations for exempt zones
// This test specifically validates the use_as_exempt_list functionality with updates
func TestAccResourceOktaNetworkZone_exempt_zone_update(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSNetworkZone, t.Name())
	resourceName := fmt.Sprintf("%s.exempt_zone_update_example", resources.OktaIDaaSNetworkZone)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSNetworkZone, doesNetworkZoneExist),
		Steps: []resource.TestStep{
			{
				Config: testOktaNetworkZoneConfig_exemptUpdate1(mgr.Seed),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%d Exempt Zone Update", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "type", "IP"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "gateways.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "usage", "POLICY"),
					resource.TestCheckResourceAttr(resourceName, "use_as_exempt_list", "true"),
				),
			},
			{
				Config: testOktaNetworkZoneConfig_exemptUpdate2(mgr.Seed),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%d Exempt Zone Update", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "type", "IP"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "gateways.#", "2"), // Updated to 2 gateways
					resource.TestCheckResourceAttr(resourceName, "usage", "POLICY"),
					resource.TestCheckResourceAttr(resourceName, "use_as_exempt_list", "true"),
				),
			},
		},
	})
}

// TestAccResourceOktaNetworkZone_exempt_validation - Test validation logic
func TestAccResourceOktaNetworkZone_exempt_validation(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSNetworkZone, t.Name())

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSNetworkZone, doesNetworkZoneExist),
		Steps: []resource.TestStep{
			{
				Config:      testOktaNetworkZoneConfig_exemptValidationFail(mgr.Seed),
				ExpectError: regexp.MustCompile("use_as_exempt_list can only be set to true for IP zones"),
			},
		},
	})
}

func testOktaNetworkZoneConfig_exemptUpdate1(rInt int) string {
	return fmt.Sprintf(`
resource "okta_network_zone" "exempt_zone_update_example" {
  name               = "testAcc_%d Exempt Zone Update"
  type               = "IP"
  gateways           = ["10.1.0.0/24"]
  usage              = "POLICY"
  status             = "ACTIVE"
  use_as_exempt_list = true
}
`, rInt)
}

func testOktaNetworkZoneConfig_exemptUpdate2(rInt int) string {
	return fmt.Sprintf(`
resource "okta_network_zone" "exempt_zone_update_example" {
  name               = "testAcc_%d Exempt Zone Update"
  type               = "IP"
  gateways           = ["10.1.0.0/24", "192.168.100.0/24"]
  usage              = "POLICY"
  status             = "ACTIVE"
  use_as_exempt_list = true
}
`, rInt)
}

func testOktaNetworkZoneConfig_exemptValidationFail(rInt int) string {
	return fmt.Sprintf(`
resource "okta_network_zone" "exempt_zone_validation_fail" {
  name               = "testAcc_%d Exempt Validation Fail"
  type               = "DYNAMIC"
  dynamic_locations  = ["US", "CA"]
  usage              = "POLICY"
  status             = "ACTIVE"
  use_as_exempt_list = true
}
`, rInt)
}
