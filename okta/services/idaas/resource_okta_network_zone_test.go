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

func TestAccResourceOktaNetworkZone_issue_2578(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSNetworkZone, t.Name())
	config := mgr.GetFixtures("basic_issue_2578.tf", t)
	updatedConfig := mgr.GetFixtures("basic_issue_2578_updated.tf", t)
	resourceName := fmt.Sprintf("%s.ip_network_zone_example", resources.OktaIDaaSNetworkZone)
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
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "proxies.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "gateways.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "usage", "POLICY"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "type", "IP"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "proxies.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "gateways.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "usage", "POLICY"),
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
