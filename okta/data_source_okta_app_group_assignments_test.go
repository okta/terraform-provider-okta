package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceAppGroupAssignments_read(t *testing.T) {
	mgr := newFixtureManager(appGroupAssignments, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_app_group_assignments.test", "groups.#"),
				),
			},
		},
	})
}
