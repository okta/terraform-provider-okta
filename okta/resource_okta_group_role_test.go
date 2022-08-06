package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaGroupAdminRole_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", groupRole)
	resourceName2 := fmt.Sprintf("%s.test_app", groupRole)
	mgr := newFixtureManager(groupRole, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	groupTarget := mgr.GetFixtures("group_targets.tf", t)
	groupTargetsUpdated := mgr.GetFixtures("group_targets_updated.tf", t)
	groupTargetsRemoved := mgr.GetFixtures("group_targets_removed.tf", t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(group, doesGroupExist),
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
