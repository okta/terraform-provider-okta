package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOktaBrands_read(t *testing.T) {
	mgr := newFixtureManager(brands, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_brands.test", "brands.#"),
					resource.TestCheckResourceAttr("data.okta_brands.test", "brands.#", "1"),
					resource.TestCheckResourceAttrSet("data.okta_brands.test", "brands.0.id"),
					resource.TestCheckResourceAttrSet("data.okta_brands.test", "brands.0.links"),
				),
			},
		},
	})
}
