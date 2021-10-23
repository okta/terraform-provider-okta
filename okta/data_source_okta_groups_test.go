package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceGroups_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(groups)
	groups := mgr.GetFixtures("okta_groups.tf", ri, t)
	config := mgr.GetFixtures("datasource.tf", ri, t)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: groups,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_group.test_1", "id"),
					resource.TestCheckResourceAttrSet("okta_group.test_2", "id"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_groups.test", "id"),
					resource.TestCheckResourceAttr("data.okta_groups.test", "groups.#", "2"),
				),
			},
		},
	})
}
