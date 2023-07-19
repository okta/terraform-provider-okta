package okta

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceGroup_read(t *testing.T) {
	mgr := newFixtureManager(group, t.Name())
	groupCreate := mgr.GetFixtures("okta_group.tf", t)
	config := mgr.GetFixtures("datasource.tf", t)
	configInvalid := mgr.GetFixtures("datasource_not_found.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
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
					resource.TestCheckResourceAttr("data.okta_group.test", "users.#", "1"),
				),
			},
			{
				Config:      configInvalid,
				ExpectError: regexp.MustCompile(`\bdoes not exist`),
			},
		},
	})
}
