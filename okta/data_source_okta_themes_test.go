package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOktaThemes_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(themes)
	config := mgr.GetFixtures("datasource.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_themes.test", "themes.#"),
					resource.TestCheckResourceAttr("data.okta_themes.test", "themes.#", "1"),
					resource.TestCheckResourceAttrSet("data.okta_themes.test", "themes.0.id"),
					resource.TestCheckResourceAttrSet("data.okta_themes.test", "themes.0.links"),
				),
			},
		},
	})
}
