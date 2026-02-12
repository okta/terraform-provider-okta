package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

func TestAccDataSourceOktaResourceSetResources_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", "okta_resource_set_resources", t.Name())
	resourceName := "data.okta_resource_set_resources.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(mgr.GetFixtures("test_basic.tf", t)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "resources.0.id"),
				),
			},
		},
	})
}

func TestAccDataSourceOktaResourceSetResources_multiple(t *testing.T) {
	mgr := newFixtureManager("data-sources", "okta_resource_set_resources", t.Name())
	resourceName := "data.okta_resource_set_resources.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(mgr.GetFixtures("test_multiple.tf", t)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "2"),
					// Verify first resource
					resource.TestCheckResourceAttrSet(resourceName, "resources.0.id"),
					// Verify second resource
					resource.TestCheckResourceAttrSet(resourceName, "resources.1.id"),
				),
			},
		},
	})
}
