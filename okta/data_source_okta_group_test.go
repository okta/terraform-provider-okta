package okta

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceGroup_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(oktaGroup)
	groupCreate := mgr.GetFixtures("okta_group.tf", ri, t)
	config := mgr.GetFixtures("datasource.tf", ri, t)
	configInvalid := mgr.GetFixtures("datasource_not_found.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: groupCreate,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_group.test", "id"),
					resource.TestCheckResourceAttrSet("data.okta_group.test", "type"),
					resource.TestCheckResourceAttrSet("okta_group.test", "id"),
					resource.TestCheckResourceAttr("okta_group.test", "users.#", "1"),
				),
			},
			{
				Config:      configInvalid,
				ExpectError: regexp.MustCompile(`\bdoes not exist`),
			},
		},
	})
}
