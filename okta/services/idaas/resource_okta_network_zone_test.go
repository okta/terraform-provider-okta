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

// TestAccResourceOktaNetworkZone_exempt_zone_update - Test for use_as_exempt_list functionality
// This test validates issue #2271 fix by demonstrating the use_as_exempt_list parameter works correctly.
// 
// NOTE: This test cannot directly test updating the DefaultExemptIpZone because:
// 1. It's a system resource that cannot be created/deleted via Terraform
// 2. Standard Terraform acceptance tests expect full lifecycle management
// 
// To manually test the DefaultExemptIpZone update functionality:
// 1. Import the zone: terraform import okta_network_zone.default nzot5r5hx70lKWqvT697
// 2. Configure with use_as_exempt_list=true and update the gateways
// 3. Run terraform apply - this will use our custom HTTP function with useAsExemptList
//
// This test validates:
// 1. Regular network zones work correctly without use_as_exempt_list
// 2. The use_as_exempt_list parameter is accepted and stored correctly
// 3. Zones can be created/updated with use_as_exempt_list=false
func TestAccResourceOktaNetworkZone_exempt_zone_update(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSNetworkZone, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSNetworkZone)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSNetworkZone, doesNetworkZoneExist),
		Steps: []resource.TestStep{
			{
				// Step 1: Create a regular zone without use_as_exempt_list
				Config: testOktaNetworkZoneConfig_regularZone(mgr.Seed),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%d Network Zone", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "type", "IP"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "gateways.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "usage", "POLICY"),
					resource.TestCheckNoResourceAttr(resourceName, "use_as_exempt_list"),
				),
			},
			{
				// Step 2: Update the regular zone (still without use_as_exempt_list)
				Config: testOktaNetworkZoneConfig_regularZoneUpdated(mgr.Seed),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%d Network Zone Updated", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "type", "IP"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "gateways.#", "2"), // Added gateway
					resource.TestCheckResourceAttr(resourceName, "usage", "POLICY"),
					resource.TestCheckNoResourceAttr(resourceName, "use_as_exempt_list"),
				),
			},
			{
				// Step 3: Create a zone with use_as_exempt_list=false (explicitly set)
				Config: testOktaNetworkZoneConfig_withExemptListFalse(mgr.Seed),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%d With Exempt False", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "type", "IP"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "gateways.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "usage", "POLICY"),
					resource.TestCheckResourceAttr(resourceName, "use_as_exempt_list", "false"),
				),
			},
			{
				// Step 4: Update the zone (keeping use_as_exempt_list=false)
				Config: testOktaNetworkZoneConfig_withExemptListFalseUpdated(mgr.Seed),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%d With Exempt False", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "type", "IP"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "gateways.#", "2"), // Added gateway
					resource.TestCheckResourceAttr(resourceName, "usage", "POLICY"),
					resource.TestCheckResourceAttr(resourceName, "use_as_exempt_list", "false"),
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


func testOktaNetworkZoneConfig_regularZone(rInt int) string {
	return fmt.Sprintf(`
resource "okta_network_zone" "test" {
  name     = "testAcc_%d Network Zone"
  type     = "IP"
  gateways = ["10.70.0.0/24"]
  usage    = "POLICY"
  status   = "ACTIVE"
}
`, rInt)
}

func testOktaNetworkZoneConfig_regularZoneUpdated(rInt int) string {
	return fmt.Sprintf(`
resource "okta_network_zone" "test" {
  name     = "testAcc_%d Network Zone Updated"
  type     = "IP"
  gateways = ["10.70.0.0/24", "10.80.0.0/24"]
  usage    = "POLICY"
  status   = "ACTIVE"
}
`, rInt)
}

func testOktaNetworkZoneConfig_withExemptListFalse(rInt int) string {
	return fmt.Sprintf(`
resource "okta_network_zone" "test" {
  name               = "testAcc_%d With Exempt False"
  type               = "IP"
  gateways           = ["10.90.0.0/24"]
  usage              = "POLICY"
  status             = "ACTIVE"
  use_as_exempt_list = false
}
`, rInt)
}

func testOktaNetworkZoneConfig_withExemptListFalseUpdated(rInt int) string {
	return fmt.Sprintf(`
resource "okta_network_zone" "test" {
  name               = "testAcc_%d With Exempt False"
  type               = "IP"
  gateways           = ["10.90.0.0/24", "10.100.0.0/24"]
  usage              = "POLICY"
  status             = "ACTIVE"
  use_as_exempt_list = false
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
