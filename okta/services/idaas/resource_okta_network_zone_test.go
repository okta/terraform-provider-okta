package idaas_test

import (
	"context"
	"fmt"
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

func TestAccResourceOktaNetworkZone_ipNormalization(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSNetworkZone, t.Name())
	config := mgr.GetFixtures("ip_normalization.tf", t)
	updatedConfig := mgr.GetFixtures("ip_normalization_updated.tf", t)

	// Define resource names following the convention from the original test
	singleIPResourceName := fmt.Sprintf("%s.ip_network_zone_single", resources.OktaIDaaSNetworkZone)
	cidrResourceName := fmt.Sprintf("%s.ip_network_zone_cidr", resources.OktaIDaaSNetworkZone)
	rangeResourceName := fmt.Sprintf("%s.ip_network_zone_range", resources.OktaIDaaSNetworkZone)
	changingSingleResourceName := fmt.Sprintf("%s.ip_network_zone_changing_single", resources.OktaIDaaSNetworkZone)
	changingCIDRResourceName := fmt.Sprintf("%s.ip_network_zone_changing_cidr", resources.OktaIDaaSNetworkZone)
	changingRangeResourceName := fmt.Sprintf("%s.ip_network_zone_changing_range", resources.OktaIDaaSNetworkZone)
	mixedResourceName := fmt.Sprintf("%s.ip_network_zone_mixed", resources.OktaIDaaSNetworkZone)
	unchangedSingleResourceName := fmt.Sprintf("%s.ip_network_zone_unchanged_single", resources.OktaIDaaSNetworkZone)
	unchangedCIDRResourceName := fmt.Sprintf("%s.ip_network_zone_unchanged_cidr", resources.OktaIDaaSNetworkZone)
	unchangedRangeResourceName := fmt.Sprintf("%s.ip_network_zone_unchanged_range", resources.OktaIDaaSNetworkZone)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSNetworkZone, doesNetworkZoneExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					// Check equivalent notation resources
					resource.TestCheckResourceAttr(singleIPResourceName, "name", acctest.BuildResourceName(mgr.Seed)+" Single"),
					resource.TestCheckResourceAttr(singleIPResourceName, "type", "IP"),
					resource.TestCheckResourceAttr(singleIPResourceName, "gateways.#", "1"),

					resource.TestCheckResourceAttr(cidrResourceName, "name", acctest.BuildResourceName(mgr.Seed)+" CIDR"),
					resource.TestCheckResourceAttr(cidrResourceName, "type", "IP"),
					resource.TestCheckResourceAttr(cidrResourceName, "gateways.#", "1"),

					resource.TestCheckResourceAttr(rangeResourceName, "name", acctest.BuildResourceName(mgr.Seed)+" Range"),
					resource.TestCheckResourceAttr(rangeResourceName, "type", "IP"),
					resource.TestCheckResourceAttr(rangeResourceName, "gateways.#", "1"),

					// Check changing resources
					resource.TestCheckResourceAttr(changingSingleResourceName, "name", acctest.BuildResourceName(mgr.Seed)+" Changing Single"),
					resource.TestCheckResourceAttr(changingSingleResourceName, "gateways.#", "1"),

					resource.TestCheckResourceAttr(changingCIDRResourceName, "name", acctest.BuildResourceName(mgr.Seed)+" Changing CIDR"),
					resource.TestCheckResourceAttr(changingCIDRResourceName, "gateways.#", "1"),

					resource.TestCheckResourceAttr(changingRangeResourceName, "name", acctest.BuildResourceName(mgr.Seed)+" Changing Range"),
					resource.TestCheckResourceAttr(changingRangeResourceName, "gateways.#", "1"),

					// Check mixed notation resource
					resource.TestCheckResourceAttr(mixedResourceName, "name", acctest.BuildResourceName(mgr.Seed)+" Mixed"),
					resource.TestCheckResourceAttr(mixedResourceName, "gateways.#", "3"),

					// Check unchanged resources
					resource.TestCheckResourceAttr(unchangedSingleResourceName, "name", acctest.BuildResourceName(mgr.Seed)+" Unchanged Single"),
					resource.TestCheckResourceAttr(unchangedSingleResourceName, "gateways.#", "1"),
					resource.TestCheckResourceAttr(unchangedSingleResourceName, "gateways.0", "192.168.2.1"),

					resource.TestCheckResourceAttr(unchangedCIDRResourceName, "name", acctest.BuildResourceName(mgr.Seed)+" Unchanged CIDR"),
					resource.TestCheckResourceAttr(unchangedCIDRResourceName, "gateways.#", "1"),
					resource.TestCheckResourceAttr(unchangedCIDRResourceName, "gateways.0", "10.1.0.0/24"),

					resource.TestCheckResourceAttr(unchangedRangeResourceName, "name", acctest.BuildResourceName(mgr.Seed)+" Unchanged Range"),
					resource.TestCheckResourceAttr(unchangedRangeResourceName, "gateways.#", "1"),
					resource.TestCheckResourceAttr(unchangedRangeResourceName, "gateways.0", "172.17.0.1-172.17.0.10"),
				),
			},
			{
				Config:             updatedConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true, // We expect changes for the non-equivalent updates
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					// Verify unchanged resources maintain their original values
					resource.TestCheckResourceAttr(unchangedSingleResourceName, "gateways.#", "1"),
					resource.TestCheckResourceAttr(unchangedSingleResourceName, "gateways.0", "192.168.2.1"),
					resource.TestCheckResourceAttr(unchangedCIDRResourceName, "gateways.#", "1"),
					resource.TestCheckResourceAttr(unchangedCIDRResourceName, "gateways.0", "10.1.0.0/24"),
					resource.TestCheckResourceAttr(unchangedRangeResourceName, "gateways.#", "1"),
					resource.TestCheckResourceAttr(unchangedRangeResourceName, "gateways.0", "172.17.0.1-172.17.0.10"),

					// Verify changing resources have new values
					resource.TestCheckResourceAttr(changingSingleResourceName, "gateways.#", "1"),
					resource.TestCheckResourceAttr(changingSingleResourceName, "gateways.0", "192.168.1.2"),
					resource.TestCheckResourceAttr(changingCIDRResourceName, "gateways.#", "1"),
					resource.TestCheckResourceAttr(changingCIDRResourceName, "gateways.0", "10.0.0.0/16"),
					resource.TestCheckResourceAttr(changingRangeResourceName, "gateways.#", "1"),
					resource.TestCheckResourceAttr(changingRangeResourceName, "gateways.0", "172.16.0.1-172.16.0.20"),
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
