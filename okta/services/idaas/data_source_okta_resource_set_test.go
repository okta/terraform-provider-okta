package idaas_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

// TODO: V6 BREAKING CHANGE - These tests will need to be updated when the datasource is split
// into separate datasources in V6. The tests will need to test both the basic resource set
// datasource and the resource set resources datasource.

func TestAccDataSourceOktaResourceSet_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", "okta_resource_set", t.Name())
	resourceName := "data.okta_resource_set.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: mgr.GetFixtures("test_basic.tf", t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "label"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					// Verify the resource contains the expected pattern
					resource.TestMatchResourceAttr(resourceName, "resources.0", regexp.MustCompile(`/api/v1/users`)),
					// Verify that resources_orn is not set when using resources
					resource.TestCheckResourceAttr(resourceName, "resources_orn.#", "0"),
				),
			},
		},
	})
}

// TestAccDataSourceOktaResourceSet_readWithResources tests data source with various resource configurations
func TestAccDataSourceOktaResourceSet_readWithResources(t *testing.T) {
	mgr := newFixtureManager("data-sources", "okta_resource_set", t.Name())
	resourceName := "data.okta_resource_set.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: mgr.GetFixtures("test_read_with_resources_step1.tf", t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "label"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "2"),
					// Verify that the resources contain valid API endpoints
					resource.TestCheckResourceAttrSet(resourceName, "resources.0"),
					resource.TestCheckResourceAttrSet(resourceName, "resources.1"),
					// Verify the resources contain the expected patterns
					resource.TestMatchResourceAttr(resourceName, "resources.0", regexp.MustCompile(`/api/v1/groups/`)),
					resource.TestMatchResourceAttr(resourceName, "resources.1", regexp.MustCompile(`/api/v1/apps/`)),
					// Verify that resources_orn is not set when using resources
					resource.TestCheckResourceAttr(resourceName, "resources_orn.#", "0"),
				),
			},
			{
				Config: mgr.GetFixtures("test_read_with_resources_step2.tf", t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "3"),
					// Verify that the resources contain valid API endpoints
					resource.TestCheckResourceAttrSet(resourceName, "resources.0"),
					resource.TestCheckResourceAttrSet(resourceName, "resources.1"),
					resource.TestCheckResourceAttrSet(resourceName, "resources.2"),
					// Verify the resources contain the expected patterns
					resource.TestMatchResourceAttr(resourceName, "resources.0", regexp.MustCompile(`/api/v1/groups/`)),
					resource.TestMatchResourceAttr(resourceName, "resources.1", regexp.MustCompile(`/api/v1/apps/`)),
					resource.TestMatchResourceAttr(resourceName, "resources.2", regexp.MustCompile(`/api/v1/users`)),
				),
			},
			{
				Config: mgr.GetFixtures("test_read_with_resources_step3.tf", t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "label"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					// Verify that both resources and resources_orn are empty
					resource.TestCheckResourceAttr(resourceName, "resources.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "resources_orn.#", "0"),
				),
			},
		},
	})
}

// TestAccDataSourceOktaResourceSet_readWithORNs tests data source with ORN references
func TestAccDataSourceOktaResourceSet_readWithORNs(t *testing.T) {
	mgr := newFixtureManager("data-sources", "okta_resource_set", t.Name())
	resourceName := "data.okta_resource_set.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: mgr.GetFixtures("test_read_with_orns_step1.tf", t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "label"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttr(resourceName, "resources_orn.#", "2"),
					// Verify that the ORNs contain the expected values
					resource.TestCheckResourceAttrSet(resourceName, "resources_orn.0"),
					resource.TestCheckResourceAttrSet(resourceName, "resources_orn.1"),
					// Verify the ORNs contain the expected patterns
					resource.TestMatchResourceAttr(resourceName, "resources_orn.0", regexp.MustCompile(`orn:okta:directory:.*:users`)),
					resource.TestMatchResourceAttr(resourceName, "resources_orn.1", regexp.MustCompile(`orn:okta:directory:.*:groups`)),
					// Verify that resources is not set when using resources_orn
					resource.TestCheckResourceAttr(resourceName, "resources.#", "0"),
				),
			},
			{
				Config: mgr.GetFixtures("test_read_with_orns_step2.tf", t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "resources_orn.#", "3"),
					// Verify that the ORNs contain the expected values
					resource.TestCheckResourceAttrSet(resourceName, "resources_orn.0"),
					resource.TestCheckResourceAttrSet(resourceName, "resources_orn.1"),
					resource.TestCheckResourceAttrSet(resourceName, "resources_orn.2"),
					// Verify the ORNs contain the expected patterns
					resource.TestMatchResourceAttr(resourceName, "resources_orn.0", regexp.MustCompile(`orn:okta:directory:.*:groups`)),
					resource.TestMatchResourceAttr(resourceName, "resources_orn.1", regexp.MustCompile(`orn:okta:directory:.*:apps`)),
					resource.TestMatchResourceAttr(resourceName, "resources_orn.2", regexp.MustCompile(`orn:okta:directory:.*:users`)),
				),
			},
		},
	})
}

// TestAccDataSourceOktaResourceSet_readResourceChanges tests that data source properly reflects resource changes
func TestAccDataSourceOktaResourceSet_readResourceChanges(t *testing.T) {
	mgr := newFixtureManager("data-sources", "okta_resource_set", t.Name())
	resourceName := "data.okta_resource_set.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: mgr.GetFixtures("test_resource_changes_step1.tf", t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestMatchResourceAttr(resourceName, "resources.0", regexp.MustCompile(`/api/v1/users`)),
				),
			},
			{
				Config: mgr.GetFixtures("test_resource_changes_step2.tf", t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "2"),
					resource.TestMatchResourceAttr(resourceName, "resources.0", regexp.MustCompile(`/api/v1/groups/`)),
					resource.TestMatchResourceAttr(resourceName, "resources.1", regexp.MustCompile(`/api/v1/users`)),
				),
			},
		},
	})
}
