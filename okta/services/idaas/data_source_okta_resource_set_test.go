package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

// NOTE: This datasource now only returns basic metadata (id, label, description, created, last_updated).
// Resources are retrieved using the separate okta_resource_set_resources datasource.

func TestAccDataSourceOktaResourceSet_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", "okta_resource_set", t.Name())
	resourceName := "data.okta_resource_set.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(mgr.GetFixtures("test_basic.tf", t)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "label"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
				),
			},
		},
	})
}

// TestAccDataSourceOktaResourceSet_readWithResources tests data source with various resource configurations
// NOTE: This test now only verifies basic metadata since resources are handled by the separate datasource
func TestAccDataSourceOktaResourceSet_readWithResources(t *testing.T) {
	mgr := newFixtureManager("data-sources", "okta_resource_set", t.Name())
	resourceName := "data.okta_resource_set.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(mgr.GetFixtures("test_read_with_resources_step1.tf", t)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "label"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
				),
			},
			{
				Config: mgr.ConfigReplace(mgr.GetFixtures("test_read_with_resources_step2.tf", t)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "label"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
				),
			},
			{
				Config: mgr.ConfigReplace(mgr.GetFixtures("test_read_with_resources_step3.tf", t)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "label"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
				),
			},
		},
	})
}

// TestAccDataSourceOktaResourceSet_readWithORNs tests data source with ORN references
// NOTE: This test now only verifies basic metadata since resources are handled by the separate datasource
func TestAccDataSourceOktaResourceSet_readWithORNs(t *testing.T) {
	mgr := newFixtureManager("data-sources", "okta_resource_set", t.Name())
	resourceName := "data.okta_resource_set.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(mgr.GetFixtures("test_read_with_orns_step1.tf", t)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "label"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
				),
			},
			{
				Config: mgr.ConfigReplace(mgr.GetFixtures("test_read_with_orns_step2.tf", t)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "label"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
				),
			},
		},
	})
}

// TestAccDataSourceOktaResourceSet_readResourceChanges tests that data source properly reflects resource changes
// NOTE: This test now only verifies basic metadata since resources are handled by the separate datasource
func TestAccDataSourceOktaResourceSet_readResourceChanges(t *testing.T) {
	mgr := newFixtureManager("data-sources", "okta_resource_set", t.Name())
	resourceName := "data.okta_resource_set.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(mgr.GetFixtures("test_resource_changes_step1.tf", t)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "label"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
				),
			},
			{
				Config: mgr.ConfigReplace(mgr.GetFixtures("test_resource_changes_step2.tf", t)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "label"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
				),
			},
		},
	})
}
