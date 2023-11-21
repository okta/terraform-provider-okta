package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOktaBrands_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", brands, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_brands.test", "brands.#"),
					resource.TestCheckResourceAttrSet("data.okta_brands.test", "brands.0.id"),
					resource.TestCheckResourceAttrSet("data.okta_brands.test", "brands.0.name"),
					resource.TestCheckResourceAttrSet("data.okta_brands.test", "brands.0.links"),
				),
			},
		},
	})
}
