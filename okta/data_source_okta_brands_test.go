package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOktaBrands_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(brands)
	config := mgr.GetFixtures("datasource.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
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
