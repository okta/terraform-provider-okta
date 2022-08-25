package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaUserGroupMemberships_crud(t *testing.T) {
	ri := acctest.RandInt()

	mgr := newFixtureManager(userGroupMemberships)
	start := mgr.GetFixtures("basic.tf", ri, t)
	update := mgr.GetFixtures("basic_update.tf", ri, t)
	remove := mgr.GetFixtures("basic_removal.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: start,
			},
			{
				Config: update,
			},
			{
				Config: remove,
			},
		},
	})
}
