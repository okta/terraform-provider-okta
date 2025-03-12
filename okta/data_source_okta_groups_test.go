package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOktaGroups_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", groups, t.Name())
	// groups := mgr.GetFixtures("okta_groups.tf", t)
	config := mgr.GetFixtures("datasource.tf", t)
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				// Config: groups,
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_group.test_1", "id"),
					resource.TestCheckResourceAttrSet("okta_group.test_2", "id"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_groups.okta_groups", "id"),
					resource.TestCheckResourceAttr("data.okta_groups.okta_groups", "groups.#", "2"),
					// the example enumeration doesn't match anything so as a string the output will be a blank string
					// resource.TestCheckOutput("special_groups", ""),
					resource.TestCheckResourceAttrSet("data.okta_groups.built_in_groups", "id"),
					resource.TestCheckResourceAttr("data.okta_groups.built_in_groups", "groups.#", "2"),
				),
			},
		},
	})
}
