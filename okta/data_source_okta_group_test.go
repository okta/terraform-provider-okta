package okta

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOktaDataSourceGroup_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(oktaGroup)
	config := mgr.GetFixtures("datasource.tf", ri, t)
	configInvalid := mgr.GetFixtures("datasource_not_found.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
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
				ExpectError: regexp.MustCompile(`\bwas not found with type\b`),
			},
		},
	})
}
