package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOktaThemes_read(t *testing.T) {
	mgr := newFixtureManager(themes, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
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
