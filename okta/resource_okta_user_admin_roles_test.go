package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaUserAdminRoles_crud(t *testing.T) {
	ri := acctest.RandInt()

	mgr := newFixtureManager(userAdminRoles)
	start := mgr.GetFixtures("basic.tf", ri, t)
	update := mgr.GetFixtures("basic_update.tf", ri, t)
	remove := mgr.GetFixtures("basic_removal.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", userAdminRoles)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: start,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "admin_roles.#", "2"),
				),
			},
			{
				Config: update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "admin_roles.#", "3"),
				),
			},
			{
				Config: remove,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "admin_roles.#", "1"),
				),
			},
		},
	})
}
