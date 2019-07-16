package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccGroupAdminRoles(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.test", groupRoles)
	mgr := newFixtureManager(groupRoles)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("all_roles.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "admin_roles.#", "1"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "admin_roles.#", "8"),
				),
			},
		},
	})
}
