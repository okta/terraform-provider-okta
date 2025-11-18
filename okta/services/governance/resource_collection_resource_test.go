package governance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

func TestAccCollectionResource_basic(t *testing.T) {
	t.Skip("Skipping Collection Resource tests - requires Okta Governance license and configured app")
	mgr := newFixtureManager("resources", "okta_collection_resource", t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := "okta_collection_resource.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "collection_id"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_id"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_orn"),
					resource.TestCheckResourceAttr(resourceName, "entitlements.#", "1"),
				),
			},
		},
	})
}

func TestAccCollectionResource_update(t *testing.T) {
	t.Skip("Skipping Collection Resource tests - requires Okta Governance license and configured app")
	mgr := newFixtureManager("resources", "okta_collection_resource", t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	configUpdated := mgr.GetFixtures("updated.tf", t)
	resourceName := "okta_collection_resource.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "entitlements.#", "1"),
				),
			},
			{
				Config: configUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "entitlements.#", "2"),
				),
			},
		},
	})
}

func TestAccCollectionResourceDataSource_basic(t *testing.T) {
	t.Skip("Skipping Collection Resource data source tests - requires Okta Governance license")
	mgr := newFixtureManager("data-sources", "okta_collection_resource", t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	dataSourceName := "data.okta_collection_resource.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "collection_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_orn"),
				),
			},
		},
	})
}
