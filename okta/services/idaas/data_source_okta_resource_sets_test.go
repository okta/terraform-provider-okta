package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

func TestAccDataSourceOktaResourceSets_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", "okta_resource_sets", t.Name())
	resourceName := "data.okta_resource_sets.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(mgr.GetFixtures("test_basic.tf", t)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					// Check that we have at least one resource set
					resource.TestCheckResourceAttrSet(resourceName, "resource_sets.#"),
					// Check that our expected resource set is present (by checking any resource set has a label)
					resource.TestCheckResourceAttrSet(resourceName, "resource_sets.0.label"),
				),
			},
		},
	})
}

func TestAccDataSourceOktaResourceSets_multiple(t *testing.T) {
	mgr := newFixtureManager("data-sources", "okta_resource_sets", t.Name())
	resourceName := "data.okta_resource_sets.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(mgr.GetFixtures("test_multiple.tf", t)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					// Check that we have at least two resource sets
					resource.TestCheckResourceAttrSet(resourceName, "resource_sets.#"),
					// Check that resource sets have expected attributes
					resource.TestCheckResourceAttrSet(resourceName, "resource_sets.0.label"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_sets.0.description"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_sets.1.label"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_sets.1.description"),
				),
			},
		},
	})
}
