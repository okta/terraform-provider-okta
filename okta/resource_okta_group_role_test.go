package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaGroupAdminRole_crud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.test", groupRole)
	resourceName2 := fmt.Sprintf("%s.test_app", groupRole)
	mgr := newFixtureManager(groupRole)
	config := mgr.GetFixtures("basic.tf", ri, t)
	groupTarget := mgr.GetFixtures("group_targets.tf", ri, t)
	groupTargetsUpdated := mgr.GetFixtures("group_targets_updated.tf", ri, t)
	groupTargetsRemoved := mgr.GetFixtures("group_targets_removed.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(oktaGroup, doesGroupExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "role_type", "READ_ONLY_ADMIN"),
					resource.TestCheckResourceAttr(resourceName2, "role_type", "APP_ADMIN"),
					resource.TestCheckResourceAttr(resourceName2, "target_app_list.#", "0"),
				),
			},
			{
				Config: groupTarget,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "role_type", "HELP_DESK_ADMIN"),
					resource.TestCheckResourceAttr(resourceName2, "role_type", "APP_ADMIN"),
					resource.TestCheckResourceAttr(resourceName2, "target_app_list.#", "1"),
				),
			},
			{
				Config: groupTargetsUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "role_type", "HELP_DESK_ADMIN"),
					resource.TestCheckResourceAttr(resourceName2, "role_type", "APP_ADMIN"),
					resource.TestCheckResourceAttr(resourceName2, "target_app_list.#", "1"),
				),
			},
			{
				Config: groupTargetsRemoved,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "role_type", "HELP_DESK_ADMIN"),
					resource.TestCheckResourceAttr(resourceName2, "role_type", "APP_ADMIN"),
					resource.TestCheckResourceAttr(resourceName2, "target_app_list.#", "0"),
				),
			},
		},
	})
}
